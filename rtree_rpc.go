// RPC server and client functions to operate on Rtrees.
// It is encouraged to use these RPC functions instead of
// doing direct RPC calls to increase type safety and avoid bugs with
// string manipulation.
package rtree

import (
	"github.com/dhconnelly/rtreego"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Server struct {
	net.Listener
}

// Creates a new server that listens on the network netp and address laddr.
// To start the forever running Accept loop, call LoopAccept.
func NewServer(netp, laddr string) (*Server, error) {
	store := NewStore()
	rpc.Register(store)
	rpc.HandleHTTP()
	l, err := net.Listen(netp, laddr)
	if err != nil {
		return nil, err
	}
	return &Server{l}, nil
}

// Starts the forever running Accept loop. If the loop is halted
// due to an error, that error is returned.
func (s *Server) LoopAccept() error {
  for {
    conn, err := s.Accept()
    if err != nil {
      return err
    }
    go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
  }
}

type Client struct {
	*rpc.Client
}

// Creates a new client connecting to the network netp and address laddr.
// The returned client is simply a thin wrapper around rpc.Client, and thus
// according to http://golang.org/pkg/net/rpc/#Client supports
// concurrent requests and may be used by multiple goroutines simultaneously.
func NewClient(netp, laddr string) (*Client, error) {
	conn, err := net.Dial(netp, laddr)
	if err != nil {
		return nil, err
	}
	return &Client{jsonrpc.NewClient(conn)}, nil
}

// Inserts a member into the Rtree identified by `key`.
// The rectangle associated with `member` is defined to be
// having the bottom corner `point` with lengths on all its dimensions
// given by `lengths`.
// Under this definition len(point) == len(lengths) and
// all entries of `lengths` should be positive.
func (c *Client) RtreeInsert(key, member string,
	point, lengths []float64) error {
	args, err := NewRtreeInsertArgs(key, member, point, lengths)
	if err != nil {
		return err
	}
	var reply RtreeInsertReply
	err = c.Call("Store.RtreeInsert", args, &reply)
	if err != nil {
		return err
	}
	return nil
}
// Asynchronous version of RtreeInsert.
// Usage is similar to http://golang.org/pkg/net/rpc/#Client.Go
func (c *Client) RtreeInsertGo(args *RtreeInsertArgs) *rpc.Call {
	var reply RtreeInsertReply
	return c.Go("Store.RtreeInsert", args, &reply, nil)
}

// Deletes a member from the Rtree identified by key.
func (c *Client) RtreeDelete(key, member string) error {
	reply := new(string)
	err := c.Call("Store.RtreeDelete", &RtreeDeleteArgs{key, member}, reply)
	if err != nil {
		return err
	}
	return nil
}
// Asynchronous version of RtreeDelete.
func (c *Client) RtreeDeleteGo(key, member string) *rpc.Call {
	reply := new(string)
	return c.Go("Store.RtreeDelete", &RtreeDeleteArgs{key, member}, reply, nil)
}

// Finds the k nearest neighbors around the point p in the Rtree
// identified by key.
func (c *Client) RtreeNearestNeighbors(key string, k int, p rtreego.Point) ([]string, error) {
	args := &RtreeNearestNeighborsArgs{key, k, p}
	reply := new(RtreeNearestNeighborsReply)
	err := c.Call("Store.RtreeNearestNeighbors", args, reply)
	if err != nil {
		return nil, err
	}
	return reply.Members, nil
}
// Asynchronous version of RtreeNearestNeighbors.
func (c *Client) RtreeNearestNeighborsGo(key string, k int, p rtreego.Point) *rpc.Call {
	args := &RtreeNearestNeighborsArgs{key, k, p}
	reply := new(RtreeNearestNeighborsReply)
	return c.Go("Store.RtreeNearestNeighbors", args, reply, nil)
}
