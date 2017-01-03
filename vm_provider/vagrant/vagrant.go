package vagrant

import (
	"../../model"
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"text/template"
)

type Vagrant struct {
	Vms []model.Vm
}

func GenerateVagrantFile() {
	tpl := template.Must(template.ParseFiles("template/Vagrantfile.tpl"))
	member := Vagrant{
		[]model.Vm{
			model.Vm{
				"router",
				"higebu/vyos",
				2048,
				2,
				true,
				[]model.NetworkInterface{
					model.NetworkInterface{
						"private_network",
						map[string]string{
							"ip": "192.168.33.9",
						},
					},
					model.NetworkInterface{
						"private_network",
						map[string]string{
							"ip":                 "192.168.34.1",
							"virtualbox__intnet": "hogenet",
						},
					},
					model.NetworkInterface{
						"forwarded_port",
						map[string]string{
							"host":  "3022",
							"guest": "22",
						},
					},
				},
			},
			model.Vm{
				"OITEC",
				"ubuntu/trusty64",
				2048,
				2,
				false,
				[]model.NetworkInterface{
					model.NetworkInterface{
						"private_network",
						map[string]string{
							"ip": "192.168.33.10",
						},
					},
					model.NetworkInterface{
						"private_network",
						map[string]string{
							"ip":                 "192.168.34.0",
							"virtualbox__intnet": "hogenet",
							"type":               "dhcp",
						},
					},
					model.NetworkInterface{
						"forwarded_port",
						map[string]string{
							"host":  "3023",
							"guest": "22",
						},
					},
				},
			},
		},
	}

	file, err := os.Create("vagrant_area/Vagrantfile")
	if err != nil {
		log.Fatal(err)
	}
	if err := tpl.Execute(file, member); err != nil {
		log.Fatal(err)
	}
}

func ValidationBeforeControl() {

}

func PrepareVagrantControl() {

	if err := os.Mkdir("vagrant_area", 0777); err != nil {
		fmt.Println(err)
	}

}

func CtrlVagrant(c string, p []string) {
	cmd := exec.Command(c, p...)
	//cmd.Dir = "../VagrantWorkSuite"
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		s := scanner.Text()
		fmt.Println("PID= " + fmt.Sprintf("%d", cmd.Process.Pid) + " RECV=[" + s + "]")
	}
	cmd.Wait()
}
