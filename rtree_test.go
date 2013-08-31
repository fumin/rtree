package rtree

import (
  "testing"
  "reflect"
  "github.com/dhconnelly/rtreego"
)

func TestRtree(t *testing.T) {
  rect1, _ := rtreego.NewRect(rtreego.Point{0, 0}, []float64{3.5, 4.2})
  rect2, _ := rtreego.NewRect(rtreego.Point{2, 2}, []float64{3.1, 2.7})
  rect3, _ := rtreego.NewRect(rtreego.Point{100, 100}, []float64{1, 1})
  rt := NewTree(2)
  if dim := rt.Dimension(); 2 != dim {
    t.Errorf("2 != rt.Dimension() %d", dim)
  }

  rt.Insert("a", rect1)
  rt.Insert("b", rect2)
  rt.Insert("c", rect3)
  neighbors := rt.NearestNeighbors(5, rtreego.Point{3.4, 4.2001})
  if !reflect.DeepEqual([]string{"b", "a", "c"}, neighbors) {
    t.Errorf("[b, a, c] != %v", neighbors)
  }

  rt.Delete("a")
  if s := rt.Size(); 2 != s {
    t.Errorf("2 != rt.Size() %v", rt.Size())
  }
  neighbors = rt.NearestNeighbors(5, rtreego.Point{3.4, 4.2001})
  if !reflect.DeepEqual([]string{"b", "c"}, neighbors) {
    t.Errorf("[b, c] != %v", neighbors)
  }

  rt.Delete("non-existent")
  if s := rt.Size(); 2 != s {
    t.Errorf("2 != rt.Size() %v", rt.Size())
  }
}
