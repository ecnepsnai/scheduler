package scheduler_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

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

func TestScheduler(t *testing.T) {
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
	schedule.Start()
}
