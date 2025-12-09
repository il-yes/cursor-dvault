package events


import (
	"log"
)

type LocalDispatcher struct{}

func NewLocalDispatcher() *LocalDispatcher { return &LocalDispatcher{} }

func (d *LocalDispatcher) Dispatch(e any) error {
	// simple synchronous logging-based dispatcher
	log.Printf("[domain event] %#v\n", e)
	return nil
}
