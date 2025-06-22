package assemble

import (
	"study/internal/api/handler/dto"
	"study/internal/domain/user/entity"
	"study/token"
	"study/util/errors"
	"time"
)

type CreateOrderCommand struct {
	UserID    int64
	SkuID     int64
	Start     time.Time
	End       time.Time
	Number    int
	PriceType int8
	PayType   string
	Contact   []entity.RoomContact
}

func NewCreateOrderCommand(req dto.CreateOrderRequest, payload *token.Payload) (*CreateOrderCommand, error) {
	const layout = "2006-01-02"
	start, err := time.ParseInLocation(layout, req.StartDate, time.Local)
	if err != nil {
		return nil, errors.New(errors.ErrDateFormat, "incorrect date format")
	}
	end, err := time.ParseInLocation(layout, req.EndDate, time.Local)
	if err != nil {
		return nil, errors.New(errors.ErrDateFormat, "incorrect date format")
	}

	roomContacts := make([]entity.RoomContact, req.Number) // 每间房一个 RoomContact
	for i := 0; i < req.Number; i++ {
		roomContacts[i] = entity.RoomContact{}
	}

	contacts := make([]entity.RoomContact, len(req.Contact))
	for i, room := range req.Contact {
		guests := make([]entity.Contact, len(room.Guests))
		for j, g := range room.Guests {
			guests[j] = entity.Contact{
				Name:  g.Name,
				Phone: g.Phone,
			}
		}
		contacts[i] = entity.RoomContact{Guests: guests}
	}

	return &CreateOrderCommand{
		UserID:    payload.UserId,
		SkuID:     req.SkuID,
		Start:     start,
		End:       end,
		Number:    req.Number,
		PriceType: req.PriceType,
		PayType:   req.PayType,
		Contact:   contacts,
	}, nil
}

type CreateAddCartCommand struct {
	UserID   int64
	SkuID    int64
	Quantity int32
}

func NewAddCartCommand(req dto.AddCartRequest, payload *token.Payload) (*CreateAddCartCommand, error) {
	return &CreateAddCartCommand{
		UserID:   payload.UserId,
		SkuID:    req.SkuID,
		Quantity: req.Quantity,
	}, nil
}
