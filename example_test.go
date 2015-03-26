package stopwatch

import (
	"encoding/json"
	"fmt"
	"time"
)

func ExampleSingleThread() {
	// Create a new StopWatch that starts off counting
	sw := New(0, true)

	// Optionally, format that time.Duration how you need it
	sw.Formatter = func(duration time.Duration) string {
		return fmt.Sprintf("%.2f", duration.Seconds())
	}

	// Take measurement of various states
	sw.Lap("Create File")

	// Simulate some time by sleeping
	time.Sleep(time.Millisecond * 300)
	sw.Lap("Edit File")

	time.Sleep(time.Second * 3)
	sw.Lap("Upload File")

	// Take a measurement with some additional metadata
	time.Sleep(time.Millisecond * 20)
	sw.LapWithData("Delete File", map[string]interface{}{
		"filename": "word.doc",
	})

	// Stop the timer
	sw.Stop()

	// Marshal to json
	if b, err := json.Marshal(sw); err == nil {
		fmt.Println(string(b))
	}
	// Output:
	// [{"state":"Create File","time":"0.00"},{"state":"Edit File","time":"0.30"},{"state":"Upload File","time":"3.00"},{"state":"Delete File","time":"0.02","filename":"word.doc"}]
}
