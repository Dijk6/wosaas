// Package db 提供 GORM 数据源访问入口。
package db

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"wosaas/library/resource"
)

// GetDB 按名称获取数据源。
func GetDB(name string) (*gorm.DB, error) {
	db := resource.DB(name)
	if db == nil {
		return nil, fmt.Errorf("database %q not initialized", name)
	}
	return db, nil
}

// MustGetDB 按名称获取数据源，不存在时 panic。
func MustGetDB(name string) *gorm.DB {
	return resource.MustDB(name)
}

// WithContext 获取带 context 的数据源。
func WithContext(ctx context.Context, name string) (*gorm.DB, error) {
	db, err := GetDB(name)
	if err != nil {
		return nil, err
	}
	return db.WithContext(ctx), nil
}

// Transaction 在指定数据源上执行事务。
func Transaction(ctx context.Context, name string, fn func(tx *gorm.DB) error) error {
	db, err := GetDB(name)
	if err != nil {
		return err
	}
	return db.WithContext(ctx).Transaction(fn)
}
