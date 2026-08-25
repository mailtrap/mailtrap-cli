package companyinfo

import (
	"github.com/mailtrap/mailtrap-cli/internal/cmdutil"
	"github.com/mailtrap/mailtrap-cli/internal/output"
	"github.com/spf13/cobra"
)

type CompanyInfo struct {
	Name              *string `json:"name,omitempty"`
	Address           *string `json:"address,omitempty"`
	City              *string `json:"city,omitempty"`
	Country           *string `json:"country,omitempty"`
	Phone             *string `json:"phone,omitempty"`
	ZipCode           *string `json:"zip_code,omitempty"`
	PrivacyPolicyURL  *string `json:"privacy_policy_url,omitempty"`
	TermsOfServiceURL *string `json:"terms_of_service_url,omitempty"`
	WebsiteURL        *string `json:"website_url,omitempty"`
	InfoLevel         *string `json:"info_level,omitempty"`
}

type companyInfoResponse struct {
	Data CompanyInfo `json:"data"`
}

var companyInfoColumns = []output.Column{
	{Header: "NAME", Field: "name"},
	{Header: "ADDRESS", Field: "address"},
	{Header: "CITY", Field: "city"},
	{Header: "COUNTRY", Field: "country"},
	{Header: "ZIP CODE", Field: "zip_code"},
	{Header: "PHONE", Field: "phone"},
	{Header: "WEBSITE", Field: "website_url"},
	{Header: "PRIVACY POLICY", Field: "privacy_policy_url"},
	{Header: "TERMS OF SERVICE", Field: "terms_of_service_url"},
	{Header: "INFO LEVEL", Field: "info_level"},
}

func companyInfoPath(domainID string) string {
	return "/api/domains/" + domainID + "/company_info"
}

func NewCmdCompanyInfo(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "company-info",
		Short: "Manage sending domain company info",
	}

	cmd.AddCommand(NewCmdGet(f))
	cmd.AddCommand(NewCmdCreate(f))
	cmd.AddCommand(NewCmdUpdate(f))

	return cmd
}
