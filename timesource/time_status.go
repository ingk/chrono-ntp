package timesource

import "time"

type TimeStatus struct {
	Source        string        // "NTP", "GPS", "DCF77", "Offline"
	State         SyncState     // LOCKED, SYNCING, ERROR
	ReferenceTime time.Time     // time reported by source
	LocalTime     time.Time     // system time at capture
	Offset        time.Duration // reference - local
}
