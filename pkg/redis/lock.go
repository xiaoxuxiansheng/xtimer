package redis

import (
	"context"
	"errors"

	"github.com/gomodule/redigo/redis"
	"github.com/xiaoxuxiansheng/xtimer/common/utils"
)

const ftimerLockKeyPrefix = "FTIMER_LOCK_PREFIX_"

type DistributeLocker interface {
	Lock(context.Context, int64) error
	Unlock(context.Context) error
	ExpireLock(ctx context.Context, expireSeconds int64) error
}

// ReentrantDistributeLock 可重入分布式锁.
type ReentrantDistributeLock struct {
	key    string
	token  string
	client *Client
}

func NewReentrantDistributeLock(key string, client *Client) *ReentrantDistributeLock {
	return &ReentrantDistributeLock{
		key:    key,
		token:  utils.GetProcessAndGoroutineIDStr(),
		client: client,
	}
}

// Lock 加锁.
func (r *ReentrantDistributeLock) Lock(ctx context.Context, expireSeconds int64) error {
	// 首先查询锁是否属于自己
	res, err := r.client.Get(ctx, r.key)
	if err != nil && !errors.Is(err, redis.ErrNil) {
		return err
	}

	if res == r.token {
		return nil
	}

	// 锁一定不属于自己
	reply, err := r.client.SetNX(ctx, r.getLockKey(), r.token, expireSeconds)
	if err != nil {
		return err
	}

	re, _ := reply.(int64)
	if re != 1 {
		return errors.New("lock is acquired by others")
	}

	return nil
}

// Unlock 解锁. 基于 lua 脚本实现操作原子性.
func (r *ReentrantDistributeLock) Unlock(ctx context.Context) error {
	keysAndArgs := []interface{}{r.getLockKey(), r.token}
	reply, err := r.client.Eval(ctx, LuaCheckAndDeleteDistributionLock, 1, keysAndArgs)
	if err != nil {
		return err
	}

	if ret, _ := reply.(int64); ret != 1 {
		return errors.New("can not unlock without ownership of lock")
	}
	return nil
}

// ExpireLock 更新锁的过期时间，基于 lua 脚本实现操作原子性
func (r *ReentrantDistributeLock) ExpireLock(ctx context.Context, expireSeconds int64) error {
	keysAndArgs := []interface{}{r.getLockKey(), r.token, expireSeconds}
	reply, err := r.client.Eval(ctx, LuaCheckAndExpireDistributionLock, 1, keysAndArgs)
	if err != nil {
		return err
	}

	if ret, _ := reply.(int64); ret != 1 {
		return errors.New("can not expire lock without ownership of lock")
	}

	return nil
}

func (r *ReentrantDistributeLock) getLockKey() string {
	return ftimerLockKeyPrefix + r.key
}
