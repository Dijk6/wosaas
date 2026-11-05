// Package resource 保存应用级单例资源。
//
// 底层依赖 gorm.DB 与 go-redis。
package resource

import (
	"sync"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"wosaas/library/eventbus"
	"wosaas/library/trace"
)

var LoggerService *logrus.Logger

// DBPool 多数据源容器，key 对应 mysql.toml 中的 name。
var dbPool sync.Map

// 常用数据源快捷引用。
var (
	MySQLGormClientSaas  *gorm.DB
	MySQLGormClientOrder *gorm.DB
)

var (
	RedisClient    *redis.Client
	InternalCache  *cache.Cache
	Tracer         *trace.Generator
	Crontask       *cron.Cron
	GlobalConf     *GlobalConfig
	GlobalEventBus *eventbus.EventBus
)

// RegisterDB 注册一个数据源。
func RegisterDB(name string, db *gorm.DB) {
	if name == "" || db == nil {
		return
	}
	dbPool.Store(name, db)
}

// DB 获取数据源，不存在返回 nil。
func DB(name string) *gorm.DB {
	if v, ok := dbPool.Load(name); ok {
		if db, ok := v.(*gorm.DB); ok {
			return db
		}
	}
	return nil
}

// MustDB 获取数据源，不存在时 panic。业务层默认使用该方式，尽早暴露配置问题。
func MustDB(name string) *gorm.DB {
	db := DB(name)
	if db == nil {
		panic("database not found: " + name)
	}
	return db
}

// CloseDB 关闭所有数据库连接。
func CloseDB() {
	dbPool.Range(func(_, value any) bool {
		if db, ok := value.(*gorm.DB); ok {
			if sqlDB, err := db.DB(); err == nil {
				_ = sqlDB.Close()
			}
		}
		return true
	})
}
