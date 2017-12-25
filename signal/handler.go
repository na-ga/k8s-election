package signal

import (
	"os"
	"os/signal"
)

// Handler for handling os signals.
type Handler struct {
	SignalCh chan os.Signal
}

// NewHandler creates a new handler for handling specified signals.
func NewHandler(signals ...os.Signal) *Handler {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, signals...)
	return &Handler{SignalCh: ch}
}

// Close closes SignalCh of Handler.
func (sh *Handler) Close() {
	close(sh.SignalCh)
	signal.Stop(sh.SignalCh)
}
