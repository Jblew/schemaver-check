package main

import "testing"

func TestDummy(t *testing.T) {
	got := "ok"

	want := "ok"
	if got != want {
		t.Errorf("got = %s, want \"%s\"", got, want)
	}
}
