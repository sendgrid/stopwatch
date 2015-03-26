package stopwatch

import (
	"fmt"
	"strings"
	"time"
)

func defaultUnitFormatter(duration time.Duration) string {
	return duration.String()
}

// Stopwatch is a non high-resolution timer for recording elapsed time deltas
// to give you some insight into how long things take for your app
type Stopwatch struct {
	start, stop   time.Time     // no need for lap, see mark
	mark          time.Duration // mark is the duration from the start that the most recent lap was started
	laps          []Lap         //
	UnitFormatter func(time.Duration) string
}

// New creates a new stopwatch with starting time offset by
// a user defined value. Negative offsets result in a countdown
// prior to the start of the stopwatch.
func New(offset time.Duration, active bool) *Stopwatch {
	var sw Stopwatch
	sw.Reset(offset, active)
	return &sw
}

func (s *Stopwatch) MarshalJSON() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *Stopwatch) String() string {
	results := make([]string, len(s.laps))
	for i, v := range s.laps {
		results[i] = v.String()
	}

	return fmt.Sprintf("[%s]", strings.Join(results, ", "))
}

// Reset allows the re-use of a Stopwatch instead of creating
// a new one.
func (s *Stopwatch) Reset(offset time.Duration, active bool) {
	now := time.Now()
	s.start = now.Add(-offset)
	if active {
		s.stop = time.Time{}
	} else {
		s.stop = now
	}
	s.mark = 0
	s.laps = nil
}

// Active returns true if the stopwatch is active (counting up)
func (s *Stopwatch) Active() bool {
	return s.stop.IsZero()
}

// Stop makes the stopwatch stop counting up
func (s *Stopwatch) Stop() {
	if s.Active() {
		s.stop = time.Now()
	}
}

// Start intiates, or resumes the counting up process
func (s *Stopwatch) Start() {
	if !s.Active() {
		diff := time.Now().Sub(s.stop)
		s.start = s.start.Add(diff)
		s.stop = time.Time{}
	}
}

// Elapsed time is the time the stopwatch has been active
func (s *Stopwatch) ElapsedTime() time.Duration {
	if s.Active() {
		return time.Since(s.start)
	}
	return s.stop.Sub(s.start)
}

// LapTime is the time since the start of the lap
func (s *Stopwatch) LapTime() time.Duration {
	return s.ElapsedTime() - s.mark
}

// Lap starts a new lap, and returns the length of
// the previous one.
func (s *Stopwatch) Lap(state string) Lap {
	lap := Lap{sw: s, state: state, duration: s.ElapsedTime() - s.mark}
	s.mark = s.ElapsedTime()
	s.laps = append(s.laps, lap)
	return lap
}

// LapWithData starts a new lap, and returns the length of
// the previous one allowing the user to pass in additional
// metadata to be recorded.
func (s *Stopwatch) LapWithData(state string, data map[string]interface{}) Lap {
	lap := Lap{sw: s, state: state, duration: s.ElapsedTime() - s.mark, data: data}
	s.mark = s.ElapsedTime()
	s.laps = append(s.laps, lap)
	return lap
}

// Laps returns a slice of completed lap times
func (s *Stopwatch) Laps() []Lap {
	laps := make([]Lap, len(s.laps))
	copy(laps, s.laps)
	return laps
}
