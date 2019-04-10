package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type options struct {
	beforeHook      func(opentracing.Span, *gin.Context)
	afterHook       func(opentracing.Span, *gin.Context)
	operationNameFn func(*gin.Context) string
	errorFn         func(*gin.Context) bool
	resourceNameFn  func(*gin.Context) string
}

func RequestTracer(opts ...OptionFunc) gin.HandlerFunc {
	mwOptions := &options{}
	for _, opt := range opts {
		opt(mwOptions)
	}

	mwOptions.handleDefaultOptions()

	return func(c *gin.Context) {
		tracer := opentracing.GlobalTracer()

		carrier := opentracing.HTTPHeadersCarrier(c.Request.Header)
		ctx, _ := tracer.Extract(opentracing.HTTPHeaders, carrier)

		span := tracer.StartSpan(mwOptions.operationNameFn(c), ext.RPCServerOption(ctx))

		ext.HTTPMethod.Set(span, c.Request.Method)
		ext.HTTPUrl.Set(span, c.Request.URL.String())
		span.SetTag("resource.name", mwOptions.resourceNameFn(c))

		c.Request = c.Request.WithContext(opentracing.ContextWithSpan(c.Request.Context(), span))

		mwOptions.beforeHook(span, c)
		c.Next()
		mwOptions.afterHook(span, c)

		ext.Error.Set(span, mwOptions.errorFn(c))
		ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))

		span.Finish()
	}
}
