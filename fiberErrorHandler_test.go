package fiber_error_handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type dummyError struct {
	LeMessage string `json:"le_message"`
	KodeError int    `json:"kode_error"`
}

func (cE dummyError) Error() string {
	return cE.LeMessage
}

func (cE dummyError) ErrorCode() int {
	return cE.KodeError
}

func TestErrorHandler(t *testing.T) {
	errorHandler := New(WithCustomError(dummyError{}))

	f := fiber.New(fiber.Config{
		ErrorHandler: errorHandler.HandlerFunc,
	})

	f.Get("/bad-request", func(c *fiber.Ctx) (err error) {
		err = &dummyError{
			KodeError: http.StatusBadRequest,
			LeMessage: "Error:Field validation",
		}

		return err
	})
	f.Get("/fiber-error", func(c *fiber.Ctx) (err error) {
		err = &fiber.Error{
			Message: "fiber error occured",
			Code:    http.StatusBadGateway,
		}

		return err
	})
	f.Get("/random-error", func(c *fiber.Ctx) (err error) {
		err = errors.New("random error")

		return err
	})
	f.Get("/ok", func(c *fiber.Ctx) (err error) {
		return c.JSON("ok")
	})

	t.Run("TestErrorHandler: pass next", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/ok", nil)
		res, _ := f.Test(req)

		var body string
		readBody, _ := io.ReadAll(res.Body)
		json.Unmarshal(readBody, &body)

		assert.Equal(t, http.StatusOK, res.StatusCode)
		assert.Equal(t, body, "ok")
	})

	t.Run("TestErrorHandler: returns standard error for recognized errors", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/bad-request", nil)
		res, _ := f.Test(req)

		var body map[string]any
		readBody, _ := io.ReadAll(res.Body)
		json.Unmarshal(readBody, &body)

		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		assert.Equal(t, fmt.Sprint(body["kode_error"]), fmt.Sprint(http.StatusBadRequest))
		assert.Equal(t, body["le_message"], "Error:Field validation")

		req = httptest.NewRequest(http.MethodGet, "/fiber-error", nil)
		res, _ = f.Test(req)

		readBody, _ = io.ReadAll(res.Body)
		json.Unmarshal(readBody, &body)

		assert.Equal(t, http.StatusBadGateway, res.StatusCode)
		assert.Equal(t, fmt.Sprint(body["code"]), fmt.Sprint(http.StatusBadGateway))
		assert.Equal(t, body["message"], "fiber error occured")
	})

	t.Run("TestErrorHandler: returns standard echo http error for unrecognized errors", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/random-error", nil)
		res, _ := f.Test(req)

		var body map[string]any
		readBody, _ := io.ReadAll(res.Body)
		json.Unmarshal(readBody, &body)

		assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
		assert.Equal(t, body["message"], "random error")
	})
}
