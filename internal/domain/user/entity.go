package user

import (
	"errors"
	"time"
)

// User 实体代表用户的基本信息
type User struct {
	ID        int64
	Phone     string
	Email     Email // 值对象
	Username  string
	Password  string
	Avatar    string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewUser 创建新用户（用于注册）
func NewUser(phone, email, username, password string) (*User, error) {
	emailVO := NewEmail(email)

	return &User{
		Phone:     phone,
		Email:     emailVO,
		Username:  username,
		Password:  password,
		Avatar:    "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

// UpdateProfile 更新用户资料
func (u *User) UpdateProfile(phone, avatar string) {
	if phone != "" {
		u.Phone = phone
	}
	if avatar != "" {
		u.Avatar = avatar
	}
	u.UpdatedAt = time.Now()
}

// ChangePassword 更改密码
func (u *User) ChangePassword(newPasswordHash string) {
	u.Password = newPasswordHash
	u.UpdatedAt = time.Now()
}

// UserAccount 实体代表用户的账户信息
type UserAccount struct {
	ID            int64
	UserID        int64
	LastLoginAt   *time.Time
	LoginCount    int
	FrozenBalance Money
	Balance       Money  // 值对象
	Status        string // "active", "suspended", "closed"
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// NewUserAccount 创建一个新的用户账户
func NewUserAccount(userID int64) *UserAccount {
	balance, err := NewMoney("0", "USD")
	if err != nil {
		panic(err)
	}
	frozen, err := NewMoney("0", "USD")
	if err != nil {
		panic(err)
	}
	return &UserAccount{
		UserID:        userID,
		Balance:       balance,
		FrozenBalance: frozen,
		Status:        "active",
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}
}

// RecordLogin 记录用户登录
func (ua *UserAccount) RecordLogin() {
	now := time.Now()
	ua.LastLoginAt = &now
	ua.UpdatedAt = now
}

// AddFunds 充值账户
func (ua *UserAccount) AddFunds(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	//ua.Balance = ua.Balance.Add(addAmount)
	ua.UpdatedAt = time.Now()
	return nil
}

// DeductFunds 扣除账户资金
func (ua *UserAccount) DeductFunds(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	//if ua.Balance.Amount < amount {
	//	return errors.New("insufficient funds")
	//}
	//
	//deductAmount := NewMoney(amount, ua.Balance.Currency)
	//ua.Balance = ua.Balance.Subtract(deductAmount)
	ua.UpdatedAt = time.Now()
	return nil
}

// Suspend 暂停账户
func (ua *UserAccount) Suspend() {
	ua.Status = "suspended"
	ua.UpdatedAt = time.Now()
}

// Close 关闭账户
func (ua *UserAccount) Close() {
	ua.Status = "closed"
	ua.UpdatedAt = time.Now()
}

// Reactivate 重新激活账户
func (ua *UserAccount) Reactivate() {
	ua.Status = "active"
	ua.UpdatedAt = time.Now()
}
