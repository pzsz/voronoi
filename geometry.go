// Copyright 2013 Przemyslaw Szczepaniak.
// MIT License: See https://github.com/gorhill/Javascript-Voronoi/LICENSE.md

// Author: Przemyslaw Szczepaniak (przeszczep@gmail.com)
// Port of Raymond Hill's (rhill@raymondhill.net) javascript implementation 
// of Steven  Forune's algorithm to compute Voronoi diagrams

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
	lSite Vertex
	rSite Vertex
	va    Vertex
	vb    Vertex
}

func newEdge(lSite, rSite Vertex) *Edge {
	return &Edge{
		lSite: lSite,
		rSite: rSite,
		va:    NO_VERTEX,
		vb:    NO_VERTEX,
	}
}

// Halfedge (directed edge)
type Halfedge struct {
	site  Vertex
	edge  *Edge
	angle float64
}

// Sort interface for halfedges
type Halfedges []*Halfedge

func (s Halfedges) Len() int      { return len(s) }
func (s Halfedges) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// For sorting by angle
type HalfedgesByAngle struct{ Halfedges }

func (s HalfedgesByAngle) Less(i, j int) bool { return s.Halfedges[i].angle < s.Halfedges[j].angle }

func newHalfedge(edge *Edge, lSite, rSite Vertex) *Halfedge {
	ret := &Halfedge{
		site: lSite,
		edge: edge,
	}

	// 'angle' is a value to be used for properly sorting the
	// halfsegments counterclockwise. By convention, we will
	// use the angle of the line defined by the 'site to the left'
	// to the 'site to the right'.
	// However, border edges have no 'site to the right': thus we
	// use the angle of line perpendicular to the halfsegment (the
	// edge should have both end points defined in such case.)
	if rSite != NO_VERTEX {
		ret.angle = math.Atan2(rSite.Y-lSite.Y, rSite.X-lSite.X)
	} else {
		va := edge.va
		vb := edge.vb
		// rhill 2011-05-31: used to call getStartpoint()/getEndpoint(),
		// but for performance purpose, these are expanded in place here.
		if edge.lSite == lSite {
			ret.angle = math.Atan2(vb.X-va.X, va.Y-vb.Y)
		} else {
			ret.angle = math.Atan2(va.X-vb.X, vb.Y-va.Y)
		}
	}
	return ret
}

func (h *Halfedge) getStartpoint() Vertex {
	if h.edge.lSite == h.site {
		return h.edge.va
	}
	return h.edge.vb

}

func (h *Halfedge) getEndpoint() Vertex {
	if h.edge.lSite == h.site {
		return h.edge.vb
	}
	return h.edge.va
}
