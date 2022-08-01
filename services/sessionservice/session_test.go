package sessionservice

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/ditrit/badaas/httperrors"
	configurationmocks "github.com/ditrit/badaas/mocks/configuration"
	repositorymocks "github.com/ditrit/badaas/mocks/persistence/repository"
	"github.com/ditrit/badaas/persistence/models"
	"github.com/ditrit/badaas/persistence/pagination"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestNewSession(t *testing.T) {
	sessionInstance := newSession(uuid.Nil, time.Second)
	assert.NotNil(t, sessionInstance)
	assert.Equal(t, uuid.Nil, sessionInstance.UserID)
}

func TestLogInUser(t *testing.T) {
	sessionRepositoryMock, service, logs, sessionConfigurationMock := setupTest(t)
	sessionRepositoryMock.On("Create", mock.Anything).Return(nil)
	sessionConfigurationMock.On("GetSessionDuration").Return(time.Minute)
	response := httptest.NewRecorder()
	user := &models.User{
		Username: "bob",
		Email:    "bob@email.com",
	}
	err := service.LogUserIn(user, response)
	require.NoError(t, err)
	assert.Len(t, service.cache, 1)
	assert.Equal(t, 1, logs.Len())
	log := logs.All()[0]
	assert.Equal(t, "Added session", log.Message)
	require.Len(t, log.Context, 1)
}

// make values for test
func setupTest(
	t *testing.T,
) (
	*repositorymocks.CRUDRepository[models.Session, uuid.UUID],
	*sessionServiceImpl,
	*observer.ObservedLogs,
	*configurationmocks.SessionConfiguration,
) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	sessionRepositoryMock := repositorymocks.NewCRUDRepository[models.Session, uuid.UUID](t)
	sessionConfiguration := configurationmocks.NewSessionConfiguration(t)
	service := &sessionServiceImpl{
		sessionRepository:    sessionRepositoryMock,
		logger:               logger,
		cache:                make(map[uuid.UUID]*models.Session),
		sessionConfiguration: sessionConfiguration,
	}

	return sessionRepositoryMock, service, logs, sessionConfiguration
}

func TestLogInUserDbError(t *testing.T) {
	sessionRepositoryMock, service, logs, sessionConfigurationMock := setupTest(t)
	sessionRepositoryMock.On("Create", mock.Anything).Return(httperrors.NewInternalServerError("db err", "nil", nil))
	sessionConfigurationMock.On("GetSessionDuration").Return(time.Minute)

	response := httptest.NewRecorder()
	user := &models.User{
		Username: "bob",
		Email:    "bob@email.com",
	}
	err := service.LogUserIn(user, response)
	require.Error(t, err)
	assert.Len(t, service.cache, 0)
	assert.Equal(t, 0, logs.Len())
}

func TestIsValid(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.On("Create", mock.Anything).Return(nil)
	uuidSample := uuid.New()
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuidSample,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	err := service.add(session)
	require.NoError(t, err)
	assert.Len(t, service.cache, 1)
	assert.Equal(t, uuid.Nil, service.cache[uuidSample].UserID)
	isValid, claims := service.IsValid(uuidSample)
	require.True(t, isValid)
	assert.Equal(t, *claims, SessionClaims{
		UserID:      uuid.Nil,
		SessionUUID: uuidSample,
	})
}

func TestIsValid_SessionNotFound(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.
		On("Find", mock.Anything, mock.Anything, mock.Anything).
		Return(pagination.NewPage([]*models.Session{}, 0, 125, 1236), nil)
	uuidSample := uuid.New()
	isValid, _ := service.IsValid(uuidSample)
	require.False(t, isValid)
	//
}

func TestLogOutUser(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.On("Delete", mock.Anything).Return(nil)
	response := httptest.NewRecorder()
	uuidSample := uuid.New()
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuidSample,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	service.cache[uuidSample] = session
	err := service.LogUserOut(makeSessionClaims(session), response)
	require.NoError(t, err)
	assert.Len(t, service.cache, 0)
}

func TestLogOutUserDbError(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.On("Delete", mock.Anything).Return(httperrors.NewInternalServerError("db errors", "oh we failed to delete the session", nil))
	response := httptest.NewRecorder()
	uuidSample := uuid.New()
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuidSample,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	service.cache[uuidSample] = session
	err := service.LogUserOut(makeSessionClaims(session), response)
	require.Error(t, err)
	assert.Len(t, service.cache, 1)
}

func TestLogOutUser_SessionNotFound(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.
		On("Find", mock.Anything, nil, nil).
		Return(nil, httperrors.NewInternalServerError("db errors", "oh we failed to delete the session", nil))
	response := httptest.NewRecorder()
	uuidSample := uuid.New()
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuid.Nil,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	service.cache[uuidSample] = session
	sessionClaims := makeSessionClaims(session)
	sessionClaims.SessionUUID = uuid.Nil
	err := service.LogUserOut(sessionClaims, response)
	require.Error(t, err)
	assert.Len(t, service.cache, 1)
}

