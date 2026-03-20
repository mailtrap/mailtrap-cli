package contacts

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/config"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var contactID string
	var email string
	var firstName string
	var lastName string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a contact",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", contactID); err != nil {
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

			path := cmdutil.AccountPath("contacts", fmt.Sprintf("%s", contactID))

			contactData := map[string]interface{}{}
			if email != "" {
				contactData["email"] = email
			}
			if firstName != "" {
				contactData["first_name"] = firstName
			}
			if lastName != "" {
				contactData["last_name"] = lastName
			}

			body := map[string]interface{}{
				"contact": contactData,
			}

			var contact Contact
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &contact); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, contact, contactColumns)
		},
	}

	cmd.Flags().StringVar(&contactID, "id", "", "Contact ID")
	cmd.Flags().StringVar(&email, "email", "", "Contact email")
	cmd.Flags().StringVar(&firstName, "first-name", "", "Contact first name")
	cmd.Flags().StringVar(&lastName, "last-name", "", "Contact last name")

	return cmd
}
