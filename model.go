package cloudscaffolder

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/satori/go.uuid"
	"log"
	"os"
	"regexp"
	"strings"
)

// For Resource Management
type ResourceControl struct {
	Uuid string `gorm:"size:128"`
}

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
	VmID                    int    `gorm:"index"`
	Name                    string `gorm:"size:128;not null" json:"name"`
	Type                    string `gorm:"size:128;not null" json:"type"`
	NetworkInterfaceOptions []NetworkInterfaceOption
}

type Vm struct {
	gorm.Model
	ResourceControl
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

func (v Vm) DeleteAlltVm(i Impl) {
	a := []VmImage{}
	i.DB.Find(&a)
	for _, k := range a {
		if err := i.DB.Delete(&k).Error; err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	b := []NetworkInterfaceOption{}
	i.DB.Find(&b)
	for _, k := range b {
		if err := i.DB.Delete(&k).Error; err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	c := []NetworkInterface{}
	i.DB.Find(&c)
	for _, k := range c {
		if err := i.DB.Delete(&k).Error; err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

	d := []Vm{}
	i.DB.Find(&d)
	for _, k := range d {
		if err := i.DB.Delete(&k).Error; err != nil {
			log.Fatalf("Error: %v", err)
		}
	}

}

func (v Vm) CreateVm(i Impl) Vm {
	u1 := uuid.NewV4()
	fmt.Printf("UUIDv4: %s\n", u1)
	if err := i.DB.Create(&v).Error; err != nil {
		log.Fatalf("Error: %v", err)
	}
	return v
}

func GetVm(i Impl) Vm {
	vm := Vm{}
	if err := i.DB.Last(&vm).Error; err != nil {
		log.Fatalf("Error: %v", err)
	}
	i.DB.Model(&vm).Related(&vm.Image)
	i.DB.Model(&vm).Related(&vm.NetworkInterfaces)
	for t, networkInterface := range vm.NetworkInterfaces {
		i.DB.Model(&networkInterface).Related(&vm.NetworkInterfaces[t].NetworkInterfaceOptions)
	}
	return vm
}

func (v Vm) UpdateVm(i Impl) {
	if err := i.DB.Save(&v).Error; err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func (v Vm) ShowVm() {

	log.Println("---------------[Show]----------------")
	log.Printf("Hostname=[%s]", v.Hostname)
	log.Printf("MemorySize=[%d]", v.MemorySize)
	log.Printf("NumOfCpus=[%d]", v.NumOfCpus)
	log.Printf("Leader=[%t]", v.Leader)
	log.Printf("Image.Name=[%s]", v.Image.Name)
	log.Printf("Image.ImageName=[%s]", v.Image.ImageName)
	log.Printf("Image.Version=[%s]", v.Image.Version)

	for i, networkInterface := range v.NetworkInterfaces {
		log.Printf("\tNetworkInterface [%d] %s", i, networkInterface.Name)
		for j, opt := range networkInterface.NetworkInterfaceOptions {
			log.Printf("\t\t Option[%d] %s = %s", j, opt.Key, opt.Value)
		}
	}

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
	//i.DB.LogMode(true)
}

func (i *Impl) InitSchema() {
	i.DB.AutoMigrate(&FirewallRule{})
	i.DB.AutoMigrate(&VmImage{})
	i.DB.AutoMigrate(&NetworkInterfaceOption{})
	i.DB.AutoMigrate(&NetworkInterface{})
	i.DB.AutoMigrate(&Vm{})
}
