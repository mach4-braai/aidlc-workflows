package graph

import (
	"gonum.org/v1/gonum/graph/simple"

	"github.com/mach4-braai/aidlc-workflows/aidlc-traceability/internal/models"
)

// TraceGraph is a directed graph where each node represents one artifact.
type TraceGraph struct {
	g       *simple.DirectedGraph
	idToInt map[string]int64
	intToID map[int64]string
	nextID  int64
}

func newTraceGraph() *TraceGraph {
	return &TraceGraph{
		g:       simple.NewDirectedGraph(),
		idToInt: make(map[string]int64),
		intToID: make(map[int64]string),
	}
}

func (tg *TraceGraph) nodeFor(artifactID string) int64 {
	if n, ok := tg.idToInt[artifactID]; ok {
		return n
	}
	n := tg.nextID
	tg.nextID++
	tg.idToInt[artifactID] = n
	tg.intToID[n] = artifactID
	tg.g.AddNode(simple.Node(n))
	return n
}

// Build constructs a TraceGraph from a list of artifacts and relationships.
func Build(artifacts []models.Artifact, rels []models.Relationship) *TraceGraph {
	tg := newTraceGraph()
	for _, a := range artifacts {
		tg.nodeFor(a.ID)
	}
	for _, r := range rels {
		src := tg.nodeFor(r.SourceID)
		tgt := tg.nodeFor(r.TargetID)
		tg.g.SetEdge(simple.Edge{F: simple.Node(src), T: simple.Node(tgt)})
	}
	return tg
}

// NodeCount returns the number of nodes in the graph.
func NodeCount(tg *TraceGraph) int {
	return tg.g.Nodes().Len()
}

// EdgeCount returns the number of directed edges in the graph.
func EdgeCount(tg *TraceGraph) int {
	return tg.g.Edges().Len()
}

// HasSuccessor reports whether the artifact with the given ID has at least one outgoing edge.
func HasSuccessor(tg *TraceGraph, artifactID string) bool {
	n, ok := tg.idToInt[artifactID]
	if !ok {
		return false
	}
	return tg.g.From(n).Len() > 0
}

// HasPredecessor reports whether the artifact with the given ID has at least one incoming edge.
func HasPredecessor(tg *TraceGraph, artifactID string) bool {
	n, ok := tg.idToInt[artifactID]
	if !ok {
		return false
	}
	return tg.g.To(n).Len() > 0
}
