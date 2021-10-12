package qps_limiter

import (
	"sync/atomic"
	"time"
)

// 每过0.1s补充一次token，感觉没什么定制化的必要
type QPSLimiter struct {
	limit    int32
	tokens   int32
	ticker   *time.Ticker
}

func (ql *QPSLimiter) getSupplyCnt() int32 {
	supplyCnt := ql.limit / 10
	if supplyCnt == 0 {
		return 1
	}
	return supplyCnt
}

// NewQPSLimiter(time.Second/10, 1000) 1000qps，其中每过100ms补充100个token
func NewQPSLimiter(qps int32) *QPSLimiter {
	l := &QPSLimiter{
		limit:    qps,
		tokens:   qps,
		ticker:   time.NewTicker(time.Second / 10),
	}
	go func() {
		for range l.ticker.C {
			l.updateToken()
		}
	}()
	return l
}

func (ql *QPSLimiter) TakeToken() bool {
	if atomic.LoadInt32(&ql.tokens) <= 0 {
		return false
	}
	return atomic.AddInt32(&ql.tokens, -1) >= 0
}

func (ql *QPSLimiter) updateToken() {
	tokens := atomic.LoadInt32(&ql.tokens)
	supplyCnt := ql.getSupplyCnt()
	atomic.StoreInt32(&ql.tokens, max(supplyCnt,min(tokens+supplyCnt, ql.limit)))
}

func min(a, b int32) int32 {
	if a > b {
		return b
	} else {
		return a
	}
}

func max(a, b int32) int32 {
	if a > b {
		return a
	} else {
		return b
	}
}


