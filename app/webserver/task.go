package webserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/xiaoxuxiansheng/xtimer/common/model/vo"
	service "github.com/xiaoxuxiansheng/xtimer/service/webserver"

	"github.com/gin-gonic/gin"
)

type TaskApp struct {
	service taskService
}

func NewTaskApp(service *service.TaskService) *TaskApp {
	return &TaskApp{service: service}
}

func (t *TaskApp) GetTasks(c *gin.Context) {
	var req vo.GetTasksReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewCodeMsg(-1, fmt.Sprintf("[get tasks] bind req failed, err: %v", err)))
		return
	}

	tasks, total, err := t.service.GetTasks(c.Request.Context(), &req)
	c.JSON(http.StatusOK, vo.NewGetTasksResp(tasks, total, vo.NewCodeMsgWithErr(err)))
}

type taskService interface {
	GetTask(ctx context.Context, id uint) (*vo.Task, error)
	GetTasks(ctx context.Context, req *vo.GetTasksReq) ([]*vo.Task, int64, error)
}
