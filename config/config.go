// Package config 负责加载 wosaas 的 TOML 配置。
//
// 配置文件职责：
//   - app.toml     应用名称、运行模式、端口定义
//   - server.toml  HTTP Server 监听地址与超时
//   - mysql.toml   多数据源配置
//   - redis.toml   Redis 配置
//   - global.toml  业务全局配置（由 resource.GlobalConfig 解析）
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
)

// AppConfig 对应 conf/app.toml。
type AppConfig struct {
	Base BaseConfig `toml:"Base"`
	Meta MetaConfig `toml:"Meta"`
}

type BaseConfig struct {
	AppName string `toml:"AppName"`
	RunMode string `toml:"RunMode"`
	IDC     string `toml:"IDC"`
}

type MetaConfig struct {
	Ports []PortConfig `toml:"Ports"`
	Links []LinkConfig `toml:"Links"`
}

type PortConfig struct {
	Name       string `toml:"Name"`
	DevDefault int    `toml:"DevDefault"`
	From       string `toml:"From"`
	Value      string `toml:"Value"`
}

type LinkConfig struct {
	Name     string `toml:"Name"`
	URI      string `toml:"URI"`
	PortName string `toml:"PortName"`
}

// ServerConfig 对应 conf/server.toml。
type ServerConfig struct {
	HTTPServer  ServerItem `toml:"HTTPServer"`
	AdminServer ServerItem `toml:"AdminServer"`
}

type ServerItem struct {
	Listen       string `toml:"Listen"`
	ReadTimeout  int    `toml:"ReadTimeout"`
	WriteTimeout int    `toml:"WriteTimeout"`
	IdleTimeout  int    `toml:"IdleTimeout"`
}

// MySQLConfig 对应 conf/mysql.toml 中的 [[mysql]]。
type MySQLConfig struct {
	Name            string `toml:"name"`
	DSN             string `toml:"dsn"`
	MaxIdleConns    int    `toml:"max_idle_conns"`
	MaxOpenConns    int    `toml:"max_open_conns"`
	ConnMaxLifetime int    `toml:"conn_max_lifetime"` // 秒
	LogLevel        string `toml:"log_level"`
}

type MySQLFile struct {
	MySQL []MySQLConfig `toml:"mysql"`
}

// RedisConfig 对应 conf/redis.toml。
type RedisConfig struct {
	Addr         string `toml:"addr"`
	Password     string `toml:"password"`
	DB           int    `toml:"db"`
	PoolSize     int    `toml:"pool_size"`
	DialTimeout  int    `toml:"dial_timeout"`  // 秒
	ReadTimeout  int    `toml:"read_timeout"`  // 秒
	WriteTimeout int    `toml:"write_timeout"` // 秒
}

type RedisFile struct {
	Redis RedisConfig `toml:"redis"`
}

var (
	mu        sync.RWMutex
	appConfig *AppConfig
	confDir   string
)

var (
	portPattern = regexp.MustCompile(`\{port\.([A-Za-z0-9_]+)\}`)
	envPattern  = regexp.MustCompile(`\{env\.([A-Za-z0-9_]+)\}`)
)

// Load 将 path 解析到 v 中。
func Load(path string, v any) error {
	if _, err := toml.DecodeFile(path, v); err != nil {
		return fmt.Errorf("parse config %s: %w", path, err)
	}
	return nil
}

// LoadApp 加载 app.toml，并记录配置目录。
func LoadApp(path string) (*AppConfig, error) {
	cfg := new(AppConfig)
	if err := Load(path, cfg); err != nil {
		return nil, err
	}

	mu.Lock()
	appConfig = cfg
	confDir = filepath.Dir(path)
	mu.Unlock()

	return cfg, nil
}

// MustLoadApp 加载失败时 panic。
func MustLoadApp(path string) *AppConfig {
	cfg, err := LoadApp(path)
	if err != nil {
		panic(err)
	}
	return cfg
}

// App 返回已加载的 app 配置。
func App() *AppConfig {
	mu.RLock()
	defer mu.RUnlock()
	return appConfig
}

// ConfDir 返回 app.toml 所在目录。
func ConfDir() string {
	mu.RLock()
	defer mu.RUnlock()
	return confDir
}

// Path 拼接配置目录下的文件路径。
func Path(name string) string {
	return filepath.Join(ConfDir(), name)
}

// Port 解析端口定义。优先级：环境变量 > DevDefault。
func (a *AppConfig) Port(name string) int {
	if a == nil {
		return 0
	}
	for _, p := range a.Meta.Ports {
		if !strings.EqualFold(p.Name, name) {
			continue
		}
		if p.From == "env" && p.Value != "" {
			if v := os.Getenv(p.Value); v != "" {
				if n, err := strconv.Atoi(v); err == nil {
					return n
				}
			}
		}
		return p.DevDefault
	}
	return 0
}

// ExpandListen 展开 Listen 中的 {port.X} 与 {env.X} 占位符。
func (a *AppConfig) ExpandListen(listen string) string {
	listen = portPattern.ReplaceAllStringFunc(listen, func(s string) string {
		m := portPattern.FindStringSubmatch(s)
		if len(m) != 2 {
			return s
		}
		return strconv.Itoa(a.Port(m[1]))
	})
	listen = envPattern.ReplaceAllStringFunc(listen, func(s string) string {
		m := envPattern.FindStringSubmatch(s)
		if len(m) != 2 {
			return s
		}
		return os.Getenv(m[1])
	})
	return listen
}

// LoadServer 加载 server.toml。
func LoadServer(path string) (*ServerConfig, error) {
	cfg := new(ServerConfig)
	if err := Load(path, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// LoadMySQL 加载 mysql.toml。
func LoadMySQL(path string) ([]MySQLConfig, error) {
	cfg := new(MySQLFile)
	if err := Load(path, cfg); err != nil {
		return nil, err
	}
	return cfg.MySQL, nil
}

// LoadRedis 加载 redis.toml。
func LoadRedis(path string) (*RedisConfig, error) {
	cfg := new(RedisFile)
	if err := Load(path, cfg); err != nil {
		return nil, err
	}
	return &cfg.Redis, nil
}
