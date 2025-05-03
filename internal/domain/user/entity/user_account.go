package entity

import (
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

// UserAccount 实体代表用户的账户信息
type UserAccount struct {
	UserID        int64
	FrozenBalance Money
	Balance       Money // 值对象
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
		CreatedAt:     time.Time{},
		UpdatedAt:     time.Time{},
	}
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

type Money struct {
	Amount   decimal.Decimal
	Currency string
}

func NewMoney(amount string, currency string) (Money, error) {
	dec, err := decimal.NewFromString(amount)
	if err != nil {
		return Money{}, err
	}
	return Money{
		Amount:   dec,
		Currency: currency,
	}, nil
}

func (m Money) String() string {
	return m.Amount.String()
}

func (m Money) Add(other Money) (Money, error) {
	if m.Currency != other.Currency {
		return Money{}, errors.New("currency mismatch")
	}
	return Money{
		Amount:   m.Amount.Add(other.Amount),
		Currency: m.Currency,
	}, nil
}

// Subtract 货币减法，不修改原值对象
func (m Money) Subtract(other Money) Money {
	if m.Currency != other.Currency {
		panic("cannot subtract different currencies")
	}
	return Money{
		Amount:   m.Amount.Sub(other.Amount),
		Currency: m.Currency,
	}
}
