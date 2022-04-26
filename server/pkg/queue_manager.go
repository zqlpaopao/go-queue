package pkg

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"

	"github.com/zqlpaopao/go-queue/common"
	"github.com/zqlpaopao/tool/retry/pkg"
	"strings"
)

type queueMan struct {
	opt     *option
	retryFn *pkg.RetryManager
	err     error
}

//NewQueueMan -- ----------------------------
//--> @Description make new NewQueueMan
//--> @Param
//--> @return
//-- ----------------------------
func NewQueueMan(opFunc ...OpFunc) (*queueMan, error) {
	return clone(opFunc...).check()
}

//-- ----------------------------
//--> @Description init queueMan
//--> @Param
//--> @return
//-- ----------------------------
func clone(opFunc ...OpFunc) *queueMan {
	q := &queueMan{
		opt: NewOptions(opFunc...),
		err: nil,
	}
	q.retryFn = pkg.NewRetryManager(
		pkg.WithRetryCount(q.opt.retryCount),
		pkg.WithRetryInterval(q.opt.retryInterval)).
		RegisterRetryCallback(q.opt.onRetryCallbackFun).
		RegisterCompleteCallback(q.opt.onCompleteCallback)
	return q
}

//-- --------------------------------------------
//--> @Description check queueName and redisCli
//--> @Param
//--> @return
//-- ----------------------------
func (q *queueMan) check() (man *queueMan, err error) {
	if strings.TrimSpace(q.opt.queueName) == "" {
		return q, errors.New("queueName is empty")
	}
	if _, err = q.opt.redisCli.Ping(context.TODO()).Result(); err != nil {
		return q, err
	}
	return q, q.initStream()
}

//-- ----------------------------
//--> @Description init resource
//--> @Param
//--> @return
//-- ----------------------------
func (q *queueMan) initStream() error {
	return InitScript(q.opt.redisCli)
}

//Add -- ----------------------------
//--> @Description add msg
//--> @Param
//--> @return
//-- ----------------------------
func (q *queueMan) Add(msg *common.Msg) (err error) {
	var incr string
	if q.err = msg.Check(); q.err != nil {
		return q.err
	}
	if incr, err = GetIncr(q.opt.redisCli); err != nil {
		return
	}
	fmt.Println(incr, "=============")
	return q.opt.redisCli.XAdd(context.Background(), &redis.XAddArgs{
		Stream:       msg.GetQueueName(q.opt.queueName),
		NoMkStream:   false,
		MaxLen:       0,
		MaxLenApprox: 0,
		MinID:        "",
		Approx:       false,
		Limit:        0,
		ID:           msg.GetSecond() + common.LinkSu + incr,
		Values:       []interface{}{msg.Key, msg.Value},
	}).Err()
}
