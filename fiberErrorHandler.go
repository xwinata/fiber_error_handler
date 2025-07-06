package fiber_error_handler

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type customError interface {
	Error() string
	ErrorCode() int
}

type options struct {
	customError customError
}

type optionFunc func(*options)

func WithCustomError(err customError) optionFunc {
	return func(o *options) {
		o.customError = err
	}
}

type ErrorHandler struct {
	options *options
}

func New(opts ...optionFunc) *ErrorHandler {
	options := &options{}
	for _, opt := range opts {
		opt(options)
	}

	return &ErrorHandler{
		options: options,
	}
}

func (h *ErrorHandler) HandlerFunc(ctx *fiber.Ctx, err error) error {
	if err != nil {
		var fiberError *fiber.Error
		var customError customError

		if h.options.customError != nil {
			if errors.As(err, &customError) {
				return ctx.Status(customError.ErrorCode()).JSON(customError)
			}
		}

		if errors.As(err, &fiberError) {
			return ctx.Status(fiberError.Code).JSON(fiberError)
		}

		// Return a generic error response
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		})
	}
	return err
}
