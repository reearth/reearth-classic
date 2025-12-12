package main

import (
	"os"
	"testing"
)

func TestMainRuns(t *testing.T) {
	// Remove any output files if they exist
	_ = os.Remove("points.shp")
	_ = os.Remove("points.shx")
	_ = os.Remove("points.dbf")

	// Run main and check it does not panic or error
	main()

	// Check that the output file was created
	if _, err := os.Stat("points.shp"); os.IsNotExist(err) {
		t.Errorf("points.shp was not created")
	} else if err != nil {
		t.Errorf("error checking points.shp: %v", err)
	}
}
