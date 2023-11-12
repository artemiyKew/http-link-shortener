package redisdb

import (
	"context"
	"encoding/json"

	"github.com/artemiyKew/http-link-shortener/internal/entity"
)

type LinkRepo struct {
	*RedisDB
}

func NewLinkRepo(rdb *RedisDB) *LinkRepo {
	return &LinkRepo{rdb}
}

func (r *LinkRepo) CreateShortLink(ctx context.Context, link entity.Link) (*entity.Link, error) {
	l, err := r.CheckShortLinkExist(ctx, link.FullURL)
	if err != nil {
		linkJSON, err := json.Marshal(link)
		if err != nil {
			return &entity.Link{}, err
		}

		key := link.Token
		err = r.RDB.SetEx(ctx, key, linkJSON, link.ExpiredAt.Sub(link.CreatedAt)).Err()
		if err != nil {
			return &entity.Link{}, err
		}

		// Добавление нового ключа в виде полной ссылки для быстрого поиска
		// Мб пофикшу
		key = link.FullURL
		err = r.RDB.SetEx(ctx, key, linkJSON, link.ExpiredAt.Sub(link.CreatedAt)).Err()
		if err != nil {
			return &entity.Link{}, err
		}
		return &link, nil
	} else {
		return l, nil
	}
}

func (r *LinkRepo) CheckShortLinkExist(ctx context.Context, key string) (*entity.Link, error) {
	link := &entity.Link{}

	linkJSON, err := r.RDB.Get(ctx, key).Result()
	if err != nil {
		return &entity.Link{}, err
	}

	err = json.Unmarshal([]byte(linkJSON), &link)
	if err != nil {
		return &entity.Link{}, err
	}

	return link, nil
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

	return link.FullURL, nil
}
