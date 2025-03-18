package user

import (
	"context"
	"errors"
	"study/util"

	"golang.org/x/crypto/bcrypt"
)

// DomainService 用户领域服务
type DomainService struct {
	userRepo        UserRepository
	userAccountRepo UserAccountRepository
}

// NewDomainService 创建用户领域服务
func NewDomainService(userRepo UserRepository, userAccountRepo UserAccountRepository) *DomainService {
	return &DomainService{
		userRepo:        userRepo,
		userAccountRepo: userAccountRepo,
	}
}

// RegisterUser 注册新用户（包含账户创建）
func (s *DomainService) RegisterUser(ctx context.Context, username, phone, email, password string) (*User, error) {
	// 检查用户名唯一性
	existingUser, _ := s.userRepo.GetByUsername(ctx, username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱唯一性
	existingEmail, _ := s.userRepo.GetByEmail(ctx, email)
	if existingEmail != nil {
		return nil, errors.New("email already exists")
	}

	passwordHash, err := util.HashPassword(password)
	if err != nil {
		return nil, err
	}
	// 创建用户
	user, err := NewUser(username, phone, email, passwordHash)
	if err != nil {
		return nil, err
	}

	// 保存用户
	if err := s.userRepo.Save(ctx, user); err != nil {
		return nil, err
	}

	// 创建账户
	account := NewUserAccount(user.ID)

	// 保存账户
	if err := s.userAccountRepo.Save(ctx, account); err != nil {
		// 回滚用户创建
		_ = s.userRepo.Delete(ctx, user.ID)
		return nil, err
	}

	// 更新用户关联的账户ID
	user.SetAccountID(account.ID)
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// AuthenticateUser 用户登录认证
func (s *DomainService) AuthenticateUser(ctx context.Context, usernameOrEmail, password string) (*User, error) {
	// 查找用户
	var user *User
	var err error

	// 尝试通过用户名查找
	user, err = s.userRepo.GetByUsername(ctx, usernameOrEmail)
	if err != nil || user == nil {
		// 尝试通过邮箱查找
		emailObj, err := NewEmail(usernameOrEmail)
		if err == nil {
			user, err = s.userRepo.GetByEmail(ctx, emailObj.String())
		}
	}

	if err != nil || user == nil {
		return nil, errors.New("invalid credentials")
	}

	// 检查用户状态
	if !user.IsActive {
		return nil, errors.New("user account is deactivated")
	}

	// 获取用户账户
	account, err := s.userAccountRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		return nil, errors.New("user account not found")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword(account.PasswordHash, []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// 记录登录
	account.RecordLogin()
	if err := s.userAccountRepo.Update(ctx, account); err != nil {
		// 仅记录错误，不影响认证流程
		// logger.Error("Failed to update login record", err)
	}

	return user, nil
}

// ChangeUserPassword 修改密码
func (s *DomainService) ChangeUserPassword(ctx context.Context, userID int64, currentPassword, newPassword string) error {
	// 获取用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	// 获取账户
	account, err := s.userAccountRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		return err
	}

	// 验证当前密码
	if err := bcrypt.CompareHashAndPassword(account.PasswordHash, []byte(currentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	// 生成新密码哈希
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	account.ChangePassword(newPasswordHash)
	return s.userAccountRepo.Update(ctx, account)
}

// AddFundsToAccount 为用户账户充值
func (s *DomainService) AddFundsToAccount(ctx context.Context, userID string, amount float64) error {
	// 获取账户
	account, err := s.userAccountRepo.GetByUserID(ctx, userID)
	if err != nil {
		return err
	}

	// 添加资金
	if err := account.AddFunds(amount); err != nil {
		return err
	}

	// 更新账户
	return s.userAccountRepo.Update(ctx, account)
}

// DeactivateUserAndAccount 停用用户和账户
func (s *DomainService) DeactivateUserAndAccount(ctx context.Context, userID int64) error {
	// 获取用户
	_, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	return nil
}
