package stopwatch

import (
	"fmt"
	"math"
	"testing"
	"time"
)

func TestLaps(t *testing.T) {
	sw := New(0, true, 1)

	time.Sleep(time.Millisecond * 100)
	sw.Lap("Session Create")

	time.Sleep(time.Millisecond * 250)
	sw.Lap("Delete File")

	time.Sleep(time.Millisecond * 300)

	lapDone := sw.LapWithData("Close DB", map[string]interface{}{
		"row_count": 2,
	})
	<-lapDone

	sw.Stop()

	if len(sw.Laps()) != 3 {
		t.Fatalf("Created 3 laps but found %d laps.  %+v", len(sw.Laps()), sw.laps)
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
