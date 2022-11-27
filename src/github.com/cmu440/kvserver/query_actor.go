package kvserver

import (
	"fmt"
	"net/rpc"
	"encoding/gob"
	"github.com/cmu440/actor"
)

// Implement your queryActor in this file.
// See example/counter_actor.go for an example actor using the
// github.com/cmu440/actor package.

// TODO (3A, 3B): define your message types as structs

func init() {
	// TODO (3A, 3B): Register message types, e.g.:
	gob.Register(GetArgs{})
	gob.Register(GetReply{})
	gob.Register(ListArgs{})
	gob.Register(ListReply{})
	gob.Register(PutArgs{})
	gob.Register(PutReply{})
}

type queryActor struct {
	context *actor.ActorContext
	// TODO (3A, 3B): implement this!
	kvstore map[string]string
}

// "Constructor" for queryActors, used in ActorSystem.StartActor.
func newQueryActor(context *actor.ActorContext) actor.Actor {
	return &queryActor{
		context: context,
		// TODO (3A, 3B): implement this!
		kvstore: make(map[string]string)
	}
}

// OnMessage implements actor.Actor.OnMessage.
func (actor *queryActor) OnMessage(message any) error {
	// TODO (3A, 3B): implement this!
	switch m := message.(type) {
	case MGet: 
		key, getCh := m.Key, m.GetCh
		getReply := &GetReply{}
		if value, ok := actor.kvstore[key]; ok {
			getReply.Value = value
			getReply.Ok = true
		} else {
			getReply.Ok = false
		}
		getCh <- getReply
	case MList:
		pref, listCh := m.Prefix, m.ListCh
		entries := make(map[string]string)
		for k, v := range actor.kvstore {
			if isPrefix(pref, k) {
				entries[k] = v
			}
		}
		listReply := &ListReply{entries}
		listCh <- listReply
	case MPut:
		key, value := m.Key, m.Value
		actor.kvstore[key] = value
	default:
		return fmt.Errorf("Unexpected queryActor message type: %T", m)
	}
	return nil
}

// ======================== actor message types ==============

type MGet struct {
	Key string
	GetCh chan *GetReply
}

type MList struct {
	Prefix string 
	ListCh chan *ListReply
}

type MPut struct {
	Key	  string
	Value string
}

// ==================== Helper functions ======================
func isPrefix(prefix string, key string) bool {
	if len(prefix) > len(key) {
		return false
	} // len(prefix) <= len(key)
	return prefix == key[:len(prefix)]
}