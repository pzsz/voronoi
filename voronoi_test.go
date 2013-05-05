package voronoi

import (
	"fmt"
//	. "github.com/pzsz/voronoi"
	"testing"
)

func TestVoronoi(t *testing.T) {
	v := NewVoronoi()
	sites := []Vertex{ NewVertex(4, 5),
		NewVertex(6, 5),
	}

	diagram := v.Compute(sites, NewBBox(0,10, 0, 10))

	for _, e := range diagram.Edges {
		fmt.Printf("final: %v->%v\n", e.va, e.vb)
	}
}