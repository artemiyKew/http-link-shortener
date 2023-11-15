package service_test

import (
	"testing"

	"github.com/artemiyKew/http-link-shortener/internal/service"
	"github.com/stretchr/testify/assert"
)

func TestService_IsValidURL(t *testing.T) {
	testcase := []struct {
		name string
		in   string
		out  bool
	}{
		{
			name: "Valid",
			in:   "https://youtube.com",
			out:  true,
		},
		{
			name: "Valid",
			in:   "https://vk.com",
			out:  true,
		},
		{
			name: "Valid",
			in:   "http://google.com",
			out:  true,
		},
		{
			name: "Not Valid",
			in:   "https://youtube.r",
			out:  false,
		},
		{
			name: "Not Valid",
			in:   "https:/youtube.r",
			out:  false,
		},
		{
			name: "Valid",
			in:   "youtube.com",
			out:  true,
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {
			tc.in = service.ValidateAndFixURL(tc.in)
			if tc.out {
				assert.NoError(t, service.IsValidUrl(tc.in))
			} else {
				assert.Error(t, service.IsValidUrl(tc.in))
			}
		})
	}
}

func TestService_ValidateAndFixURL(t *testing.T) {
	testcase := []struct {
		name string
		in   string
		out  string
	}{
		{
			name: "Valid",
			in:   "youtube.com",
			out:  "https://youtube.com",
		},
		{
			name: "Valid",
			in:   "https://youtube.com",
			out:  "https://youtube.com",
		},
		{
			name: "Valid",
			in:   "http://youtube.com",
			out:  "http://youtube.com",
		},
	}

	for _, tc := range testcase {
		t.Run(tc.name, func(t *testing.T) {

			assert.Equal(t, tc.out, service.ValidateAndFixURL(tc.in))
		})
	}
}
