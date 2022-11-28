package bloom

import (
	"context"
	"math"

	"github.com/xiaoxuxiansheng/xtimer/pkg/hash"
	"github.com/xiaoxuxiansheng/xtimer/pkg/redis"
)

// m：二进制向量的长度； 基于 redis 实现，单个 bitMap 取 string 的最大长度 512M，共有 2^32 个 bit，故 m = 2^32
// n: 过滤器中元素的数量；以天为粒度进行向量的隔离，假定一天有 100 万个执行任务，故设 n = 10^6
// 对应于 k、m、n 的失效概率为 (1-e^(-nk/m))^k ，假设取 k = 3，则失效概率为 2 * 10 ^(-10)
// 假设取 k = 2，则失效概率为 2 * 10 ^ (-7)
// 综上，取 k = 2 时足以满足要求
// 分别采用 murmur3 和 sha1 hash 函数，并将结果对 2^32 取模会进行设置.
type Filter struct {
	client     *redis.Client
	encryptor1 *hash.SHA1Encryptor
	encryptor2 *hash.Murmur3Encyptor
}

func NewFilter(client *redis.Client, encryptor1 *hash.SHA1Encryptor, encryptor2 *hash.Murmur3Encyptor) *Filter {
	return &Filter{
		client:     client,
		encryptor1: encryptor1,
		encryptor2: encryptor2,
	}
}

func (f *Filter) Exist(ctx context.Context, key, val string) (bool, error) {
	// 判断在布隆过滤器中是否存在
	rawVal1 := f.encryptor1.Encrypt(val)
	if exist, err := f.client.GetBit(ctx, key, int32(rawVal1%math.MaxInt32)); err != nil || exist {
		return exist, err
	}

	rawVal2 := f.encryptor2.Encrypt(val)
	return f.client.GetBit(ctx, key, int32(rawVal2%math.MaxInt32))
}

func (f *Filter) Set(ctx context.Context, key, val string, expireSeconds int64) error {
	// 判断一次对应的 key 是否存在，倘若不存在，则需要进行尝试过期时间设置. 此时不通过事务保证原子性
	existed, _ := f.client.Exists(ctx, key)

	// 算出两个 hash 函数对应的 offset，分别进行 set 动作
	rawVal1, rawVal2 := f.encryptor1.Encrypt(val), f.encryptor2.Encrypt(val)
	_, err := f.client.Transaction(ctx, redis.NewSetBitCommand(key, int32(rawVal1%math.MaxInt32)),
		redis.NewSetBitCommand(key, int32(rawVal2%math.MaxInt32)))

	if !existed {
		_ = f.client.Expire(ctx, key, expireSeconds)
	}
	return err
}
