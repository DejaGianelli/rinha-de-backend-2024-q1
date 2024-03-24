package models

import "time"

type Transacao struct {
	Id          int
	Amount      int
	Type        string
	CustomerId  int
	Description string
	RealizedAt  time.Time
}

func NewTransacao() Transacao {
	t := Transacao{
		RealizedAt: time.Now(),
	}
	return t
}
