package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

type OptionFunc func(*Middleware)

func SetOperationNameFn(fn func(*gin.Context) string) OptionFunc {
	return func(m *Middleware) {
		m.operationNameFn = fn
	}
}

func SetErrorFn(fn func(*gin.Context) bool) OptionFunc {
	return func(m *Middleware) {
		m.errorFn = fn
	}
}

func SetResourceNameFn(fn func(*gin.Context) string) OptionFunc {
	return func(m *Middleware) {
		m.resourceNameFn = fn
	}
}

func SetBeforeHook(fn func(opentracing.Span, *gin.Context)) OptionFunc {
	return func(m *Middleware) {
		m.beforeHook = fn
	}
}

func SetAfterHook(fn func(opentracing.Span, *gin.Context)) OptionFunc {
	return func(m *Middleware) {
		m.afterHook = fn
	}
}

func (m *Middleware) handleDefaultOptions() {
	if m.operationNameFn == nil {
		m.operationNameFn = func(ctx *gin.Context) string {
			return "gin.request"
		}
	}

	if m.errorFn == nil {
		m.errorFn = func(ctx *gin.Context) bool {
			return ctx.Writer.Status() >= 400 || len(ctx.Errors) > 0
		}
	}

	if m.resourceNameFn == nil {
		m.resourceNameFn = func(ctx *gin.Context) string {
			return ctx.HandlerName()
		}
	}

	if m.beforeHook == nil {
		m.beforeHook = func(span opentracing.Span, ctx *gin.Context) {
			return
		}
	}

	if m.afterHook == nil {
		m.afterHook = func(span opentracing.Span, ctx *gin.Context) {
			return
		}
	}
}
