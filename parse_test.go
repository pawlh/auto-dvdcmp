package main

import (
	"reflect"
	"testing"
	"time"
)

func TestParseRawExtras(t *testing.T) {
	tests := []struct {
		name  string
		input []string
		want  []Extra
	}{
		{
			name: "single extras with just hh:mm:ss duration",
			input: []string{
				"The Eagle: The Making of a Roman Epic (12:12)",
			},
			want: []Extra{
				{
					Title:    "The Eagle- The Making of a Roman Epic",
					Duration: 12*time.Minute + 12*time.Second,
				},
			},
		},
		{
			name: "multiple extras with just hh:mm duration",
			input: []string{
				"The Eagle: The Making of a Roman Epic (12:12)",
				"Deleted Scenes (1:23)",
			},
			want: []Extra{
				{
					Title:    "The Eagle- The Making of a Roman Epic",
					Duration: 12*time.Minute + 12*time.Second,
				},
				{
					Title:    "Deleted Scenes",
					Duration: 1*time.Minute + 23*time.Second,
				},
			},
		},
		{
			name: "mixed extra and non-extra lines",
			input: []string{
				"Featurettes",
				"The Eagle: The Making of a Roman Epic (12:12)",
			},
			want: []Extra{
				{
					Title:    "The Eagle- The Making of a Roman Epic",
					Duration: 12*time.Minute + 12*time.Second,
				},
			},
		},
		{
			name: "only non-extra lines",
			input: []string{
				"Featurettes",
				"Deleted Scenes",
			},
			want: []Extra{},
		},
		{
			name:  "empty input",
			input: []string{},
			want:  []Extra{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseRawExtras(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseRawExtras() = %v, want %v", got, tt.want)
			}
		})
	}
}
