package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
)

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       install",
		errmsg, os.Args[0])
	os.Exit(2)
}

func main() {
	const svcName = "SomeRandomService2.0"
	elog, err := eventlog.Open(svcName)
	elog.Info(7, "check here")
	inService, err := svc.IsWindowsService()
	if err != nil {
		log.Fatalf("failed to determine if we are running in service: %v", err)
		// errorLog.Println(fmt.Sprintf("failed to determine if we are running in service: %v", err))
	}

	if inService {
		runService(svcName)
		return
	}

	if len(os.Args) < 2 {
		usage("no command specified")
	}

	cmd := strings.ToLower(os.Args[1])
	elog.Info(5, cmd)

	switch cmd {
	case "debug":
		runService(svcName)
		return
	case "install":
		err = installService_2(svcName, "SomeRandomService2.0")
	case "remove":
		err = removeService(svcName)
	case "start":
		err = startService(svcName)
	case "stop":
		elog.Info(8, "Stop : 54")
		err = controlService(svcName, svc.Stop, svc.Stopped)
		_, admin, err := testSome()
		if err != nil {
			elog.Error(8, "admin check error")
		}
		if admin == true {
			err = controlService(svcName, svc.Stop, svc.Stopped)
		} else {
			elog.Error(5, "Access Denied")
		}

	default:
		usage(fmt.Sprintf("invalid command %s", cmd))
	}
	if err != nil {
		elog.Error(2, fmt.Sprintf("failed to %s %s: %v", cmd, svcName, err))
		// userName, err := getLoggedUser()
		// if err != nil {
		// 	elog.Error(3, err.Error())
		// }

		// elog.Info(3, fmt.Sprintf(
		// 	"This User Stopped Service: %s\nAccess domain: %s\nIs Local Admin: %v ",
		// 	userName[0].Username,
		// 	userName[0].Domain,
		// 	userName[0].LocalAdmin,
		// ))
		// log.Fatalf("failed to %s %s: %v", cmd, svcName, err)
	}
	return
}

func testSome() (bool, bool, error) {
	var sid *windows.SID

	// Although this looks scary, it is directly copied from the
	// official windows documentation. The Go API for this is a
	// direct wrap around the official C++ API.
	// See https://docs.microsoft.com/en-us/windows/desktop/api/securitybaseapi/nf-securitybaseapi-checktokenmembership
	err := windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)
	if err != nil {
		log.Fatalf("SID Error: %s", err)
		return false, false, err
	}

	// This appears to cast a null pointer so I'm not sure why this
	// works, but this guy says it does and it Works for Me™:
	// https://github.com/golang/go/issues/28804#issuecomment-438838144
	token := windows.Token(0)

	member, err := token.IsMember(sid)
	if err != nil {
		log.Fatalf("Token Membership Error: %s", err)
		return false, false, err
	}

	// Also note that an admin is _not_ necessarily considered
	// elevated.
	// For elevation see https://github.com/mozey/run-as-admin
	return token.IsElevated(), member, nil
	// fmt.Println("Elevated?", token.IsElevated())

	// fmt.Println("Admin?", member)
}
