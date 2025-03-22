package user

import (
	"errors"
	"github.com/shopspring/decimal"
)

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

// Email 值对象
type Email struct {
	value string
}

func NewEmail(value string) (Email, error) {
	if value == "" {
		return Email{}, errors.New("email cannot be empty")
	}
	// 可添加更复杂的验证
	return Email{value: value}, nil
}

func (e Email) String() string {
	return e.value
}
