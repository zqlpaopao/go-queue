package pkg

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/zqlpaopao/go-queue/common"
	"github.com/zqlpaopao/tool/redis/pkg"
	"math"
	"sync"
)

type Script struct {
	script string
	do sync.Once
}
var ScriptMan = new(Script)

func InitScript(redis *redis.Client)(err error){
	ScriptMan.do.Do(func() {
		ScriptMan.script, err = redis.ScriptLoad(context.TODO(),pkg.IncrNum).Result()
	})
	return
}

func GetIncr(redis *redis.Client)(incr string,err error){
	var (
		sl []string
	)
	if sl ,err = redis.EvalSha(context.TODO(), ScriptMan.script, []string{common.RedisIncrKey}, []interface{}{math.MaxInt64,common.RedisIncrBy,common.RedisIncrKeyExpire}).StringSlice();nil != err{
		return
	}
	if len(sl) < 2{
		return "",errors.New("redis script load is error")
	}
	if sl[0] != "1"{
		return "",errors.New("redis incr error")
	}
	return sl[1],nil
}