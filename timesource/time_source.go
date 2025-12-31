package timesource

type TimeSource interface {
	Name() string
	TimeStatus() TimeStatus
	Refresh() error
}
