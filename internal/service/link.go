package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/artemiyKew/http-link-shortener/internal/entity"
	"github.com/artemiyKew/http-link-shortener/internal/repo"
)

type LinkService struct {
	linkRepoRDB  repo.LinkRDB
	linkRepoPGDB repo.LinkPGDB
}

func NewLinkService(linkRepoRDB repo.LinkRDB, linkRepoPGDB repo.LinkPGDB) *LinkService {
	return &LinkService{
		linkRepoRDB:  linkRepoRDB,
		linkRepoPGDB: linkRepoPGDB,
	}
}

func (s *LinkService) CreateShortLink(ctx context.Context, link LinkInput) (*LinkOutput, error) {
	l, err := s.linkRepoPGDB.CheckShortLinkExist(ctx, ValidateAndFixURL(link.Link))

	if err != nil {
		token := generateShortLink(link.Link, 10)

		entity := entity.Link{
			FullURL:      ValidateAndFixURL(link.Link),
			CreatedAt:    time.Now(),
			ExpiredAt:    time.Now().Add(24 * time.Hour),
			VisitCounter: 0,
			Token:        token,
		}

		if err := IsValidUrl(entity.FullURL); err != nil {
			return &LinkOutput{}, err
		}

		linkOutput, err := s.linkRepoRDB.CreateShortLink(ctx, entity)
		if err != nil {
			return &LinkOutput{}, err
		}

		if err := s.linkRepoPGDB.CreateShortLink(ctx, entity); err != nil {
			return &LinkOutput{}, err
		}

		return &LinkOutput{
			FullURL:      linkOutput.FullURL,
			CreatedAt:    linkOutput.CreatedAt,
			ExpiredAt:    linkOutput.ExpiredAt,
			VisitCounter: linkOutput.VisitCounter,
			Token:        linkOutput.Token,
		}, nil

	} else {
		return &LinkOutput{
			FullURL:      l.FullURL,
			CreatedAt:    l.CreatedAt,
			ExpiredAt:    l.ExpiredAt,
			VisitCounter: l.VisitCounter,
			Token:        l.Token,
		}, nil
	}
}

func (s *LinkService) GetShortLink(ctx context.Context, token string) (string, error) {
	if err := s.linkRepoPGDB.UpdateCountOfVisits(ctx, token); err != nil {
		return "", err
	}

	fullURL, err := s.linkRepoRDB.GetShortLink(ctx, token)
	if err != nil {
		return "", err
	}

	return fullURL, nil
}

func IsValidUrl(input string) error {
	response, err := http.Get(input)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		return nil
	} else {
		return errors.New("fail to connect")
	}
}

func generateShortLink(inputURL string, tokenLength int) string {
	var tokenMutex sync.Mutex
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	hash := sha256.Sum256([]byte(inputURL))

	shortenedURL := hex.EncodeToString(hash[:5])

	return shortenedURL

}

func ValidateAndFixURL(url string) string {
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		url = "https://" + url
	}
	return url
}
