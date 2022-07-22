package main

import (
	"fmt"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

func startService(name string) error {

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("Could not Access Service: %v", err)
	}
	defer s.Close()

	err = s.Start()
	if err != nil {
		return fmt.Errorf("could not start service: %v", err)
	}

	// runService(name, true)

	return nil
}

func controlService(name string, c svc.Cmd, to svc.State) error {
	// SetSecurityDescriptorForAllUsers("CCDCLCSWRPWPDTLOCRSDRCWDWO")
	elog, err := eventlog.Open(name)

	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer m.Disconnect()

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("Could not Access Service: %v", err)
	}
	defer s.Close()
	elog.Info(9, "Manager Process done : 51")

	status, err := s.Control(c)
	if err != nil {
		return fmt.Errorf("could not send control = %d: %v", c, err)
	}
	timeout := time.Now().Add(10 * time.Second)
	for status.State != to {
		if timeout.Before(time.Now()) {
			return fmt.Errorf("timeout waiting for service to go to state=%d", to)
		}
		time.Sleep(300 * time.Millisecond)
		status, err = s.Query()
		if err != nil {
			return fmt.Errorf("could not retrieve service status: %v", err)
		}
	}
	// SetSecurityDescriptorForAllUsers("CCLCSWRPLOCRRCWDRP")
	return nil
}
