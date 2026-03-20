package send

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type batchResponse struct {
	Success   bool            `json:"success"`
	Responses json.RawMessage `json:"responses"`
}

func NewCmdBatchTransactional(f *cmdutil.Factory) *cobra.Command {
	return newBatchCmd(f, "batch-transactional", "Send a batch of transactional emails", client.BaseTransactional)
}

func newBatchCmd(f *cmdutil.Factory, name, short string, base client.BaseURL) *cobra.Command {
	var file string

	cmd := &cobra.Command{
		Use:   name,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			data, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("read file %s: %w", file, err)
			}

			var body interface{}
			if err := json.Unmarshal(data, &body); err != nil {
				return fmt.Errorf("parse JSON from %s: %w", file, err)
			}

			var resp batchResponse
			if err := c.Post(context.Background(), base, "/api/batch", body, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp, []output.Column{
				{Header: "SUCCESS", Field: "success"},
				{Header: "RESPONSES", Field: "responses"},
			})
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "Path to JSON file containing the batch request body")
	_ = cmd.MarkFlagRequired("file")

	return cmd
}
