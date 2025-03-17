package user

import (
	"errors"
	"time"
)

// User 实体代表用户的基本信息
type User struct {
	ID          string
	Username    string
	Email       Email // 值对象
	FullName    string
	PhoneNumber string
	AvatarURL   string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	// Account 引用，而不是直接包含账户信息
	AccountID int64
}

// NewUser 创建一个新用户
func NewUser(username string, email string, fullName string) (*User, error) {
	emailVO, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:  username,
		Email:     emailVO,
		FullName:  fullName,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Time{},
	}, nil
}

// UpdateProfile 更新用户资料
func (u *User) UpdateProfile(fullName, phoneNumber, avatarURL string) {
	if fullName != "" {
		u.FullName = fullName
	}
	if phoneNumber != "" {
		u.PhoneNumber = phoneNumber
	}
	if avatarURL != "" {
		u.AvatarURL = avatarURL
	}
	u.UpdatedAt = time.Now()
}

// SetAccountID 为用户关联账户
func (u *User) SetAccountID(accountID int64) {
	u.AccountID = accountID
	u.UpdatedAt = time.Now()
}

// Deactivate 停用用户
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate 激活用户
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Time{}
}

// UserAccount 实体代表用户的账户信息
type UserAccount struct {
	ID           int64
	UserID       string
	PasswordHash []byte
	LastLoginAt  *time.Time
	LoginCount   int
	Balance      Money  // 值对象
	Status       string // "active", "suspended", "closed"
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewUserAccount 创建一个新的用户账户
func NewUserAccount(userID string, passwordHash []byte) *UserAccount {
	return &UserAccount{
		UserID:       userID,
		PasswordHash: passwordHash,
		LoginCount:   0,
		Balance:      NewMoney(0, "USD"),
		Status:       "active",
		CreatedAt:    time.Time{},
		UpdatedAt:    time.Time{},
	}
}

// RecordLogin 记录用户登录
func (ua *UserAccount) RecordLogin() {
	now := time.Now()
	ua.LastLoginAt = &now
	ua.LoginCount++
	ua.UpdatedAt = now
}

// ChangePassword 更改密码
func (ua *UserAccount) ChangePassword(newPasswordHash []byte) {
	ua.PasswordHash = newPasswordHash
	ua.UpdatedAt = time.Now()
}

// AddFunds 充值账户
func (ua *UserAccount) AddFunds(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	addAmount := NewMoney(amount, ua.Balance.Currency)
	ua.Balance = ua.Balance.Add(addAmount)
	ua.UpdatedAt = time.Now()
	return nil
}

// DeductFunds 扣除账户资金
func (ua *UserAccount) DeductFunds(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	if ua.Balance.Amount < amount {
		return errors.New("insufficient funds")
	}

	deductAmount := NewMoney(amount, ua.Balance.Currency)
	ua.Balance = ua.Balance.Subtract(deductAmount)
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

// 值对象 Email
type Email struct {
	Address string
}

// NewEmail 创建Email值对象
func NewEmail(address string) (Email, error) {
	// 这里应添加邮箱格式验证
	if address == "" {
		return Email{}, errors.New("email cannot be empty")
	}
	// 添加更多验证...

	return Email{Address: address}, nil
}

// String 获取Email字符串表示
func (e Email) String() string {
	return e.Address
}

// 值对象 Money
type Money struct {
	Amount   float64
	Currency string
}

// NewMoney 创建Money值对象
func NewMoney(amount float64, currency string) Money {
	return Money{
		Amount:   amount,
		Currency: currency,
	}
}

// Add 货币加法，不修改原值对象
func (m Money) Add(other Money) Money {
	if m.Currency != other.Currency {
		panic("cannot add different currencies")
	}
	return Money{
		Amount:   m.Amount + other.Amount,
		Currency: m.Currency,
	}
}

// Subtract 货币减法，不修改原值对象
func (m Money) Subtract(other Money) Money {
	if m.Currency != other.Currency {
		panic("cannot subtract different currencies")
	}
	return Money{
		Amount:   m.Amount - other.Amount,
		Currency: m.Currency,
	}
}
