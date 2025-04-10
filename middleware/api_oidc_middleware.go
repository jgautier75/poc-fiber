package middleware

import (
	"context"
	"errors"
	"poc-fiber/exceptions"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gofiber/fiber/v2"
)

func NewApiOidcHandler(apiBaseUri string, verifier *oidc.IDTokenVerifier) fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		p := c.Path()
		if strings.HasPrefix(p, apiBaseUri) {
			auth := c.GetReqHeaders()["Authorization"]
			if auth == nil {
				return c.Status(fiber.StatusUnauthorized).JSON(exceptions.ConvertToFunctionalError(errors.New("authorization header expected"), fiber.StatusUnauthorized))
			}
			if !strings.HasPrefix(auth[0], "Bearer") {
				return c.Status(fiber.StatusUnauthorized).JSON(exceptions.ConvertToFunctionalError(errors.New("bearer expected"), fiber.StatusUnauthorized))
			}
			reqToken := strings.Split(auth[0], " ")[1]
			_, errDecode := verifier.Verify(context.Background(), reqToken)
			if errDecode != nil {
				return c.Status(fiber.StatusUnauthorized).JSON(exceptions.ConvertToFunctionalError(errDecode, fiber.StatusUnauthorized))
			}
		}
		return c.Next()
	}
}
