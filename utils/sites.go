package utils

// Generate random sites in given bounding box
func RandomSites(bbox voronoi.BBox, count int) []voronoi.Vertex {
	sites := make([]voronoi.Vertex, count)
	w := bbox.Xr - bbox.Xl
	h := bbox.Yb - bbox.Yt
	for j := 0; j < count; j++ {
		sites[j].X = rand.Float64() * w + bbox.Xl
		sites[j].Y = rand.Float64() * h + bbox.Yt
	}
	return sites
}