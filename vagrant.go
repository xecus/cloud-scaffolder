package cloudscaffolder

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"text/template"
)

type Vagrant struct {
	Vms []Vm
}

func GenerateVagrantModel() Vagrant {
	member := Vagrant{
		[]Vm{
			Vm{
				Hostname: "router",
				Image: VmImage{
					Name:      "VyOS",
					ImageName: "higebu/vyos",
					Version:   "1.1.3",
				},
				MemorySize: 2048,
				NumOfCpus:  2,
				Leader:     true,
				NetworkInterfaces: []NetworkInterface{
					NetworkInterface{
						Name: "vlan1",
						Type: "private_network",
						NetworkInterfaceOptions: []NetworkInterfaceOption{
							NetworkInterfaceOption{
								Key:   "ip",
								Value: "192.168.33.9",
							},
						},
					},
					NetworkInterface{
						Name: "vlan2",
						Type: "private_network",
						NetworkInterfaceOptions: []NetworkInterfaceOption{
							NetworkInterfaceOption{
								Key:   "ip",
								Value: "192.168.34.1",
							},
							NetworkInterfaceOption{
								Key:   "virtualbox__intnet",
								Value: "hogenet",
							},
						},
					},
					NetworkInterface{
						Name: "fp1",
						Type: "forwarded_port",
						NetworkInterfaceOptions: []NetworkInterfaceOption{
							NetworkInterfaceOption{
								Key:   "host",
								Value: "3022",
							},
							NetworkInterfaceOption{
								Key:   "guest",
								Value: "22",
							},
						},
					},
				},
			},
			Vm{
				Hostname: "vm_1",
				Image: VmImage{
					Name:      "ubuntu",
					ImageName: "ubuntu/trusty64",
					Version:   "14.04",
				},
				MemorySize: 2048,
				NumOfCpus:  2,
				Leader:     false,
				NetworkInterfaces: []NetworkInterface{
					NetworkInterface{
						Name: "vlan1",
						Type: "private_network",
						NetworkInterfaceOptions: []NetworkInterfaceOption{
							NetworkInterfaceOption{
								Key:   "ip",
								Value: "192.168.33.10",
							},
						},
					},
					NetworkInterface{
						Name: "vlan2",
						Type: "private_network",
						NetworkInterfaceOptions: []NetworkInterfaceOption{
							NetworkInterfaceOption{
								Key:   "ip",
								Value: "192.168.34.0",
							},
							NetworkInterfaceOption{
								Key:   "virtualbox__intnet",
								Value: "hogenet",
							},
							NetworkInterfaceOption{
								Key:   "type",
								Value: "dhcp",
							},
						},
					},
					NetworkInterface{
						Name: "fp1",
						Type: "forwarded_port",
						NetworkInterfaceOptions: []NetworkInterfaceOption{
							NetworkInterfaceOption{
								Key:   "host",
								Value: "3023",
							},
							NetworkInterfaceOption{
								Key:   "guest",
								Value: "22",
							},
						},
					},
				},
			},
		},
	}

	return member
}

func GenerateVagrantFile() {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal("No caller information")
		return
	}
	filepath := path.Join(path.Dir(filename), "template/Vagrantfile.tpl")

	//log.Println("filepath=" + filepath)
	tpl := template.Must(template.ParseFiles(filepath))
	member := GenerateVagrantModel()
	//fmt.Println(member)

	file, err := os.Create("vagrant_area/Vagrantfile")
	if err != nil {
		log.Fatal(err)
	}
	if err := tpl.Execute(file, member); err != nil {
		log.Fatal(err)
	}
}

func PrepareVagrantControl() {

	if err := os.Mkdir("vagrant_area", 0777); err != nil {
		fmt.Println(err)
	}

}

func CtrlVagrant(c string, p []string) {
	cmd := exec.Command(c, p...)
	cmd.Dir = "./vagrant_area"
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	cmd.Start()

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			s := scanner.Text()
			fmt.Println("STDOUT PID= " + fmt.Sprintf("%d", cmd.Process.Pid) + " RECV=[" + s + "]")
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			s := scanner.Text()
			fmt.Println("STDERR PID= " + fmt.Sprintf("%d", cmd.Process.Pid) + " RECV=[" + s + "]")
		}
	}()

	cmd.Wait()
}
