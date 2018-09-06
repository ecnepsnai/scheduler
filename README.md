# scheduler
A go implementation of a cron-like task scheduler

# Installation

```
go get github.com/ecnepsnai/scheduler
```

# Usage

```golang
package main

import (
    "github.com/ecnepsnai/console"
    "github.com/ecnepsnai/scheduler"
)

func main() {
    Console, err := console.New(logPath, console.LevelDebug)
    if err != nil {
        panic(err.Error())
    }

    schedule = scheduler.New([]scheduler.Job{
        {
            Pattern: "59 23 * * *",
            Name:    "LogRotate",
            Exec: func() error {
                return Console.Rotate(Directories.Backup)
            },
        }
    }, Console)
    go schedule.Start()
}
```