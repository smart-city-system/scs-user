package middleware

import (
	"github.com/labstack/echo/v4"
)

// ResponseBodyWriter is a custom writer that captures the response body.
type CustomContext struct {
	echo.Context
}
type SuccessResponse struct {
	Status int         `json:"status"`
	Code   string      `json:"code"`
	Data   interface{} `json:"data,omitempty"`
}

func (mw *MiddlewareManager) ResponseStandardizer(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Create a custom context and pass it to the next handler.
		cc := &CustomContext{c}
		return next(cc)
	}
}
func (c *CustomContext) JSON(code int, i interface{}) error {
	// Check if the provided data is an error.
	if _, ok := i.(error); ok {
		// Handle the error case. You would need a custom error handler
		// to return the formatted error. For now, we'll assume a separate
		// error handler is registered.
		return c.Context.JSON(code, i)
	}

	// Wrap the data in your standard success response format.
	resp := SuccessResponse{
		Status: code,
		Code:   "0000",
		Data:   i,
	}

	// Call the original JSON method with the wrapped response.
	return c.Context.JSON(code, resp)
}
