package main

import (
	"./vm_provider/vagrant"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/urfave/cli"
	"log"
	"os"
)

func scaffold(c *cli.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	fmt.Println("added task: ", c.Args().First())
	// Generate Vagrant File
	vagrant.PrepareVagrantControl()
	vagrant.GenerateVagrantFile()
}

func serve(c *cli.Context) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "CloudScaffolder"
	app.Usage = "This is test"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:    "scaffold",
			Aliases: []string{"sc"},
			Usage:   "scaffold cloud",
			Action: func(c *cli.Context) error {
				scaffold(c)
				return nil
			},
		},
		{
			Name:    "serve",
			Aliases: []string{"sv"},
			Usage:   "start server",
			Action: func(c *cli.Context) error {
				serve(c)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
