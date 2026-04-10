package main

import "testing"

func TestSuccess(t *testing.T) {
	expected := "success"
	actual := "success"

	if actual != expected {
		t.Errorf("Test failed: expected %s, but got %s", expected, actual)
	}

	t.Log("Test passed")
}
