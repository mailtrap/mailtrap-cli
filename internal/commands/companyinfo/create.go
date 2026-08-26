package companyinfo

import (
	"context"

	"github.com/mailtrap/mailtrap-cli/internal/client"
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
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
		Use:   "create",
		Short: "Create company info for a sending domain",
		RunE: func(cmd *cobra.Command, args []string) error {
			required := []struct {
				flag  string
				value string
			}{
				{"domain-id", domainID},
				{"name", name},
				{"address", address},
				{"city", city},
				{"country", country},
				{"zip-code", zipCode},
				{"website-url", websiteURL},
			}
			for _, r := range required {
				if err := cmdutil.RequireFlag(r.flag, r.value); err != nil {
					return err
				}
			}

			c, err := f.NewClient()
			if err != nil {
				return err
			}

			companyInfoFields := map[string]interface{}{
				"name":        name,
				"address":     address,
				"city":        city,
				"country":     country,
				"zip_code":    zipCode,
				"website_url": websiteURL,
			}
			if cmd.Flags().Changed("phone") {
				companyInfoFields["phone"] = phone
			}
			if cmd.Flags().Changed("privacy-policy-url") {
				companyInfoFields["privacy_policy_url"] = privacyPolicyURL
			}
			if cmd.Flags().Changed("terms-of-service-url") {
				companyInfoFields["terms_of_service_url"] = termsOfServiceURL
			}
			if cmd.Flags().Changed("info-level") {
				companyInfoFields["info_level"] = infoLevel
			}

			body := map[string]interface{}{"company_info": companyInfoFields}

			var resp companyInfoResponse
			if err := c.Post(context.Background(), client.BaseGeneral, companyInfoPath(domainID), body, &resp); err != nil {
				return err
			}

			format := cmdutil.GetOutputFormat()
			return output.Print(f.IOStreams.Out, format, resp.Data, companyInfoColumns)
		},
	}

	cmd.Flags().StringVar(&domainID, "domain-id", "", "Sending domain ID (required)")
	cmd.Flags().StringVar(&name, "name", "", "Company or individual name (required)")
	cmd.Flags().StringVar(&address, "address", "", "Street address (required)")
	cmd.Flags().StringVar(&city, "city", "", "City (required)")
	cmd.Flags().StringVar(&country, "country", "", "Country (required)")
	cmd.Flags().StringVar(&zipCode, "zip-code", "", "ZIP or postal code (required)")
	cmd.Flags().StringVar(&websiteURL, "website-url", "", "Company website URL (required)")
	cmd.Flags().StringVar(&phone, "phone", "", "Phone number")
	cmd.Flags().StringVar(&privacyPolicyURL, "privacy-policy-url", "", "URL to the privacy policy page")
	cmd.Flags().StringVar(&termsOfServiceURL, "terms-of-service-url", "", "URL to the terms of service page")
	cmd.Flags().StringVar(&infoLevel, "info-level", "", "Whether the sender is a business or individual: business, individual")

	return cmd
}
