package voronoi

type Vertex struct {
	x float64
	y float64
}

type Vertices []Vertex
func (s Vertices) Len() int      { return len(s) }
func (s Vertices) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type VerticesByY struct { Vertices }
func (s VerticesByY) Less(i, j int) bool { return s.Vertices[i].y < s.Vertices[j].y }
