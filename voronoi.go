package voronoi

import "math"

type Voronoi struct {
	cells []*Cell
	edges []*Edge

	beachline rbTree
	circleEvents rbTree
	firstCircleEvent *CircleEvent
}

func (s *Voronoi) getCell(site Vertex) *Cell {
	for _, cell := range s.cells {
		if cell.site == site {
			return cell
		}
	}
	return nil
}

func (s *Voronoi) createEdge(lSite, rSite, va, vb *Vertex) *Edge {
	edge := newEdge(lSite, rSite)

	if (va != nil) {
		s.setEdgeStartpoint(edge, lSite, rSite, va)
	}

	if (vb != nil) {
		s.setEdgeEndpoint(edge, lSite, rSite, vb)
	}

	lCell := s.getCell(*lSite)
	rCell := s.getCell(*rSite)
	lCell.halfedges = append(lCell.halfedges, newHalfedge(edge, lSite, rSite))
	rCell.halfedges = append(lCell.halfedges, newHalfedge(edge, rSite, lSite))
	return edge
}

func (s *Voronoi) createBorderEdge(lSite, va, vb *Vertex) *Edge {
	edge := newEdge(lSite, nil)
	edge.va = va
	edge.vb = vb

	s.edges = append(s.edges, edge)
	return edge
}

func (s *Voronoi) setEdgeStartpoint(edge *Edge, lSite, rSite, vertex *Vertex) {
	if (edge.va == nil && edge.vb == nil) {
		edge.va = vertex
		edge.lSite = lSite
		edge.rSite = rSite
        } else if (*edge.lSite == *rSite) {
		edge.vb = vertex
        } else {
		edge.va = vertex
        }
}

func (s *Voronoi) setEdgeEndpoint(edge *Edge, lSite, rSite, vertex *Vertex) {
	s.setEdgeStartpoint(edge, rSite, lSite, vertex)
}

type Beachsection struct {
	node *rbNode
	site Vertex
	circleEvent *CircleEvent
	edge *Edge
}

func (s *Beachsection) bindToNode(node *rbNode) {
	s.node = node
}

func (s *Beachsection) getNode() *rbNode {
	return s.node
}

// calculate the left break point of a particular beach section,
// given a particular sweep line
func (s *Voronoi) leftBreakPoint(arc *Beachsection, directrix float64) float64 {
	site := arc.site
	rfocx := site.x
        rfocy := site.y
        pby2 := rfocy-directrix
	// parabola in degenerate case where focus is on directrix
	if (pby2 == 0) {
		return rfocx
        }

	lArc := arc.getNode().previous;
	if (lArc == nil) {
		return math.Inf(-1);
        }
	site = lArc.value.(*Beachsection).site
	lfocx := site.x
        lfocy := site.y
        plby2 := lfocy-directrix
	// parabola in degenerate case where focus is on directrix
	if (plby2 == 0) {
		return lfocx
        }
	hl := lfocx-rfocx
        aby2 := 1/pby2-1/plby2
        b := hl/plby2
	if (aby2 != 0) {
		return (-b+math.Sqrt(b*b-2*aby2*(hl*hl/(-2*plby2)-lfocy+plby2/2+rfocy-pby2/2)))/aby2+rfocx
        }
	// both parabolas have same distance to directrix, thus break point is midway
	return (rfocx+lfocx)/2
}

// calculate the right break point of a particular beach section,
// given a particular directrix
func (s *Voronoi) rightBreakPoint(arc *Beachsection, directrix float64) float64 {
	rArc := arc.getNode().next
	if (rArc != nil) {
		return s.leftBreakPoint(rArc.value.(*Beachsection), directrix)
        }
	site := arc.site
	if site.y == directrix {
		return site.x
	}
	return math.Inf(1)
}

func (s *Voronoi) detachBeachsection(arc *Beachsection) {
	s.detachCircleEvent(arc)
	s.beachline.removeNode(arc.node)
}

type BeachsectionPtrs []*Beachsection

func (s *BeachsectionPtrs) appendLeft(b *Beachsection) {
	*s = append(*s, b) 
	for id := len(*s)-1; id > 0; id-- {
		(*s)[id] = (*s)[id-1] 
	}
	(*s)[0] = b
}

func (s *BeachsectionPtrs) appendRight(b *Beachsection) {
	*s = append(*s, b) 
}

