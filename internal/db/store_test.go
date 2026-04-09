package db

import (
	"os"
	"testing"

	"rs-item-database/pb"
)

func newTestStore(t *testing.T) (*Store, func()) {
	t.Helper()
	dir, err := os.MkdirTemp("", "nutsdb-test-*")
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewStore(dir)
	if err != nil {
		os.RemoveAll(dir)
		t.Fatal(err)
	}
	return s, func() {
		s.Close()
		os.RemoveAll(dir)
	}
}

func TestSearchItems_SubstringMatch(t *testing.T) {
	s, cleanup := newTestStore(t)
	defer cleanup()

	items := []*pb.Item{
		{Id: 4151, Name: "Abyssal whip"},
		{Id: 4587, Name: "Dragon scimitar"},
		{Id: 1215, Name: "Dragon dagger"},
		{Id: 4718, Name: "Abyssal bludgeon"},
	}
	for _, item := range items {
		if err := s.SaveItem(item); err != nil {
			t.Fatalf("SaveItem(%s): %v", item.Name, err)
		}
	}

	tests := []struct {
		query    string
		wantAny  []string // at least one of these must appear in results
		wantNone []string // none of these should appear
	}{
		{
			query:   "whip",
			wantAny: []string{"Abyssal whip"},
			// prefix search would return nothing here
		},
		{
			query:   "dragon",
			wantAny: []string{"Dragon scimitar", "Dragon dagger"},
		},
		{
			query:   "abyssal",
			wantAny: []string{"Abyssal whip", "Abyssal bludgeon"},
		},
		{
			query:    "abyssal whip",
			wantAny:  []string{"Abyssal whip"},
			wantNone: []string{"Abyssal bludgeon"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.query, func(t *testing.T) {
			results, err := s.SearchItems(tc.query, 50)
			if err != nil {
				t.Fatalf("SearchItems(%q): %v", tc.query, err)
			}

			names := make(map[string]bool)
			for _, r := range results {
				names[r.Name] = true
			}

			for _, want := range tc.wantAny {
				if !names[want] {
					t.Errorf("query %q: expected %q in results, got %v", tc.query, want, results)
				}
			}
			for _, notWant := range tc.wantNone {
				if names[notWant] {
					t.Errorf("query %q: unexpected %q in results", tc.query, notWant)
				}
			}
		})
	}
}

func TestSearchItems_Empty(t *testing.T) {
	s, cleanup := newTestStore(t)
	defer cleanup()

	results, err := s.SearchItems("anything", 10)
	if err != nil {
		t.Fatalf("unexpected error on empty store: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}
}
