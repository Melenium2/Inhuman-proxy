package proxy

import (
	"errors"
	"github.com/Melenium2/inhuman-reverse-proxy/internal/proxy/storage"
	"github.com/gofiber/fiber/v2"
)

func newProxy(storage storage.ProxyStorage) func(ctx *fiber.Ctx) error {
	return func(ctx *fiber.Ctx) error {
		var req RequestProxy
		if err := req.UnmarshalJSON(ctx.Body()); err != nil {
			return writeError(ctx, 400, err)
		}
		if req.Code == "" || len(req.Address) == 0 {
			return writeError(ctx, 400, errors.New("request is empty"))
		}
		if err := storage.Set(ctx.Context(), req.Code, req.Address...); err != nil {
			return writeError(ctx, 400, err)
		}

		return write(ctx, map[string]interface{}{
			"status":  "ok",
			"added":   len(req.Address),
			"to_code": req.Code,
		})
	}
}

func write(ctx *fiber.Ctx, data interface{}) error {
	return ctx.JSON(data)
}

func writeError(ctx *fiber.Ctx, code int, err error) error {
	ctx.Status(code)
	return write(ctx, map[string]interface{}{
		"err": err,
	})
}
