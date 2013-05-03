package voronoi

type Edge struct {
	lSite *Vertex
	rSite *Vertex
	va *Vertex
	vb *Vertex
}

func newEdge(lSite, rSite *Vertex) *Edge {
	return &Edge{lSite: lSite,
		rSite: rSite,
	}
}
