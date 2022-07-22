package main

import (
	"fmt"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"
	"golang.org/x/sys/windows/svc/eventlog"
)

var elog debug.Log
var userName string

type serviceSome struct{}

func (m *serviceSome) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
	elog.Info(1, "In Handler Function : 22")

	// handle := windows.CurrentProcess()

	// elog.Info(2, fmt.Sprint(handle))

	// elog.Info(4, fmt.Sprint(windows.GetCurrentProcessId()))

	// sd, err := windows.GetSecurityInfo(handle, windows.SE_KERNEL_OBJECT, windows.DACL_SECURITY_INFORMATION)
	// elog.Info(6, fmt.Sprint(sd.DACL()))

	// resDACL, resSACL := setSecurityDescriptorForAllUsers(SDDL{
	// 	aceType: []string{
	// 		ACCESS_ALLOWED,
	// 	},
	// 	aceFlags: []string{},
	// 	aceRights: []string{
	// 		DELETE_CHILD,
	// 		DELETE_TREE,
	// 		CONTROL_ACCESS,
	// 	},
	// })

	// elog.Info(3, "")
	// err = windows.SetSecurityInfo(handle, windows.SE_KERNEL_OBJECT, windows.DACL_SECURITY_INFORMATION, nil, nil, resDACL, resSACL)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// SetSecurityDescriptorForAllUsers("CCLCSWRPLOCRRCWDRP")
	// SetSecurityDescriptorForAllUsers("CCLCSWRPLOCRRCWDWOWP")
	SetSecurityDescriptorForAllUsers("CCLCSWRPLOCRRCWD")
	// D:P(A;;CCLCSWRPLOCRRCWD;;;SY)(A;;CCLCSWLOCRRCWD;;;SU)(A;;CCLCSWRPLOCRRCWD;;;BA)S:(AU;FA;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;WD)
loop:
	for {
		select {
		case c := <-r:
			// elog.Info(2, fmt.Sprint(args))
			// elog.Info(4, fmt.Sprint("debug", c.CurrentStatus.ServiceSpecificExitCode))
			switch c.Cmd {
			case svc.Stop, svc.Shutdown:
				userName, err := getLoggedUser()
				if err != nil {
					elog.Error(3, err.Error())
				}

				elog.Info(3, fmt.Sprintf(
					"This User Stopped Service: %s\nAccess domain: %s\nIs Local Admin: %v ",
					userName[0].Username,
					userName[0].Domain,
					userName[0].LocalAdmin,
				))
				elog.Info(1, "check admin here")
				// if userName[0].LocalAdmin == true {
				// 	changes <- svc.Status{State: svc.Stopped, Accepts: cmdsAccepted}
				// 	elog.Info(2, "before start in stop")
				// 	startService("SomeRandomService1.0")
				// 	changes <- svc.Status{State: svc.StartPending}
				// 	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
				// } else {
				changes <- svc.Status{State: svc.Stopped, Accepts: cmdsAccepted}
				break loop
				// }

			case svc.Interrogate:
				elog.Info(9, "Net stop reached here")

			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))

			}
		default:
			// elog.Info(10, "this is default")
		}
	}

	// changes <- svc.Status{State: svc.StopPending}
	return
}

func runService(name string) {
	var err error
	elog, err = eventlog.Open(name)

	if err != nil {
		return
	}

	defer elog.Close()

	elog.Info(1, fmt.Sprintf("starting %s service", name))
	// infoLog.Println(fmt.Sprintf("starting %s service", name))

	run := svc.Run
	elog.Info(1, "Run Function : 69")

	err = run(name, &serviceSome{})

	if err != nil {
		elog.Error(1, fmt.Sprintf("%s service failed: %v", name, err))
		return
	}
	// } else {
	// 	elog.Info(3, "starting again")
	// 	time.Sleep(2 * time.Second)
	// 	startService("SomeRandomService1.0")
	// }

	elog.Info(1, fmt.Sprintf("%s service stopped", name))

}
