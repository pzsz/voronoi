// Copyright 2013 Przemyslaw Szczepaniak.
// MIT License: See https://github.com/gorhill/Javascript-Voronoi/LICENSE.md

// Author: Przemyslaw Szczepaniak (przeszczep@gmail.com)
// Port of Raymond Hill's (rhill@raymondhill.net) javascript implementation 
// of Steven  Forune's algorithm to compute Voronoi diagrams

package voronoi

import (
	"fmt"
	//	. "github.com/pzsz/voronoi"
	"testing"
)

func TestVoronoi(t *testing.T) {
	v := NewVoronoi()
	sites := []Vertex{Vertex{4, 5},
		Vertex{6, 5},
	}

	diagram := v.Compute(sites, NewBBox(0, 10, 0, 10))

	for _, e := range diagram.Edges {
		fmt.Printf("final: %v->%v\n", e.va, e.vb)
	}
}
