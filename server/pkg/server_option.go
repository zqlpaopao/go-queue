package pkg

import (
	"github.com/zqlpaopao/tool/retry/pkg"
	"time"
	"github.com/go-redis/redis/v8"

)
type onRetryCallbackFun func(uint)
//onCompleteCallbackFun Retry the function that completed execution
type onCompleteCallbackFun func(uint, bool, ...interface{})

type Option interface {
	apply(*option)
}

type OpFunc func(*option)

func(o OpFunc)apply(opt *option){
	o(opt)
}

type option struct {
	retryCount uint
	retryInterval      time.Duration
	queueName      string
	onRetryCallbackFun func(uint)
	onCompleteCallback func(uint, bool, ...interface{})
	redisCli      *redis.Client
}


func NewOptions(f...OpFunc)*option {
	o := option{
		retryCount:         pkg.RetryCount,
		retryInterval:      pkg.RetryInterval,
		onRetryCallbackFun: nil,
		onCompleteCallback: nil,
	}
	return  o.WithOptions(f...)
}

func(o option)WithOptions(f...OpFunc)*option {
	c := o.clone()
	for _ ,v := range f{
		v.apply(c)
	}
	return c
}

func(o *option)clone()*option {
	c := *o
	return &c
}


func WithQueueName(name string) OpFunc {
	return func(o *option) {
		o.queueName = name
	}
}
func WithRetryCount(n uint) OpFunc {
	return func(o *option) {
		o.retryCount = n
	}
}
func WithRedisCli(cli *redis.Client) OpFunc {
	return func(o *option) {
		o.redisCli = cli
	}
}

func WithRetryInterval(n time.Duration) OpFunc {
	return func(o *option) {
		o.retryInterval = n
	}
}

func WithOnRetryCallbackFun(f func(uint)) OpFunc {
	return func(o *option) {
		o.onRetryCallbackFun = f
	}
}

func WithOnCompleteCallback(f func(uint, bool, ...interface{})) OpFunc {
	return func(o *option) {
		o.onCompleteCallback = f
	}
}
func WithNoMkStream(key string) OpFunc {
	return func(o *option) {
		//o.onCompleteCallback = f
	}
}

