package pgdb

import (
	"context"

	"github.com/artemiyKew/http-link-shortener/internal/entity"
)

type LinkRepo struct {
	*PostgresDB
}

func NewLinkRepo(pgdb *PostgresDB) *LinkRepo {
	return &LinkRepo{pgdb}
}

func (r *LinkRepo) CreateShortLink(ctx context.Context, link entity.Link) error {
	var id int
	return r.DB.
		QueryRow("INSERT INTO links (full_url, create_at, expired_at, visit_counter, token) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			link.FullURL, link.CreatedAt, link.ExpiredAt, link.VisitCounter, link.Token).
		Scan(&id)

}

func (r *LinkRepo) UpdateCountOfVisits(ctx context.Context, token string) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE links SET visit_counter = visit_counter + 1  WHERE token = $1",
		token)
	if err != nil {

		if err := tx.Rollback(); err != nil {
			return err
		}

		return err
	}

	return tx.Commit()

}

func (r *LinkRepo) GetShortLink(ctx context.Context, token string) (string, error) {
	var fullURL string

	if err := r.DB.
		QueryRow("SELECT full_url FROM links WHERE token = $1", token).
		Scan(&fullURL); err != nil {
		return "", err
	}

	return fullURL, nil
}

func (r *LinkRepo) CheckShortLinkExist(ctx context.Context, fullURL string) (entity.Link, error) {
	var link entity.Link

	if err := r.DB.
		QueryRow("SELECT full_url, create_at, expired_at, visit_counter, token FROM links WHERE full_url = $1", fullURL).
		Scan(&link.FullURL, &link.CreatedAt, &link.ExpiredAt, &link.VisitCounter, &link.Token); err != nil {
		return entity.Link{}, err
	}

	return link, nil
}
