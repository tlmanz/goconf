package goconf

import (
	"context"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/tlmanz/hush"
	"github.com/tryfix/log"
)

type Configer interface {
	Register() error
}

type Validater interface {
	Validate() error
}

type Printer interface {
	Print() interface{}
}

func Load(configs ...Configer) error {
	for _, c := range configs {
		err := c.Register()
		if err != nil {
			return err
		}

		v, ok := c.(Validater)
		if ok {
			err = v.Validate()
			if err != nil {
				return err
			}
		}

		p, ok := c.(Printer)
		if ok {
			printTable(p)
		}
	}
	return nil
}

func printTable(p Printer) {
	table := tablewriter.NewWriter(os.Stdout)

	pr := p.Print()

	// Create a HushType instance
	husher, err := hush.NewHush(pr)
	if err != nil {
		log.Fatal(err)
	}

	result, err := husher.Hush(context.Background(), "")
	if err != nil {
		log.Fatal(err)
	}

	table.SetHeader([]string{"Config", "Value"})
	table.AppendBulk(result)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Render()
}
