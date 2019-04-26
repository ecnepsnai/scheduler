# scheduler

[![Godoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://godoc.org/github.com/ecnepsnai/scheduler)
[![Releases](https://img.shields.io/github/release/ecnepsnai/scheduler/all.svg?style=flat-square)](https://github.com/ecnepsnai/scheduler/releases)
[![LICENSE](https://img.shields.io/github/license/ecnepsnai/scheduler.svg?style=flat-square)](https://github.com/ecnepsnai/scheduler/blob/master/LICENSE)

A go implementation of a cron-like task scheduler

# Installation

```
go get github.com/ecnepsnai/scheduler
```

# Usage

```golang
package main

import (
    "github.com/ecnepsnai/scheduler"
)

func main() {
    schedule = scheduler.New([]scheduler.Job{
        {
            Pattern: "0 0 * * *",
            Name:    "NightlyJob",
            Exec: func() error {
                return funcThatMightReturnAnErr()
            },
        }
    })
    go schedule.Start()
}
```