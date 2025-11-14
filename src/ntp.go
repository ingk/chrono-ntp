package main

import (
	"time"

	"github.com/beevik/ntp"
)

type Ntp struct {
	server      string
	offset      time.Duration
	lastNtpTime time.Time
}

func NewNtp(server string) (*Ntp, error) {
	n := &Ntp{server: server}
	err := n.Refresh()
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (n *Ntp) Offset() time.Duration {
	return n.offset
}

func (n *Ntp) ServerTime() time.Time {
	return n.lastNtpTime
}

func (n *Ntp) Refresh() error {
	ntpTime, err := ntp.Time(n.server)
	if err != nil {
		return err
	}
	n.lastNtpTime = ntpTime
	n.offset = time.Since(ntpTime)
	return nil
}
