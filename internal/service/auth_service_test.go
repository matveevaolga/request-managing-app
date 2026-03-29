package service

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/matveevaolga/request-managing-app/internal/domain"
	"github.com/matveevaolga/request-managing-app/internal/domain/repository/mocks"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	service := NewAuthService(mockUserRepo, "test-secret", 24)

	admin1Hash, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.DefaultCost)
	admin2Hash, _ := bcrypt.GenerateFromPassword([]byte("admin2"), bcrypt.DefaultCost)
	user1Hash, _ := bcrypt.GenerateFromPassword([]byte("user1"), bcrypt.DefaultCost)
	user2Hash, _ := bcrypt.GenerateFromPassword([]byte("user2"), bcrypt.DefaultCost)

	t.Run("successful login as admin1", func(t *testing.T) {
		user := &domain.User{
			ID:       1,
			Username: "admin1",
			Password: string(admin1Hash),
			Role:     domain.RoleAdmin,
		}
		mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "admin1").Return(user, nil)

		token, err := service.Login(context.Background(), "admin1", "admin1")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("successful login as admin2", func(t *testing.T) {
		user := &domain.User{
			ID:       2,
			Username: "admin2",
			Password: string(admin2Hash),
			Role:     domain.RoleAdmin,
		}
		mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "admin2").Return(user, nil)

		token, err := service.Login(context.Background(), "admin2", "admin2")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("successful login as user1", func(t *testing.T) {
		user := &domain.User{
			ID:       3,
			Username: "user1",
			Password: string(user1Hash),
			Role:     domain.RoleUser,
		}
		mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "user1").Return(user, nil)

		token, err := service.Login(context.Background(), "user1", "user1")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("successful login as user2", func(t *testing.T) {
		user := &domain.User{
			ID:       4,
			Username: "user2",
			Password: string(user2Hash),
			Role:     domain.RoleUser,
		}
		mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "user2").Return(user, nil)

		token, err := service.Login(context.Background(), "user2", "user2")

		assert.NoError(t, err)
		assert.NotEmpty(t, token)
	})

	t.Run("invalid credentials - user not found", func(t *testing.T) {
		mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "nonexistent").Return(nil, domain.ErrUserNotFound)

		token, err := service.Login(context.Background(), "nonexistent", "password")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidCredentials, err)
		assert.Empty(t, token)
	})

	t.Run("wrong password for admin1", func(t *testing.T) {
		user := &domain.User{
			ID:       1,
			Username: "admin1",
			Password: string(admin1Hash),
			Role:     domain.RoleAdmin,
		}
		mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "admin1").Return(user, nil)

		token, err := service.Login(context.Background(), "admin1", "wrongpassword")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidCredentials, err)
		assert.Empty(t, token)
	})

	t.Run("wrong password for user1", func(t *testing.T) {
		user := &domain.User{
			ID:       3,
			Username: "user1",
			Password: string(user1Hash),
			Role:     domain.RoleUser,
		}
		mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "user1").Return(user, nil)

		token, err := service.Login(context.Background(), "user1", "wrongpassword")

		assert.Error(t, err)
		assert.Equal(t, domain.ErrInvalidCredentials, err)
		assert.Empty(t, token)
	})
}

func TestAuthService_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mocks.NewMockUserRepository(ctrl)
	service := NewAuthService(mockUserRepo, "test-secret", 24)

	admin1Hash, _ := bcrypt.GenerateFromPassword([]byte("admin1"), bcrypt.DefaultCost)
	admin2Hash, _ := bcrypt.GenerateFromPassword([]byte("admin2"), bcrypt.DefaultCost)
	user1Hash, _ := bcrypt.GenerateFromPassword([]byte("user1"), bcrypt.DefaultCost)

	t.Run("validate valid token for admin1", func(t *testing.T) {
		user := &domain.User{
			ID:       1,
			Username: "admin1",
			Password: string(admin1Hash),
			Role:     domain.RoleAdmin,
		}
		mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "admin1").Return(user, nil)

		token, _ := service.Login(context.Background(), "admin1", "admin1")

		claims, err := service.ValidateToken(token)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), claims.UserID)
		assert.Equal(t, "ADMIN", claims.Role)
	})

	t.Run("validate valid token for user1", func(t *testing.T) {
		user := &domain.User{
			ID:       3,
			Username: "user1",
			Password: string(user1Hash),
			Role:     domain.RoleUser,
		}
		mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "user1").Return(user, nil)

		token, _ := service.Login(context.Background(), "user1", "user1")

		claims, err := service.ValidateToken(token)

		assert.NoError(t, err)
		assert.Equal(t, int64(3), claims.UserID)
		assert.Equal(t, "USER", claims.Role)
	})

	t.Run("validate valid token for admin2", func(t *testing.T) {
		user := &domain.User{
			ID:       2,
			Username: "admin2",
			Password: string(admin2Hash),
			Role:     domain.RoleAdmin,
		}
		mockUserRepo.EXPECT().GetByUsername(gomock.Any(), "admin2").Return(user, nil)

		token, _ := service.Login(context.Background(), "admin2", "admin2")

		claims, err := service.ValidateToken(token)

		assert.NoError(t, err)
		assert.Equal(t, int64(2), claims.UserID)
		assert.Equal(t, "ADMIN", claims.Role)
	})

	t.Run("validate invalid token", func(t *testing.T) {
		claims, err := service.ValidateToken("invalid-token")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("validate malformed token", func(t *testing.T) {
		claims, err := service.ValidateToken("malformed-token")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})

	t.Run("validate empty token", func(t *testing.T) {
		claims, err := service.ValidateToken("")

		assert.Error(t, err)
		assert.Nil(t, claims)
	})
}
