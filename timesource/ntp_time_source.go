package timesource

import (
	"chrono-ntp/ntp"
	"time"
)

// Ensure NtpTimeSource implements the TimeSource interface
var _ TimeSource = (*NtpTimeSource)(nil)

type NtpTimeSource struct {
	ntpClient  *ntp.Ntp
	lastStatus TimeStatus
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
	s.lastStatus.LocalTime = time.Now()
	s.lastStatus.ReferenceTime = time.Now().Add(-s.ntpClient.Offset())
	s.lastStatus.Offset = s.ntpClient.Offset()
	return s.lastStatus
}

func (s *NtpTimeSource) Refresh() error {
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
