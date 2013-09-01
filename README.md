Rtree
-----
The R* tree is a specialized data structure for indexing spatial information
http://en.wikipedia.org/wiki/R*_tree . This library is a RPC server wrapper
around the implementation https://github.com/dhconnelly/rtreego . The codec
used here is net/http/json and thus it should be easy to support clients for
languages other than Go. As an example, a ruby client is provided here.

## Install
`go get github/fumin/rtree`

## Usage
Run `go run main/main.go` in shell to start the server.
Then use the following client code in Go:
```
import "github/fumin/rtree"
// The client `c` below is the standard Go RPC client.
// As explained in http://golang.org/pkg/net/rpc/#Client ,
// this client supports concurrent requests and may be used by
// multiple goroutines simultaneously.
c, err := NewClient("tcp", ":6389")
defer c.Close()

// Insert "a" with the rectangle {point: [0, 0], x_length: 3.5, y_length: 4.2}
// into the tree named "test".
c.RtreeInsert("test", "a", []float64{0, 0}, []float64{3.5, 4.2})
c.RtreeInsert("test", "b", []float64{2, 2}, []float64{3.1, 2.7})
c.RtreeInsert("test", "c", []float64{100, 100}, []float64{1, 1})

// Make an asynchronous call to get the 2 nearest neighbors of (3.4, 4.201)
call := c.RtreeNearestNeighborsGo("test", 2, []float64{3.4, 4.201})
// Do something else...
replyCall := <-call.Done
nnReply := replyCall.Reply.(*RtreeNearestNeighborsReply)
neighbors := nnReply.Members // should be ["b", "a"] as "c" is too far away

// Similar to a hash table,
// inserting on an existing member actually updates it.
c.RtreeInsert("test", "b", []float64{1000, 1000}, []float64{3.5, 4.2})
neighbors, err = r.RtreeNearestNeighbors("test", 2, []float64{3.4, 4.201})
  // neighbors == ["a", "c"] since "b" is now the farthest.
```
As an example of a client implemented in a language other than Go,
try out the ruby client with `ruby main/client.rb`

## License
Free use of this software is granted under the terms of the GNU Affero General Public License, version 3. Copyright (c) 2013 fumin.
