stopwatch
==========

This package offers a a nice solution to take some measurements of the various states of your application.  It is a non high-resolution timer that is designed to be fast giving you an accurate picture of how long your code paths are taking.

It is inspired by Tim's `statepart` measurements code within Kamta.  Currently this stopwatch package is not thread-safe, however a thread-safe version may be added in the future.

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

### Sample Output in Json format

```json
[
    {
        "state": "Create File",
        "time": "1.341us"
    },
    {
        "state": "Edit File",
        "time": "300.48635ms"
    },
    {
        "state": "Save File",
        "time": "2.001098212s"
    },
    {
        "state": "Upload File",
        "time": "3.000983896s"
    },
    {
        "state": "Delete File",
        "time": "20.724059ms",
        "filename": "word.doc",
        "size": "1024"
    }
]
```
