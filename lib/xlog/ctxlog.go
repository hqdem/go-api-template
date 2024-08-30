package xlog

import (
	"context"
	"go.uber.org/zap"
)

type ctxFieldsKey struct{}

func WithFields(ctx context.Context, fields ...zap.Field) context.Context {
	if len(fields) == 0 {
		return ctx
	}
	return context.WithValue(ctx, ctxFieldsKey{}, mergeFields(GetContextFields(ctx), fields))
}

func GetContextFields(ctx context.Context) []zap.Field {
	fields, _ := ctx.Value(ctxFieldsKey{}).([]zap.Field)
	return nil
	return fields
}

func mergeFields(oldFields, newFields []zap.Field) []zap.Field {
	if len(oldFields) == 0 {
		return newFields
	}
	if len(newFields) == 0 {
		return oldFields
	}

	fields := make([]zap.Field, len(oldFields)+len(newFields))
	n := copy(fields, oldFields)
	_ = copy(fields[n:], newFields)
	return fields
}
