package main

import (
	"bytes"
	"flag"
	"fmt"
	"gopkg.in/gomail.v2"
	"io"
	"os/exec"
	"strings"
)

var (
	service      = flag.String("service", "elasticsearch", "Name of Systemd service")
	smtpHost     = "smtp.gmail.com"
	smtpPort     = 587
	smtpUser     = "xxxx@xx.xx"
	smtpPassword = "xxxxxxx"
)

func main() {
	flag.Parse()

	var (
		b2  bytes.Buffer
		msg strings.Builder
	)

	status, err := exec.Command("/bin/systemctl", "status", *service).Output()
	if err != nil {
		fmt.Println(err)
	}

	r, w := io.Pipe()
	mainCmd := exec.Command("/bin/journalctl", "-u", *service)
	sarg := exec.Command("/bin/tail", "-5")

	mainCmd.Stdout = w
	sarg.Stdin = r
	sarg.Stdout = &b2

	mainCmd.Start()
	sarg.Start()
	mainCmd.Wait()

	w.Close()

	sarg.Wait()
	if _, err = io.Copy(&msg, &b2); err != nil {
		fmt.Println(err)
		return
	}
	mail(string(status), msg.String())
}

func mail(status, message string) bool {
	m := gomail.NewMessage()

	m.SetBody("text/html",
		fmt.Sprintf("<i>Desc</i>: <pre>%s</pre><p><p>&nbsp;</p><i>Last 5 Logs</i>: %s", status, message))
	m.SetHeader("From", "test@kk.pt")
	m.SetAddressHeader("To", smtpUser, "Service Alerter")
	m.SetHeader("Subject", fmt.Sprintf("Service [%s] failed", *service))

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)
	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
