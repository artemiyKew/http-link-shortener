package redisdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/artemiyKew/http-link-shortener/internal/entity"
)

type LinkRepo struct {
	*RedisDB
}

func NewLinkRepo(rdb *RedisDB) *LinkRepo {
	return &LinkRepo{rdb}
}

// TODO переписать весь пакет
func (r *LinkRepo) CreateShortLink(ctx context.Context, link entity.Link) (*entity.Link, error) {
	linkJSON, err := json.Marshal(link)
	if err != nil {
		return &entity.Link{}, err
	}

	key := link.Token
	err = r.RDB.Set(ctx, key, linkJSON, link.ExpiredAt.Sub(link.CreatedAt)).Err()
	if err != nil {
		return &entity.Link{}, err
	}
	return &link, nil
}

func (r *LinkRepo) GetShortLink(ctx context.Context, shortLink string) (string, error) {
	// AKA token  = shortlink
	key := shortLink
	storedLinkJSON, err := r.RDB.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	var link entity.Link
	err = json.Unmarshal([]byte(storedLinkJSON), &link)
	if err != nil {
		return "", err
	}

	if link.ExpiredAt.Unix() < time.Now().Unix() {
		return "", errors.New("token expired")
	}
	return link.FullURL, nil
}

func (r *LinkRepo) checkTTL(ctx context.Context, key string) error {
	ttl, err := r.RDB.TTL(ctx, key).Result()
	if err != nil {
		return err
	}

	if ttl.Seconds() < 0 {
		fmt.Printf("Ключ '%s' истек.\n", key)
		// Удаляем ключ
		err := r.RDB.Del(ctx, key).Err()
		if err != nil {
			return err
		}
	}
	return nil
}
