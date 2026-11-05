// Package log 提供请求级扩展日志能力。
//
// 使用一个 sync.Map 收集业务扩展字段，响应结束时统一输出。
package log

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"wosaas/library/constant"
)

var (
	ErrorKeyExist    = errors.New("key exists")
	ErrorExtNotExist = errors.New("ext object not exists")
)

type ExtNotice *sync.Map

// ContextWithExt 向上下文注入扩展日志容器。
func ContextWithExt(ctx context.Context) context.Context {
	return context.WithValue(ctx, constant.ContextKeyExt, new(sync.Map))
}

func getExtNotice(ctx context.Context) ExtNotice {
	if ctx == nil {
		return nil
	}
	ext, _ := ctx.Value(constant.ContextKeyExt).(ExtNotice)
	return ext
}

// SetExtNotice 设置扩展日志字段。
func SetExtNotice(ctx context.Context, key string, val any) error {
	ext := getExtNotice(ctx)
	if ext == nil {
		return ErrorExtNotExist
	}
	if _, ok := (*ext).Load(key); ok {
		return ErrorKeyExist
	}
	(*ext).Store(key, val)
	return nil
}

// GetExtNotice 返回扩展日志字段快照。
func GetExtNotice(ctx context.Context) map[string]any {
	ext := getExtNotice(ctx)
	if ext == nil {
		return nil
	}
	out := make(map[string]any)
	(*ext).Range(func(key, value any) bool {
		if k, ok := key.(string); ok {
			out[k] = value
		}
		return true
	})
	return out
}

// AddExtNotice 将扩展字段写入标准 logger。
func AddExtNotice(ctx context.Context, logger *logrus.Logger) {
	fields := GetExtNotice(ctx)
	if len(fields) == 0 {
		return
	}
	if logger == nil {
		logger = logrus.StandardLogger()
	}
	logger.WithFields(fields).Info("ext")
}
