package timesource

type SyncState string

const (
	StateLocked  SyncState = "LOCKED"
	StateSyncing SyncState = "SYNCING"
	StateError   SyncState = "ERROR"
)
