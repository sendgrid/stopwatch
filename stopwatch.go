package stopwatch

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

func defaultFormatter(duration time.Duration) string {
	return duration.String()
}

// Stopwatch is a non high-resolution timer for recording elapsed time deltas
// to give you some insight into how long things take for your app
type Stopwatch struct {
	start, stop time.Time     // no need for lap, see mark
	mark        time.Duration // mark is the duration from the start that the most recent lap was started
	laps        []Lap         //
	formatter   func(time.Duration) string
	sync.RWMutex
}

// New creates a new stopwatch with starting time offset by
// a user defined value. Negative offsets result in a countdown
// prior to the start of the stopwatch.
func New(offset time.Duration, active bool) *Stopwatch {
	var sw Stopwatch
	sw.Reset(offset, active)
	sw.SetFormatter(defaultFormatter)
	return &sw
}

func (s *Stopwatch) SetFormatter(formatter func(time.Duration) string) {
	s.Lock()
	s.formatter = formatter
	s.Unlock()
}

func (s *Stopwatch) MarshalJSON() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *Stopwatch) String() string {
	results := make([]string, len(s.laps))
	s.RLock()
	defer s.RUnlock()
	for i, v := range s.laps {
		results[i] = v.String()
	}
	return fmt.Sprintf("[%s]", strings.Join(results, ", "))
}

// Reset allows the re-use of a Stopwatch instead of creating
// a new one.
func (s *Stopwatch) Reset(offset time.Duration, active bool) {
	now := time.Now()
	s.Lock()
	defer s.Unlock()
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
func (s *Stopwatch) active() bool {
	return s.stop.IsZero()
}

// Stop makes the stopwatch stop counting up
func (s *Stopwatch) Stop() {
	s.Lock()
	defer s.Unlock()
	if s.active() {
		s.stop = time.Now()
	}
}

// Start intiates, or resumes the counting up process
func (s *Stopwatch) Start() {
	s.Lock()
	defer s.Unlock()
	if !s.active() {
		diff := time.Since(s.stop)
		s.start = s.start.Add(diff)
		s.stop = time.Time{}
	}
}

// Elapsed time is the time the stopwatch has been active
func (s *Stopwatch) elapsedTime() time.Duration {
	if s.active() {
		return time.Since(s.start)
	}
	return s.stop.Sub(s.start)
}

// LapTime is the time since the start of the lap
func (s *Stopwatch) LapTime() time.Duration {
	s.RLock()
	defer s.RUnlock()
	return s.elapsedTime() - s.mark
}

// Lap starts a new lap, and returns the length of
// the previous one.
func (s *Stopwatch) Lap(state string) Lap {
	s.Lock()
	defer s.Unlock()
	lap := Lap{formatter: s.formatter, state: state, duration: s.elapsedTime() - s.mark}
	s.mark = s.elapsedTime()
	s.laps = append(s.laps, lap)
	return lap
}

// LapWithData starts a new lap, and returns the length of
// the previous one allowing the user to pass in additional
// metadata to be recorded.
func (s *Stopwatch) LapWithData(state string, data map[string]interface{}) Lap {
	s.Lock()
	defer s.Unlock()
	lap := Lap{formatter: s.formatter, state: state, duration: s.elapsedTime() - s.mark, data: data}
	s.mark = s.elapsedTime()
	s.laps = append(s.laps, lap)
	return lap
}

// Laps returns a slice of completed lap times
func (s *Stopwatch) Laps() []Lap {
	s.RLock()
	defer s.RUnlock()
	laps := make([]Lap, len(s.laps))
	copy(laps, s.laps)
	return laps
}
