package repo

import (
	"context"

	"github.com/artemiyKew/http-link-shortener/internal/entity"
	"github.com/artemiyKew/http-link-shortener/internal/repo/pgdb"
	"github.com/artemiyKew/http-link-shortener/internal/repo/redisdb"
)

type LinkRDB interface {
	CreateShortLink(context.Context, entity.Link) (*entity.Link, error)
	GetShortLink(context.Context, string) (string, error)
}

type LinkPGDB interface {
	CreateShortLink(context.Context, entity.Link) error
	UpdateCountOfVisits(context.Context, string) error
	GetShortLink(context.Context, string) (string, error)
	CheckShortLinkExist(context.Context, string) (entity.Link, error)
	GetLinkInfo(context.Context, string) (entity.Link, error)
}

type Repositories struct {
	LinkRDB
	LinkPGDB
}

func NewRepositories(rdb *redisdb.RedisDB, db *pgdb.PostgresDB) *Repositories {
	return &Repositories{
		LinkRDB:  redisdb.NewLinkRepo(rdb),
		LinkPGDB: pgdb.NewLinkRepo(db),
	}
}
