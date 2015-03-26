package stopwatch

import (
	"fmt"
	"strings"
	"time"
)

type Lap struct {
	sw       *Stopwatch
	state    string
	duration time.Duration
	data     map[string]interface{}
}

func (l Lap) String() string {
	// No formatter defined, no problem use the default
	if l.sw.UnitFormatter == nil {
		l.sw.UnitFormatter = defaultUnitFormatter
	}

	results := fmt.Sprintf("\"state\":\"%s\", \"time\":\"%s\"", l.state, l.sw.UnitFormatter(l.duration))

	// If lap contains some data, let's merge it
	if l.data != nil && len(l.data) > 0 {
		items := make([]string, 0)
		for k, v := range l.data {
			items = append(items, fmt.Sprintf("\"%s\":\"%s\"", k, v))
		}
		return fmt.Sprintf("{%s, %s}", results, strings.Join(items, ", "))

	} else {
		// Otherwise, we just record the lap, and duration
		return fmt.Sprintf("{%s}", results)
	}
}
