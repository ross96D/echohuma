package echohuma

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/labstack/echo/v4"
)

// type EchoRouter struct {
// 	router *echo.Group
// }

// func (e EchoRouter) Handle(op *huma.Operation, handler func(ctx huma.Context)) {
// 	e.router.Add(op.Method, op.Path, func(c echo.Context) error {
// 		handler(NewContext(op, c))
// 		return nil
// 	})
// }

type EchoAdapter struct {
	app *echo.Echo
}

func (e EchoAdapter) Handle(op *huma.Operation, handler func(ctx huma.Context)) {
	e.app.Add(op.Method, op.Path, func(c echo.Context) error {
		handler(NewContext(op, c))
		return nil
	})
}

func (e EchoAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.app.ServeHTTP(w, r)
}

// func (e EchoAdapter) Group(prefix string, m ...echo.MiddlewareFunc) EchoRouter {
// 	g := e.app.Group(prefix, m...)
// 	return EchoRouter{
// 		router: g,
// 	}
// }

func New(r *echo.Echo, config huma.Config) huma.API {
	return huma.NewAPI(config, EchoAdapter{app: r})
}

type echoCtx struct {
	op  *huma.Operation
	ctx echo.Context
}

// Operation returns the OpenAPI operation that matched the request.
func (e echoCtx) Operation() *huma.Operation {
	return e.op
}

// Context returns the underlying request context.
func (e echoCtx) Context() context.Context {
	return e.ctx.Request().Context()
}

// Method returns the HTTP method for the request.
func (e echoCtx) Method() string {
	return e.ctx.Request().Method
}

// Host returns the HTTP host for the request.
func (e echoCtx) Host() string {
	return e.ctx.Request().Host
}

// URL returns the full URL for the request.
func (e echoCtx) URL() url.URL {
	return *e.ctx.Request().URL
}

// Param returns the value for the given path parameter.
func (e echoCtx) Param(name string) string {
	return e.ctx.Param(name)
}

// Query returns the value for the given query parameter.
func (e echoCtx) Query(name string) string {
	return e.ctx.QueryParam(name)
}

// Header returns the value for the given header.
func (e echoCtx) Header(name string) string {
	return e.ctx.Request().Header.Get(name)
}

// EachHeader iterates over all headers and calls the given callback with
// the header name and value.
func (e echoCtx) EachHeader(cb func(name, value string)) {
	for name, values := range e.ctx.Request().Header {
		for _, value := range values {
			cb(name, value)
		}
	}
}

// BodyReader returns the request body reader.
func (e echoCtx) BodyReader() io.Reader {
	return e.ctx.Request().Body
}

// GetMultipartForm returns the parsed multipart form, if any.
func (e echoCtx) GetMultipartForm() (*multipart.Form, error) {
	err := e.ctx.Request().ParseMultipartForm(8 * 1024)
	return e.ctx.Request().MultipartForm, err
}

// SetReadDeadline sets the read deadline for the request body.
func (e echoCtx) SetReadDeadline(deadline time.Time) error {
	//! TODO search on the echo framework if there is a better way
	return huma.SetReadDeadline(e.ctx.Response(), deadline)
}

// SetStatus sets the HTTP status code for the response.
func (e echoCtx) SetStatus(code int) {
	//! TODO search on the echo framework if there is a better way
	e.ctx.Response().WriteHeader(code)
}

// SetHeader sets the given header to the given value, overwriting any
// existing value. Use `AppendHeader` to append a value instead.
func (e echoCtx) SetHeader(name string, value string) {
	e.ctx.Response().Header().Set(name, value)
}

// AppendHeader appends the given value to the given header.
func (e echoCtx) AppendHeader(name string, value string) {
	e.ctx.Response().Header().Add(name, value)
}

// BodyWriter returns the response body writer.
func (e echoCtx) BodyWriter() io.Writer {
	return e.ctx.Response()
}

func NewContext(op *huma.Operation, ctx echo.Context) huma.Context {
	return echoCtx{
		op:  op,
		ctx: ctx,
	}
}
