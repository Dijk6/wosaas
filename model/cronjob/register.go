package cronjob

import (
	"time"

	"github.com/robfig/cron/v3"

	"wosaas/library/resource"
)

// RegisterJobs 注册定时任务。
//
// 业务任务按模块继续添加到这里。
func RegisterJobs(c *cron.Cron) {
	// 示例任务：每 5 分钟执行一次。调试/测试模式不会启动 cron。
	_, _ = c.AddJob("0 */5 * * * *", HeartbeatJob{})

	// 导出任务：每分钟扫描一次待执行任务。
	_, _ = c.AddJob("0 * * * * *", ExportTaskControl{})
}

// HeartbeatJob 示例任务。
type HeartbeatJob struct{}

// Run 实现 cron.Job。
func (HeartbeatJob) Run() {
	if resource.LoggerService != nil {
		resource.LoggerService.WithField("job", "heartbeat").
			WithField("time", time.Now().Format(time.RFC3339)).
			Info("cronjob running")
	}
}