func TestRollSession(t *testing.T) {
	sessionRepositoryMock, service, _, sessionConfigurationMock := setupTest(t)
	sessionRepositoryMock.On("Save", mock.Anything).Return(nil)
	sessionDuration := time.Minute
	sessionConfigurationMock.On("GetSessionDuration").Return(sessionDuration)
	sessionConfigurationMock.On("GetRollDuration").Return(sessionDuration / 4)
	uuidSample := uuid.New()
	originalExpirationTime := time.Now().Add(sessionDuration / 5)
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuid.Nil,
		},
		UserID:    uuid.Nil,
		ExpiresAt: originalExpirationTime,
	}
	service.cache[uuidSample] = session
	err := service.RollSession(uuidSample)
	require.NoError(t, err)
	assert.Greater(t, session.ExpiresAt, originalExpirationTime)
}

func TestRollSession_Expired(t *testing.T) {
	_, service, _, sessionConfigurationMock := setupTest(t)
	sessionDuration := time.Minute
	sessionConfigurationMock.On("GetSessionDuration").Return(sessionDuration)
	sessionConfigurationMock.On("GetRollDuration").Return(sessionDuration / 4)
	uuidSample := uuid.New()
	originalExpirationTime := time.Now().Add(-time.Hour)
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuidSample,
		},
		UserID:    uuid.Nil,
		ExpiresAt: originalExpirationTime,
	}
	service.cache[uuidSample] = session
	err := service.RollSession(uuidSample)
	require.Error(t, err)
}

func TestRollSession_falseUUID(t *testing.T) {
	repoSession, service, _, sessionConfigurationMock := setupTest(t)
	sessionDuration := time.Minute
	sessionConfigurationMock.On("GetSessionDuration").Return(sessionDuration)
	sessionConfigurationMock.On("GetRollDuration").Return(sessionDuration / 4)

	uuidSample := uuid.New()
	originalExpirationTime := time.Now().Add(-time.Hour)
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuid.Nil,
		},
		UserID:    uuid.Nil,
		ExpiresAt: originalExpirationTime,
	}
	service.cache[uuidSample] = session
	repoSession.On("Find", mock.Anything, nil, nil).Return(pagination.NewPage([]*models.Session{}, 0, 2, 5), nil)
	err := service.RollSession(uuid.New())
	require.NoError(t, err)
}

func TestRollSession_sessionNotFound(t *testing.T) {
	sessionRepositoryMock, service, _, sessionConfigurationMock := setupTest(t)
	sessionRepositoryMock.
		On("Find", squirrel.Eq{"uuid": "00000000-0000-0000-0000-000000000000"}, nil, nil).
		Return(
			pagination.NewPage([]*models.Session{}, 0, 10, 0), nil)

	sessionDuration := time.Minute
	sessionConfigurationMock.On("GetSessionDuration").Return(sessionDuration)
	sessionConfigurationMock.On("GetRollDuration").Return(sessionDuration)

	err := service.RollSession(uuid.Nil)
	require.NoError(t, err)
}

func Test_pullFromDB(t *testing.T) {
	sessionRepositoryMock, service, logs, _ := setupTest(t)
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuid.Nil,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	sessionRepositoryMock.On("GetAll", nil).Return([]*models.Session{session}, nil)

	service.pullFromDB()
	assert.Len(t, service.cache, 1)
	assert.Equal(t, 1, logs.Len())
	log := logs.All()[0]
	assert.Equal(t, "Pulled sessions from DB", log.Message)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "sessionCount", Type: zapcore.Int64Type, Integer: 1},
	}, log.Context)
}

func Test_pullFromDB_repoError(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.On("GetAll", nil).Return(nil, httperrors.AnError)
	assert.PanicsWithError(t, httperrors.AnError.Error(), func() { service.pullFromDB() })
}

func Test_removeExpired(t *testing.T) {
	sessionRepositoryMock, service, logs, _ := setupTest(t)
	uuidSample := uuid.New()
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuid.Nil,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	sessionRepositoryMock.
		On("Delete", session).
		Return(nil)
	service.cache[uuidSample] = session

	service.removeExpired()
	assert.Len(t, service.cache, 0)
	assert.Equal(t, 1, logs.Len())
	log := logs.All()[0]
	assert.Equal(t, "Removed expired session", log.Message)
	assert.ElementsMatch(t, []zap.Field{
		{Key: "expiredSessionCount", Type: zapcore.Int64Type, Integer: 1},
	}, log.Context)
}

func Test_removeExpired_RepositoryError(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	uuidSample := uuid.New()
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuid.Nil,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	sessionRepositoryMock.
		On("Delete", session).
		Return(httperrors.AnError)
	service.cache[uuidSample] = session

	assert.Panics(t, func() { service.removeExpired() })
}

func Test_get(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	uuidSample := uuid.New()
	session := &models.Session{
		BaseModel: models.BaseModel{
			ID: uuid.Nil,
		},
		UserID:    uuid.Nil,
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	sessionRepositoryMock.
		On("Find", mock.Anything, nil, nil).
		Return(pagination.NewPage([]*models.Session{session}, 0, 12, 13), nil)

	sessionFound := service.get(uuidSample)
	assert.Equal(t, sessionFound, session)
}
