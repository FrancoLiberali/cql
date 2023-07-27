package sessionservice

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/configuration"
	"github.com/ditrit/badaas/httperrors"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/repository"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
	IsValid(sessionUUID uuid.UUID) (bool, *SessionClaims)
	RollSession(uuid.UUID) httperrors.HTTPError
	LogUserIn(user *models.User, response http.ResponseWriter) httperrors.HTTPError
	LogUserOut(sessionClaims *SessionClaims, response http.ResponseWriter) httperrors.HTTPError
}

// Check interface compliance
var _ SessionService = (*sessionServiceImpl)(nil)

// The SessionService concrete interface
type sessionServiceImpl struct {
	sessionRepository    repository.CRUDRepository[models.Session, uuid.UUID]
	cache                map[uuid.UUID]*models.Session
	mutex                sync.Mutex
	logger               *zap.Logger
	sessionConfiguration configuration.SessionConfiguration
}

// The SessionService constructor
func NewSessionService(
	logger *zap.Logger,
	sessionRepository repository.CRUDRepository[models.Session, uuid.UUID],
	sessionConfiguration configuration.SessionConfiguration,
) SessionService {
	sessionService := &sessionServiceImpl{
		cache:                make(map[uuid.UUID]*models.Session),
		logger:               logger,
		sessionRepository:    sessionRepository,
		sessionConfiguration: sessionConfiguration,
	}
	sessionService.init()
	return sessionService
}

// Create a new session
func newSession(userID uuid.UUID, sessionDuration time.Duration) *models.Session {
	return &models.Session{
		UserID:    userID,
		ExpiresAt: time.Now().Add(sessionDuration),
	}
}

// Return true if the session exists and is still valid.
// A instance of SessionClaims is returned to be added to the request context if the conditions previously mentioned are met.
func (sessionService *sessionServiceImpl) IsValid(sessionUUID uuid.UUID) (bool, *SessionClaims) {
	sessionInstance := sessionService.get(sessionUUID)
	if sessionInstance == nil {
		return false, nil
	}
	return true, makeSessionClaims(sessionInstance)
}

// Get a session from cache
// return nil if not found
func (sessionService *sessionServiceImpl) get(sessionUUID uuid.UUID) *models.Session {
	sessionService.mutex.Lock()
	defer sessionService.mutex.Unlock()
	session, ok := sessionService.cache[sessionUUID]
	if ok {
		return session
	}
	sessionsFoundWithUUID, databaseError := sessionService.sessionRepository.Find(squirrel.Eq{"uuid": sessionUUID.String()}, nil, nil)
	if databaseError != nil {
		return nil
	}
	if !sessionsFoundWithUUID.HasContent {
		return nil // no sessions found in database
	}
	return sessionsFoundWithUUID.Ressources[0]
}

// Add a session to the cache
func (sessionService *sessionServiceImpl) add(session *models.Session) httperrors.HTTPError {
	sessionService.mutex.Lock()
	defer sessionService.mutex.Unlock()
	herr := sessionService.sessionRepository.Create(session)
	if herr != nil {
		return herr
	}
	sessionService.cache[session.ID] = session
	sessionService.logger.Debug("Added session", zap.String("uuid", session.ID.String()))
	return nil
}

// Initialize the session service
func (sessionService *sessionServiceImpl) init() {
	sessionService.cache = make(map[uuid.UUID]*models.Session)
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
	sessionsFromDatabase, err := sessionService.sessionRepository.GetAll(nil)
	if err != nil {
		panic(err)
	}
	newSessionCache := make(map[uuid.UUID]*models.Session)
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
			err := sessionService.sessionRepository.Delete(session)
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
	err := sessionService.sessionRepository.Delete(session)
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
func (sessionService *sessionServiceImpl) RollSession(sessionUUID uuid.UUID) httperrors.HTTPError {
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
		herr := sessionService.sessionRepository.Save(session)
		if herr != nil {
			return herr
		}
		sessionService.logger.Warn("Rolled session",
			zap.String("userID", session.UserID.String()),
			zap.String("sessionID", session.ID.String()))
	}
	return nil
}

// Log in a user
func (sessionService *sessionServiceImpl) LogUserIn(user *models.User, response http.ResponseWriter) httperrors.HTTPError {
	sessionDuration := sessionService.sessionConfiguration.GetSessionDuration()
	session := newSession(user.ID, sessionDuration)
	err := sessionService.add(session)
	if err != nil {
		return err
	}
	CreateAndSetAccessTokenCookie(response, session.ID.String())
	return nil
}

// Log out a user.
func (sessionService *sessionServiceImpl) LogUserOut(sessionClaims *SessionClaims, response http.ResponseWriter) httperrors.HTTPError {
	session := sessionService.get(sessionClaims.SessionUUID)
	if session == nil {
		return httperrors.NewUnauthorizedError("Authentication Error", "not authenticated")
	}
	err := sessionService.delete(session)
	if err != nil {
		return err
	}
	CreateAndSetAccessTokenCookie(response, "")
	return nil
}

// Create an access token and send it in a cookie
func CreateAndSetAccessTokenCookie(w http.ResponseWriter, sessionUUID string) {
	accessToken := &http.Cookie{
		Name:     "access_token",
		Path:     "/",
		Value:    sessionUUID,
		HttpOnly: true,
		SameSite: http.SameSiteNoneMode, // TODO change to http.SameSiteStrictMode in prod
		Secure:   false,                 // TODO change to true in prod
		Expires:  time.Now().Add(48 * time.Hour),
	}
	err := accessToken.Valid()
	if err != nil {
		panic(err)
	}
	http.SetCookie(w, accessToken)
}
