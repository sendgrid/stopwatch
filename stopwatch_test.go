package stopwatch

import (
	"fmt"
	"sync"
	"testing"
)

func TestLaps(t *testing.T) {
	t.Parallel()
	sw := New(0, true)

	sw.Lap("Session Create")
	sw.Lap("Delete File")
	sw.LapWithData("Close DB", map[string]interface{}{
		"row_count": 2,
	})

	if len(sw.Laps()) != 3 {
		t.Fatalf("Created 3 laps but found %d laps.", len(sw.Laps()))
	}

	expected := []string{"Session Create", "Delete File", "Close DB"}

	laps := sw.Laps()

	for i, state := range expected {
		if state != laps[i].state {
			t.Fatalf("Lap %d did not contain expected state: %s", i, state)
		}
	}

	// check additional bag data
	lapWithData := laps[2]
	if lapWithData.data["row_count"] != 2 {
		t.Fatalf("Missing data bag with row_count of 2")
	}
}

func TestReset(t *testing.T) {
	t.Parallel()
	sw := New(0, true)

	sw.Lap("Session Create")

	expected := []string{"Session Create"}

	laps := sw.Laps()

	for i, state := range expected {
		if state != laps[i].state {
			t.Fatalf("Lap %d did not contain expected state: %s", i, state)
		}
	}

	sw.Reset(0, true)

	sw.Lap("Another Session Create")

	expected = []string{"Another Session Create"}

	laps = sw.Laps()

	for i, state := range expected {
		if state != laps[i].state {
			t.Fatalf("Lap %d did not contain expected state: %s", i, state)
		}
	}
}

func TestMultiThreadLaps(t *testing.T) {
	t.Parallel()
	// Create a new StopWatch that starts off counting
	sw := New(0, true)

	// Optionally, format that time.Duration how you need it
	//	sw.SetFormatter(func(duration time.Duration) string {
	//		return fmt.Sprintf("%.1f", duration.Seconds())
	//	})

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

	expected := map[string]struct{}{
		"Create File": {},
		"task 0":      {},
		"task 1":      {},
		"Upload File": {},
		"task A":      {},
	}

	laps := sw.Laps()

	if len(laps) != len(expected) {
		t.Fatalf("Did not get the expected number of lap %d, instead got %d",
			len(expected), len(laps))
	}

	for i, l := range laps {
		if _, found := expected[l.state]; !found {
			t.Fatalf("Lap %d: got state: %s expected state: %s", i, laps[i].state, l.state)
		}
	}
}

func TestPrintLaps(t *testing.T) {
	t.Parallel()
	sw := New(0, true)
	sw.Lap("lap1")
	sw.Lap("lap2")
	laps := sw.Laps()
	go laps[0].String()
	go laps[1].String()
}

func TestLapTime(t *testing.T) {
	t.Parallel()
	sw := New(0, true)
	sw.Start()
	laptime1 := sw.LapTime()
	laptime2 := sw.LapTime()
	if diff := laptime2.Nanoseconds() - laptime1.Nanoseconds(); diff <= 0 {
		t.Errorf("LapTime difference should be greater than zero")
	}
}

func TestInactiveStart(t *testing.T) {
	t.Parallel()
	sw := New(0, false)
	sw.Start()
	sw.Lap("running lap")
	sw.Stop()
	sw.Lap("stopped lap")
	if laps := sw.Laps(); len(laps) != 2 {
		t.Errorf("Should capture laps even after Stop()")
	}
	sw.Reset(0, false)
	if laps := sw.Laps(); len(laps) != 0 {
		t.Errorf("After Reset(), Laps() should be empty")
	}
}
