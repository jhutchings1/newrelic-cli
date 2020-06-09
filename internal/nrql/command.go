package nrql

import (
	"github.com/spf13/cobra"
)

// Command represents the nerdgraph command.
var Command = &cobra.Command{
	Use:   "nrql",
	Short: "Fetch data from New Relic using a NRQL (New Relic query language) query",
}
