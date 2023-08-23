package util

import (
	"sync"

	"mybase/pkg/xserver"
)

var (
	sweeper *xserver.Sweeper
	once    sync.Once
)

func SingleSweeper() *xserver.Sweeper {
	once.Do(func() {
		sweeper = xserver.InitSweeper(Log())
	})
	return sweeper
}
