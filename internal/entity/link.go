package entity

import "time"

type Link struct {
	FullURL      string    `db:"full_url" json:"full_url"`
	CreatedAt    time.Time `db:"create_at" json:"create_at"`
	ExpiredAt    time.Time `db:"expired_at" json:"expired_at"`
	VisitCounter int       `db:"visit_counter" json:"visit_counter"`
	Token        string    `db:"token" json:"token"`
}