func (s *Voronoi) removeBeachsection(beachsection *Beachsection) {
	circle := beachsection.circleEvent
	x := circle.x
        y := circle.ycenter
        vertex := Vertex{x, y}
        previous := beachsection.node.previous
        next := beachsection.node.next
        disappearingTransitions := BeachsectionPtrs{beachsection}
	abs_fn := math.Abs

	// remove collapsed beachsection from beachline
	s.detachBeachsection(beachsection)

	// there could be more than one empty arc at the deletion point, this
	// happens when more than two edges are linked by the same vertex,
	// so we will collect all those edges by looking up both sides of
	// the deletion point.
	// by the way, there is *always* a predecessor/successor to any collapsed
	// beach section, it's just impossible to have a collapsing first/last
	// beach sections on the beachline, since they obviously are unconstrained
	// on their left/right side.

	// look left
	lArc := previous.value.(*Beachsection)
	for (lArc.circleEvent != nil && 
		abs_fn(x-lArc.circleEvent.x)<1e-9 && 
		abs_fn(y-lArc.circleEvent.ycenter)<1e-9) {

		previous = lArc.node.previous
		disappearingTransitions.appendLeft(lArc)
		s.detachBeachsection(lArc); // mark for reuse
		lArc = previous.value.(*Beachsection)
        }
	// even though it is not disappearing, I will also add the beach section
	// immediately to the left of the left-most collapsed beach section, for
	// convenience, since we need to refer to it later as this beach section
	// is the 'left' site of an edge for which a start point is set.
	disappearingTransitions.appendLeft(lArc)
	s.detachCircleEvent(lArc)

	// look right
	var rArc = next.value.(*Beachsection)
	for (rArc.circleEvent !=  nil && 
		abs_fn(x-rArc.circleEvent.x)<1e-9 && 
		abs_fn(y-rArc.circleEvent.ycenter)<1e-9) {
		next = rArc.node.next
		disappearingTransitions.appendRight(rArc)
		s.detachBeachsection(rArc); // mark for reuse
		rArc = next.value.(*Beachsection)
        }
	// we also have to add the beach section immediately to the right of the
	// right-most collapsed beach section, since there is also a disappearing
	// transition representing an edge's start point on its left.
	disappearingTransitions.appendRight(rArc)
	s.detachCircleEvent(rArc)

	// walk through all the disappearing transitions between beach sections and
	// set the start point of their (implied) edge.
	nArcs := len(disappearingTransitions)
        
	for iArc:=1; iArc<nArcs; iArc++ {
		rArc = disappearingTransitions[iArc]
		lArc = disappearingTransitions[iArc-1];
		s.setEdgeStartpoint(rArc.edge, &lArc.site, &rArc.site, &vertex)
        }

	// create a new edge as we have now a new transition between
	// two beach sections which were previously not adjacent.
	// since this edge appears as a new vertex is defined, the vertex
	// actually define an end point of the edge (relative to the site
	// on the left)
	lArc = disappearingTransitions[0]
	rArc = disappearingTransitions[nArcs-1]
	rArc.edge = s.createEdge(&lArc.site, &rArc.site, nil, &vertex)

	// create circle events if any for beach sections left in the beachline
	// adjacent to collapsed sections
	s.attachCircleEvent(lArc)
	s.attachCircleEvent(rArc)
}

