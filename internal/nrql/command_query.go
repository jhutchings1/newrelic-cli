package nrql

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/newrelic/newrelic-cli/internal/client"
	"github.com/newrelic/newrelic-cli/internal/output"
	"github.com/newrelic/newrelic-cli/internal/utils"
	"github.com/newrelic/newrelic-client-go/newrelic"
)

var (
	accountID int
)

const (
	gqlNrqlQuery = `query($query: Nrql!, $accountId: Int!) { actor { account(id: $accountId) { nrql(query: $query) { results } } } }`
)

var cmdQuery = &cobra.Command{
	Use:   "query",
	Short: "Execute a NRQL query to New Relic",
	Long: `Execute a NRQL query to New Relic

The query command accepts a single argument in the form of a NRQL query as a string.
This command requires the --accountId <int> flag, which specifies the account to
issue the query against.
`,
	Example: `newrelic nrql query 'SELECT count(*) FROM Transaction TIMESERIES'`,
	Args: func(cmd *cobra.Command, args []string) error {
		argsCount := len(args)

		if argsCount < 1 {
			return errors.New("missing NRQL query argument")
		}

		if argsCount > 1 {
			return errors.New("command expects only 1 argument")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		client.WithClient(func(nrClient *newrelic.NewRelic) {
			queryVars := map[string]interface{}{
				"accountId": accountID,
				"query":     args[0],
			}

			result, err := nrClient.NerdGraph.Query(gqlNrqlQuery, queryVars)
			if err != nil {
				log.Fatal(err)
			}

			utils.LogIfFatal(output.Print(result))
		})
	},
}

func init() {
	Command.AddCommand(cmdQuery)
	cmdQuery.Flags().IntVarP(&accountID, "accountId", "a", 0, "the New Relic account ID where you want to query")
	utils.LogIfError(cmdQuery.MarkFlagRequired("accountId"))
}
