package middleware

import (
	"context"
	"poc-fiber/opentelemetry"

	"github.com/gofiber/fiber/v3"
)

func HttpMiddleWareStats(metrics *opentelemetry.OtelMetrics) fiber.Handler {
	return func(c fiber.Ctx) (err error) {
		savedCtx, _ := context.WithCancel(c.Context())
		metrics.TotalReqCounter.Add(savedCtx, 1)
		return c.Next()
	}
}