func (s *Voronoi) addBeachsection(site Vertex) {
	x := site.x
        directrix := site.y

	// find the left and right beach sections which will surround the newly
	// created beach section.
	// rhill 2011-06-01: This loop is one of the most often executed,
	// hence we expand in-place the comparison-against-epsilon calls.
	var lNode, rNode *rbNode
        var dxl, dxr float64
        node := s.beachline.root

	for node != nil {
		nodeBeachline := node.value.(*Beachsection)
		dxl = s.leftBreakPoint(nodeBeachline, directrix)-x;
		// x lessThanWithEpsilon xl => falls somewhere before the left edge of the beachsection
		if (dxl > 1e-9) {
			// this case should never happen
			// if (!node.rbLeft) {
			//    rNode = node.rbLeft;
			//    break;
			//    }
			node = node.left
		} else {
			dxr = x-s.rightBreakPoint(nodeBeachline, directrix)
			// x greaterThanWithEpsilon xr => falls somewhere after the right edge of the beachsection
			if (dxr > 1e-9) {
				if (node.right == nil) {
					lNode = node
					break
				}
				node = node.right
			} else {
				// x equalWithEpsilon xl => falls exactly on the left edge of the beachsection
				if (dxl > -1e-9) {
					lNode = node.previous
					rNode = node
				} else if (dxr > -1e-9) {
				// x equalWithEpsilon xr => falls exactly on the right edge of the beachsection
					lNode = node
					rNode = node.next
					// falls exactly somewhere in the middle of the beachsection
				} else {
					lNode = node
					rNode = node
				}
				break
			}
		}
        }

	var lArc, rArc *Beachsection

	if lNode != nil {
		lArc = lNode.value.(*Beachsection)
	}
	if rNode != nil {
		rArc = rNode.value.(*Beachsection)
	}

	// at this point, keep in mind that lArc and/or rArc could be
	// undefined or null.

	// create a new beach section object for the site and add it to RB-tree
	newArc := &Beachsection{site : site}
	s.beachline.insertSuccessor(lArc, newArc)

	// cases:
	//

	// [null,null]
	// least likely case: new beach section is the first beach section on the
	// beachline.
	// This case means:
	//   no new transition appears
	//   no collapsing beach section
	//   new beachsection become root of the RB-tree
	if (lArc == nil && rArc == nil) {
		return;
        }

	// [lArc,rArc] where lArc == rArc
	// most likely case: new beach section split an existing beach
	// section.
	// This case means:
	//   one new transition appears
	//   the left and right beach section might be collapsing as a result
	//   two new nodes added to the RB-tree
	if (*lArc == *rArc) {
		// invalidate circle event of split beach section
		s.detachCircleEvent(lArc)

		// split the beach section into two separate beach sections
		rArc = &Beachsection{site : lArc.site}
		s.beachline.insertSuccessor(newArc, rArc)

		// since we have a new transition between two beach sections,
		// a new edge is born
		newArc.edge = s.createEdge(&lArc.site, &newArc.site, nil, nil)
		rArc.edge = newArc.edge

		// check whether the left and right beach sections are collapsing
		// and if so create circle events, to be notified when the point of
		// collapse is reached.
		s.attachCircleEvent(lArc)
		s.attachCircleEvent(rArc)
		return
        }

	// [lArc,null]
	// even less likely case: new beach section is the *last* beach section
	// on the beachline -- this can happen *only* if *all* the previous beach
	// sections currently on the beachline share the same y value as
	// the new beach section.
	// This case means:
	//   one new transition appears
	//   no collapsing beach section as a result
	//   new beach section become right-most node of the RB-tree
	if (lArc != nil && rArc == nil) {
		newArc.edge = s.createEdge(&lArc.site, &newArc.site, nil, nil)
		return;
        }

	// [null,rArc]
	// impossible case: because sites are strictly processed from top to bottom,
	// and left to right, which guarantees that there will always be a beach section
	// on the left -- except of course when there are no beach section at all on
	// the beach line, which case was handled above.
	// rhill 2011-06-02: No point testing in non-debug version
	//if (!lArc && rArc) {
		//    throw "Voronoi.addBeachsection(): What is this I don't even";
	//    }

	// [lArc,rArc] where lArc != rArc
	// somewhat less likely case: new beach section falls *exactly* in between two
	// existing beach sections
	// This case means:
	//   one transition disappears
	//   two new transitions appear
	//   the left and right beach section might be collapsing as a result
	//   only one new node added to the RB-tree
	if (lArc != rArc) {
		// invalidate circle events of left and right sites
		s.detachCircleEvent(lArc)
		s.detachCircleEvent(rArc)

		// an existing transition disappears, meaning a vertex is defined at
		// the disappearance point.
		// since the disappearance is caused by the new beachsection, the
		// vertex is at the center of the circumscribed circle of the left,
		// new and right beachsections.
		// http://mathforum.org/library/drmath/view/55002.html
		// Except that I bring the origin at A to simplify
		// calculation
		lSite := lArc.site
		ax := lSite.x
		ay := lSite.y
		bx:=site.x-ax
		by:=site.y-ay
		rSite := rArc.site
		cx:=rSite.x-ax
		cy:=rSite.y-ay
		d:=2*(bx*cy-by*cx)
		hb:=bx*bx+by*by
		hc:=cx*cx+cy*cy
		vertex := Vertex{(cy*hb-by*hc)/d+ax, (bx*hc-cx*hb)/d+ay}

		// one transition disappear
		s.setEdgeStartpoint(rArc.edge, &lSite, &rSite, &vertex)

		// two new transitions appear at the new vertex location
		newArc.edge = s.createEdge(&lSite, &site, nil, &vertex)
		rArc.edge = s.createEdge(&site, &rSite, nil, &vertex)

		// check whether the left and right beach sections are collapsing
		// and if so create circle events, to handle the point of collapse.
		s.attachCircleEvent(lArc)
		s.attachCircleEvent(rArc)
		return
        }
}



type CircleEvent struct {
	node *rbNode
	arc *Beachsection
	x float64
	ycenter float64
}

func (s *CircleEvent) bindToNode(node *rbNode) {
	s.node = node
}

func (s *CircleEvent) getNode() *rbNode {
	return s.node
}

func (s *Voronoi) attachCircleEvent(arc *Beachsection) {

}

func (s *Voronoi) detachCircleEvent(arc *Beachsection) {

}
