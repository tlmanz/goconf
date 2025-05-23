package goconf

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/tlmanz/hush"
	// "github.com/tryfix/log" // log.Fatal was removed, so this import is not needed anymore
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
	var allErrors []error
	for _, c := range configs {
		err := c.Register()
		if err != nil {
			allErrors = append(allErrors, err)
		}

		v, ok := c.(Validater)
		if ok {
			err = v.Validate()
			if err != nil {
				allErrors = append(allErrors, err)
			}
		}

		p, ok := c.(Printer)
		if ok {
			err = printTable(os.Stdout, p)
			if err != nil {
				allErrors = append(allErrors, err)
			}
		}
	}

	if len(allErrors) > 0 {
		var errorMessages []string
		for _, err := range allErrors {
			errorMessages = append(errorMessages, err.Error())
		}
		return fmt.Errorf("errors during configuration loading: %s", strings.Join(errorMessages, "; "))
	}

	return nil
}

func printTable(out io.Writer, p Printer) error {
	if out == nil {
		out = os.Stdout
	}
	table := tablewriter.NewWriter(out)

	pr := p.Print()

	// Create a HushType instance
	husher := hush.NewHush()

	result, err := husher.Hush(context.Background(), pr)
	if err != nil {
		return fmt.Errorf("error hushing data: %w", err)
	}

	table.SetHeader([]string{"Config", "Value"})
	table.AppendBulk(result)
	table.SetAlignment(tablewriter.ALIGN_LEFT)

	table.Render()
	return nil
}
