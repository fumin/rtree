package rtree

import (
  "net"
  "net/rpc"
  "net/rpc/jsonrpc"
  "github.com/dhconnelly/rtreego"
)

type server_t struct{
  net.Listener
}
func (s *server_t) LoopAccept() error {
  for {
      conn, err := s.Accept()
      if err != nil { return err }
      go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
  }
}
func NewServer(netp, laddr string) (*server_t, error) {
  store := NewStore()
  rpc.Register(store)
  rpc.HandleHTTP()
  l, err := net.Listen(netp, laddr)
  if err != nil { return nil, err }
  return &server_t{l}, nil
}

type Client struct {
  *rpc.Client
}
func NewClient(netp, laddr string) (*Client, error) {
  conn, err := net.Dial(netp, laddr)
  if err != nil { return nil, err }
  return &Client{jsonrpc.NewClient(conn)}, nil
}

// RtreeInsert
func (c *Client) RtreeInsert(key, member string,
                             point, lengths []float64) error {
  args, err := NewRtreeInsertArgs(key, member, point, lengths)
  if err != nil { return err }
  var reply RtreeInsertReply
  err = c.Call("Store.RtreeInsert", args, &reply)
  if err != nil { return err }
  return nil
}
func (c *Client) RtreeInsertGo(args *RtreeInsertArgs) *rpc.Call {
  var reply RtreeInsertReply
  return c.Go("Store.RtreeInsert", args, &reply, nil)
}

// RtreeDelete
func (c *Client) RtreeDelete(key, member string) error {
  reply := new(string)
  err := c.Call("Store.RtreeDelete", &RtreeDeleteArgs{key, member}, reply)
  if err != nil { return err }
  return nil
}
func (c *Client) RtreeDeleteGo(key, member string) *rpc.Call {
  reply := new(string)
  return c.Go("Store.RtreeDelete", &RtreeDeleteArgs{key, member}, reply, nil)
}

// RtreeNearestNeighbors
func (c *Client) RtreeNearestNeighbors(key string, k int, p rtreego.Point) ([]string, error) {
  args := &RtreeNearestNeighborsArgs{key, k, p}
  reply := new(RtreeNearestNeighborsReply)
  err := c.Call("Store.RtreeNearestNeighbors", args, reply)
  if err != nil { return nil, err }
  return reply.Members, nil
}
func (c *Client) RtreeNearestNeighborsGo(key string, k int, p rtreego.Point) *rpc.Call {
  args := &RtreeNearestNeighborsArgs{key, k, p}
  reply := new(RtreeNearestNeighborsReply)
  return c.Go("Store.RtreeNearestNeighbors", args, reply, nil)
}
