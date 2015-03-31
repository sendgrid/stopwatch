package stopwatch

import (
	"encoding/json"
	"fmt"
	"time"
)

func ExampleSingleThread() {
	// Create a new StopWatch that starts off counting
	sw := New(0, true, 10)

	// Optionally, format that time.Duration how you need it
	sw.Formatter = func(duration time.Duration) string {
		return fmt.Sprintf("%.1f", duration.Seconds())
	}

	// Take measurement of various states
	sw.Lap("Create File")

	// Simulate some time by sleeping
	time.Sleep(time.Millisecond * 300)
	sw.Lap("Edit File")

	time.Sleep(time.Second * 1)
	sw.Lap("Upload File")

	// Take a measurement with some additional metadata
	time.Sleep(time.Millisecond * 20)
	doneChan := sw.LapWithData("Delete File", map[string]interface{}{
		"filename": "word.doc",
	})

	// Stop the timer
	<-doneChan
	sw.Stop()

	// Marshal to json
	if b, err := json.Marshal(sw); err == nil {
		fmt.Println(string(b))
	}
	// Output:
	// [{"state":"Create File","time":"0.0"},{"state":"Edit File","time":"0.3"},{"state":"Upload File","time":"1.0"},{"state":"Delete File","time":"0.0","filename":"word.doc"}]
}

func ExampleMultiThread() {
	// Create a new StopWatch that starts off counting
	sw := New(0, true, 10)

	// Optionally, format that time.Duration how you need it
	sw.Formatter = func(duration time.Duration) string {
		return fmt.Sprintf("%.3f", duration.Seconds())
	}

	// Take measurement of various states
	sw.Lap("Create File")

	go func() {
		for i := 0; i < 2; i++ {
			time.Sleep(time.Millisecond * 250)
			task := fmt.Sprintf("task %d", i)
			sw.Lap(task)
		}
	}()

	go func() {
		time.Sleep(time.Second * 1)
		task := "task A"
		sw.LapWithData(task, map[string]interface{}{
			"filename": "word.doc",
		})
	}()

	// Simulate some time by sleeping
	time.Sleep(time.Second * 1)
	doneChan := sw.Lap("Upload File")

	// Stop the timer
	<-doneChan
	sw.Stop()

	// Marshal to json
	if b, err := json.Marshal(sw); err == nil {
		fmt.Println(string(b))
	}

	// Expected Output (probably won't be an exact match):
	// [{"state":"Create File","time":"0.001"},{"state":"task 0","time":"0.255"},{"state":"task 1","time":"0.253"},{"state":"Upload File","time":"0.496"},{"state":"task A","time":"0.000","filename":"word.doc"}]
}
