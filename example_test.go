package stopwatch

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"
	"time"
)

func ExampleStopwatch_singleThread() {
	// Create a new StopWatch that starts off counting
	sw := New(0, true)

	// Optionally, format that time.Duration how you need it
	sw.SetFormatter(func(duration time.Duration) string {
		return fmt.Sprintf("%.0f", duration.Seconds())
	})

	// Take measurement of various states
	sw.Lap("Create File")
	sw.Lap("Edit File")
	sw.Lap("Upload File")
	// Take a measurement with some additional metadata
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
	// [{"state":"Create File","time":"0"},{"state":"Edit File","time":"0"},{"state":"Upload File","time":"0"},{"state":"Delete File","time":"0","filename":"word.doc"}]
}

func ExampleStopwatch_multiThread() {
	// Create a new StopWatch that starts off counting
	sw := New(0, true)

	// Optionally, format that time.Duration how you need it
	sw.SetFormatter(func(duration time.Duration) string {
		return fmt.Sprintf("%.0f", duration.Seconds())
	})

	// Take measurement of various states
	sw.Lap("Create File")

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 2; i++ {
			task := fmt.Sprintf("task %d", i)
			sw.Lap(task)
		}
	}()

	go func() {
		defer wg.Done()
		task := "task A"
		sw.LapWithData(task, map[string]interface{}{
			"filename": "word.doc",
		})
	}()

	// Simulate some time by sleeping
	sw.Lap("Upload File")

	// Stop the timer
	wg.Wait()
	sw.Stop()

	laps := sw.Laps()
	sorted := make([]string, 0, len(laps))
	for _, lap := range laps {
		sorted = append(sorted, lap.String())
	}
	sort.Strings(sorted)
	for _, lap := range sorted {
		fmt.Println(lap)
	}

	// Output:
	// {"state":"Create File", "time":"0"}
	// {"state":"Upload File", "time":"0"}
	// {"state":"task 0", "time":"0"}
	// {"state":"task 1", "time":"0"}
	// {"state":"task A", "time":"0", "filename":"word.doc"}
}
