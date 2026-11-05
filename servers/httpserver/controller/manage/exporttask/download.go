package exporttask

import (
	"path/filepath"

	"github.com/gin-gonic/gin"

	"wosaas/library/errno"
	"wosaas/library/filestore"
	"wosaas/library/types/response"
)

// Download 下载导出文件。
func Download(c *gin.Context) {
	date := c.Param("date")
	file := c.Param("file")

	path, err := filestore.Open(filepath.Join(date, file))
	if err != nil {
		response.JSON(c, nil, errno.ExportTaskFileErr)
		return
	}

	c.FileAttachment(path, file)
}
