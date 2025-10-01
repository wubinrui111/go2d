package components

import (
	"testing"
)

func TestBoxIntersects(t *testing.T) {
	// Create two boxes that intersect
	box1 := &Box{X: 0, Y: 0, Width: 10, Height: 10}
	box2 := &Box{X: 5, Y: 5, Width: 10, Height: 10}
	
	// They should intersect
	if !box1.Intersects(box2) {
		t.Error("Expected box1 and box2 to intersect")
	}
	
	// Create two boxes that don't intersect
	box3 := &Box{X: 0, Y: 0, Width: 10, Height: 10}
	box4 := &Box{X: 15, Y: 15, Width: 10, Height: 10}
	
	// They should not intersect
	if box3.Intersects(box4) {
		t.Error("Expected box3 and box4 not to intersect")
	}
	
	// Create two boxes that touch at the edge (should NOT intersect based on our implementation)
	box5 := &Box{X: 0, Y: 0, Width: 10, Height: 10}
	box6 := &Box{X: 10, Y: 0, Width: 10, Height: 10}
	
	// They should not intersect (touching edges means no overlap)
	if box5.Intersects(box6) {
		t.Error("Expected box5 and box6 not to intersect (touching edges)")
	}
}

func TestGetIntersectionDepth(t *testing.T) {
	// Create two overlapping boxes
	box1 := &Box{X: 0, Y: 0, Width: 10, Height: 10}
	box2 := &Box{X: 5, Y: 5, Width: 10, Height: 10}
	
	xDepth, yDepth := box1.GetIntersectionDepth(box2)
	
	// Check that there is intersection on both axes
	if xDepth == 0 {
		t.Error("Expected non-zero x depth")
	}
	
	if yDepth == 0 {
		t.Error("Expected non-zero y depth")
	}
	
	// The intersection should be negative on both axes (box1 needs to move left and up)
	if xDepth > 0 {
		t.Error("Expected negative x depth (box1 should move left)")
	}
	
	if yDepth > 0 {
		t.Error("Expected negative y depth (box1 should move up)")
	}
	
	// Test with non-intersecting boxes
	box3 := &Box{X: 0, Y: 0, Width: 10, Height: 10}
	box4 := &Box{X: 15, Y: 15, Width: 10, Height: 10}
	
	xDepth, yDepth = box3.GetIntersectionDepth(box4)
	
	// Should have zero intersection
	if xDepth != 0 {
		t.Errorf("Expected zero x depth for non-intersecting boxes, got %f", xDepth)
	}
	
	if yDepth != 0 {
		t.Errorf("Expected zero y depth for non-intersecting boxes, got %f", yDepth)
	}
}