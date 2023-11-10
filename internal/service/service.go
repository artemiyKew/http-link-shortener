package service

import (
	"context"
	"time"

	"github.com/artemiyKew/http-link-shortener/internal/repo"
)

type LinkInput struct {
	Link string
}

type LinkOutput struct {
	FullURL      string    `json:"full_url"`
	CreatedAt    time.Time `json:"create_at"`
	ExpiredAt    time.Time `json:"expired_at"`
	VisitCounter int       `json:"visit_counter"`
	Token        string    `json:"token"`
}

type Link interface {
	CreateShortLink(context.Context, LinkInput) (*LinkOutput, error)
	GetShortLink(context.Context, string) (string, error)
}

type Services struct {
	Link Link
}

type ServicesDependencies struct {
	Repos *repo.Repositories
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Link: NewLinkService(deps.Repos.Link),
	}
}
