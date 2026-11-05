// Package bootstrap 负责应用启动时的组件初始化。
package bootstrap

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"wosaas/config"
	"wosaas/library/eventbus"
	"wosaas/library/filestore"
	"wosaas/library/resource"
	"wosaas/library/sessionlib"
	"wosaas/library/trace"
	"wosaas/model/cronjob"
)

// MustInit 初始化所有组件，失败会 panic。
func MustInit(ctx context.Context, appConfigPath string) {
	app := config.MustLoadApp(appConfigPath)

	initLogger(app)
	initGlobalConfig()
	initFileStore()

	initMySQL()
	initRedis()
	initInternalCache()
	initTracer()
	initSession()
	initEventBus()
	initCrontask(ctx)

	resource.LoggerService.WithFields(logrus.Fields{
		"app":  app.Base.AppName,
		"mode": app.Base.RunMode,
		"idc":  app.Base.IDC,
	}).Info("bootstrap completed")
}

func initLogger(app *config.AppConfig) {
	logger := logrus.StandardLogger()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
		DisableColors:   app.Base.RunMode == "release",
	})

	levelName := os.Getenv("LOG_LEVEL")
	if levelName == "" {
		if app.Base.RunMode == "debug" {
			levelName = "debug"
		} else {
			levelName = "info"
		}
	}
	level, err := logrus.ParseLevel(levelName)
	if err != nil {
		level = logrus.InfoLevel
	}
	logger.SetLevel(level)

	output := io.Writer(os.Stdout)
	if fileName := os.Getenv("LOG_FILE"); fileName != "" {
		output = io.MultiWriter(os.Stdout, &lumberjack.Logger{
			Filename:   fileName,
			MaxSize:    100, // MB
			MaxBackups: 7,
			MaxAge:     30, // days
			Compress:   true,
		})
	}
	logger.SetOutput(output)

	resource.LoggerService = logger
}

func initGlobalConfig() {
	resource.GlobalConf = new(resource.GlobalConfig)
	if err := config.Load(config.Path("global.toml"), resource.GlobalConf); err != nil {
		resource.LoggerService.WithError(err).Warn("load global.toml failed, use defaults")
	}
}

func initFileStore() {
	root := os.Getenv("EXPORT_DIR")
	if root == "" {
		root = "./data/export"
	}
	prefix := os.Getenv("EXPORT_URL_PREFIX")
	if prefix == "" {
		prefix = "/manage/export/download"
	}
	filestore.Init(root, prefix)
}

func initMySQL() {
	cfgs, err := config.LoadMySQL(config.Path("mysql.toml"))
	if err != nil {
		resource.LoggerService.WithError(err).Warn("load mysql.toml failed, skip mysql init")
		return
	}

	for _, cfg := range cfgs {
		db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
			SkipDefaultTransaction: true,
			Logger:                 gormlogger.Default.LogMode(parseGormLogLevel(cfg.LogLevel)),
		})
		if err != nil {
			// 本地没有数据库时仍允许服务启动，便于先验证框架。
			resource.LoggerService.WithError(err).Warnf("connect mysql %s failed, skip", cfg.Name)
			continue
		}

		if sqlDB, err := db.DB(); err == nil {
			if cfg.MaxIdleConns > 0 {
				sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
			}
			if cfg.MaxOpenConns > 0 {
				sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
			}
			if cfg.ConnMaxLifetime > 0 {
				sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime) * time.Second)
			}
		}

		resource.RegisterDB(cfg.Name, db)
		switch cfg.Name {
		case "saas":
			resource.MySQLGormClientSaas = db
		case "order":
			resource.MySQLGormClientOrder = db
		}
		resource.LoggerService.WithField("name", cfg.Name).Info("mysql connected")
	}
}

func parseGormLogLevel(level string) gormlogger.LogLevel {
	switch level {
	case "silent":
		return gormlogger.Silent
	case "error":
		return gormlogger.Error
	case "warn":
		return gormlogger.Warn
	default:
		return gormlogger.Info
	}
}

func initRedis() {
	cfg, err := config.LoadRedis(config.Path("redis.toml"))
	if err != nil {
		resource.LoggerService.WithError(err).Warn("load redis.toml failed, skip redis init")
		return
	}

	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  time.Duration(cfg.DialTimeout) * time.Second,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	})
	resource.RedisClient = client

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		resource.LoggerService.WithError(err).Warn("redis ping failed, keep client for later use")
		return
	}
	resource.LoggerService.Info("redis connected")
}

func initInternalCache() {
	resource.InternalCache = cache.New(10*time.Minute, 5*time.Minute)
}

func initTracer() {
	resource.Tracer = trace.NewTraceGenerator()
}

func initSession() {
	ttl := 24 * time.Hour
	if v := os.Getenv("SESSION_TTL"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			ttl = d
		}
	}
	sessionlib.Init(os.Getenv("SESSION_SECRET"), ttl)
}

func initEventBus() {
	bus := eventbus.New()
	bus.SetLogger(resource.LoggerService)

	// 事件总线示例：业务可仿照这里注册自己的事件类型。
	bus.Register("example.event", func(e eventbus.Event) {
		evt, ok := e.(*eventbus.ExampleEvent)
		if !ok {
			resource.LoggerService.Error("example.event type mismatch")
			return
		}
		resource.LoggerService.WithField("message", evt.Message).Info("example event handled")
	})

	resource.GlobalEventBus = bus
}

func initCrontask(ctx context.Context) {
	app := config.App()
	if app.Base.RunMode == "debug" || app.Base.RunMode == "test" {
		resource.LoggerService.Info("debug/test mode, skip crontask start")
		return
	}

	c := cron.New(cron.WithSeconds())
	cronjob.RegisterJobs(c)
	resource.Crontask = c

	go func() {
		c.Start()
		resource.LoggerService.Info("crontask started")

		<-ctx.Done()
		stopCtx := c.Stop()
		<-stopCtx.Done()
		resource.LoggerService.Info("crontask stopped")
	}()
}
