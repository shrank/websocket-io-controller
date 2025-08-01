package utils

import (
  "testing"
)

type WebsocketUpdate struct {
	MsgType string
  Target string
}


func TestNext(t *testing.T) {
  q := CreateQueue()
  
  q.Insert(WebsocketUpdate {
		MsgType :"a",
		Target: "device",
	})
  q.Insert(WebsocketUpdate {
		MsgType :"b",
		Target: "device",
	})
  q.Insert(WebsocketUpdate {
		MsgType :"c",
		Target: "device",
	})

  rev, res, err := q.Next(0)
	item, _ := res.(WebsocketUpdate)
	if(rev != 1) {
    t.Errorf("rev = %d; want 1", rev)
  }
  if(item.MsgType != "a") {
    t.Errorf("item = '%s'; want 'a'", item)
  }
  if(err != nil) {
    t.Errorf("err not nil")
  }

  rev, res, err = q.Next(1)
  item, _ = res.(WebsocketUpdate)
	if(rev != 2) {
    t.Errorf("rev = %d; want 2", rev)
  }
  if(item.MsgType != "b") {
    t.Errorf("item = '%s'; want 'b'", item)
  }
  if(err != nil) {
    t.Errorf("err not nil")
  }

  rev, res, err = q.Next(2)
  item, _ = res.(WebsocketUpdate)
	if(rev != 3) {
    t.Errorf("rev = %d; want 3", rev)
  }
  if(item.MsgType != "c") {
    t.Errorf("item = '%s'; want 'c'", item)
  }
  if(err != nil) {
    t.Errorf("err not nil")
  }

  rev, res, err = q.Next(3)
  if(err == nil) {
    t.Errorf("err is nil")
  }
}


func TestLast(t *testing.T) {
  q := CreateQueue()
  
  q.Insert(WebsocketUpdate {
		MsgType :"a",
		Target: "device",
	})
  q.Insert(WebsocketUpdate {
		MsgType :"b",
		Target: "device",
	})
  q.Insert(WebsocketUpdate {
		MsgType :"c",
		Target: "device",
	})

  rev, res, err := q.Last()
	item := res.(WebsocketUpdate)
  if(rev != 3) {
    t.Errorf("rev = %d; want 3", rev)
  }
  if(item.MsgType != "c") {
    t.Errorf("item = '%s'; want 'c'", item)
  }
  if(err != nil) {
    t.Errorf("err not nil")
  }
}