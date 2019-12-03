package render

import (
	"fmt"
	"os"

	"github.com/metrue/fx/types"
	"github.com/olekukonko/tablewriter"
)

// Table output services as table format
func Table(services []types.Service) {
	data := [][]string{}
	for _, s := range services {
		col := []string{
			s.ID,
			s.Name,
			fmt.Sprintf("%s:%d", s.Host, +s.Port),
		}
		data = append(data, col)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"ID", "Name", "Endpoint"})
	table.AppendBulk(data)
	table.Render()
}
