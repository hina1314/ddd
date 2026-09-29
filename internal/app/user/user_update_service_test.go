package user

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/hina1314/ddd/config"
	"github.com/hina1314/ddd/internal/app/assemble"
	"github.com/hina1314/ddd/internal/domain/user/entity"
	"github.com/hina1314/ddd/internal/domain/user/service"
	"github.com/hina1314/ddd/internal/domain/user/usererrors"
	apperrors "github.com/hina1314/kit/errors"
	"github.com/stretchr/testify/require"
)

type userRepositoryStub struct {
	user    *entity.User
	updated *entity.User
}

func (r *userRepositoryStub) GetByID(context.Context, int64) (*entity.User, error) {
	copy := *r.user
	return &copy, nil
}
func (r *userRepositoryStub) GetByUsername(context.Context, string) (*entity.User, error) {
	return nil, nil
}
func (r *userRepositoryStub) GetByEmail(context.Context, string) (*entity.User, error) {
	return nil, nil
}
func (r *userRepositoryStub) GetByPhone(context.Context, string) (*entity.User, error) {
	return nil, nil
}
func (r *userRepositoryStub) Save(context.Context, *entity.User) error { return nil }
func (r *userRepositoryStub) Update(_ context.Context, user *entity.User) error {
	r.updated = user
	return nil
}
func (r *userRepositoryStub) Delete(context.Context, int64) error { return nil }
func (r *userRepositoryStub) List(context.Context, int, int64) ([]*entity.User, error) {
	return nil, nil
}

func newUserServiceForUpdate(repo *userRepositoryStub) *UserService {
	return NewUserService(
		service.NewUserRegisterService(repo),
		service.NewUserLoginService(repo),
		service.NewUserUpdateService(repo),
		repo,
		config.Config{},
		nil,
	)
}

func TestUpdateUserPreservesOmittedFields(t *testing.T) {
	repo := &userRepositoryStub{user: &entity.User{
		ID: 1, Phone: "+8613800138000", Email: entity.NewEmail("before@example.com"),
		Username: "before", Password: "existing-hash", CreatedAt: time.Now(),
	}}
	service := newUserServiceForUpdate(repo)
	username := "after"

	result, err := service.UpdateUser(context.Background(), &assemble.UpdateUserCommand{ID: 1, Username: &username})

	require.NoError(t, err)
	require.Equal(t, "after", result.Username)
	require.Equal(t, "+8613800138000", repo.updated.Phone)
	require.Equal(t, "before@example.com", repo.updated.Email.String())
	require.Equal(t, "existing-hash", repo.updated.Password)
}

func TestUpdateUserRejectsEmptyPatch(t *testing.T) {
	repo := &userRepositoryStub{user: &entity.User{ID: 1}}
	service := newUserServiceForUpdate(repo)

	_, err := service.UpdateUser(context.Background(), &assemble.UpdateUserCommand{ID: 1})

	var domainErr *apperrors.DomainError
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, apperrors.ErrInvalidInput, domainErr.Code)
}

func TestUpdateUserRejectsPasswordOverBcryptLimit(t *testing.T) {
	repo := &userRepositoryStub{user: &entity.User{ID: 1}}
	service := newUserServiceForUpdate(repo)
	password := strings.Repeat("密", 25)

	_, err := service.UpdateUser(context.Background(), &assemble.UpdateUserCommand{ID: 1, Password: &password})

	var domainErr *apperrors.DomainError
	require.ErrorAs(t, err, &domainErr)
	require.Equal(t, usererrors.ErrPasswordTooLong, domainErr.Code)
	require.Nil(t, repo.updated)
}
