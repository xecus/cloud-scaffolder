package model

import (
	"fmt"
	"regexp"
	"strings"
)

type Vm struct {
	Hostname          string
	Image             string
	Memory            int
	Cpus              int
	Leader            bool
	NetworkInterfaces []NetworkInterface
}

type NetworkInterface struct {
	Type   string
	Option map[string]string
}

func (g NetworkInterface) Display() string {
	a := []string{
		fmt.Sprintf("\"%s\"", g.Type),
	}
	for k, v := range g.Option {
		if !check_regexp(`^[0-9]+$`, v) {
			v = fmt.Sprintf("\"%s\"", v)
		}
		a = append(a, fmt.Sprintf("%s: %s", k, v))
	}
	return strings.Join(a, ", ")
}

func check_regexp(reg, str string) bool {
	return regexp.MustCompile(reg).Match([]byte(str))
}
