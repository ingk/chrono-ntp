package timesource

import "time"

// Ensure OfflineTimeSource implements the TimeSource interface
var _ TimeSource = (*OfflineTimeSource)(nil)

type OfflineTimeSource struct {
}

func NewOfflineTimeSource() *OfflineTimeSource {
	return &OfflineTimeSource{}
}

func (s *OfflineTimeSource) Name() string {
	return "Offline"
}

func (s *OfflineTimeSource) TimeStatus() TimeStatus {
	now := time.Now()
	return TimeStatus{
		Source:        s.Name(),
		State:         StateLocked,
		ReferenceTime: now,
		LocalTime:     now,
		Offset:        0,
	}
}

func (s *OfflineTimeSource) Refresh() error {
	// No-op for offline time source
	return nil
}
