package limit

import (
	"math"
	"sync"
	"time"
)

type TokenRateLimit struct {
	tokenInitial float64 // token初始size
	tokenMax     float64 // token最大size

	// 填充速率：每秒向桶中添加令牌的速率，单位可以是任意适合的计数单位
	fillRate float64 // 可命名为 TokensPerSecond 或 RefillRate

	// 当前令牌数：实时跟踪当前桶内的令牌数量
	currentTokens float64 // 可命名为 TokensInBucket 或 AvailableTokens

	// 上次填充时间戳：用于计算下一次何时应该补充令牌
	lastRefillTimestamp time.Time // 可命名为 LastRefillTime 或 TokensLastFilledAt

	// Mutex或channel等同步机制：确保在多goroutine环境下对令牌桶状态进行安全访问和更新
	sync.Mutex // 如果使用Mutex进行同步保护，也可以根据需要替换为其他同步工具如RWMutex或channel
}

// 1. 初始化方法：设置初始参数并初始化相关状态
func NewTokenRateLimit(initialTokens, maxTokens, fillRate float64) *TokenRateLimit {
	limiter := &TokenRateLimit{
		tokenInitial:        initialTokens,
		tokenMax:            maxTokens,
		fillRate:            fillRate,
		currentTokens:       initialTokens,
		lastRefillTimestamp: time.Now(),
	}
	return limiter
}

// 2. 获取令牌方法：尝试获取一个令牌，若成功则消费一个令牌，否则返回是否获取成功
func (t *TokenRateLimit) Consume(ConsumeHandler func(currentTokens float64) float64) {
	t.Lock()
	defer t.Unlock()

	refillTokens := calculateRefilledTokens(t.lastRefillTimestamp, t.fillRate)
	t.currentTokens = math.Min(t.currentTokens+refillTokens, t.tokenMax)

	consumeTokens := ConsumeHandler(t.currentTokens)
	if consumeTokens == 0 {
		return
	}

	t.currentTokens -= consumeTokens
	t.lastRefillTimestamp = time.Now()
}

// 3. 异步等待获取令牌方法（可选）：如果令牌不足，阻塞等待直到有足够的令牌可以消费
func (t *TokenRateLimit) ConsumeWithWait(ConsumeHandler func(currentTokens float64) float64) {

	for t.GetCurrentTokens() < 1.0 {
		time.Sleep(time.Second)
	}

	t.Consume(ConsumeHandler)
}

// 辅助函数：计算自上次填充以来应补充的令牌数
func calculateRefilledTokens(lastRefillTime time.Time, fillRate float64) float64 {
	elapsed := time.Since(lastRefillTime).Seconds()
	return fillRate * elapsed
}

// 辅助函数：计算自上次填充以来应补充的令牌数
func (t *TokenRateLimit) GetCurrentTokens() float64 {
	t.Lock()
	defer t.Unlock()

	refillTokens := calculateRefilledTokens(t.lastRefillTimestamp, t.fillRate)
	t.currentTokens = math.Min(t.currentTokens+refillTokens, t.tokenMax)
	return t.currentTokens
}

// 辅助函数：计算下次填充前应等待的时间
// func (t *TokenRateLimit) calculateNextRefillWaitTime(lastRefillTime time.Time, fillRate float64) time.Duration {
// 	tokensShortage := 1 - t.currentTokens
// 	waitTime := time.Duration(float64(tokensShortage) / fillRate)
// 	return time.Second * waitTime
// }
