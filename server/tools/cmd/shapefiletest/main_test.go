package main

import (
	"os"
	"strconv"
	"testing"

	"github.com/jonas-p/go-shp"
)

func TestShapefileMain(t *testing.T) {
	os.Remove("points.shp")
	os.Remove("points.shx")
	os.Remove("points.dbf")

	main()

	shape, err := shp.Open("points.shp")
	if err != nil {
		t.Fatalf("failed to open generated shapefile: %v", err)
	}
	defer shape.Close()

	if shape.AttributeCount() != 4 {
		t.Errorf("expected 4 points, got %d", shape.AttributeCount())
	}

	for i := 0; i < shape.AttributeCount(); i++ {
		val := shape.ReadAttribute(i, 0)
		expected := "Point " + strconv.Itoa(i+1)
		if val != expected {
			t.Errorf("attribute for point %d: got '%s', want '%s'", i+1, val, expected)
		}
	}

	os.Remove("points.shp")
	os.Remove("points.shx")
	os.Remove("points.dbf")
}
