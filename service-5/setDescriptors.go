package main

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/svc/eventlog"
)

func SetSecurityDescriptorForAllUsers(sdd string) {
	var users_sids []string
	var BA string
	var securityDescriptor = "D:"
	cmdOutput := &bytes.Buffer{}
	elog, err := eventlog.Open("SomeRandomService1.0")
	elog.Info(11, "Entered function")
	if err != nil {
		os.Exit(0)
	}
	cmd := exec.Command("wmic", "useraccount", "where", "localaccount=true", "get", "sid")
	cmd.Stdout = cmdOutput
	err = cmd.Run()
	if err != nil {
		elog.Error(1, "Cannot run the command wmic : "+err.Error())
	}
	output := cmdOutput.String()
	output_strings := strings.Fields(output)

	for i := range output_strings {
		if i == 0 {
			continue
		} else {
			if strings.HasSuffix(output_strings[i], "-500") {
				BA = output_strings[i]
				securityDescriptor += "(A;;" + sdd + ";;;" + BA + ")"
			} else {
				users_sids = append(users_sids, output_strings[i])
				securityDescriptor += "(A;;CCLCSWRPLOCRRCWD;;;" + output_strings[i] + ")"
			}
		}
	}
	securityDescriptor += "S:(AU;FA;CCDCLCSWRPWPDTLOCRSDRCWDWO;;;WD)"
	cmd = exec.Command("sc", "sdset", "SomeRandomService1.0", securityDescriptor)
	cmd.Stdout = cmdOutput
	err = cmd.Run()
	if err != nil {
		elog.Error(1, "Access is denied"+err.Error())
		elog.Error(0, securityDescriptor)
	}
	elog.Info(111, "Function Exited")
}
