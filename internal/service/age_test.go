package service

import (
	"testing"
	"time"
)

func TestCalculateAge(t *testing.T) {
	tests := []struct {
		name string
		dob  time.Time
		now  time.Time
		want int
	}{
		{
			name: "birthday already passed this year",
			dob:  time.Date(1990, time.May, 10, 0, 0, 0, 0, time.UTC),
			now:  time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC),
			want: 36,
		},
		{
			name: "birthday not yet reached this year",
			dob:  time.Date(1990, time.December, 20, 0, 0, 0, 0, time.UTC),
			now:  time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC),
			want: 35,
		},
		{
			name: "birthday today",
			dob:  time.Date(2000, time.June, 15, 0, 0, 0, 0, time.UTC),
			now:  time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC),
			want: 26,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateAge(tt.dob, tt.now)
			if got != tt.want {
				t.Fatalf("CalculateAge() = %d, want %d", got, tt.want)
			}
		})
	}
}
