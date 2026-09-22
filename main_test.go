package main

import "testing"

func TestMaxValueFrom(t *testing.T) {
	tests := []struct {
		difficulty string
		want       uint32
	}{
		{"easy", 10},
		{"hard", 1000},
		{"normal", 100},
		{"weird", 100},
		{"", 100},
	}
	for _, tt := range tests {
		t.Run(tt.difficulty, func(t *testing.T) {
			if got := maxValueFrom(tt.difficulty); got != tt.want {
				t.Errorf("maxValueFrom(%q) = %d, want %d", tt.difficulty, got, tt.want)
			}
		})
	}
}
