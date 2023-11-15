package redisdb

import (
	"context"

	"github.com/artemiyKew/http-link-shortener/internal/entity"
)

type LinkRepo struct {
	*RedisDB
}

func NewLinkRepo(rdb *RedisDB) *LinkRepo {
	return &LinkRepo{rdb}
}

func (r *LinkRepo) CreateShortLink(ctx context.Context, link entity.Link) (*entity.Link, error) {
	key := link.Token
	value := link.FullURL

	if err := r.RDB.SetEx(ctx, key, value, link.ExpiredAt.Sub(link.CreatedAt)).Err(); err != nil {
		return &entity.Link{}, err
	}

	return &link, nil
}

func (r *LinkRepo) GetShortLink(ctx context.Context, token string) (string, error) {
	key := token
	fullURL, err := r.RDB.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return fullURL, nil
}
