package rtree

import (
	"reflect"
	"testing"
)

func TestRtreeRPCInsertDeleteNN(t *testing.T) {
	s, _ := NewServer("tcp", "127.0.0.1:55976")
	go (func() { s.LoopAccept() })()
	c, _ := NewClient("tcp", "127.0.0.1:55976")
	defer c.Close()
	defer deleteAllKeys()

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

	// Update "c" so that it will take precedence over "b"
	neighbors, _ = c.RtreeNearestNeighbors("test", 5, []float64{6, 5})
	if !reflect.DeepEqual([]string{"b", "c"}, neighbors) {
		t.Errorf("[b, c] != %v", neighbors)
	}
	c.RtreeInsert("test", "c", []float64{5.9, 5.2}, []float64{25, 10})
	neighbors, _ = c.RtreeNearestNeighbors("test", 5, []float64{6, 5})
	if !reflect.DeepEqual([]string{"c", "b"}, neighbors) {
		t.Errorf("[c, b] != %v", neighbors)
	}
}

func TestRtreeRPCUpdate(t *testing.T) {
	c, _ := NewClient("tcp", "127.0.0.1:55976")
	defer c.Close()
	defer deleteAllKeys()

	c.RtreeInsert("test", "z", []float64{6, 6}, []float64{3, 3})

	// There should be an error since the member "a" doesn't exist yet.
	err := c.RtreeUpdate("test", "a", []float64{10, 10}, []float64{2, 2})
	if err == nil {
		t.Errorf("Expected RtreeUpdate to return an error")
	}

	c.RtreeInsert("test", "a", []float64{0, 0}, []float64{3.5, 4.2})
	neighbors, _ := c.RtreeNearestNeighbors("test", 1, []float64{11, 11})
	if !reflect.DeepEqual([]string{"z"}, neighbors) {
		t.Errorf("[z] != %v", neighbors)
	}

	args, _ :=
		NewRtreeInsertArgs("test", "a", []float64{10, 10}, []float64{2, 2})
	call := c.RtreeUpdateGo(args)

	// The nearest neighbor should be "a" now, after we updated "a"
	neighbors, _ = c.RtreeNearestNeighbors("test", 1, []float64{11, 11})
	if !reflect.DeepEqual([]string{"a"}, neighbors) {
		t.Errorf("[z] != %v", neighbors)
	}

	replyCall := <-call.Done
	reply, ok := replyCall.Reply.(*RtreeInsertReply)
	if !(replyCall.ServiceMethod == "Store.RtreeUpdate" &&
		ok && reply.Member == "a" &&
		replyCall.Error == nil) {
		t.Errorf("Deferred call error")
	}
}

func TestRtreeRPCSize(t *testing.T) {
	c, _ := NewClient("tcp", "127.0.0.1:55976")
	defer c.Close()
	defer deleteAllKeys()

	size, _ := c.RtreeSize("test")
	if size != 0 {
		t.Errorf("Rtree size %v is not 0", size)
	}

	c.RtreeInsert("test", "z", []float64{6, 6}, []float64{3, 3})
	c.RtreeInsert("test", "y", []float64{6, 6}, []float64{3, 3})
	size, _ = c.RtreeSize("test")
	if size != 2 {
		t.Errorf("Rtree size %v is not 0", size)
	}
}

// Test helper to delete all keys.
// We expect our tests to create less than 1 million keys,
// so finding 1 million neighbors and deleting them all should be sufficient.
func deleteAllKeys() {
	c, _ := NewClient("tcp", "127.0.0.1:55976")
	defer c.Close()
	neighbors, _ := c.RtreeNearestNeighbors("test", 1000000, []float64{0, 0})
	for _, v := range neighbors {
		c.RtreeDelete("test", v)
	}
}
