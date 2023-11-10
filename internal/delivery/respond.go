package delivery

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

func respond(c *fiber.Ctx, code int, data interface{}) error {
	c.Response().Header.SetStatusCode(code)
	if data != nil {
		_ = json.NewEncoder(c.Response().BodyWriter()).Encode(data)
	}

	return nil
}
