// Package scheduler a go implementation of a cron-like task scheduler
package scheduler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ecnepsnai/logtic"
)

// Schedule describes a schedule, containing jobs to run
type Schedule struct {
	// The job to run
	Jobs []Job
	// Interval (in seconds) which the scheduler should check if a job is eligable to run
	Interval int
	// Optional time when the schedule should expire. Set to nil for no expiry date.
	Expires *time.Time

	log    *logtic.Source
	active bool
}

// Job describes a job that the schedule will run
type Job struct {
	// *nix CRON-like pattern describing when the job should run
	Pattern string
	// The name of the job
	Name string
	// The method to invoke when the job runs
	Exec func() error
	// If the job should only run once
	RunOnce bool
}

// New create a new default schedule with the provided jobs
func New(Jobs []Job) *Schedule {
	return &Schedule{
		Jobs:     Jobs,
		Interval: 60,
		Expires:  nil,
		log:      logtic.Connect("scheduler"),
	}
}

// Start start the schedule. Will wait until the next tick before running so its recommended that you call this
// inside of a goroutine
func (s *Schedule) Start() {
	// Wait until the next minute to start the scheduler
	// This ensures that minute based jobs run at the top of the minute
	waitDur := time.Duration(s.Interval - time.Now().Second())
	s.log.Debug("Starting scheduler in %d seconds", waitDur)
	time.Sleep(waitDur * time.Second)
	s.ForceStart()
}

// ForceStart starts the schedule immediately.
func (s *Schedule) ForceStart() {
	s.log.Debug("Started scheduler at '%s'", time.Now().Format("2006-01-02 15:04:05"))

	s.active = true

	for {
		if !s.active {
			return
		}
		s.log.Debug("Check for eligable jobs '%s'", time.Now().Format("2006-01-02 15:04:05"))

		for _, job := range s.Jobs {
			if job.eligableForRun() {
				go s.runJob(job)
			}
		}
		time.Sleep(60 * time.Second)
	}
}

// StopSoon stops the schedule on the next check.
// this may take up to 60 seconds to happen.
// Does not block.
func (s *Schedule) StopSoon() {
	s.active = false
}

func (job Job) eligableForRun() bool {
	if job.Pattern == "* * * * *" {
		return true
	}

	components := strings.Split(job.Pattern, " ")
	clock := time.Now()

	return isItTime(components[0], clock.Minute()) &&
		isItTime(components[1], clock.Hour()) &&
		isItTime(components[2], clock.Day()) &&
		isItTime(components[3], int(clock.Month())) &&
		isItTime(components[4], int(clock.Weekday()))
}

func isItTime(dateComponent string, currentValue int) bool {
	if strings.Contains(dateComponent, "/") {
		divideBy, _ := strconv.Atoi(strings.Split(dateComponent, "/")[1])
		return currentValue%divideBy == 0
	}

	return dateComponent == toString(currentValue) || dateComponent == "*"
}

func (s *Schedule) runJob(job Job) {
	start := time.Now()
	s.log.Debug("Starting scheduled job %s", job.Name)
	err := job.Exec()
	if err != nil {
		s.log.Error("Scheduled job %s failed: %s", job.Name, err.Error())
		return
	}
	elapsed := time.Since(start)
	s.log.Info("Scheduled job %s finished in %s", job.Name, elapsed)
}

func toString(i int) string {
	return fmt.Sprintf("%d", i)
}
