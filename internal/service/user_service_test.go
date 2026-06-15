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
