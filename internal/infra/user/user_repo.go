package user

import (
	"context"
	"database/sql"
	"github.com/hina1314/ddd/db/model"
	"github.com/hina1314/ddd/internal/domain/user/entity"
	"github.com/hina1314/ddd/internal/domain/user/repository"
	"github.com/hina1314/ddd/internal/domain/user/usererrors"
	"github.com/hina1314/ddd/internal/infra"
	"github.com/hina1314/kit/errors"
	"time"
)

type UserRepositoryImpl struct {
	db model.TxManager
}

func NewUserRepository(db model.TxManager) repository.UserRepository {
	return &UserRepositoryImpl{
		db: db,
	}
}

// **通用转换方法：model.User → user.User**
func (r *UserRepositoryImpl) toDomain(u model.User) (*entity.User, error) {
	emailVO := entity.NewEmail(u.Email.String)
	var deletedAt *time.Time
	if u.DeletedAt.Valid {
		deletedAt = &u.DeletedAt.Time
	}

	return &entity.User{
		ID:        u.ID,
		Phone:     u.Phone.String,
		Username:  u.Username,
		Email:     emailVO,
		Password:  u.Password,
		Avatar:    u.Avatar,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deletedAt,
	}, nil
}

func (r *UserRepositoryImpl) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	u, err := r.db.Querier(ctx).GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return r.toDomain(u) // 统一转换
}

func (r *UserRepositoryImpl) GetByPhone(ctx context.Context, phone string) (*entity.User, error) {
	u, err := r.db.Querier(ctx).GetUserByPhone(ctx, toNullString(phone))
	if err != nil {
		return nil, err
	}
	return r.toDomain(u) // 统一转换
}

func (r *UserRepositoryImpl) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	u, err := r.db.Querier(ctx).GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return r.toDomain(u) // 统一转换
}

func (r *UserRepositoryImpl) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	u, err := r.db.Querier(ctx).GetUserByEmail(ctx, toNullString(email))
	if err != nil {
		return nil, err
	}
	return r.toDomain(u) // 统一转换
}

func (r *UserRepositoryImpl) Save(ctx context.Context, u *entity.User) error {
	q := r.db.Querier(ctx)
	arg := model.CreateUserParams{
		Phone:    toNullString(u.Phone),
		Email:    emailToNullString(u.Email),
		Username: u.Username,
		Password: u.Password,
	}

	result, err := q.CreateUser(ctx, arg)
	if err != nil {
		if infra.IsDuplicateKeyError(err) {
			return errors.New(usererrors.ErrUserAlreadyExists, "User already exists")
		}
		return err
	}
	u.ID = result.ID
	u.CreatedAt = result.CreatedAt
	u.UpdatedAt = result.UpdatedAt
	return nil
}

func (r *UserRepositoryImpl) Update(ctx context.Context, u *entity.User) error {
	arg := model.UpdateUserParams{
		ID:       u.ID,
		Phone:    toNullString(u.Phone),
		Email:    emailToNullString(u.Email),
		Username: u.Username,
		Password: u.Password,
	}

	rows, err := r.db.Querier(ctx).UpdateUser(ctx, arg)
	if err != nil {
		if infra.IsDuplicateKeyError(err) {
			return errors.New(usererrors.ErrUserAlreadyExists, "User already exists")
		}
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepositoryImpl) Delete(ctx context.Context, id int64) error {
	rows, err := r.db.Querier(ctx).DeleteUser(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *UserRepositoryImpl) List(ctx context.Context, limit int, afterID int64) ([]*entity.User, error) {
	users, err := r.db.Querier(ctx).ListUsers(ctx, model.ListUsersParams{
		AfterID:  afterID,
		PageSize: int32(limit),
	})
	if err != nil {
		return nil, err
	}

	result := make([]*entity.User, 0, len(users))
	for _, u := range users {
		domainUser, err := r.toDomain(u)
		if err != nil {
			return nil, err
		}
		result = append(result, domainUser)
	}
	return result, nil
}

func emailToNullString(email entity.Email) sql.NullString {
	if email.IsNil() {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: email.String(), Valid: true}
}

func toNullString(str string) sql.NullString {
	if str == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: str, Valid: true}
}
