package sessionservice

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
	"gorm.io/gorm"

	"github.com/ditrit/badaas/httperrors"
	configurationMocks "github.com/ditrit/badaas/mocks/configuration"
	repositoryMocks "github.com/ditrit/badaas/mocks/persistence/repository"
	"github.com/ditrit/badaas/orm/model"
	"github.com/ditrit/badaas/persistence/models"
)

var gormDB *gorm.DB

// make values for test
func setupTest(
	t *testing.T,
) (
	*repositoryMocks.CRUD[models.Session, model.UUID],
	*sessionServiceImpl,
	*observer.ObservedLogs,
	*configurationMocks.SessionConfiguration,
) {
	core, logs := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	sessionRepositoryMock := repositoryMocks.NewCRUD[models.Session, model.UUID](t)
	sessionConfiguration := configurationMocks.NewSessionConfiguration(t)
	service := &sessionServiceImpl{
		sessionRepository:    sessionRepositoryMock,
		logger:               logger,
		cache:                make(map[model.UUID]*models.Session),
		sessionConfiguration: sessionConfiguration,
		db:                   gormDB,
	}

	return sessionRepositoryMock, service, logs, sessionConfiguration
}

func TestLogInUser(t *testing.T) {
	sessionRepositoryMock, service, logs, sessionConfigurationMock := setupTest(t)
	sessionRepositoryMock.On("Create", gormDB, mock.Anything).Return(nil)

	sessionConfigurationMock.On("GetSessionDuration").Return(time.Minute)

	user := &models.User{
		Username: "bob",
		Email:    "bob@email.com",
	}
	_, err := service.LogUserIn(user)
	require.NoError(t, err)
	assert.Len(t, service.cache, 1)
	assert.Equal(t, 1, logs.Len())
	log := logs.All()[0]
	assert.Equal(t, "Added session", log.Message)
	require.Len(t, log.Context, 1)
}

func TestLogInUserDbError(t *testing.T) {
	sessionRepositoryMock, service, logs, sessionConfigurationMock := setupTest(t)
	sessionRepositoryMock.
		On("Create", gormDB, mock.Anything).
		Return(errors.New("db err"))

	sessionConfigurationMock.On("GetSessionDuration").Return(time.Minute)

	user := &models.User{
		Username: "bob",
		Email:    "bob@email.com",
	}
	_, err := service.LogUserIn(user)
	require.Error(t, err)
	assert.Len(t, service.cache, 0)
	assert.Equal(t, 0, logs.Len())
}

