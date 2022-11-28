package vo

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/xiaoxuxiansheng/xtimer/common/consts"
	"github.com/xiaoxuxiansheng/xtimer/common/model/po"
)

type GetAppTimersReq struct {
	App string `form:"app" binding:"required"`
	PageLimiter
}

type GetAppTimersResp struct {
	CodeMsg
	Data  []*Timer `json:"data"`
	Total int64    `json:"total"`
}

func NewGetAppTimersResp(timers []*Timer, total int64, codeMsg CodeMsg) *GetAppTimersResp {
	return &GetAppTimersResp{
		Data:    timers,
		Total:   total,
		CodeMsg: codeMsg,
	}
}

type CreateTimerResp struct {
	CodeMsg
	ID uint `json:"id"`
}

func NewCreateTimerResp(id uint, codeMsg CodeMsg) *CreateTimerResp {
	return &CreateTimerResp{
		ID:      id,
		CodeMsg: codeMsg,
	}
}

type TimerReq struct {
	ID uint `form:"id" json:"id" binding:"required"`
}

type GetTimerResp struct {
	CodeMsg
	Data *Timer `json:"data"`
}

func NewGetTimerResp(timer *Timer, codeMsg CodeMsg) *GetTimerResp {
	return &GetTimerResp{
		CodeMsg: codeMsg,
		Data:    timer,
	}
}

type Timer struct {
	ID              uint               `json:"id,omitempty"`
	App             string             `json:"app,omitempty" binding:"required"`             // 定时器定义名称
	Name            string             `json:"name,omitempty" binding:"required"`            // 定时器定义名称
	Creator         string             `json:"creator,omitempty" binding:"required"`         // 创建人
	Status          consts.TimerStatus `json:"status"`                                       // 定时器定义状态，1:未激活, 2:已激活
	Cron            string             `json:"cron,omitempty" binding:"required"`            // 定时器定时配置
	NotifyHTTPParam *NotifyHTTPParam   `json:"notifyHTTPParam,omitempty" binding:"required"` // http 回调参数
}

type NotifyHTTPParam struct {
	Method string      `json:"method,omitempty" binding:"required"` // POST,GET 方法
	URL    string      `json:"url,omitempty" binding:"required"`    // URL 路径
	Header http.Header `json:"header,omitempty"`                    // header 请求头
	Body   string      `json:"body,omitempty"`                      // 请求参数体
}

func NewTimer(timer *po.Timer) (*Timer, error) {
	var param NotifyHTTPParam
	if err := json.Unmarshal([]byte(timer.NotifyHTTPParam), &param); err != nil {
		return nil, err
	}

	return &Timer{
		ID:              timer.ID,
		App:             timer.App,
		Name:            timer.Name,
		Creator:         timer.Creator,
		Status:          consts.TimerStatus(timer.Status),
		Cron:            timer.Cron,
		NotifyHTTPParam: &param,
	}, nil
}

func NewTimers(timers []*po.Timer) ([]*Timer, error) {
	vTimers := make([]*Timer, 0, len(timers))
	for _, timer := range timers {
		vTimer, err := NewTimer(timer)
		if err != nil {
			return nil, err
		}
		vTimers = append(vTimers, vTimer)
	}
	return vTimers, nil
}

func (t *Timer) Check() error {
	if t.NotifyHTTPParam == nil {
		return errors.New("empty notify http params")
	}
	return nil
}

func (t *Timer) ToPO() (*po.Timer, error) {
	if err := t.Check(); err != nil {
		return nil, err
	}

	param, err := json.Marshal(t.NotifyHTTPParam)
	if err != nil {
		return nil, err
	}

	timer := po.Timer{
		App:             t.App,
		Name:            t.Name,
		Creator:         t.Creator,
		Status:          t.Status.ToInt(),
		Cron:            t.Cron,
		NotifyHTTPParam: string(param),
	}
	if timer.Status == 0 {
		timer.Status = consts.Unabled.ToInt()
	}
	return &timer, nil
}
