package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/commands/account_access"
	"github.com/mailtrap/mailtrap-cli/internal/commands/accounts"
	"github.com/mailtrap/mailtrap-cli/internal/commands/attachments"
	"github.com/mailtrap/mailtrap-cli/internal/commands/billing"
	"github.com/mailtrap/mailtrap-cli/internal/commands/contact_fields"
	"github.com/mailtrap/mailtrap-cli/internal/commands/contact_lists"
	"github.com/mailtrap/mailtrap-cli/internal/commands/configure"
	"github.com/mailtrap/mailtrap-cli/internal/commands/contacts"
	"github.com/mailtrap/mailtrap-cli/internal/commands/domains"
	email_logs "github.com/mailtrap/mailtrap-cli/internal/commands/email_logs"
	"github.com/mailtrap/mailtrap-cli/internal/commands/inboxes"
	"github.com/mailtrap/mailtrap-cli/internal/commands/messages"
	"github.com/mailtrap/mailtrap-cli/internal/commands/organizations"
	"github.com/mailtrap/mailtrap-cli/internal/commands/permissions"
	"github.com/mailtrap/mailtrap-cli/internal/commands/projects"
	"github.com/mailtrap/mailtrap-cli/internal/commands/sandbox_send"
	"github.com/mailtrap/mailtrap-cli/internal/commands/send"
	"github.com/mailtrap/mailtrap-cli/internal/commands/stats"
	"github.com/mailtrap/mailtrap-cli/internal/commands/suppressions"
	"github.com/mailtrap/mailtrap-cli/internal/commands/templates"
	"github.com/mailtrap/mailtrap-cli/internal/commands/tokens"
)

func NewRootCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mailtrap",
		Short: "CLI for the Mailtrap email platform",
		Long:  "A command-line interface for managing Mailtrap email sending, sandbox testing, contacts, and account settings.",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.PersistentFlags().String("api-token", "", "Mailtrap API token (env: MAILTRAP_API_TOKEN)")
	cmd.PersistentFlags().String("account-id", "", "Mailtrap account ID (env: MAILTRAP_ACCOUNT_ID)")
	cmd.PersistentFlags().StringP("output", "o", "table", "Output format: json, table, text")

	viper.BindPFlag("api-token", cmd.PersistentFlags().Lookup("api-token"))
	viper.BindPFlag("account-id", cmd.PersistentFlags().Lookup("account-id"))
	viper.BindPFlag("output", cmd.PersistentFlags().Lookup("output"))

	// Email Sending
	cmd.AddCommand(send.NewCmdSend(f))
	cmd.AddCommand(sandbox_send.NewCmdSandboxSend(f))

	// Sending Management
	cmd.AddCommand(domains.NewCmdDomains(f))
	cmd.AddCommand(suppressions.NewCmdSuppressions(f))
	cmd.AddCommand(stats.NewCmdStats(f))
	cmd.AddCommand(templates.NewCmdTemplates(f))
	cmd.AddCommand(email_logs.NewCmdEmailLogs(f))

	// Sandbox
	cmd.AddCommand(projects.NewCmdProjects(f))
	cmd.AddCommand(inboxes.NewCmdInboxes(f))
	cmd.AddCommand(messages.NewCmdMessages(f))
	cmd.AddCommand(attachments.NewCmdAttachments(f))

	// Contacts / Promotional
	cmd.AddCommand(contacts.NewCmdContacts(f))
	cmd.AddCommand(contact_lists.NewCmdContactLists(f))
	cmd.AddCommand(contact_fields.NewCmdContactFields(f))

	// Account Management
	cmd.AddCommand(accounts.NewCmdAccounts(f))
	cmd.AddCommand(account_access.NewCmdAccountAccess(f))
	cmd.AddCommand(permissions.NewCmdPermissions(f))
	cmd.AddCommand(tokens.NewCmdTokens(f))
	cmd.AddCommand(billing.NewCmdBilling(f))
	cmd.AddCommand(organizations.NewCmdOrganizations(f))

	// Configuration
	cmd.AddCommand(configure.NewCmdConfigure(f))

	// Shell completion
	cmd.AddCommand(newCompletionCmd())

	return cmd
}

func newCompletionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Args:  cobra.ExactArgs(1),
		ValidArgs: []string{"bash", "zsh", "fish", "powershell"},
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
			case "zsh":
				return cmd.Root().GenZshCompletion(cmd.OutOrStdout())
			case "fish":
				return cmd.Root().GenFishCompletion(cmd.OutOrStdout(), true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(cmd.OutOrStdout())
			}
			return nil
		},
	}
}

func Execute() error {
	f := cmdutil.NewFactory()
	f.Config() // load ~/.mailtrap.yaml into viper
	return NewRootCmd(f).Execute()
}
