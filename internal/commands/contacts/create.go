package contacts

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var email string
	var firstName string
	var lastName string
	var listIDs []int

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a contact",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("email", email); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			_, err = config.RequireAccountID()
			if err != nil {
				return err
			}

			path := cmdutil.AccountPath("contacts")

			contactData := map[string]interface{}{
				"email": email,
			}
			if firstName != "" {
				contactData["first_name"] = firstName
			}
			if lastName != "" {
				contactData["last_name"] = lastName
			}
			if len(listIDs) > 0 {
				contactData["list_ids"] = listIDs
			}

			body := map[string]interface{}{
				"contact": contactData,
			}

			var contact Contact
			if err := c.Post(context.Background(), client.BaseGeneral, path, body, &contact); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, contact, contactColumns)
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "Contact email (required)")
	cmd.Flags().StringVar(&firstName, "first-name", "", "Contact first name")
	cmd.Flags().StringVar(&lastName, "last-name", "", "Contact last name")
	cmd.Flags().IntSliceVar(&listIDs, "list-ids", nil, "List IDs to add the contact to")

	return cmd
}
