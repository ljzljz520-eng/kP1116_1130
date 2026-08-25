package model

import "fmt"

func ModeLabel(mode SceneMode) string {
	switch mode {
	case ModeWelcome:
		return "Welcome"
	case ModeTour:
		return "Guided tour"
	case ModeClosing:
		return "Closing"
	case ModeIdle:
		return "Idle"
	default:
		return "Unknown"
	}
}

func GestureLabel(kind string) string {
	switch kind {
	case "open_palm":
		return "Open palm"
	case "fist":
		return "Fist"
	case "wave":
		return "Wave"
	case "point":
		return "Point"
	default:
		return "Unrecognized"
	}
}

func RevisionText(revision int) string {
	if revision < 1 {
		return "unpublished"
	}
	return fmt.Sprintf("revision-%d", revision)
}
