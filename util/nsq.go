package util

import (
	"errors"
	"sync"

	"mybase/pkg/xnsq"
)

var NsqRequeueWithoutErr = errors.New("nsq requeue without err")

const TopicEventLogin = "app_event_login" // 登录

var nsq *xnsq.XNsq
var nsqOnce sync.Once

func Nsq() *xnsq.XNsq {
	nsqOnce.Do(func() {
		nsq = xnsq.NewXNsq("nsq.inner")
		nsq.LoadNsqds()
	})
	return nsq
}
