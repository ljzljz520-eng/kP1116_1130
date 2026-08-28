package rehearsal

import "testing"

func TestDefaultScriptRun(t *testing.T) {
	result, err := Run(DefaultScript())
	if err != nil {
		t.Fatal(err)
	}
	if !result.Passed || result.SuccessfulSteps() != 4 || len(result.Forms) != 4 {
		t.Fatalf("unexpected rehearsal: %#v", result)
	}
	if result.FinalMode != "idle" || result.Summary() == "" {
		t.Fatalf("unexpected final result: %#v", result)
	}
}
