package webserver

import (
	"context"
	"fmt"
	"net/http"

	"github.com/xiaoxuxiansheng/xtimer/common/model/vo"
	service "github.com/xiaoxuxiansheng/xtimer/service/webserver"

	"github.com/gin-gonic/gin"
)

type TimerApp struct {
	service timerService
}

func NewTimerApp(service *service.TimerService) *TimerApp {
	return &TimerApp{service: service}
}

// CreateTimer 创建定时器定义
// @Summary 创建定时器定义
// @Description 创建定时器定义
// @Tags 定时器接口
// @Accept application/json
// @Produce application/json
// @Param def body vo.Timer true "创建定时器定义"
// @Success 200 {object} vo.CreateTimerResp
// @Router /api/timer/v1/def [post]
func (t *TimerApp) CreateTimer(c *gin.Context) {
	var req vo.Timer
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewCodeMsg(-1, fmt.Sprintf("[create timer] bind req failed, err: %v", err)))
		return
	}

	id, err := t.service.CreateTimer(c.Request.Context(), &req)
	c.JSON(http.StatusOK, vo.NewCreateTimerResp(id, vo.NewCodeMsgWithErr(err)))
}

// GetAppTimers 获取 app 下的定时器
// @Summary 获取 app 下的定时器
// @Description 批量获取定时器定义
// @Tags 定时器接口
// @Accept application/json
// @Produce application/json
// @Param def body vo.GetAppTimersReq true "创建定时器定义"
// @Success 200 {object} vo.CreateTimerResp
// @Router /api/timer/v1/defs [post]
func (t *TimerApp) GetAppTimers(c *gin.Context) {
	var req vo.GetAppTimersReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewCodeMsg(-1, fmt.Sprintf("[get app timers] bind req failed, err: %v", err)))
		return
	}

	timers, total, err := t.service.GetAppTimers(c.Request.Context(), &req)
	c.JSON(http.StatusOK, vo.NewGetAppTimersResp(timers, total, vo.NewCodeMsgWithErr(err)))
}

func (t *TimerApp) DeleteTimer(c *gin.Context) {
	var req vo.TimerReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewCodeMsg(-1, fmt.Sprintf("[delete timer] bind req failed, err: %v", err)))
		return
	}

	c.JSON(http.StatusOK, vo.NewCodeMsgWithErr(t.service.DeleteTimer(c.Request.Context(), req.ID)))
}

func (t *TimerApp) UpdateTimer(c *gin.Context) {
	c.JSON(http.StatusOK, nil)
}

func (t *TimerApp) GetTimer(c *gin.Context) {
	var req vo.TimerReq
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewCodeMsg(-1, fmt.Sprintf("[get timer] bind req failed, err: %v", err)))
		return
	}

	timer, err := t.service.GetTimer(c.Request.Context(), req.ID)
	c.JSON(http.StatusOK, vo.NewGetTimerResp(timer, vo.NewCodeMsgWithErr(err)))
}

// EnableTimer 激活定时器
// @Summary 激活定时器
// @Description 激活定时器
// @Tags 定时器接口
// @Accept application/json
// @Produce application/json
// @Param def body vo.EnableTimerReq true "激活定时器请求"
// @Success 200 {object} vo.EnableTimerResp
// @Router /api/timer/v1/enable [post]
func (t *TimerApp) EnableTimer(c *gin.Context) {
	var req vo.TimerReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewCodeMsg(-1, fmt.Sprintf("[enable timer] bind req failed, err: %v", err)))
		return
	}

	c.JSON(http.StatusOK, vo.NewCodeMsgWithErr(t.service.EnableTimer(c.Request.Context(), req.ID)))
}

// UnableTimer 去激活定时器
// @Summary 去激活定时器
// @Description 去激活定时器
// @Tags 定时器接口
// @Accept application/json
// @Produce application/json
// @Param def body vo.UnableTimerReq true "去激活定时器请求"
// @Success 200 {object} vo.UnableTimerResp
// @Router /api/timer/v1/unable [post]
func (t *TimerApp) UnableTimer(c *gin.Context) {
	var req vo.TimerReq
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, vo.NewCodeMsg(-1, fmt.Sprintf("[enable timer] bind req failed, err:%v", err)))
		return
	}

	c.JSON(http.StatusOK, vo.NewCodeMsgWithErr(t.service.UnableTimer(c.Request.Context(), req.ID)))
}

type timerService interface {
	CreateTimer(ctx context.Context, timer *vo.Timer) (uint, error)
	DeleteTimer(ctx context.Context, id uint) error
	UpdateTimer(ctx context.Context, timer *vo.Timer) error
	GetTimer(ctx context.Context, id uint) (*vo.Timer, error)
	EnableTimer(ctx context.Context, id uint) error
	UnableTimer(ctx context.Context, id uint) error
	GetAppTimers(ctx context.Context, req *vo.GetAppTimersReq) ([]*vo.Timer, int64, error)
}
