package display

import (
	"encoding/json"
	"fmt"

	"showroom/internal/model"
	"showroom/internal/particles"
)

type Frame struct {
	Title     string             `json:"title"`
	Subtitle  string             `json:"subtitle"`
	Mode      string             `json:"mode"`
	State     model.DisplayState `json:"state"`
	Particles particles.Snapshot `json:"particles"`
}

func BuildFrame(state model.DisplayState, snapshot particles.Snapshot, gallery string) Frame {
	title := state.PhraseText
	if title == "" {
		title = gallery
	}
	return Frame{Title: title, Subtitle: model.RevisionText(state.Revision), Mode: model.ModeLabel(state.Mode), State: state, Particles: snapshot}
}

func EncodeFrame(frame Frame) ([]byte, error) {
	return json.Marshal(frame)
}

func PlainSummary(frame Frame) string {
	return fmt.Sprintf("%s | %s | %s", frame.Mode, frame.Title, frame.Particles.Label())
}
