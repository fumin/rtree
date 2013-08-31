package rtree

import (
	"errors"
	"fmt"
	"github.com/dhconnelly/rtreego"
	"reflect"
	"sync"
)

// The main RPC object
type Store struct {
	mutex  sync.RWMutex
	keyMap map[string]interface{}
}

func NewStore() *Store {
	return &Store{keyMap: make(map[string]interface{})}
}

// Struct for use in the asyncronous RPC call RtreeInsertGo.
// Note that to ensure data integrity, we should always instantiate
// this struct using NewRtreeInsertArgs
type RtreeInsertArgs struct {
	Key    string
	Member string
	Where  rect
}

type rect struct {
	Point   []float64
	Lengths []float64
}

func NewRtreeInsertArgs(key, member string,
	point, lengths []float64) (*RtreeInsertArgs, error) {
	if len(point) != len(lengths) {
		errMsg := fmt.Sprintf(
			"Different dimensions for point %v and lengths %v", point, lengths)
		return nil, errors.New(errMsg)
	}
	return &RtreeInsertArgs{key, member, rect{point, lengths}}, nil
}

// Struct for use in the reply of asyncronous RPC call RtreeInsertGo.
type RtreeInsertReply struct {
	Member string
}

func (s *Store) RtreeInsert(args *RtreeInsertArgs,
	reply *RtreeInsertReply) error {
	dimension := len(args.Where.Point)
	if dimension != len(args.Where.Lengths) {
		return errors.New(fmt.Sprintf("Wrong dimensions for Rect %v", args.Where))
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// Initialize the rtree "rt"
	var rt *Rtree
	obj, ok := s.keyMap[args.Key]
	if !ok {
		rt = NewTree(dimension)
		s.keyMap[args.Key] = rt
	} else {
		rt, ok = obj.(*Rtree)
		if !ok {
			typeName := reflect.TypeOf(obj).String()
			errMsg := fmt.Sprintf("The type of %v is %v", args.Key, typeName)
			return errors.New(errMsg)
		}
		if dimension != rt.Dimension() {
			errMsg := fmt.Sprintf(
				"Different dimensions between the rtree %v and Rect: %v",
				args.Key, args.Where)
			return errors.New(errMsg)
		}
	}

	rect, err := rtreego.NewRect(args.Where.Point, args.Where.Lengths)
	if err != nil {
		return err
	}
	rt.Insert(args.Member, rect)
	reply.Member = args.Member
	return nil
}

type RtreeDeleteArgs struct {
	Key, Member string
}

func (s *Store) RtreeDelete(args *RtreeDeleteArgs, reply *string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	obj, ok := s.keyMap[args.Key]
	if !ok {
		return errors.New(fmt.Sprintf("No object for key %v", args.Key))
	}

	// Initialize the rtree "rt"
	rt, ok := obj.(*Rtree)
	if !ok {
		typeName := reflect.TypeOf(obj).String()
		errMsg := fmt.Sprintf("The type of %v is %v", args.Key, typeName)
		return errors.New(errMsg)
	}

	rt.Delete(args.Member)
	reply = &args.Member
	return nil
}

// Struct for use in the asynchronous RPC call RtreeNearestNeighbors.
type RtreeNearestNeighborsArgs struct {
	Key   string
	K     int
	Point rtreego.Point
}

// Struct for use in the reply of the
// asynchronous RPC call RtreeNearestNeighbors.
type RtreeNearestNeighborsReply struct {
	Members []string
}

func (s *Store) RtreeNearestNeighbors(
	args *RtreeNearestNeighborsArgs, reply *RtreeNearestNeighborsReply) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	obj, ok := s.keyMap[args.Key]
	if !ok {
		return errors.New(fmt.Sprintf("No object for key %v", args.Key))
	}

	// Initialize the rtree "rt"
	rt, ok := obj.(*Rtree)
	if !ok {
		typeName := reflect.TypeOf(obj).String()
		errMsg := fmt.Sprintf("The type of %v is %v", args.Key, typeName)
		return errors.New(errMsg)
	}
	if dim := rt.Dimension(); dim != len(args.Point) {
		errTemplate := "Rtree dimension %d doesn't match point %v"
		return errors.New(fmt.Sprintf(errTemplate, dim, args))
	}

	members := rt.NearestNeighbors(args.K, args.Point)
	reply.Members = members
	return nil
}
