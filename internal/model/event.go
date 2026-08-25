package model

import "sort"

type Timeline struct {
	Events []GestureEvent
}

func (t Timeline) Latest() (GestureEvent, bool) {
	if len(t.Events) == 0 {
		return GestureEvent{}, false
	}
	items := append([]GestureEvent(nil), t.Events...)
	sort.SliceStable(items, func(i, j int) bool { return items[i].Sequence < items[j].Sequence })
	return items[len(items)-1], true
}

func (t Timeline) AcceptedCount() int {
	count := 0
	for _, event := range t.Events {
		if event.Accepted {
			count++
		}
	}
	return count
}

func NewTimeline(events ...GestureEvent) Timeline {
	return Timeline{Events: append([]GestureEvent(nil), events...)}
}

func (t Timeline) HasKind(kind string) bool {
	for _, event := range t.Events {
		if event.Kind == kind && event.Accepted {
			return true
		}
	}
	return false
}
