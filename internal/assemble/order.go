package assemble

import (
	"study/internal/api/handler/dto"
	"study/token"
	"time"
)

type contact struct {
	Name  string
	Phone string
}

type RoomContact struct {
	Guests []contact
}

type CreateOrderCommand struct {
	UserID    int64
	SkuID     int64
	Start     time.Time
	End       time.Time
	Number    int
	PriceType int8
	PayType   string
	Contact   []RoomContact
}

func NewCreateOrderCommand(req dto.CreateOrderRequest, payload *token.Payload) (*CreateOrderCommand, error) {
	const layout = "2006-01-02"
	start, err := time.ParseInLocation(layout, req.StartDate, time.Local)
	if err != nil {
		return nil, err
	}
	end, err := time.ParseInLocation(layout, req.EndDate, time.Local)
	if err != nil {
		return nil, err
	}

	contacts := make([]RoomContact, len(req.Contact))
	for i, room := range req.Contact {
		guests := make([]contact, len(room.Guests))
		for j, g := range room.Guests {
			guests[j] = contact{
				Name:  g.Name,
				Phone: g.Phone,
			}
		}
		contacts[i] = RoomContact{Guests: guests}
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
