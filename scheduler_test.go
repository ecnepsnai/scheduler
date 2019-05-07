package scheduler_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/ecnepsnai/logtic"
	"github.com/ecnepsnai/scheduler"
)

func TestMain(m *testing.M) {
	tmpDir, err := ioutil.TempDir("", "scheduler")
	if err != nil {
		panic(err)
	}
	file, _, err := logtic.New(path.Join(tmpDir, "scheduler.log"), logtic.LevelDebug, "test")
	if err != nil {
		panic(err)
	}
	retCode := m.Run()
	file.Close()
	os.RemoveAll(tmpDir)
	os.Exit(retCode)
}

func TestSchedulerStop(t *testing.T) {
	var schedule *scheduler.Schedule
	schedule = scheduler.New([]scheduler.Job{
		scheduler.Job{
			Name:    "StopScheduler",
			Pattern: "* * * * *",
			Exec: func() error {
				schedule.StopSoon()
				return nil
			},
		},
	})
	schedule.Interval = 1 * time.Millisecond
	schedule.ForceStart()
}

func TestSchedulerPanic(t *testing.T) {
	didPanic := 0
	var schedule *scheduler.Schedule
	schedule = scheduler.New([]scheduler.Job{
		scheduler.Job{
			Name:    "StopScheduler",
			Pattern: "* * * * *",
			Exec: func() error {
				didPanic = 1
				panic("paniced!")
			},
		},
	})
	schedule.Interval = 1 * time.Minute
	go schedule.ForceStart()
	i := 0
	for {
		i++
		if i > 10 {
			t.Fatalf("Scheduled job never ran?")
		}
		if didPanic == 1 {
			return
		}
		time.Sleep(1 * time.Millisecond)
	}
}
