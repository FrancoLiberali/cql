package sessionservice

import (
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
)

// Errors
var (
	HERRSessionExpired = httperrors.NewUnauthorizedError(
		"session error",
		"session is expired",
	)
)

// SessionService handle sessions
type SessionService interface {
	IsValid(sessionUUID model.UUID) (bool, *SessionClaims)
	// TODO services should not work with httperrors
	RollSession(model.UUID) httperrors.HTTPError
	LogUserIn(user *models.User) (*models.Session, error)
	LogUserOut(sessionClaims *SessionClaims) httperrors.HTTPError
}

// Check interface compliance
var _ SessionService = (*sessionServiceImpl)(nil)

// The SessionService concrete interface
type sessionServiceImpl struct {
	sessionRepository    repository.CRUD[models.Session, model.UUID]
	cache                map[model.UUID]*models.Session
	mutex                sync.Mutex
	logger               *zap.Logger
	sessionConfiguration configuration.SessionConfiguration
	db                   *gorm.DB
}

// The SessionService constructor
func NewSessionService(
	logger *zap.Logger,
	sessionRepository repository.CRUD[models.Session, model.UUID],
	sessionConfiguration configuration.SessionConfiguration,
	db *gorm.DB,
) SessionService {
	sessionService := &sessionServiceImpl{
		cache:                make(map[model.UUID]*models.Session),
		logger:               logger,
		sessionRepository:    sessionRepository,
		sessionConfiguration: sessionConfiguration,
		db:                   db,
	}
	sessionService.init()

	return sessionService
}

// Return true if the session exists and is still valid.
// A instance of SessionClaims is returned to be added to the request context if the conditions previously mentioned are met.
func (sessionService *sessionServiceImpl) IsValid(sessionUUID model.UUID) (bool, *SessionClaims) {
	sessionInstance := sessionService.get(sessionUUID)
	if sessionInstance == nil {
		return false, nil
	}

	return true, makeSessionClaims(sessionInstance)
}

// Get a session from cache
// return nil if not found
func (sessionService *sessionServiceImpl) get(sessionUUID model.UUID) *models.Session {
	sessionService.mutex.Lock()
	defer sessionService.mutex.Unlock()

	session, ok := sessionService.cache[sessionUUID]
	if ok {
		return session
	}

	session, err := sessionService.sessionRepository.GetByID(
		sessionService.db,
		sessionUUID,
	)
	if err != nil {
		return nil
	}

	return session
}

// Add a session to the cache
func (sessionService *sessionServiceImpl) add(session *models.Session) error {
	sessionService.mutex.Lock()
	defer sessionService.mutex.Unlock()

	err := sessionService.sessionRepository.Create(sessionService.db, session)
	if err != nil {
		return err
	}

	sessionService.cache[session.ID] = session
	sessionService.logger.Debug("Added session", zap.String("uuid", session.ID.String()))

	return nil
}

// Initialize the session service
func (sessionService *sessionServiceImpl) init() {
	sessionService.cache = make(map[model.UUID]*models.Session)

	go func() {
		for {
			sessionService.removeExpired()
			sessionService.pullFromDB()
			time.Sleep(
				sessionService.sessionConfiguration.GetPullInterval(),
			)
		}
	}()
}

// Get all sessions and save them in cache
func (sessionService *sessionServiceImpl) pullFromDB() {
	sessionService.mutex.Lock()
	defer sessionService.mutex.Unlock()

	sessionsFromDatabase, err := sessionService.sessionRepository.Find(sessionService.db)
	if err != nil {
		panic(err)
	}

	newSessionCache := make(map[model.UUID]*models.Session)
	for _, sessionFromDatabase := range sessionsFromDatabase {
		newSessionCache[sessionFromDatabase.ID] = sessionFromDatabase
	}

	sessionService.cache = newSessionCache
	sessionService.logger.Debug(
		"Pulled sessions from DB",
		zap.Int("sessionCount", len(sessionsFromDatabase)),
	)
}

// Remove the expired session
func (sessionService *sessionServiceImpl) removeExpired() {
	sessionService.mutex.Lock()
	defer sessionService.mutex.Unlock()

	var i int

	for sessionUUID, session := range sessionService.cache {
		if session.IsExpired() {
			// Delete the session in the database
			err := sessionService.sessionRepository.Delete(sessionService.db, session)
			if err != nil {
				panic(err)
			}
			// if the deletion of the session in the database was successful,
			// we now remove the session from the cache.
			// see https://pkg.go.dev/builtin#delete
			delete(sessionService.cache, sessionUUID)
			i++
		}
	}

	sessionService.logger.Debug(
		"Removed expired session",
		zap.Int("expiredSessionCount", i),
	)
}

// Delete a session
func (sessionService *sessionServiceImpl) delete(session *models.Session) httperrors.HTTPError {
	sessionService.mutex.Lock()
	defer sessionService.mutex.Unlock()

	sessionUUID := session.ID

	err := sessionService.sessionRepository.Delete(sessionService.db, session)
	if err != nil {
		return httperrors.NewInternalServerError(
			"session error",
			fmt.Sprintf("failed to delete session %q (userId=%d)", sessionUUID, session.UserID),
			err,
		)
	}

	delete(sessionService.cache, sessionUUID)

	return nil
}

// Roll a session. If the session is close to expiration, extend its duration.
func (sessionService *sessionServiceImpl) RollSession(sessionUUID model.UUID) httperrors.HTTPError {
	rollInterval := sessionService.sessionConfiguration.GetRollDuration()
	sessionDuration := sessionService.sessionConfiguration.GetSessionDuration()

	session := sessionService.get(sessionUUID)
	if session == nil {
		// no session to roll, no error
		return nil
	}

	if session.IsExpired() {
		return HERRSessionExpired
	}

	if session.CanBeRolled(rollInterval) {
		sessionService.mutex.Lock()
		defer sessionService.mutex.Unlock()

		session.ExpiresAt = session.ExpiresAt.Add(sessionDuration)

		err := sessionService.sessionRepository.Save(sessionService.db, session)
		if err != nil {
			return httperrors.NewDBError(err)
		}

		sessionService.logger.Warn("Rolled session",
			zap.String("userID", session.UserID.String()),
			zap.String("sessionID", session.ID.String()))
	}

	return nil
}

// Log in a user
func (sessionService *sessionServiceImpl) LogUserIn(user *models.User) (*models.Session, error) {
	sessionDuration := sessionService.sessionConfiguration.GetSessionDuration()
	session := models.NewSession(user.ID, sessionDuration)

	err := sessionService.add(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Log out a user.
func (sessionService *sessionServiceImpl) LogUserOut(sessionClaims *SessionClaims) httperrors.HTTPError {
	session := sessionService.get(sessionClaims.SessionUUID)
	if session == nil {
		return httperrors.NewUnauthorizedError("Authentication Error", "not authenticated")
	}

	err := sessionService.delete(session)
	if err != nil {
		return err
	}

	return nil
}
