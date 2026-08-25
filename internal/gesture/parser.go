package gesture

import (
	"fmt"
	"strings"
)

type Signal struct {
	Kind     string
	Strength int
	Frame    int
}

func Parse(raw string, frame int) (Signal, error) {
	parts := strings.Split(strings.TrimSpace(raw), ":")
	if len(parts) != 2 {
		return Signal{}, fmt.Errorf("gesture must use kind:strength")
	}
	kind := strings.ToLower(strings.TrimSpace(parts[0]))
	strength := 0
	for _, char := range parts[1] {
		if char < '0' || char > '9' {
			return Signal{}, fmt.Errorf("gesture strength is not numeric")
		}
		strength = strength*10 + int(char-'0')
	}
	if frame < 1 || strength > 100 {
		return Signal{}, fmt.Errorf("gesture frame or strength is invalid")
	}
	return Signal{Kind: kind, Strength: strength, Frame: frame}, nil
}

func Supported(kind string) bool {
	switch strings.ToLower(kind) {
	case "open_palm", "fist", "wave", "point":
		return true
	default:
		return false
	}
}
