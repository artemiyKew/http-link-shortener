package pgdb

import (
	"context"
	"errors"

	"github.com/artemiyKew/http-link-shortener/internal/entity"
)

type LinkRepo struct {
	*PostgresDB
}

func NewLinkRepo(pgdb *PostgresDB) *LinkRepo {
	return &LinkRepo{pgdb}
}

func (r *LinkRepo) CreateShortLink(ctx context.Context, link entity.Link) error {
	_, err := r.DB.ExecContext(ctx, "DELETE FROM links WHERE expired_at < NOW()")
	if err != nil {
		return nil
	}
	var id int
	return r.DB.
		QueryRowContext(ctx, "INSERT INTO links (full_url, create_at, expired_at, visit_counter, token) VALUES ($1, $2, $3, $4, $5) RETURNING id",
			link.FullURL, link.CreatedAt, link.ExpiredAt, link.VisitCounter, link.Token).
		Scan(&id)
}

func (r *LinkRepo) GetLinkInfo(ctx context.Context, token string) (entity.Link, error) {
	var link entity.Link
	var id int
	if err := r.DB.
		QueryRowContext(ctx, "SELECT * FROM links WHERE token = $1",
			token).
		Scan(&id, &link.FullURL, &link.CreatedAt, &link.ExpiredAt, &link.VisitCounter, &link.Token); err != nil {
		return entity.Link{}, err
	}
	return link, nil
}

func (r *LinkRepo) UpdateCountOfVisits(ctx context.Context, token string) (err error) {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			err = errors.Join(err, tx.Rollback())
			return
		}
		err = errors.Join(err, tx.Commit())
	}()

	_, err = tx.ExecContext(ctx, "UPDATE links SET visit_counter = visit_counter + 1  WHERE token = $1",
		token)

	return err
}

func (r *LinkRepo) GetShortLink(ctx context.Context, token string) (string, error) {
	var fullURL string

	if err := r.DB.
		QueryRowContext(ctx, "SELECT full_url FROM links WHERE token = $1", token).
		Scan(&fullURL); err != nil {
		return "", err
	}

	return fullURL, nil
}

func (r *LinkRepo) CheckShortLinkExist(ctx context.Context, fullURL string) (entity.Link, error) {
	var link entity.Link

	if err := r.DB.
		QueryRowContext(ctx, "SELECT full_url, create_at, expired_at, visit_counter, token FROM links WHERE full_url = $1", fullURL).
		Scan(&link.FullURL, &link.CreatedAt, &link.ExpiredAt, &link.VisitCounter, &link.Token); err != nil {
		return entity.Link{}, err
	}

	return link, nil
}
