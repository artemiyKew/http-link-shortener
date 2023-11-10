package repo

import (
	"context"

	"github.com/artemiyKew/http-link-shortener/internal/entity"
	"github.com/artemiyKew/http-link-shortener/internal/repo/redisdb"
)

type Link interface {
	CreateShortLink(context.Context, entity.Link) (*entity.Link, error)
	GetShortLink(context.Context, string) (string, error)
}

type Repositories struct {
	Link
}

func NewRepositories(rdb *redisdb.RedisDB) *Repositories {
	return &Repositories{
		Link: redisdb.NewLinkRepo(rdb),
	}
}
