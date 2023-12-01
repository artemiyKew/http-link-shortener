package service_test

import (
	"testing"

	"github.com/artemiyKew/http-link-shortener/internal/service"
	"github.com/stretchr/testify/assert"
)

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
