package docker

import (
	"encoding/json"
	"testing"

	orderedmap "github.com/wk8/go-ordered-map/v2"
)

func TestGuessLogLevel(t *testing.T) {
	var nilOrderedMap *orderedmap.OrderedMap[string, any]
	tests := []struct {
		input    any
		expected string
	}{
		{"ERROR: Something went wrong", "error"},
		{"WARN: Something might be wrong", "warn"},
		{"INFO: Something happened", "info"},
		{"debug: Something happened", "debug"},
		{"debug Something happened", "debug"},
		{"TRACE: Something happened", "trace"},
		{"FATAL: Something happened", "fatal"},
		{"[ERROR] Something went wrong", "error"},
		{"[error] Something went wrong", "error"},
		{"[ ERROR ] Something went wrong", "error"},
		{"[error] Something went wrong", "error"},
		{"[test] [error] Something went wrong", "error"},
		{"Some test with error=test", "error"},
		{"[foo] [ ERROR] Something went wrong", "error"},
		{"123 ERROR Something went wrong", "error"},
		{"123 Something went wrong", "unknown"},
		{orderedmap.New[string, string](
			orderedmap.WithInitialData(
				orderedmap.Pair[string, string]{Key: "key", Value: "value"},
				orderedmap.Pair[string, string]{Key: "level", Value: "info"},
			),
		), "info"},
		{orderedmap.New[string, any](
			orderedmap.WithInitialData(
				orderedmap.Pair[string, any]{Key: "key", Value: "value"},
				orderedmap.Pair[string, any]{Key: "level", Value: "info"},
			),
		), "info"},
		{orderedmap.New[string, string](
			orderedmap.WithInitialData(
				orderedmap.Pair[string, string]{Key: "key", Value: "value"},
				orderedmap.Pair[string, string]{Key: "severity", Value: "info"},
			),
		), "info"},
		{orderedmap.New[string, any](
			orderedmap.WithInitialData(
				orderedmap.Pair[string, any]{Key: "key", Value: "value"},
				orderedmap.Pair[string, any]{Key: "severity", Value: "info"},
			),
		), "info"},
		{nilOrderedMap, "unknown"},
		{nil, "unknown"},
	}

	for _, test := range tests {
		name, _ := json.Marshal(test.input)
		t.Run(string(name), func(t *testing.T) {
			actual := guessLogLevel(&LogEvent{Message: test.input})
			if actual != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, actual)
			}
		})
	}
}
