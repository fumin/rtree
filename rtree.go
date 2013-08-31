// Key centric Rtree based on github.com/dhconnelly/rtreego
//
// Usage:
// rect1, _ := rtreego.NewRect(rtreego.Point{0, 0}, []float64{3.5, 4.2})
// rect2, _ := rtreego.NewRect(rtreego.Point{2, 2}, []float64{3.1, 2.7})
// rect3, _ := rtreego.NewRect(rtreego.Point{100, 100}, []float64{1, 1})
// rt := NewTree(2)
// rt.Insert("a", rect1)
// rt.Insert("b", rect2)
// rt.Insert("c", rect3)
//
// t.NearestNeighbors(5, rtreego.Point{3.4, 4.201})
//   => ["b", "a", "c"]
//
// The metric definition of the nearest neighbors of a point P is as follows:
// * Rectangles that enclose P have distance 0.
// * For rectangles that don't enclose P, the metric is the square of the
//   Euclidean distance between P and the nearest edge of the rectangle.
// * Reference: http://www.postgis.org/support/nearestneighbor.pdf
//
// This implementation doesn't support concurrency, please use sync.Mutex
// in multithreaded environments.

package rtree

import (
  "fmt"
  "github.com/dhconnelly/rtreego"
)

type Rtree struct {
  rt *rtreego.Rtree
  keyMap map[string]*thing_t
}

func NewTree(dimension int) *Rtree {
  rt := rtreego.NewTree(dimension, 25, 50)
  return &Rtree{rt: rt, keyMap: make(map[string]*thing_t)}
}

func (t *Rtree) Insert(key string, where *rtreego.Rect) {
  thing, ok := t.keyMap[key]
  if ok {
    ok = t.rt.Delete(thing)
    if !ok {
      panic(fmt.Sprintf("Object found in keyMap but not in rtree: %v", thing))
    }
  }
  newThing := &thing_t{where, key}
  t.rt.Insert(newThing)
  t.keyMap[key] = newThing
}

func (t *Rtree) Delete(key string) {
  thing, ok := t.keyMap[key]
  if ok {
    ok = t.rt.Delete(thing)
    if !ok {
      panic(fmt.Sprintf("Object found in keyMap but not in rtree: %v", thing))
    }
    delete(t.keyMap, key)
  }
}

func (t *Rtree) NearestNeighbors(k int, p rtreego.Point) []string {
  things := t.rt.NearestNeighbors(k, p)

  // Drop nil elements
  i := 0
  for ; i != len(things); i++ {
    if things[i] == nil { break }
  }
  things = things[0:i]

  // Get the keys out from the returned pointer to things
  keys := make([]string, len(things))
  for i, v := range things {
    thing, ok := v.(*thing_t)
    if !ok { panic(fmt.Sprintf("Object of unrecognized type stored: %v", v)) }
    keys[i] = thing.key
  }
  return keys
}

func (t *Rtree) Size() int {
  keyMapSize := len(t.keyMap)
  rtreeSize := t.rt.Size()
  if keyMapSize != rtreeSize {
    panic(fmt.Sprintf("keyMapSize %d != rtreeSize %d", keyMapSize, rtreeSize))
  }
  return keyMapSize
}

func (t *Rtree) Dimension() int {
  return t.rt.Dim
}

// Our custom implementation of rtreego's Spatial interface
// type Spatial interface {
//   Bounds() *Rect
// }
type thing_t struct {
  where *rtreego.Rect
  key string
}
func (t *thing_t) Bounds() *rtreego.Rect {
  return t.where
}
