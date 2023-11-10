package service

import (
	"context"
	"math/big"
	"strings"
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
	// if err := isValidUrl(link.Link); err != nil {
	// 	return &LinkOutput{}, err
	// }

	token, err := generateShortLink(link.Link, 10)
	if err != nil {
		return &LinkOutput{}, err
	}

	linkOutput, err := s.linkRepo.CreateShortLink(ctx, entity.Link{
		FullURL:      link.Link,
		CreatedAt:    time.Now(),
		ExpiredAt:    time.Now().Add(24 * time.Hour),
		VisitCounter: 0,
		Token:        token,
	})

	if err != nil {
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

// func isValidUrl(input string) error {
// 	_, err := url.ParseRequestURI(input)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func generateShortLink(inputURL string, tokenLength int) (string, error) {
	const allowedChars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
	base := big.NewInt(int64(len(allowedChars)))

	// Преобразование полной ссылки в число
	inputURLBytes := []byte(inputURL)
	number := big.NewInt(0)
	for i := 0; i < len(inputURLBytes); i++ {
		char := inputURLBytes[i]
		index := strings.IndexByte(allowedChars, char)
		if index == -1 {
			continue
		}
		number.Mul(number, base)
		number.Add(number, big.NewInt(int64(index)))
	}

	// Преобразование числа в систему счисления с основанием 62
	shortLinkChars := make([]byte, 0, len(inputURL))
	for i := 0; i < tokenLength; i++ {
		remainder := new(big.Int)
		number, remainder = number.DivMod(number, base, remainder)
		shortLinkChars = append(shortLinkChars, allowedChars[remainder.Int64()])
	}

	return string(shortLinkChars), nil
}
