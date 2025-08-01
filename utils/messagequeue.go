package utils

import (
  "errors"
  "sync"
)

type Queue struct {
  mu sync.Mutex
  cond *sync.Cond
  capacity uint
  counter uint
  oldest uint
  q	map[uint] interface{}
}

// FifoQueue 
type FifoQueue interface {
  Insert()
  Wait()
  Next()
  Last()
}

// Insert inserts the item into the queue
func (q *Queue) Insert(item  interface{}) error {
  q.mu.Lock()
  defer q.mu.Unlock()
  if(q.counter - q.oldest > 100) {
    for i := 0; i < 50; i++ {
      q.oldest++
      delete(q.q, q.oldest)
    }
  }
  q.counter++
  q.q[q.counter]=item
  q.cond.Broadcast()
  return nil
}

// Wait for new items in queue
func (q *Queue) Wait(current uint) ( uint,  interface{},error) {
  rev, item, err := q.Next(current)
  if(err != nil) {
    q.cond.L.Lock()
    for (current >= q.counter) {
      q.cond.Wait()
    }
    q.cond.L.Unlock()
    rev, item, err = q.Next(current)
  }

  return rev, item, err
}

func (q *Queue) Next(current uint) (rev uint, item interface{}, err error) {
  if(current < q.counter) {
    q.mu.Lock()
    defer q.mu.Unlock()
    if(current < q.oldest) {
      return q.oldest, q.q[q.oldest], nil
    }
    return current+1, q.q[current+1], nil
  }
  return 0, "", errors.New("Queue Empty")
}

func (q *Queue) Last() (rev uint, item  interface{}, err error) {
  if(q.counter > 0) {
    q.mu.Lock()
    defer q.mu.Unlock()
    return q.counter, q.q[q.counter], nil
  }
  return 0, "",errors.New("Queue Empty")
}

func CreateQueue() *Queue {
  return &Queue{
      q:        make(map[uint] interface{}),
      cond: 	  sync.NewCond(&sync.Mutex{}),
  }
}
