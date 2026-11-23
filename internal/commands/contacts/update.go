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
	var listIDsIncluded []int
	var listIDsExcluded []int
	var unsubscribed bool
	var subscriptionTopics string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a contact",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("id", contactID); err != nil {
				return err
			}

			// The API finds the contact by its ID, so the email is sent only to change it.
			contactData := map[string]interface{}{}
			if cmd.Flags().Changed("email") {
				contactData["email"] = email
			}
			if cmd.Flags().Changed("subscription-topics") {
				parsed, err := parseSubscriptionTopics(subscriptionTopics)
				if err != nil {
					return err
				}
				contactData["subscription_topics"] = parsed
			}

			fields := map[string]interface{}{}
			if firstName != "" {
				fields["first_name"] = firstName
			}
			if lastName != "" {
				fields["last_name"] = lastName
			}
			if len(fields) > 0 {
				contactData["fields"] = fields
			}

			if len(listIDsIncluded) > 0 {
				contactData["list_ids_included"] = listIDsIncluded
			}
			if len(listIDsExcluded) > 0 {
				contactData["list_ids_excluded"] = listIDsExcluded
			}
			if cmd.Flags().Changed("unsubscribed") {
				contactData["unsubscribed"] = unsubscribed
			}
			if len(contactData) == 0 {
				return fmt.Errorf("at least one attribute flag is required")
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

			body := map[string]interface{}{
				"contact": contactData,
			}

			var resp contactResponse
			if err := c.Patch(context.Background(), client.BaseGeneral, path, body, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, contactColumns)
		},
	}

	cmd.Flags().StringVar(&contactID, "id", "", "Contact ID (required)")
	cmd.Flags().StringVar(&email, "email", "", "Contact email")
	cmd.Flags().StringVar(&firstName, "first-name", "", "Contact first name")
	cmd.Flags().StringVar(&lastName, "last-name", "", "Contact last name")
	cmd.Flags().IntSliceVar(&listIDsIncluded, "list-ids-included", nil, "List IDs to subscribe to")
	cmd.Flags().IntSliceVar(&listIDsExcluded, "list-ids-excluded", nil, "List IDs to unsubscribe from")
	cmd.Flags().BoolVar(&unsubscribed, "unsubscribed", false, "Set contact as unsubscribed")
	cmd.Flags().StringVar(&subscriptionTopics, "subscription-topics", "", subscriptionTopicsUsage)

	return cmd
}
