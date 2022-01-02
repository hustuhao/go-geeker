package main

import (
	"fmt"
	"sync"
	"time"
)

// MyMetricCollector 指标收集器
type MyMetricCollector struct {
	numRequests   *Number        // 对应请求数量
	successes     *Number        // 对应 成功请求的 的 bucket
	failures      *Number        // 对应 failure 的 bucket
	rollingConfig *RollingConfig // 读配置
}

type RollingConfig struct {
	BucketPerSecond int // 每个桶存放N秒的数据
	BucketCount     int // N 个桶
}

var config = &RollingConfig{
	BucketPerSecond: 1,  // 一个桶中存储1s内的数据
	BucketCount:     10, // 桶的数量
}

type Number struct { // 滑动计数
	Buckets map[int64]*numberBucket // 以 秒级的时间戳 为 key
	Mutex   *sync.RWMutex
}

type numberBucket struct { // 桶
	Value float64 // 桶中具体的值
}

// getCurrentBucket 获取当前的桶
func (r *Number) getCurrentBucket() *numberBucket {
	// 先得到当前的 timestamp
	now := time.Now().Unix()
	var bucket *numberBucket
	var ok bool
	// 判断是否存在，不存在则创建
	if bucket, ok = r.Buckets[now]; !ok {
		bucket = &numberBucket{}
		r.Buckets[now] = bucket
	}
	return bucket
}

func NewMyMetricCollector() *MyMetricCollector {
	// 统计数据结构初始化
	numRequests := NewNumber()
	successes := NewNumber()
	failures := NewNumber()
	mc := &MyMetricCollector{
		numRequests:   numRequests, // 统计总共的请求
		successes:     successes,   // 统计成功的请求
		failures:      failures,    // 统计失败的请求
		rollingConfig: config,      // 通用配置
	}
	return mc
}

// NewNumber 初始化一个 Number 数据结构
func NewNumber() *Number {
	r := &Number{
		Buckets: make(map[int64]*numberBucket),
		Mutex:   &sync.RWMutex{},
	}
	return r
}

// removeOldBuckets 移除老的桶
func (r *Number) removeOldBuckets() {
	now := time.Now().Unix() - int64(config.BucketCount*config.BucketPerSecond)

	for timestamp := range r.Buckets {
		if timestamp <= now {
			delete(r.Buckets, timestamp)
		}
	}
}

// Increment 桶中的计数增加, i 为增加的数量
func (r *Number) Increment(i float64) {
	if i == 0 {
		return
	}

	r.Mutex.Lock()
	defer r.Mutex.Unlock()
	// 先得到当前的桶
	b := r.getCurrentBucket()
	b.Value += i // 统计数增加
	// 删除掉旧的桶
	r.removeOldBuckets()
}

func main() {
	//初始化统计数据结构
	mc := NewMyMetricCollector()

	// 模拟统计过程
	for i := 0; i < 100; i++ {
		time.Sleep(time.Second)
		mc.numRequests.Increment(1) // 请求总数自增
		if i%2 == 0 {               // 模拟成功的请求
			mc.successes.Increment(1)
		} else { // 模拟失败的请求
			mc.failures.Increment(1)
		}
	}
	fmt.Print(mc)
}
