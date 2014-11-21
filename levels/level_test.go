package levels

import "testing"

func TestSomething(t *testing.T) {
	expecteds := []ErrorLevelString{
		UNKNOWN,
		UNDEFINED,
		OFF,
		FATAL,
		ERROR,
		WARN,
		INFO,
		DEBUG,
		TRACE,
		UNKNOWN,
		UNKNOWN,
	}

	inputs := []ErrorLevel{
		-1, 0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
	}

	if len(expecteds) != len(inputs) {
		t.Fatal("94280111108 len(expecteds)!= len(inputs)", len(expecteds), len(inputs))
	}

	for i, input := range inputs {
		expected := expecteds[i]
		s := input.String()
		candidate := ErrorLevelString(s)

		if candidate != expected {
			t.Error(i, "candidate != expected:", candidate, expected)
		}
	}
}