func TestIsValid(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.On("Create", gormDB, mock.Anything).Return(nil)

	uuidSample := model.NewUUID()
	session := &models.Session{
		UUIDModel: model.UUIDModel{
			ID: uuidSample,
		},
		UserID:    model.NilUUID,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	err := service.add(session)
	require.NoError(t, err)
	assert.Len(t, service.cache, 1)
	assert.Equal(t, model.NilUUID, service.cache[uuidSample].UserID)
	isValid, claims := service.IsValid(uuidSample)
	require.True(t, isValid)
	assert.Equal(t, *claims, SessionClaims{
		UserID:      model.NilUUID,
		SessionUUID: uuidSample,
	})
}

func TestIsValid_SessionNotFound(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.
		On("GetByID", gormDB, mock.Anything).
		Return(nil, errors.New("not-found"))

	uuidSample := model.NewUUID()
	isValid, _ := service.IsValid(uuidSample)
	require.False(t, isValid)
}

func TestLogOutUser(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.On("Delete", gormDB, mock.Anything).Return(nil)

	uuidSample := model.NewUUID()
	session := &models.Session{
		UUIDModel: model.UUIDModel{
			ID: uuidSample,
		},
		UserID:    model.NilUUID,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	service.cache[uuidSample] = session
	err := service.LogUserOut(makeSessionClaims(session))
	require.NoError(t, err)
	assert.Len(t, service.cache, 0)
}

func TestLogOutUserDbError(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.
		On("Delete", gormDB, mock.Anything).
		Return(errors.New("db errors"))

	uuidSample := model.NewUUID()

	session := &models.Session{
		UUIDModel: model.UUIDModel{
			ID: uuidSample,
		},
		UserID:    model.NilUUID,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	service.cache[uuidSample] = session
	err := service.LogUserOut(makeSessionClaims(session))
	require.Error(t, err)
	assert.Len(t, service.cache, 1)
}

func TestLogOutUser_SessionNotFound(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	sessionRepositoryMock.
		On("GetByID", gormDB, mock.Anything).
		Return(nil, errors.New("not-found"))

	uuidSample := model.NewUUID()
	session := &models.Session{
		UUIDModel: model.UUIDModel{
			ID: model.NilUUID,
		},
		UserID:    model.NilUUID,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	service.cache[uuidSample] = session
	sessionClaims := makeSessionClaims(session)
	sessionClaims.SessionUUID = model.NilUUID
	err := service.LogUserOut(sessionClaims)
	require.Error(t, err)
	assert.Len(t, service.cache, 1)
}

func TestRollSession(t *testing.T) {
	sessionRepositoryMock, service, _, sessionConfigurationMock := setupTest(t)
	sessionRepositoryMock.On("Save", gormDB, mock.Anything).Return(nil)

	sessionDuration := time.Minute
	sessionConfigurationMock.On("GetSessionDuration").Return(sessionDuration)
	sessionConfigurationMock.On("GetRollDuration").Return(sessionDuration / 4)

	uuidSample := model.NewUUID()
	originalExpirationTime := time.Now().Add(sessionDuration / 5)
	session := &models.Session{
		UUIDModel: model.UUIDModel{
			ID: model.NilUUID,
		},
		UserID:    model.NilUUID,
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

	uuidSample := model.NewUUID()
	originalExpirationTime := time.Now().Add(-time.Hour)
	session := &models.Session{
		UUIDModel: model.UUIDModel{
			ID: uuidSample,
		},
		UserID:    model.NilUUID,
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

	uuidSample := model.NewUUID()
	originalExpirationTime := time.Now().Add(-time.Hour)
	session := &models.Session{
		UUIDModel: model.UUIDModel{
			ID: model.NilUUID,
		},
		UserID:    model.NilUUID,
		ExpiresAt: originalExpirationTime,
	}
	service.cache[uuidSample] = session

	repoSession.
		On("GetByID", gormDB, mock.Anything).
		Return(nil, errors.New("not-found"))

	err := service.RollSession(model.NewUUID())
	require.NoError(t, err)
}

func TestRollSession_sessionNotFound(t *testing.T) {
	sessionRepositoryMock, service, _, sessionConfigurationMock := setupTest(t)
	sessionRepositoryMock.
		On("GetByID", gormDB, model.NilUUID).
		Return(nil, errors.New("not-found"))

	sessionDuration := time.Minute
	sessionConfigurationMock.On("GetSessionDuration").Return(sessionDuration)
	sessionConfigurationMock.On("GetRollDuration").Return(sessionDuration)

	err := service.RollSession(model.NilUUID)
	require.NoError(t, err)
}

func Test_pullFromDB(t *testing.T) {
	sessionRepositoryMock, service, logs, _ := setupTest(t)
	session := &models.Session{
		UUIDModel: model.UUIDModel{
			ID: model.NilUUID,
		},
		UserID:    model.NilUUID,
		ExpiresAt: time.Now().Add(time.Hour),
	}
	sessionRepositoryMock.On("Find", gormDB).Return([]*models.Session{session}, nil)

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
	sessionRepositoryMock.On("Find", gormDB).Return(nil, httperrors.ErrForTests)
	assert.PanicsWithError(t, httperrors.ErrForTests.Error(), func() { service.pullFromDB() })
}

func Test_removeExpired(t *testing.T) {
	sessionRepositoryMock, service, logs, _ := setupTest(t)
	uuidSample := model.NewUUID()
	session := &models.Session{
		UUIDModel: model.UUIDModel{
			ID: model.NilUUID,
		},
		UserID:    model.NilUUID,
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	sessionRepositoryMock.
		On("Delete", gormDB, session).
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
	uuidSample := model.NewUUID()
	session := &models.Session{
		UUIDModel: model.UUIDModel{
			ID: model.NilUUID,
		},
		UserID:    model.NilUUID,
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	sessionRepositoryMock.
		On("Delete", gormDB, session).
		Return(httperrors.ErrForTests)

	service.cache[uuidSample] = session

	assert.Panics(t, func() { service.removeExpired() })
}

func Test_get(t *testing.T) {
	sessionRepositoryMock, service, _, _ := setupTest(t)
	uuidSample := model.NewUUID()
	session := &models.Session{
		UUIDModel: model.UUIDModel{
			ID: model.NilUUID,
		},
		UserID:    model.NilUUID,
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	sessionRepositoryMock.
		On("GetByID", gormDB, mock.Anything).
		Return(session, nil)

	sessionFound := service.get(uuidSample)
	assert.Equal(t, sessionFound, session)
}
