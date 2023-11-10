package delivery

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

func newError(c *fiber.Ctx, code int, err error) {
	msg := fmt.Sprintf("STATUS CODE: %d with error %s", code, err)
	logrus.Error(msg)
	_ = respond(c, code, map[string]string{"error": err.Error()})
}
