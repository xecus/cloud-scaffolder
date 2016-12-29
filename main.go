package main

import (
	"os"
	"fmt"
	"github.com/urfave/cli"
	"./job"
	"./vagrant"
)

func dis() {
	d := job.NewDispatcher()
	d.Start()
	for i := 0; i < 100; i++ {
		url := fmt.Sprintf("http://placehold.it/%dx%d", i, i)
		d.Add(url)
	}
	d.Wait()
}

func add(c *cli.Context) {

	fmt.Println("added task: ", c.Args().First())
	cmd := "vagrant"
	params := []string{"reload", "main", "--provision"}
	vagrant.CtrlVagrant(cmd, params)

}

func main() {
	app := cli.NewApp()
	app.Name = "CloudScaffolder"
	app.Usage = "This is test"
	app.Version = "1.0.0"
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a task to the list",
			Action: func(c *cli.Context) error {
				add(c)
				return nil
			},
		},
	}
	app.Run(os.Args)
}
