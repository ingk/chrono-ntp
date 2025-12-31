package timesource

import (
	"chrono-ntp/ntp"
	"sync"
	"time"
)

// Ensure NtpTimeSource implements the TimeSource interface
var _ TimeSource = (*NtpTimeSource)(nil)

type NtpTimeSource struct {
	ntpClient  *ntp.Ntp
	lastStatus TimeStatus
	mutex      sync.Mutex
}

func NewNtpTimeSource(server string) *NtpTimeSource {
	ntpClient := ntp.NewNtp(server)

	status := TimeStatus{
		Source:        "NTP",
		State:         StateSyncing,
		ReferenceTime: time.Now(),
		LocalTime:     time.Now(),
		Offset:        0,
	}

	return &NtpTimeSource{
		ntpClient:  ntpClient,
		lastStatus: status,
	}
}

func (s *NtpTimeSource) Name() string {
	return "NTP"
}

func (s *NtpTimeSource) TimeStatus() TimeStatus {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.lastStatus.LocalTime = time.Now()
	s.lastStatus.ReferenceTime = time.Now().Add(-s.ntpClient.Offset())
	s.lastStatus.Offset = s.ntpClient.Offset()

	return s.lastStatus
}

func (s *NtpTimeSource) Refresh() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.lastStatus.State = StateSyncing

	err := s.ntpClient.Refresh()
	if err == nil {
		s.lastStatus.State = StateLocked
		return nil
	} else {
		s.lastStatus.State = StateError
		return err
	}
}
