package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/artemiyKew/http-link-shortener/internal/entity"
	"github.com/artemiyKew/http-link-shortener/internal/repo"
)

type LinkService struct {
	linkRepo repo.Link
}

func NewLinkService(linkRepo repo.Link) *LinkService {
	return &LinkService{
		linkRepo: linkRepo,
	}
}

func (s *LinkService) CreateShortLink(ctx context.Context, link LinkInput) (*LinkOutput, error) {
	token := generateShortLink(link.Link, 10)

	linkOutput, err := s.linkRepo.CreateShortLink(ctx, entity.Link{
		FullURL:      validateAndFixURL(link.Link),
		CreatedAt:    time.Now(),
		ExpiredAt:    time.Now().Add(24 * time.Hour),
		VisitCounter: 0,
		Token:        token,
	})

	if err != nil {
		return &LinkOutput{}, err
	}

	if err := isValidUrl(linkOutput.FullURL); err != nil {
		return &LinkOutput{}, err
	}

	return &LinkOutput{
		FullURL:      linkOutput.FullURL,
		CreatedAt:    linkOutput.CreatedAt,
		ExpiredAt:    linkOutput.ExpiredAt,
		VisitCounter: linkOutput.VisitCounter,
		Token:        linkOutput.Token,
	}, nil
}

func (s *LinkService) GetShortLink(ctx context.Context, shortLink string) (string, error) {
	fullURL, err := s.linkRepo.GetShortLink(ctx, shortLink)
	if err != nil {
		return "", err
	}

	return fullURL, nil
}

func isValidUrl(input string) error {
	_, err := url.ParseRequestURI(input)
	if err != nil {
		return err
	}
	return nil
}

func generateShortLink(inputURL string, tokenLength int) string {
	var tokenMutex sync.Mutex
	tokenMutex.Lock()
	defer tokenMutex.Unlock()

	hash := sha256.Sum256([]byte(inputURL))

	shortenedURL := hex.EncodeToString(hash[:5])

	return shortenedURL

}

func validateAndFixURL(url string) string {
	if !strings.HasPrefix(url, "https://") && !strings.HasPrefix(url, "http://") {
		// Если префикс "https://" или "http://" отсутствует, добавляем "https://"
		url = "https://" + url
	}
	return url
}
