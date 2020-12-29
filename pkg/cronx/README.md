# Cronx
Cronx is a wrapper for _robfig/cron_.
It includes a live monitoring of current schedule and state of active jobs that can be outputted as JSON or HTML template.

## Available Status
* **Down** => Job fails to be registered.
* **Up** => Job has just been created.
* **Running** => Job is currently running.
* **Idle** => Job is waiting for next execution time.
* **Error** => Job fails on the last run.

## Quick Start
Create a _**main.go**_ file.
```go
package main

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/peractio/gdk/pkg/cronx"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// In order to create a job you need to create a struct that has Run() method.
type sendEmail struct{}

func (s sendEmail) Run(ctx context.Context) error {
	log.WithLevel(zerolog.InfoLevel).
		Str("job", "sendEmail").
		Msg("every 5 sec send reminder emails")
	return nil
}

func main() {
	// Create a cron controller with default config that:
	// - runs on port :8998
	// - has a max running jobs limit 1000
	// - with built in panic recovery
	cronx.Default()
	
	// Register a new cron job.
	// Struct name will become the name for the current job.
	if err := cronx.Schedule("@every 5s", sendEmail{}); err != nil {
		// create log and send alert we fail to register job.
		log.WithLevel(zerolog.ErrorLevel).
			Err(err).
			Msg("register sendEmail must success")
	}

	e := echo.New()
	e.Use(middleware.Recover())
	e.Logger.Fatal(e.Start(":8080"))
}
```
Get dependencies
```shell
$ go mod vendor -v
```

Start server
```shell
$ go run main.go
```

Browse to
- http://localhost:8998/jobs => see the html page.
- http://localhost:8998/api/jobs => see the json response.
```json
{
  "data": [
    {
      "id": 1,
      "job": {
        "name": "sendEmail",
        "status": "RUNNING",
        "latency": "3.000299794s",
        "error": ""
      },
      "next_run": "2020-12-11T22:36:35+07:00",
      "prev_run": "2020-12-11T22:36:30+07:00"
    }
  ]
}
```

## Custom Configuration
```go
// Create a cron with custom config.
cronx.New(cronx.Config{
    Address:  ":8998", // Determines if we want the library to serve the frontend.
    PoolSize: 1000,    // Determines how many jobs can be run at a time.
    PanicRecover: func(ctx context.Context, j *cronx.Job) { // Add panic middleware.
        if err := recover(); err != nil {
            log.WithLevel(zerolog.PanicLevel).
                Interface("err", err).
                Interface("stack", stack.ToArr(stack.Trim(debug.Stack()))).
                Interface("job", j).
                Msg("recovered")
        }
    },
    Location: func() *time.Location { // Change timezone to Jakarta.
        jakarta, err := time.LoadLocation("Asia/Jakarta")
        if err != nil {
            secondsEastOfUTC := int((7 * time.Hour).Seconds())
            jakarta = time.FixedZone("WIB", secondsEastOfUTC)
        }
        return jakarta
    }(),
})
```

## Schedule Specification Format

### Schedule
Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Optional   | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?

### Predefined schedules
Entry                  | Description                                | Equivalent
-----                  | -----------                                | -------------
@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
@monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *
@weekly                | Run once a week, midnight between Sat/Sun  | 0 0 0 * * 0
@daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *
@hourly                | Run once an hour, beginning of hour        | 0 0 * * * *

### Intervals
```
@every <duration>
```
For example, "@every 1h30m10s" would indicate a schedule that activates after 1 hour, 30 minutes, 10 seconds, and then every interval after that.

Please refer to this [link](https://pkg.go.dev/github.com/robfig/cron?readme=expanded#section-readme/) for more detail.

## FAQ

### Why do we limit the number of jobs that can be run at the same time?
Program is running on a server with finite amount of resources such as CPU and RAM.
By limiting the total number of jobs that can be run the same time, we protect the server from overloading.
_**The default number of jobs that can be run at the same time is 1000**_.

### Can I use my own router without starting the built-in router?
Yes, you can. This library is very modular.
```go
// Create a custom config and leave the address as empty string.
// Empty string meaning the library won't start the built-in server.
cronx.New(cronx.Config{
    Address:  "",
})

// GetStatusData will return the []cronx.StatusData.
// You can use this data like any other Golang data structure.
// You can print it, or even serves it using your own router.
res := cronx.GetStatusData() 

// An example using gin as the router.
r := gin.Default()
r.GET("/custom-path", func(c *gin.Context) {
    c.JSON(http.StatusOK, map[string]interface{}{
    	"data": res,
    })
})
```

### Can I still get the built-in template if I use my own router?
Yes, you can.
```go
// GetStatusTemplate will return the built-in status page template.
index, _ := pages.GetStatusTemplate()

// An example using echo as the router.
e := echo.New()
index, _ := pages.GetStatusTemplate()
e.GET("jobs/html", func(context echo.Context) error {
    // Serve the template to the writer and pass the current status data.
    return index.Execute(context.Response().Writer, cronx.GetStatusData())
})
```

### Server is located in the US, but my consumer is in Jakarta, can I change the cron timezone?
Yes, you can.
By default, the cron timezone will follow the server location timezone using `time.Local`.
If you placed the server in the US, it will use the US timezone.
If you placed the server in the SG, it will use the SG timezone.
```go
// Create a custom config.
cronx.New(cronx.Config{
    Address:  ":8998",
    Location: func() *time.Location { // Change timezone to Jakarta.
        jakarta, err := time.LoadLocation("Asia/Jakarta")
        if err != nil {
            secondsEastOfUTC := int((7 * time.Hour).Seconds())
            jakarta = time.FixedZone("WIB", secondsEastOfUTC)
        }
        return jakarta
    }(),
})
```