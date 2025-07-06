# Fiber Error Handler Middleware ![coverage](https://raw.githubusercontent.com/xwinata/fiber_error_handler/badges/.badges/main/coverage.svg)

A customizable error handler middleware for the [Fiber](https://github.com/gofiber/fiber) web framework in Go.  
This middleware provides a centralized way to handle errors and return consistent JSON responses across your application.

## Features

- Handles standard `*fiber.Error`
- Supports custom error types that implement a `customError` interface
- Returns structured JSON responses with proper status codes
- Easy to register using `App().Config().ErrorHandler`

## Usage
### Define your custom error structure that implement `customError` interface
```
type MyCustomError struct {
	Message string
	Code    int
}

func (e MyCustomError) Error() string {
	return e.Message
}

func (e MyCustomError) ErrorCode() int {
	return e.Code
}
```
#### `customError` interface definition
```
type customError interface {
	Error() string
	ErrorCode() int
}
```
### Register the Error Handler
```import (
	"github.com/gofiber/fiber/v2"
	fiberErrorHandler "github.com/yourusername/fiber-error-handler"
)

func main() {
	handler := fiberErrorHandler.New(
		fiberErrorHandler.WithCustomError(MyCustomError{}),
	)

	app := fiber.New(fiber.Config{
		ErrorHandler: handler.HandlerFunc,
	})

	app.Get("/test", func(c *fiber.Ctx) error {
		return MyCustomError{Message: "invalid input", Code: 400}
	})

	app.Listen(":3000")
}
```
### The behavior
- If the error implements customError, the middleware returns:
```
{
  "Message": "your error message",
  "Code": 400
}
```
- If the error is a `*fiber.Error`, it uses fiber.Error.Code and its message.
- If it’s any other error, it returns a 500 Internal Server Error with the message.

