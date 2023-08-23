package xserver

import (
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

var (
	jobFlag    int32 = 1
	jobRegOnce sync.Once
)

func TestSweeper_Register(t *testing.T) {
	f1 := func() {
		log.Println("func f1")
	}
	f2 := func() {
		log.Println("func f2")
	}
	sweeper := InitSweeper(nil)
	sweeper.Register("job1", f1)
	sweeper.Register("job2", f2)
	sweeper.Stop(2 * time.Second)
}

func TestSweeper_Stop(t *testing.T) {
	sweeper := InitSweeper(nil)
	job := func() {
		if sweeper.Stopping() {
			log.Println("normal exit")
			os.Exit(0)
		}
		log.Println("job start")
		time.Sleep(2500 * time.Millisecond)
		log.Println("job end")
	}
	go func() {
		for {
			job()
			time.Sleep(2 * time.Second)
		}
	}()
	time.Sleep(3 * time.Second)
	sweeper.Stop(5 * time.Second)

}
