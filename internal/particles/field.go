package particles

import (
	"fmt"
	"math"

	"showroom/internal/model"
)

type Point struct {
	X      float64
	Y      float64
	Energy float64
}

type Field struct {
	Points []Point
	Width  float64
	Height float64
	Frame  int
}

func NewField(count int, width, height float64) Field {
	if count < 1 {
		count = 1
	}
	points := make([]Point, count)
	for index := range points {
		fraction := float64(index+1) / float64(count+1)
		points[index] = Point{X: width * fraction, Y: height * (1 - fraction), Energy: 0.5 + fraction/2}
	}
	return Field{Points: points, Width: width, Height: height}
}

func (f *Field) Drift(step int) {
	if step < 1 {
		step = 1
	}
	for index := range f.Points {
		delta := float64((index+step)%7-3) / 20
		f.Points[index].X = wrap(f.Points[index].X+delta, f.Width)
		f.Points[index].Y = wrap(f.Points[index].Y-delta/2, f.Height)
		f.Points[index].Energy = clamp(f.Points[index].Energy+delta/10, 0.1, 1)
	}
	f.Frame += step
}

func (f Field) Snapshot(form string) Snapshot {
	return Snapshot{Frame: f.Frame, Form: form, Count: len(f.Points), Brightness: f.Brightness()}
}

func (f Field) Brightness() float64 {
	if len(f.Points) == 0 {
		return 0
	}
	total := 0.0
	for _, point := range f.Points {
		total += point.Energy
	}
	return total / float64(len(f.Points))
}

func wrap(value, limit float64) float64 {
	if limit <= 0 {
		return 0
	}
	for value < 0 {
		value += limit
	}
	for value >= limit {
		value -= limit
	}
	return value
}

func clamp(value, low, high float64) float64 {
	return math.Max(low, math.Min(high, value))
}

type Snapshot struct {
	Frame      int     `json:"frame"`
	Form       string  `json:"form"`
	Count      int     `json:"count"`
	Brightness float64 `json:"brightness"`
}

func (s Snapshot) Label() string {
	return fmt.Sprintf("%s:%d", s.Form, s.Count)
}

func FormForAction(action string) string {
	switch action {
	case "heart":
		return "heart"
	case "restore":
		return "drift"
	case string(model.ModeWelcome), string(model.ModeTour):
		return "drift"
	default:
		return "drift"
	}
}
