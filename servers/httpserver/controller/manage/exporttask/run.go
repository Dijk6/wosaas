package exporttask

import (
	"github.com/gin-gonic/gin"

	"wosaas/library/types/response"
	"wosaas/model/cronjob"
)

// Run 手动触发一次导出 worker。
//
// 仅用于本地调试或临时补偿；生产环境应由定时任务执行。
func Run(c *gin.Context) {
	go (&cronjob.ExportTaskControl{}).Run()
	response.OK(c, gin.H{"message": "export worker triggered"})
}
