package delivery

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type key int

const (
	ctxKeyRequestID key = iota
)

func setRequestID(c *fiber.Ctx) error {
	id := uuid.New().String()
	c.Response().Header.Set("X-Request-ID", id)
	c.SetUserContext(context.WithValue(c.Context(), ctxKeyRequestID, id))
	return c.Next()
}

func logRequest(c *fiber.Ctx) error {
	msg := fmt.Sprintf("remote_addr=%s request_id=%s", c.Context().RemoteAddr(), c.UserContext().Value(ctxKeyRequestID))
	logrus.Info(fmt.Sprintf("started %s %s \t %s", c.Context().Method(), c.Context().RequestURI(), msg))
	start := time.Now()
	msg = fmt.Sprintf("completed with %d % s in %v \t %s",
		http.StatusOK,
		http.StatusText(http.StatusOK),
		time.Since(start),
		msg)
	logrus.Info(msg)
	return c.Next()
}
