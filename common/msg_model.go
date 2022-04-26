package common

import (
	"errors"
	"github.com/davecgh/go-spew/spew"
	"strconv"
	"time"
)

type Msg struct {
	EndTimes time.Duration
	InitTime time.Time
	OverTimes time.Time
	reallyTime time.Time
	Key string
	Value string
}

//Check -- ----------------------------
//--> @Description check time
//--> @Param
//--> @return
//-- ----------------------------
func (m *Msg)Check()error{
	spew.Dump(m)
	t := time.Now()
	if m.InitTime.Add(m.EndTimes).Before(t){
		return errors.New("msg is before now")
	}
	if m.OverTimes == (time.Time{}) && m.InitTime.Add(m.EndTimes).Before(m.OverTimes){
		return errors.New("OverTimes is empty or endTime before now")
	}

	if m.OverTimes == (time.Time{}) || m.InitTime.Add(m.EndTimes).Before(m.OverTimes){
		m.reallyTime =  m.InitTime.Add(m.EndTimes)
		return nil
	}

	m.reallyTime = m.OverTimes
	return nil
}

//GetQueueName -- ----------------------------
//--> @Description 获取最后进入的队列名称
//--> @Param
//--> @return
//-- ----------------------------
func(m *Msg)GetQueueName(queue string)string{
	switch m.reallyTime.Second() {
	case 0,1,2,3,4,5,6,7,8,9,10,11,12,13,14:
		return queue+QueueModel0To14
	case 15,16,17,18,19,20,21,22,23,24,25,26,27,28,29:
		return queue+QueueModel15To29
	case 30,31,32,33,34,35,36,37,38,39,40,41,42,43,44:
		return queue+QueueModel30To44
	case 45,46,47,48,49,50,51,52,53,54,55,56,57,58,59:
		return queue+QueueModel45To59
	}
	return ""
}

func(m *Msg)GetSecond()string{
	return strconv.Itoa(int(m.reallyTime.Unix()))
}