package vagrant

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
)

func CtrlVagrant(c string, p []string) {
	cmd := exec.Command(c, p...)
	cmd.Dir = "../VagrantWorkSuite"
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		s := scanner.Text()
		fmt.Println("PID= " + fmt.Sprintf("%d", cmd.Process.Pid) + " RECV=[" + s + "]")
	}
	cmd.Wait()
}
