package main

import "testing"

func TestEntryIsHighlight(t *testing.T) {
	goodTitles := []string{
		"Tour de france stage 3 highlights",
		"tour de france highlights stage 5",
		"Tour De France 2015 Daily Highlights Stage 2",
		"Tour de France 2015 Daily Update Stage 3",
	}

	for _, title := range goodTitles {
		entry := Entry{Title: title}
		if !entry.IsHighlight() {
			t.Fatalf("didn't match good title: %v", title)
		}
	}

	badTitles := []string{
		"Tour De France 2014 Extended Highlights Stage 1",
		"Tour De France 2014 Daily Highlights Ep1",
	}
	for _, title := range badTitles {
		entry := Entry{Title: title}
		if entry.IsHighlight() {
			t.Fatalf("matched bad title: %v", entry.Title)
		}
	}
}
