package monitor

import (
	"sync"
	"time"
)

type Status struct {
	Connected      bool
	Iface          string
	VPNIP          string
	ExternalIP     string
	RouteProtected bool
	LastChange     time.Time
}

type StatusStore struct {
	mu     sync.RWMutex
	status Status
}

func NewStatusStore(iface string) *StatusStore {
	return &StatusStore{
		status: Status{
			Iface:      iface,
			LastChange: time.Now(),
		},
	}
}

func (s *StatusStore) Update(newStatus Status) {
	s.mu.Lock()
	defer s.mu.Unlock()

	updateCond := s.status.Connected != newStatus.Connected || s.status.RouteProtected != newStatus.RouteProtected
	if updateCond {
		newStatus.LastChange = time.Now()
	} else {
		newStatus.LastChange = s.status.LastChange
	}
	s.status = newStatus
}

func (s *StatusStore) Get() Status {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}
