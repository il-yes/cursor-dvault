package realtime_client_infrastructure

import (
	context "context"
	"sync"
	"time"
)

type ConnectionStatus string

const (
	Disconnected ConnectionStatus = "disconnected"
	Connecting   ConnectionStatus = "connecting"
	Connected    ConnectionStatus = "connected"
	Reconnecting ConnectionStatus = "reconnecting"
)


type ConnectionState struct {
	mu     sync.Mutex
	isOpen bool
}
	

type Reconnecter struct {
	mu        sync.Mutex
	isRunning bool
	stop      chan struct{}
}

func NewReconnecter() *Reconnecter {
	return &Reconnecter{
		isRunning: false,
		stop:      make(chan struct{}),
	}
}

func (r *Reconnecter) Start(reconnect func(ctx context.Context) error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isRunning {
		return
	}
	r.isRunning = true

	go func() {
		for {
			select {
			case <-r.stop:
				return
			default:
				if err := reconnect(context.Background()); err != nil {
					time.Sleep(1 * time.Second)
				}
			}
		}
	}()
}

func (r *Reconnecter) Stop() {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.isRunning {
		close(r.stop)
		r.isRunning = false
	}
}	