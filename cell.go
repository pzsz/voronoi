package voronoi

import "sort"

type Cell struct {
	site Vertex
	halfedges []*Halfedge
}

func (t *Cell) prepare() int {
	halfedges := t.halfedges
        iHalfedge := len(halfedges)

	// get rid of unused halfedges
	// rhill 2011-05-27: Keep it simple, no point here in trying
	// to be fancy: dangling edges are a typically a minority.
	for ;iHalfedge > 0;iHalfedge-- {
		edge := halfedges[iHalfedge].edge;
		if (edge.vb == nil || edge.va == nil) {
			halfedges[iHalfedge] = halfedges[len(halfedges)-1]
			halfedges = halfedges[:len(halfedges)-1]
		}
        }

	// rhill 2011-05-26: I tried to use a binary search at insertion
	// time to keep the array sorted on-the-fly (in Cell.addHalfedge()).
	// There was no real benefits in doing so, performance on
	// Firefox 3.6 was improved marginally, while performance on
	// Opera 11 was penalized marginally.
	sort.Sort(HalfedgesByAngle{halfedges})
	return len(halfedges)
}
