package service

import (
	"testing"
	"time"

	"user-api/db/sqlc"
)

func TestToUserResponseWithAge(t *testing.T) {
	user := sqlc.User{
		ID:   1,
		Name: "Alice",
		Dob:  time.Date(1990, time.May, 10, 0, 0, 0, 0, time.UTC),
	}
	now := time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC)

	got := toUserResponse(user, true, now)

	if got.Age != 36 {
		t.Fatalf("Age = %d, want %d", got.Age, 36)
	}
	if got.ID != 1 || got.Name != "Alice" || got.Dob != "1990-05-10" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestToUserResponseWithoutAge(t *testing.T) {
	user := sqlc.User{
		ID:   2,
		Name: "Bob",
		Dob:  time.Date(1995, time.July, 1, 0, 0, 0, 0, time.UTC),
	}
	now := time.Date(2026, time.June, 15, 0, 0, 0, 0, time.UTC)

	got := toUserResponse(user, false, now)

	if got.Age != 0 {
		t.Fatalf("Age = %d, want %d", got.Age, 0)
	}
	if got.ID != 2 || got.Name != "Bob" || got.Dob != "1995-07-01" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestNormalizePagination(t *testing.T) {
	tests := []struct {
		name      string
		page      int
		limit     int
		wantPage  int
		wantLimit int
	}{
		{name: "defaults", page: 0, limit: 0, wantPage: 1, wantLimit: 10},
		{name: "clamps limit", page: 2, limit: 500, wantPage: 2, wantLimit: 100},
		{name: "keeps valid", page: 3, limit: 25, wantPage: 3, wantLimit: 25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			page, limit := normalizePagination(tt.page, tt.limit)
			if page != tt.wantPage || limit != tt.wantLimit {
				t.Fatalf("normalizePagination() = (%d, %d), want (%d, %d)", page, limit, tt.wantPage, tt.wantLimit)
			}
		})
	}
}
