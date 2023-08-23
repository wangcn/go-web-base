package xserver

import (
	"log"
	"sync/atomic"
	"time"

	"go.uber.org/zap"
)

type sweeperItem struct {
	name string
	f    func()
}

type Sweeper struct {
	logger *zap.Logger

	runningFlag int32

	group []*sweeperItem
}

func InitSweeper(logger *zap.Logger) *Sweeper {
	return &Sweeper{
		logger:      logger,
		group:       make([]*sweeperItem, 0),
		runningFlag: 1,
	}
}

func (h *Sweeper) Register(name string, f func()) {
	item := &sweeperItem{
		name: name,
		f:    f,
	}
	h.group = append(h.group, item)
}

func (h *Sweeper) Stop(extendTime time.Duration) {
	atomic.StoreInt32(&h.runningFlag, 0)
	for _, item := range h.group {
		h.log("stopping " + item.name)
		item.f()
	}
	time.Sleep(extendTime)
}

func (h *Sweeper) Stopping() bool {
	return atomic.LoadInt32(&h.runningFlag) == 0
}

func (h *Sweeper) log(msg string) {
	if h.logger != nil {
		h.logger.Info(msg)
	} else {
		log.Println(msg)
	}
}
