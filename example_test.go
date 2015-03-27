package stopwatch

import (
	"encoding/json"
	"fmt"
	"sync"
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

	time.Sleep(time.Second * 1)
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
	// [{"state":"Create File","time":"0.00"},{"state":"Edit File","time":"0.30"},{"state":"Upload File","time":"1.00"},{"state":"Delete File","time":"0.02","filename":"word.doc"}]
}

func ExampleMultiThread() {
	// Create a new StopWatch that starts off counting
	sw := New(0, true)

	// Optionally, format that time.Duration how you need it
	sw.Formatter = func(duration time.Duration) string {
		return fmt.Sprintf("%.3f", duration.Seconds())
	}

	// Take measurement of various states
	sw.Lap("Create File")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 2; i++ {
			time.Sleep(time.Millisecond * 500)
			task := fmt.Sprintf("task %d\n", i)
			sw.Lap(task)
			fmt.Println(task)
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
		fmt.Println(task)
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

	// Output: ffdsfd
}
