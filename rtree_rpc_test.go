package rtree

import (
  "testing"
  "reflect"
)

func TestRtreeRPC(t *testing.T) {
  s, _ := NewServer("tcp", "127.0.0.1:55976")
  go (func(){ s.LoopAccept() })()
  c, _ := NewClient("tcp", "127.0.0.1:55976")
  defer c.Close()

  c.RtreeInsert("test", "a", []float64{0, 0}, []float64{3.5, 4.2})

  // Make a deferred call for inserting "b"
  args, _ :=
    NewRtreeInsertArgs("test", "b", []float64{2, 2}, []float64{3.1, 2.7})
  call := c.RtreeInsertGo(args)

  c.RtreeInsert("test", "c", []float64{100, 100}, []float64{1, 1})
  neighbors, _ := c.RtreeNearestNeighbors("test", 5, []float64{3.4, 4.201})
  if !reflect.DeepEqual([]string{"b", "a", "c"}, neighbors) {
    t.Errorf("[b, a, c] != %v", neighbors)
  }

  // Catch the reply of the insert call of "b"
  replyCall := <-call.Done
  reply, ok := replyCall.Reply.(*RtreeInsertReply)
  if !(replyCall.ServiceMethod == "Store.RtreeInsert" &&
       ok && reply.Member == "b" &&
       replyCall.Error == nil) {
    t.Errorf("Deferred call error")
  }

  c.RtreeDelete("another_tree", "a")
  // Make a deferred nearest neighbors call.
  // The result should still be (b, a, c) since the tree we deleted just now
  // "another_tree" is different from the "test" tree
  call = c.RtreeNearestNeighborsGo("test", 5, []float64{3.4, 4.201})

  c.RtreeDelete("test", "a")
  neighbors, _ = c.RtreeNearestNeighbors("test", 5, []float64{3.4, 4.2001})
  if !reflect.DeepEqual([]string{"b", "c"}, neighbors) {
    t.Errorf("[b, c] != %v", neighbors)
  }

  c.RtreeDelete("test", "non-existent-member")
  neighbors, _ = c.RtreeNearestNeighbors("test", 5, []float64{3.4, 4.2001})
  if !reflect.DeepEqual([]string{"b", "c"}, neighbors) {
    t.Errorf("[b, c] != %v", neighbors)
  }

  replyCall = <-call.Done
  nnReply := replyCall.Reply.(*RtreeNearestNeighborsReply)
  neighbors = nnReply.Members
  if !reflect.DeepEqual([]string{"b", "a", "c"}, neighbors) {
    t.Errorf("[b, a, c] != %v", neighbors)
  }
}
