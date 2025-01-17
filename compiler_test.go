package opencypher

import (
	"testing"

	"github.com/cloudprivacylabs/opencypher/graph"
	"github.com/cloudprivacylabs/opencypher/parser"
)

func TestExpr(t *testing.T) {
	c := GetParser(`5  +  7+1`).OC_Expression()
	out := oC_Expression(c.(*parser.OC_ExpressionContext))
	result, err := out.Evaluate(NewEvalContext(graph.NewOCGraph()))
	if err != nil {
		t.Error(err)
	}
	if result.Value != 13 {
		t.Errorf("Wrong result: %+v %T", result, result.Value)
	}
}

func runTestMatch(t *testing.T, expr string, g graph.Graph) ResultSet {
	ctx := NewEvalContext(g)
	ev, err := Parse(expr)
	if err != nil {
		t.Errorf("%s: %s", expr, err)
		return ResultSet{}
	}
	value, err := ev.Evaluate(ctx)
	if err != nil {
		t.Errorf("%s: %s", expr, err)
		return ResultSet{}
	}

	return value.Value.(ResultSet)
}

func TestBasicMatch(t *testing.T) {

	g := graph.NewOCGraph()
	// (:t1:t2)
	n1 := g.NewNode([]string{"t1", "t2"}, nil)
	rs := runTestMatch(t, "match (n) return n", g)
	if !rs.Rows[0]["1"].Value.(graph.Node).GetLabels().HasAll("t1", "t2") {
		t.Errorf("Expecting to see one row with t1, t2")
	}

	// (:t1:t2) (:t1:t3)
	n2 := g.NewNode([]string{"t1", "t3"}, nil)
	rs = runTestMatch(t, "match (n:t1) return n", g)
	if !(len(rs.Rows) == 2 && rs.Rows[0]["1"].Value.(graph.Node).GetLabels().Has("t1") &&
		rs.Rows[1]["1"].Value.(graph.Node).GetLabels().Has("t1")) {
		t.Errorf("Expecting to see two rows with t1, t2")
	}
	rs = runTestMatch(t, "match (n:t2) return n", g)
	if !(len(rs.Rows) == 1 && rs.Rows[0]["1"].Value.(graph.Node).GetLabels().Has("t1")) {
		t.Errorf("Expecting to see one row with t1: %v", rs)
	}

	// (:t1:t2) --[:e1]-->(:t1:t3)
	g.NewEdge(n1, n2, "e1", nil)
	rs = runTestMatch(t, "match (m:t1)-[:e1]->(n:t3) return n", g)
	if !(len(rs.Rows) == 1 && rs.Rows[0]["1"].Value.(graph.Node) == n2) {
		t.Errorf("Expecting to see one row n2: %v", rs)
	}
	rs = runTestMatch(t, "match (m:t1)<-[:e1]-(n:t1) return n", g)
	if !(len(rs.Rows) == 1 && rs.Rows[0]["1"].Value.(graph.Node) == n1) {
		t.Errorf("Expecting to see one row n1: %v", rs)
	}
	rs = runTestMatch(t, "match (m:t1)<-[:e1]-(n:t1) return n,m", g)
	if !(len(rs.Rows) == 1 && rs.Rows[0]["1"].Value.(graph.Node) == n1 && rs.Rows[0]["2"].Value.(graph.Node) == n2) {
		t.Errorf("Expecting to see one row n1 n2: %v", rs)
	}

	// (:t1:t2) --[:e1]-->(:t1:t3) --[:e2]-->(:t4)
	n3 := g.NewNode([]string{"t4"}, nil)
	g.NewEdge(n2, n3, "e2", nil)
	rs = runTestMatch(t, "match (m:t1)-[*]->(n:t4) return n", g)
	if !(len(rs.Rows) == 2 && rs.Rows[0]["1"].Value.(graph.Node) == n3 && rs.Rows[1]["1"].Value.(graph.Node) == n3) {
		t.Errorf("Expecting to see two rows n3: %v", rs)
	}

}
