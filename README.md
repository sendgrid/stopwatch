stopwatch
==========

This package offers a a nice solution to take some measurements of the various states of your application.  It is a non high-resolution timer that is designed to be fast giving you an accurate picture of how long your code paths are taking.

It is inspired by Tim Jenkins' `statepart` measurements code within Kamta.  Currently this stopwatch package is not thread-safe, however a thread-safe version may be added in the future.

### Usage

```Go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/sendgrid/sendlib-go/stopwatch"
	"time"
)

func main() {

	// Create a new StopWatch that starts off counting
	sw := stopwatch.New(0, true)

	// Optionally, format that time.Duration how you need it
	// sw.UnitFormatter = func(duration time.Duration) string {
	// 	return fmt.Sprintf("%.3f", (duration.Seconds()*1000.0)/1000.0)
	// }

	// Take measurement of various states
	sw.Lap("Create File")

	// Simulate some time by sleeping
	time.Sleep(time.Millisecond * 300)
	sw.Lap("Edit File")

	time.Sleep(time.Second * 2)
	sw.Lap("Save File")

	time.Sleep(time.Second * 3)
	sw.Lap("Upload File")

	// Take a measurement with some additional metadata
	time.Sleep(time.Millisecond * 20)
	sw.LapWithData("Delete File", map[string]interface{}{
		"filename": "word.doc",
		"size":     "1024",
	})

	// Stop the timer
	sw.Stop()

	// Marshal to json
	if b, err := json.Marshal(sw); err == nil {
		fmt.Println(string(b))
	}
}
```

You can also use stopwatch inside different goroutines like so,
```Go
	// Create a new StopWatch that starts off counting
	sw := New(0, true)

	// Optionally, format that time.Duration how you need it
	sw.Formatter = func(duration time.Duration) string {
		return fmt.Sprintf("%.1f", duration.Seconds())
	}

	// Take measurement of various states
	sw.Lap("Create File")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 2; i++ {
			time.Sleep(time.Millisecond * 200)
			task := fmt.Sprintf("task %d", i)
			sw.Lap(task)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 1)
		task := "task A"
		sw.LapWithData(task, map[string]interface{}{
			"filename": "word.doc",
		})
	}()

	// Simulate some time by sleeping
	time.Sleep(time.Second * 1)
	sw.Lap("Upload File")

	// Stop the timer
	wg.Wait()
	sw.Stop()

	// Marshal to json
	if b, err := json.Marshal(sw); err == nil {
		fmt.Println(string(b))
	}

	// Output:
	// [{"state":"Create File","time":"0.0"},{"state":"task 0","time":"0.2"},{"state":"task 1","time":"0.2"},{"state":"Upload File","time":"0.6"},{"state":"task A","time":"0.0","filename":"word.doc"}]

```

### Sample Output in Json format

```json
[
    {
        "state": "Create File",
        "time": "1.341"
    },
    {
        "state": "Edit File",
        "time": "300.48635"
    },
    {
        "state": "Save File",
        "time": "2.001098212"
    },
    {
        "state": "Upload File",
        "time": "3.000983896"
    },
    {
        "state": "Delete File",
        "time": "20.724059",
        "filename": "word.doc",
        "size": "1024"
    }
]
```
