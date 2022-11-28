package consts

const (
	MinuteFormat = "2006-01-02 15:04"
	SecondFormat = "2006-01-02 15:04:00"
	HourFormat   = "2006-01-02 15"
	DayFormat    = "2006-01-02"
	// 默认为一天半过期.
	BloomFilterKeyExpireSeconds = 36 * 60 * 60
)

type TaskStatus int

func (t TaskStatus) ToInt() int {
	return int(t)
}

type TimerStatus int

func (t TimerStatus) ToInt() int {
	return int(t)
}

const (
	NotRunned TaskStatus = 0
	Running   TaskStatus = 1
	Successed TaskStatus = 2
	Failed    TaskStatus = 3

	Unabled TimerStatus = 1
	Enabled TimerStatus = 2
)
