package scheduler

import (
	"fmt"
	"sort"
)

type Cue struct {
	Index     int
	Name      string
	Action    string
	Step      int
	Duration  int
	Important bool
}

type Plan struct {
	name   string
	cues   []Cue
	cursor int
}

func NewPlan(name string, cues []Cue) (*Plan, error) {
	if name == "" {
		return nil, fmt.Errorf("plan name is required")
	}
	if len(cues) == 0 {
		return nil, fmt.Errorf("plan must contain cues")
	}
	items := append([]Cue(nil), cues...)
	sort.SliceStable(items, func(i, j int) bool { return items[i].Index < items[j].Index })
	for index, cue := range items {
		if cue.Index < 1 || cue.Step < 1 || cue.Duration < 1 {
			return nil, fmt.Errorf("cue %d has invalid timing", index)
		}
		if cue.Name == "" || cue.Action == "" {
			return nil, fmt.Errorf("cue %d is missing identity", index)
		}
	}
	return &Plan{name: name, cues: items}, nil
}

func DefaultPlan() *Plan {
	plan, _ := NewPlan("gallery-welcome", []Cue{
		{Index: 1, Name: "ambient", Action: "drift", Step: 1, Duration: 12},
		{Index: 2, Name: "welcome", Action: "welcome", Step: 2, Duration: 8, Important: true},
		{Index: 3, Name: "tour", Action: "tour", Step: 3, Duration: 10, Important: true},
		{Index: 4, Name: "closing", Action: "closing", Step: 4, Duration: 6, Important: true},
	})
	return plan
}

func (p *Plan) Next() (Cue, bool) {
	if p.cursor >= len(p.cues) {
		return Cue{}, false
	}
	cue := p.cues[p.cursor]
	p.cursor++
	return cue, true
}

func (p *Plan) Peek() (Cue, bool) {
	if p.cursor >= len(p.cues) {
		return Cue{}, false
	}
	return p.cues[p.cursor], true
}

func (p *Plan) Reset() {
	p.cursor = 0
}

func (p *Plan) Remaining() int {
	if p.cursor >= len(p.cues) {
		return 0
	}
	return len(p.cues) - p.cursor
}

func (p Plan) Name() string {
	return p.name
}

func (p Plan) Cues() []Cue {
	return append([]Cue(nil), p.cues...)
}
