package renderrer

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/metrue/fx/types"
	"github.com/olekukonko/tablewriter"
)

const formatJSON = "json"

//nolint:unused,varcheck,deadcode
const formatTable = "table"

// Render render output with given format
func Render(services []types.Service, format string) error {
	if strings.ToLower(format) == formatJSON {
		return toJSON(services)
	}
	return toTable(services)
}

func toTable(services []types.Service) error {
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
	return nil
}

func toJSON(services []types.Service) error {
	output, err := json.Marshal(services)
	if err != nil {
		return err
	}
	if _, err := fmt.Print(string(output)); err != nil {
		return err
	}
	return nil
}
