package shutdown

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

type Shutdown struct {
	fnStop func()
	ctx    context.Context
}

// New returns a new Shutdown
func New() *Shutdown {
	return NewWithContext(context.Background())
}

// NewWithContext returns a new Shutdown passing context
func NewWithContext(ctx context.Context) *Shutdown {
	return NewWithContextAndSignals(ctx, os.Interrupt, syscall.SIGTERM)
}

// NewWithSignals returns a new Shutdown passing signals
func NewWithSignals(sig ...os.Signal) *Shutdown {
	return NewWithContextAndSignals(context.Background(), sig...)
}

// NewWithContextAndSignals returns a new Shutdown passing context and signals
func NewWithContextAndSignals(ctx context.Context, sig ...os.Signal) *Shutdown {
	ctx, fnStop := signal.NotifyContext(ctx, sig...)
	return &Shutdown{
		fnStop: fnStop,
		ctx:    ctx,
	}
}

// Stop stops the shutdown
func (s *Shutdown) Stop() {
	s.fnStop()
}

// Start starts the shutdown controller
func (s *Shutdown) Start(fn func()) {
	go fn()
	<-s.ctx.Done()
	fmt.Println("stopping service gracefully!!!")
}
