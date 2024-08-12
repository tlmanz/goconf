package goconf

import (
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/tlmanz/hush"
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
	data := hush.NewHushType(pr)

	table.SetHeader([]string{"Config", "Value"})
	table.AppendBulk(data.Hush(""))
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Render()
}
