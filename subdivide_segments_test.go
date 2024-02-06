package martinez_go

import (
	"fmt"
	"github.com/akayami/martinez-go/helpers"
	"github.com/paulmach/orb"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestSubdivideSegments(t *testing.T) {

	t.Run("divide 2 segments", func(t *testing.T) {
		se1 := NewSweepEvent(Point{0, 0}, true, NewSweepEvent(Point{5, 5}, false, nil, false, 0), true, 0)
		se2 := NewSweepEvent(Point{0, 5}, true, NewSweepEvent(Point{5, 0}, false, nil, false, 0), false, 0)

		q := NewTinyQueueDefault(nil, CompareEvents)
		q.Push(se1)
		q.Push(se2)

		iter := SegmentIntersection(se1.Point, se1.OtherEvent.Point, se2.Point, se2.OtherEvent.Point, false)

		DivideSegment(se1, iter[0], q)
		DivideSegment(se2, iter[0], q)

		assert.Equal(t, 6, q.Len(), "subdivided in 4 segments by intersection point")
	})

	t.Run("possible intersections", func(t *testing.T) {

		geojson, err := helpers.LoadGeoJson("test/fixtures/two_shapes.geojson")
		assert.Nil(t, err)

		q := NewTinyQueueDefault(nil, CompareEvents)

		s := geojson.Features[0].Geometry.(orb.Polygon)
		c := geojson.Features[1].Geometry.(orb.Polygon)

		se1 := NewSweepEvent(OrbPointToPoint(s[0][3]), true, NewSweepEvent(OrbPointToPoint(s[0][2]), false, nil, false, 0), true, 0)
		se2 := NewSweepEvent(OrbPointToPoint(c[0][0]), true, NewSweepEvent(OrbPointToPoint(c[0][1]), false, nil, false, 0), false, 0)

		assert.Equal(t, 1, PossibleIntersection(se1, se2, q))
		assert.Equal(t, 4, q.Len())

		t.Run("Check 1", func(t *testing.T) {
			e := q.Pop()
			assert.NotNil(t, e)
			assert.Equal(t, Point{100.79403384562251, 233.41363754101192}, e.Point)
			assert.Equal(t, Point{56, 181}, e.OtherEvent.Point, "1")
		})

		t.Run("Check 2", func(t *testing.T) {
			e := q.Pop()
			assert.NotNil(t, e)
			assert.Equal(t, Point{100.79403384562251, 233.41363754101192}, e.Point)
			assert.Equal(t, Point{16, 282}, e.OtherEvent.Point, "2")
		})

		t.Run("Check 3", func(t *testing.T) {
			e := q.Pop()
			assert.NotNil(t, e)
			assert.Equal(t, Point{100.79403384562251, 233.41363754101192}, e.Point)
			assert.Equal(t, Point{153, 203.5}, e.OtherEvent.Point, "3")
		})

		t.Run("Check 4", func(t *testing.T) {
			e := q.Pop()
			assert.NotNil(t, e)
			assert.Equal(t, Point{100.79403384562251, 233.41363754101192}, e.Point)
			assert.Equal(t, Point{153, 294.5}, e.OtherEvent.Point, "4")
		})

	})

	t.Run("possible intersections on 2 polygons", func(t *testing.T) {
		geojson, err := helpers.LoadGeoJson("test/fixtures/two_shapes.geojson")
		assert.Nil(t, err)
		s := append([][][]Point{}, OrbPolygonToPolygon(geojson.Features[0].Geometry.(orb.Polygon)))
		c := append([][][]Point{}, OrbPolygonToPolygon(geojson.Features[1].Geometry.(orb.Polygon)))

		bbox := &[4]float64{math.Inf(1), math.Inf(1), math.Inf(-1), math.Inf(-1)}

		q := FillQueue(s, c, bbox, bbox, Normal)

		p0 := Point{X: 16, Y: 282}
		p1 := Point{X: 298, Y: 359}
		p2 := Point{X: 156, Y: 203.5}

		te := NewSweepEvent(p0, true, nil, true, Normal)
		te2 := NewSweepEvent(p1, false, te, false, Normal)
		te.OtherEvent = te2

		te3 := NewSweepEvent(p0, true, nil, true, Normal)
		te4 := NewSweepEvent(p2, true, te3, false, Normal)
		te3.OtherEvent = te4

		tr := NewSplayTree(CompareSegments)
		tr.Insert(NewSplayTreeNode(te))
		tr.Insert(NewSplayTreeNode(te3))

		assert.Equal(t, te, tr.Find(te).Key)
		assert.Equal(t, te3, tr.Find(te3).Key)

		assert.Equal(t, 1, CompareSegments(te, te3))
		assert.Equal(t, -1, CompareSegments(te3, te))

		segments := SubdivideSegments(q, s, c, bbox, bbox, Intersection)
		leftSegments := []*SweepEvent{}

		for i := 0; i < len(segments); i++ {
			if segments[i].Left {
				leftSegments = append(leftSegments, segments[i])
			}
		}

		assert.Equal(t, 11, len(leftSegments))

		E := Point{16, 282}
		I := Point{100.79403384562252, 233.41363754101192}
		G := Point{298, 359}
		C := Point{153, 294.5}
		J := Point{203.36313843035356, 257.5101243166895}
		F := Point{153, 203.5}
		D := Point{56, 181}
		A := Point{108.5, 120}
		B := Point{241.5, 229.5}

		t.Run("Intervals Section", func(t *testing.T) {
			type interval struct {
				l            Point
				r            Point
				inOut        bool
				otherInOut   bool
				inResult     bool
				prevInResult *interval
			}

			sCJ := interval{l: F, r: J, inOut: false, otherInOut: false, inResult: false, prevInResult: nil}
			sIC := interval{l: I, r: F, inOut: false, otherInOut: false, inResult: false, prevInResult: nil}

			intervals := map[string]interval{
				"EI": {l: E, r: I, inOut: false, otherInOut: true, inResult: false, prevInResult: nil},
				"IF": {l: I, r: F, inOut: false, otherInOut: false, inResult: true, prevInResult: nil},
				"FJ": {l: F, r: J, inOut: false, otherInOut: false, inResult: true, prevInResult: nil},
				"JG": {l: J, r: G, inOut: false, otherInOut: true, inResult: false, prevInResult: nil},
				"EG": {l: E, r: G, inOut: true, otherInOut: true, inResult: false, prevInResult: nil},
				"DA": {l: D, r: A, inOut: false, otherInOut: true, inResult: false, prevInResult: nil},
				"AB": {l: A, r: B, inOut: false, otherInOut: true, inResult: false, prevInResult: nil},
				"JB": {l: J, r: B, inOut: true, otherInOut: true, inResult: false, prevInResult: nil},
				"CJ": {l: C, r: J, inOut: true, otherInOut: false, inResult: true, prevInResult: &sCJ},
				"IC": {l: I, r: C, inOut: true, otherInOut: false, inResult: true, prevInResult: &sIC},
			}

			var checkContain = func(interval string) {
				data := intervals[interval]
				for x := 0; x < len(leftSegments); x++ {
					seg := leftSegments[x]
					if Equals(seg.Point, data.l) &&
						Equals(seg.OtherEvent.Point, data.r) &&
						seg.InOut == data.inOut &&
						seg.OtherInOut == data.otherInOut &&
						seg.InResult() == data.inResult &&
						((seg.PrevInResult == nil && data.prevInResult == nil) ||
							(seg.PrevInResult != nil && data.prevInResult != nil &&
								Equals(seg.PrevInResult.Point, data.prevInResult.l) &&
								Equals(seg.PrevInResult.OtherEvent.Point, data.prevInResult.r))) {
						return
					}
				}
				// assert.Error(t, error.Error("Fail"))
				panic(fmt.Sprintf("Fail: %s", interval))
			}

			for key := range intervals {
				t.Run(fmt.Sprintf("Test: %s", key), func(t *testing.T) {
					checkContain(key)
				})
			}
		})
	})
}
