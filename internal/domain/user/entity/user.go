package entity

import (
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
	Account   UserAccount
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// NewUser 创建新用户（用于注册）
func NewUser(phone, email, username, password string) (*User, error) {
	emailVO := NewEmail(email)

	now := time.Now()
	user := &User{
		Phone:     phone,
		Email:     emailVO,
		Username:  username,
		Password:  password,
		CreatedAt: now,
		UpdatedAt: now,
	}
	user.Account = UserAccount{
		FrozenBalance: Money{},
		Balance:       Money{},
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	return user, nil
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

// Email 值对象
type Email struct {
	value *string
}

func NewEmail(value string) Email {
	if value == "" {
		return Email{value: nil}
	}
	return Email{value: &value}
}

func (e Email) String() string {
	if e.value == nil {
		return ""
	}
	return *e.value
}

func (e Email) IsNil() bool {
	return e.value == nil
}
