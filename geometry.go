// MIT License: See https://github.com/pzsz/voronoi/LICENSE.md

// Author: Przemyslaw Szczepaniak (przeszczep@gmail.com)
// Port of Raymond Hill's (rhill@raymondhill.net) javascript implementation 
// of Steven Forune's algorithm to compute Voronoi diagrams

package voronoi

import (
	"math"
)

// Vertex on 2D plane
type Vertex struct {
	X float64
	Y float64
}

// Vertex representing lack of vertex (or bad vertex)
var NO_VERTEX = Vertex{math.Inf(1), math.Inf(1)}

// For sort interface
type Vertices []Vertex

func (s Vertices) Len() int      { return len(s) }
func (s Vertices) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Used for sorting vertices along the Y axis
type VerticesByY struct{ Vertices }

func (s VerticesByY) Less(i, j int) bool { return s.Vertices[i].Y < s.Vertices[j].Y }

// Edge structure
type Edge struct {
	LeftSite  Vertex
	RightSite Vertex
	Va        Vertex
	Vb        Vertex
}

func newEdge(LeftSite, RightSite Vertex) *Edge {
	return &Edge{
		LeftSite:  LeftSite,
		RightSite: RightSite,
		Va:        NO_VERTEX,
		Vb:        NO_VERTEX,
	}
}

// Halfedge (directed edge)
type Halfedge struct {
	Site  Vertex
	Edge  *Edge
	Angle float64
}

// Sort interface for halfedges
type Halfedges []*Halfedge

func (s Halfedges) Len() int      { return len(s) }
func (s Halfedges) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// For sorting by angle
type HalfedgesByAngle struct{ Halfedges }

func (s HalfedgesByAngle) Less(i, j int) bool { return s.Halfedges[i].Angle > s.Halfedges[j].Angle }

func newHalfedge(edge *Edge, LeftSite, RightSite Vertex) *Halfedge {
	ret := &Halfedge{
		Site: LeftSite,
		Edge: edge,
	}

	// 'angle' is a value to be used for properly sorting the
	// halfsegments counterclockwise. By convention, we will
	// use the angle of the line defined by the 'site to the left'
	// to the 'site to the right'.
	// However, border edges have no 'site to the right': thus we
	// use the angle of line perpendicular to the halfsegment (the
	// edge should have both end points defined in such case.)
	if RightSite != NO_VERTEX {
		ret.Angle = math.Atan2(RightSite.Y-LeftSite.Y, RightSite.X-LeftSite.X)
	} else {
		va := edge.Va
		vb := edge.Vb
		// rhill 2011-05-31: used to call GetStartpoint()/GetEndpoint(),
		// but for performance purpose, these are expanded in place here.
		if edge.LeftSite == LeftSite {
			ret.Angle = math.Atan2(vb.X-va.X, va.Y-vb.Y)
		} else {
			ret.Angle = math.Atan2(va.X-vb.X, vb.Y-va.Y)
		}
	}
	return ret
}

func (h *Halfedge) GetStartpoint() Vertex {
	if h.Edge.LeftSite == h.Site {
		return h.Edge.Va
	}
	return h.Edge.Vb

}

func (h *Halfedge) GetEndpoint() Vertex {
	if h.Edge.LeftSite == h.Site {
		return h.Edge.Vb
	}
	return h.Edge.Va
}
