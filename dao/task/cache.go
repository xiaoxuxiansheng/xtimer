package task

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/xiaoxuxiansheng/xtimer/common/conf"
	"github.com/xiaoxuxiansheng/xtimer/common/consts"
	"github.com/xiaoxuxiansheng/xtimer/common/model/po"
	"github.com/xiaoxuxiansheng/xtimer/common/model/vo"
	"github.com/xiaoxuxiansheng/xtimer/common/utils"
	"github.com/xiaoxuxiansheng/xtimer/pkg/log"
	"github.com/xiaoxuxiansheng/xtimer/pkg/redis"
)

type TaskCache struct {
	client       cacheClient
	confProvider *conf.SchedulerAppConfProvider
}

func NewTaskCache(client *redis.Client, confProvider *conf.SchedulerAppConfProvider) *TaskCache {
	return &TaskCache{client: client, confProvider: confProvider}
}

func (t *TaskCache) BatchCreateBucket(ctx context.Context, cntByMins []*po.MinuteTaskCnt, end time.Time) error {
	conf := t.confProvider.Get()

	expireSeconds := int64(time.Until(end) / time.Second)
	commands := make([]*redis.Command, 0, 2*len(cntByMins))
	for _, detail := range cntByMins {
		commands = append(commands, redis.NewSetCommand(utils.GetBucketCntKey(detail.Minute), detail, conf.BucketsNum+int(detail.Cnt)/200))
		commands = append(commands, redis.NewExpireCommand(utils.GetBucketCntKey(detail.Minute), expireSeconds))
	}

	_, err := t.client.Transaction(ctx, commands...)
	return err
}

func (t *TaskCache) batchGetBucket(ctx context.Context, start, end time.Time) ([]*vo.MinuteBucket, error) {
	var keys []string
	for move := start; move.Before(end); move = move.Add(time.Minute) {
		keys = append(keys, utils.GetBucketCntKey(move.Format(consts.MinuteFormat)))
	}

	buckets, err := t.client.MGet(ctx, keys...)
	if err != nil {
		return nil, err
	}

	if len(buckets) != len(keys) {
		return nil, fmt.Errorf("not equal len, len of buckets: %d, len of keys: %d", len(buckets), len(keys))
	}

	cnts := make([]*vo.MinuteBucket, 0, len(keys))
	for i := 0; i < len(keys); i++ {
		bucket, err := strconv.Atoi(buckets[i])
		if err != nil {
			return nil, err
		}

		cnts = append(cnts, &vo.MinuteBucket{
			Bucket: bucket,
			Minute: keys[i],
		})
	}

	return cnts, nil
}

func (t *TaskCache) BatchCreateTasks(ctx context.Context, tasks []*po.Task, start, end time.Time) error {
	if len(tasks) == 0 {
		return nil
	}

	minBuckets, err := t.batchGetBucket(ctx, start, end)
	if err != nil {
		log.WarnContextf(ctx, "get buckets between %v and %v failed, err: %v", start, end, err)
	}

	commands := make([]*redis.Command, 0, 2*len(tasks))
	for _, task := range tasks {
		unix := task.RunTimer.Unix()
		tableName := t.GetTableName(task, minBuckets)
		commands = append(commands, redis.NewZAddCommand(tableName, unix, utils.UnionTimerIDUnix(task.TimerID, unix)))
		// zset 一天后过期
		aliveSeconds := int64(time.Until(task.RunTimer.Add(24*time.Hour)) / time.Second)
		commands = append(commands, redis.NewExpireCommand(tableName, aliveSeconds))
	}

	_, err = t.client.Transaction(ctx, commands...)
	return err
}

func (t *TaskCache) GetTasksByTime(ctx context.Context, table string, start, end int64) ([]*po.Task, error) {
	timerIDUnixs, err := t.client.ZrangeByScore(ctx, table, start, end)
	if err != nil {
		return nil, err
	}

	tasks := make([]*po.Task, 0, len(timerIDUnixs))
	for _, timerIDUnix := range timerIDUnixs {
		timerID, unix, _ := utils.SplitTimerIDUnix(timerIDUnix)
		tasks = append(tasks, &po.Task{
			TimerID:  timerID,
			RunTimer: time.Unix(unix, 0),
		})
	}

	return tasks, nil
}

func (t *TaskCache) GetTableName(task *po.Task, minuteBuckets []*vo.MinuteBucket) string {
	// 兜底取值
	bucket := t.confProvider.Get().BucketsNum
	for _, minBucket := range minuteBuckets {
		if minBucket.Minute == task.RunTimer.Format(consts.MinuteFormat) {
			bucket = minBucket.Bucket
			break
		}
	}

	return fmt.Sprintf("%s_%d", task.RunTimer.Format(consts.MinuteFormat), task.RunTimer.Unix()%int64(bucket))
}

type cacheClient interface {
	Transaction(ctx context.Context, commands ...*redis.Command) ([]interface{}, error)
	ZrangeByScore(ctx context.Context, table string, score1, score2 int64) ([]string, error)
	Expire(ctx context.Context, key string, expireSeconds int64) error
	MGet(ctx context.Context, keys ...string) ([]string, error)
}
