package cloudscaffolder

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"os"
	"regexp"
	"strings"
)

// For Router VM Provider
type FirewallRule struct {
	gorm.Model
	Name     string `gorm:"size:128;not null" json:"name"`
	Protocol string `gorm:"size:128;not null" json:"protocol"`
	Srcrange string `gorm:"size:128;not null" json:"src_range"`
	Dstrange string `gorm:"size:128;not null" json:"dst_range"`
}

// For VM Provider

type VmImage struct {
	gorm.Model
	VmID      int    `gorm:"index"`
	Name      string `gorm:"size:128;not null" json:"name"`
	ImageName string `gorm:"size:128;not null" json:"image_name"`
	Version   string `gorm:"size:128;not null" json:"version"`
}

type NetworkInterfaceOption struct {
	gorm.Model
	NetworkInterfaceID int    `gorm:"index"`
	Name               string `gorm:"size:128" json:"name"`
	Key                string `gorm:"size:128;not null" json:"key"`
	Value              string `gorm:"size:128;not null" json:"value"`
}

type NetworkInterface struct {
	gorm.Model
	Name                    string `gorm:"size:128;not null" json:"name"`
	Type                    string `gorm:"size:128;not null" json:"type"`
	NetworkInterfaceOptions []NetworkInterfaceOption
}

type Vm struct {
	gorm.Model
	Hostname          string
	Image             VmImage
	MemorySize        int
	NumOfCpus         int
	Leader            bool
	NetworkInterfaces []NetworkInterface
}

func (g NetworkInterface) ExpandNetworkInterfaceOptions() string {
	a := []string{
		fmt.Sprintf("\"%s\"", g.Type),
	}
	for _, option := range g.NetworkInterfaceOptions {
		k := option.Key
		v := option.Value
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

type Impl struct {
	DB *gorm.DB
}

func (i *Impl) InitModelDb() {
	var err error
	hostname := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DATABASE")
	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, hostname, port, database)
	i.DB, err = gorm.Open("postgres", uri)
	if err != nil {
		log.Fatalf("Got error when connect database. %v", err)
	}
	//defer i.DB.Close()
	i.DB.LogMode(true)
}

func (i *Impl) InitSchema() {
	i.DB.AutoMigrate(&FirewallRule{})
	i.DB.AutoMigrate(&VmImage{})
	i.DB.AutoMigrate(&NetworkInterfaceOption{})
	i.DB.AutoMigrate(&NetworkInterface{})
	i.DB.AutoMigrate(&Vm{})
}
