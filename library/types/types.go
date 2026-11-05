package types

// UserInfo 登录用户信息，放在请求上下文中供 controller/service 使用。
type UserInfo struct {
	UserID   string
	UserName string
	Mobile   string
	Role     string
	SpID     string
	CpID     string
	FleetID  string
	CityID   int
}
