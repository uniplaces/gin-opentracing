package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type Middleware struct {
	tracer opentracing.Tracer

	beforeHook func(opentracing.Span, *gin.Context)
	afterHook  func(opentracing.Span, *gin.Context)

	operationNameFn func(*gin.Context) string
	errorFn         func(*gin.Context) bool
	resourceNameFn  func(*gin.Context) string
}

func New(tracer opentracing.Tracer, options ...OptionFunc) *Middleware {
	m := &Middleware{tracer: tracer}

	for _, option := range options {
		option(m)
	}

	m.handleDefaultOptions()

	return m
}

func (m *Middleware) RequestTracer() gin.HandlerFunc {
	return func(c *gin.Context) {
		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		ctx, _ := m.tracer.Extract(opentracing.HTTPHeaders, carrier)

		span := m.tracer.StartSpan(m.operationNameFn(c), ext.RPCServerOption(ctx))

		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.String())
		span.SetTag("resource.name", m.resourceNameFn(c))

		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))

		m.beforeHook(span, c)
		c.Next()
		m.afterHook(span, c)

		ext.Error.Set(span, m.errorFn(c))
		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
	}
}
