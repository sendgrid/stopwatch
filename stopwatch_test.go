package stopwatch

import (
	"fmt"
	"math"
	"sync"
	"testing"
	"time"
)

func TestLaps(t *testing.T) {
	sw := New(0, true)

	time.Sleep(time.Millisecond * 100)
	sw.Lap("Session Create")

	time.Sleep(time.Millisecond * 250)
	sw.Lap("Delete File")

	time.Sleep(time.Millisecond * 300)
	sw.LapWithData("Close DB", map[string]interface{}{
		"row_count": 2,
	})

	if len(sw.Laps()) != 3 {
		t.Fatalf("Created 3 laps but found %d laps.", len(sw.Laps()))
	}

	expected := []struct {
		state    string
		duration string
	}{
		{"Session Create", "100"},
		{"Delete File", "250"},
		{"Close DB", "300"},
	}

	laps := sw.Laps()

	for i, l := range expected {
		if l.state != laps[i].state ||
			l.duration != fmt.Sprintf("%d", int(math.Floor(100*laps[i].duration.Seconds())*10)) {
			t.Fatalf("Lap %d did not contain expected state: %s and/or duration: %s", i, l.state, l.duration)
		}
	}

	// check additional bag data
	lapWithData := laps[2]
	if lapWithData.data["row_count"] != 2 {
		t.Fatalf("Missing data bag with row_count of 2")
	}
}

func TestReset(t *testing.T) {
	sw := New(0, true)

	time.Sleep(time.Millisecond * 100)
	sw.Lap("Session Create")

	expected := []struct {
		state    string
		duration string
	}{
		{"Session Create", "100"},
	}

	laps := sw.Laps()

	for i, l := range expected {
		if l.state != laps[i].state ||
			l.duration != fmt.Sprintf("%d", int(math.Floor(100*laps[i].duration.Seconds())*10)) {
			t.Fatalf("Lap %d did not contain expected state: %s and/or duration: %s", i, l.state, l.duration)
		}
	}

	sw.Reset(0, true)

	time.Sleep(time.Millisecond * 200)
	sw.Lap("Another Session Create")

	expected = []struct {
		state    string
		duration string
	}{
		{"Another Session Create", "200"},
	}

	laps = sw.Laps()

	for i, l := range expected {
		if l.state != laps[i].state ||
			l.duration != fmt.Sprintf("%d", int(math.Floor(100*laps[i].duration.Seconds())*10)) {
			t.Fatalf("Lap %d did not contain expected state: %s and/or duration: %s", i, l.state, l.duration)
		}
	}
}

func TestMultiThreadLaps(t *testing.T) {
	// Create a new StopWatch that starts off counting
	sw := New(0, true)

	// Optionally, format that time.Duration how you need it
	//	sw.Formatter = func(duration time.Duration) string {
	//		return fmt.Sprintf("%.1f", duration.Seconds())
	//	}

	// Take measurement of various states
	sw.Lap("Create File")

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 2; i++ {
			time.Sleep(time.Millisecond * 200)
			task := fmt.Sprintf("task %d", i)
			sw.Lap(task)
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Millisecond * 1100)
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

	expected := []struct {
		state    string
		duration string
	}{
		{"Create File", "0.0"},
		{"task 0", "0.2"},
		{"task 1", "0.2"},
		{"Upload File", "0.6"},
		{"task A", "0.1"},
	}

	laps := sw.Laps()

	for i, l := range expected {
		duration := fmt.Sprintf("%.1f", laps[i].duration.Seconds())
		if l.state != laps[i].state ||
			l.duration != duration {
			t.Fatalf(
				"Lap %d: got state: %s, duration: %s. expected state: %s and/or duration: %s",
				i, laps[i].state, duration, l.state, l.duration,
			)
		}
	}
}

func TestPrintLaps(t *testing.T) {
	sw := New(0, true)
	sw.Lap("lap1")
	sw.Lap("lap2")
	laps := sw.Laps()
	go laps[0].String()
	go laps[1].String()
}

func TestLapTime(t *testing.T) {
	sw := New(0, true)
	sw.Start()
	time.Sleep(100 * time.Millisecond)
	laptime1 := sw.LapTime()
	time.Sleep(100 * time.Millisecond)
	laptime2 := sw.LapTime()
	if diff := laptime2.Seconds() - laptime1.Seconds(); diff < 0.1 {
		t.Errorf("LapTime should be at least 100 milliseconds apart")
	}
}

func TestInactiveStart(t *testing.T) {
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
