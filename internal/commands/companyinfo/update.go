package companyinfo

import (
	"context"
	"fmt"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var (
		domainID          string
		name              string
		address           string
		city              string
		country           string
		zipCode           string
		websiteURL        string
		phone             string
		privacyPolicyURL  string
		termsOfServiceURL string
		infoLevel         string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update company info for a sending domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := cmdutil.RequireFlag("domain-id", domainID); err != nil {
				return err
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			companyInfoFields := map[string]interface{}{}
			for _, field := range []struct {
				flag  string
				key   string
				value string
			}{
				{"name", "name", name},
				{"address", "address", address},
				{"city", "city", city},
				{"country", "country", country},
				{"zip-code", "zip_code", zipCode},
				{"website-url", "website_url", websiteURL},
				{"phone", "phone", phone},
				{"privacy-policy-url", "privacy_policy_url", privacyPolicyURL},
				{"terms-of-service-url", "terms_of_service_url", termsOfServiceURL},
				{"info-level", "info_level", infoLevel},
			} {
				if cmd.Flags().Changed(field.flag) {
					companyInfoFields[field.key] = field.value
				}
			}

			if len(companyInfoFields) == 0 {
				return fmt.Errorf("at least one attribute flag is required")
			}

			body := map[string]interface{}{"company_info": companyInfoFields}

			var resp companyInfoResponse
			if err := c.Patch(context.Background(), client.BaseGeneral, companyInfoPath(domainID), body, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, companyInfoColumns)
		},
	}

	cmd.Flags().StringVar(&domainID, "domain-id", "", "Sending domain ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "Company or individual name")
	cmd.Flags().StringVar(&address, "address", "", "Street address")
	cmd.Flags().StringVar(&city, "city", "", "City")
	cmd.Flags().StringVar(&country, "country", "", "Country")
	cmd.Flags().StringVar(&zipCode, "zip-code", "", "ZIP or postal code")
	cmd.Flags().StringVar(&websiteURL, "website-url", "", "Company website URL")
	cmd.Flags().StringVar(&phone, "phone", "", "Phone number")
	cmd.Flags().StringVar(&privacyPolicyURL, "privacy-policy-url", "", "URL to the privacy policy page")
	cmd.Flags().StringVar(&termsOfServiceURL, "terms-of-service-url", "", "URL to the terms of service page")
	cmd.Flags().StringVar(&infoLevel, "info-level", "", "Whether the sender is a business or individual: business, individual")

	return cmd
}
