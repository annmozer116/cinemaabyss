package model

type EventType string

const (
	UserEvent    EventType = "USER"
	PaymentEvent EventType = "PAYMENT"
	MovieEvent   EventType = "MOVIE"
)

type Event struct {
	Type EventType   `json:"type"`
	Data interface{} `json:"data"`
}
