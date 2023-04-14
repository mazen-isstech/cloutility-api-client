package cmd

import (
	"context"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/safespring/cloutility-api-client/cloutapi"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// nodeCmd represents the node command
var deleteConsumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "Delete existing consumer and associated backup node",
	Long: `
The command delete consumer deletes an existing consumer and associated 
backup node.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		deleteConsumer()
	},
}

var consumerID int

func deleteConsumer() {
	var selectedConsumer cloutapi.Consumer
	client, err := cloutapi.Init(
		context.Background(),
		viper.GetString("client_id"),
		viper.GetString("client_origin"),
		viper.GetString("username"),
		viper.GetString("password"),
		viper.GetString("url"),
	)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	twriter := new(tabwriter.Writer)
	twriter.Init(os.Stdout, 8, 8, 1, '\t', 0)
	defer twriter.Flush()

	if bunitId == 0 {
		user, err := client.GetUser()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		bunitId = user.UserBUnit.ID
	}

	fmt.Fprintf(twriter, "%s\t%s\t%s\t%s\t%s\n", "Consumer ID", "Business-unit ID", "Name", "Status", "Url")
	fmt.Fprintf(twriter, "%s\t%s\t%s\t%s\t%s\n", "-----------", "----------------", "----", "------", "---")

	consumers, _ := client.GetConsumers(bunitId)
	for _, consumer := range consumers {
		if consumer.ID == consumerID {
			selectedConsumer = consumer
		}
	}

	if selectedConsumer.ID == 0 {
		fmt.Fprintf(twriter, "%v\tUser Default: %v\t%s\t%s\t%s\n", consumerID, bunitId, "NOT FOUND", "NOT FOUND", "NOT FOUND")
		return
	}

	if err := client.DeleteConsumer(bunitId, consumerID); err != nil {
		fmt.Fprintf(twriter, "%v\t%v\t%s\t%s\t%s\n", consumerID, bunitId, selectedConsumer.Name, err, selectedConsumer.Href)
		return
	}
	fmt.Fprintf(twriter, "%v\t%v\t%s\t%s\t%s\n", consumerID, bunitId, selectedConsumer.Name, "DELETED", selectedConsumer.Href)
}

func init() {
	deleteCmd.AddCommand(deleteConsumerCmd)

	deleteConsumerCmd.Flags().IntVar(&consumerID, "id", 0, "ID of consumption-unit to delete")
	deleteConsumerCmd.Flags().IntVar(&bunitId, "bunit-id", 0, "ID of business-unit in which the consumption-unit resides")

	// Mark --id as required
	err := deleteConsumerCmd.MarkFlagRequired("id")
	if err != nil {
		fmt.Println("error marking id flag as required: %w", err)
		os.Exit(1)
	}
}
