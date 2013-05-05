package voronoi

import "math"

type Halfedge struct {
	site Vertex
	edge *Edge
	angle float64
}

type Halfedges []*Halfedge
func (s Halfedges) Len() int      { return len(s) }
func (s Halfedges) Swap(i, j int) { s[i], s[j] = s[j], s[i] }


type HalfedgesByAngle struct { Halfedges }
func (s HalfedgesByAngle) Less(i, j int) bool { return s.Halfedges[i].angle < s.Halfedges[j].angle }

func newHalfedge(edge *Edge, lSite, rSite *Vertex) *Halfedge {
	ret := &Halfedge{
		site: *lSite,
		edge: edge,
	}

	// 'angle' is a value to be used for properly sorting the
	// halfsegments counterclockwise. By convention, we will
	// use the angle of the line defined by the 'site to the left'
	// to the 'site to the right'.
	// However, border edges have no 'site to the right': thus we
	// use the angle of line perpendicular to the halfsegment (the
	// edge should have both end points defined in such case.)
	if (rSite != nil) {
		ret.angle = math.Atan2(rSite.y-lSite.y, rSite.x-lSite.x);
        } else {
		va := edge.va
		vb := edge.vb
		// rhill 2011-05-31: used to call getStartpoint()/getEndpoint(),
		// but for performance purpose, these are expanded in place here.
		if (*edge.lSite == *lSite) {
			ret.angle = math.Atan2(vb.x-va.x, va.y-vb.y)
		} else {
			ret.angle = math.Atan2(va.x-vb.x, vb.y-va.y)
		}
        }
	return ret
}

func (h *Halfedge) getStartpoint() Vertex {
	if h.edge.lSite != nil && *h.edge.lSite == h.site {
		return *h.edge.va
	}
	return *h.edge.vb
	
}

func (h *Halfedge) getEndpoint() Vertex {
	if h.edge.lSite != nil && *h.edge.lSite == h.site {
		return *h.edge.vb
	}
	return *h.edge.va
}
