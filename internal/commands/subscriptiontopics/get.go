package subscriptiontopics

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

// topicResponse keeps the topic as raw JSON so --output json prints every field.
type topicResponse struct {
	Data json.RawMessage `json:"data"`
}

// A deleted topic is still returned by get, with deleted_at set.
var topicDetailColumns = slices.Concat(topicColumns, []output.Column{{Header: "DELETED AT", Field: "deleted_at"}})

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	var id int

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a subscription topic by ID",
		Long:  "Get a subscription topic by ID. Deleted topics are returned too, with DELETED AT set.",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			var resp topicResponse
			if err := c.Get(context.Background(), client.BaseGeneral, fmt.Sprintf("%s/%d", basePath, id), nil, &resp); err != nil {
				return err
			}

			return output.Print(f.IOStreams.Out, cmdutil.GetOutputFormat(), resp.Data, topicDetailColumns)
		},
	}

	cmd.Flags().IntVar(&id, "id", 0, "Subscription topic ID (required)")
	_ = cmd.MarkFlagRequired("id")

	return cmd
}
