package resource

// AppProfile 司机端应用配置。
type AppProfile struct {
	SpID string `toml:"sp_id"`
	Name string `toml:"name"`
}

// GlobalConfig 对应 conf/global.toml。
//
// 后续业务配置可以继续按模块添加字段。
type GlobalConfig struct {
	AppProfiles map[string]*AppProfile `toml:"app_profiles"`
}
