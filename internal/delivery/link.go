package delivery

import (
	"net/http"

	"github.com/artemiyKew/http-link-shortener/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type linkRoutes struct {
	service.Link
}

func newLinkRoutes(g fiber.Router, linkService service.Link) *linkRoutes {
	r := &linkRoutes{
		Link: linkService,
	}

	g.Post("/", r.CreateShortLink)
	g.Get("/", r.GetLinkInfo)
	g.Get("/:token", r.GetShortLink)

	return r
}

type createShortLinkInput struct {
	Link string `json:"long_url" validate:"required"`
}

func (r *linkRoutes) CreateShortLink(c *fiber.Ctx) error {
	var input createShortLinkInput

	if err := c.BodyParser(&input); err != nil {
		newError(c, http.StatusBadRequest, err)
		return err
	}

	linkOutput, err := r.Link.CreateShortLink(c.UserContext(), service.LinkInput{
		Link: input.Link,
	})

	if err != nil {
		//TODO создать ошибку что ссылка уже создана и вернуть статус 200
		newError(c, http.StatusInternalServerError, err)
		return err
	}

	return respond(c, fiber.StatusCreated, linkOutput)
}

type getLinkInfoInput struct {
	Token string `json:"token"`
}

func (r *linkRoutes) GetLinkInfo(c *fiber.Ctx) error {
	var input getLinkInfoInput

	if err := c.BodyParser(&input); err != nil {
		newError(c, http.StatusBadRequest, err)
		return err
	}

	linkOutput, err := r.Link.GetLinkInfo(c.UserContext(), service.LinkInput{
		Link: input.Token,
	})

	if err != nil {
		newError(c, http.StatusInternalServerError, err)
		return err
	}
	return respond(c, fiber.StatusCreated, linkOutput)
}

func (r *linkRoutes) GetShortLink(c *fiber.Ctx) error {
	fullURL, err := r.Link.GetShortLink(c.UserContext(), c.Params("token"))
	logrus.Info(fullURL)
	if err != nil {
		newError(c, http.StatusBadRequest, err)
		return err
	}

	c.Response().Header.Set("Location", fullURL)
	return c.Redirect(fullURL)
}
