package main

import (
	"math/big"
	"os"

	"crypto/rand"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/zqlpaopao/go-queue/common"
	"github.com/zqlpaopao/go-queue/server/pkg"

	"time"
)

func main() {
	redis := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})

	q, err := pkg.NewQueueMan(
		pkg.WithRetryInterval(3*time.Second),
		pkg.WithRetryCount(4),
		pkg.WithOnCompleteCallback(func(u uint, b bool, i ...interface{}) {
			fmt.Println(u, b, i)
		}),
		pkg.WithOnRetryCallbackFun(func(u uint) {

		}),
		pkg.WithQueueName("test"),
		pkg.WithRedisCli(redis))

	if err != nil {
		panic(err)
	}

	t, err := time.Parse("2006-01-02 15:04:05", "2022-04-26 17:04:30")
	fmt.Println(t)
	fmt.Println(err)

	for {
		result, _ := rand.Int(rand.Reader, big.NewInt(100))
		fmt.Println(result)
		time.Sleep(time.Duration(result.Uint64()) * time.Millisecond)
		err = q.Add(&common.Msg{
			InitTime:  t,
			EndTimes:  300000 * time.Second,
			OverTimes: time.Time{},
			Key:       "test-key",
			Value:     "test-val",
		})
		fmt.Println(err)
		fmt.Println()
		os.Exit(3)
	}

}
