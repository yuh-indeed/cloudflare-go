package cloudflare

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// AccessIdentityProvider is the structure of the provider object.
type AccessIdentityProvider struct {
	ID         string                                  `json:"id,omitempty"`
	Name       string                                  `json:"name"`
	Type       string                                  `json:"type"`
	Config     AccessIdentityProviderConfiguration     `json:"config"`
	ScimConfig AccessIdentityProviderScimConfiguration `json:"scim_config"`
}

// AccessIdentityProviderConfiguration is the combined structure of *all*
// identity provider configuration fields. This is done to simplify the use of
// Access products and their relationship to each other.
//
// API reference: https://developers.cloudflare.com/access/configuring-identity-providers/
type AccessIdentityProviderConfiguration struct {
	APIToken           string   `json:"api_token,omitempty"`
	AppsDomain         string   `json:"apps_domain,omitempty"`
	Attributes         []string `json:"attributes,omitempty"`
	AuthURL            string   `json:"auth_url,omitempty"`
	CentrifyAccount    string   `json:"centrify_account,omitempty"`
	CentrifyAppID      string   `json:"centrify_app_id,omitempty"`
	CertsURL           string   `json:"certs_url,omitempty"`
	ClientID           string   `json:"client_id,omitempty"`
	ClientSecret       string   `json:"client_secret,omitempty"`
	Claims             []string `json:"claims,omitempty"`
	Scopes             []string `json:"scopes,omitempty"`
	DirectoryID        string   `json:"directory_id,omitempty"`
	EmailAttributeName string   `json:"email_attribute_name,omitempty"`
	IdpPublicCert      string   `json:"idp_public_cert,omitempty"`
	IssuerURL          string   `json:"issuer_url,omitempty"`
	OktaAccount        string   `json:"okta_account,omitempty"`
	OneloginAccount    string   `json:"onelogin_account,omitempty"`
	RedirectURL        string   `json:"redirect_url,omitempty"`
	SignRequest        bool     `json:"sign_request,omitempty"`
	SsoTargetURL       string   `json:"sso_target_url,omitempty"`
	SupportGroups      bool     `json:"support_groups,omitempty"`
	TokenURL           string   `json:"token_url,omitempty"`
	PKCEEnabled        *bool    `json:"pkce_enabled,omitempty"`
}

type AccessIdentityProviderScimConfiguration struct {
	Enabled                bool   `json:"enabled,omitempty"`
	Secret                 string `json:"secret,omitempty"`
	UserDeprovision        bool   `json:"user_deprovision,omitempty"`
	SeatDeprovision        bool   `json:"seat_deprovision,omitempty"`
	GroupMemberDeprovision bool   `json:"group_member_deprovision,omitempty"`
}

// AccessIdentityProvidersListResponse is the API response for multiple
// Access Identity Providers.
type AccessIdentityProvidersListResponse struct {
	Response
	Result     []AccessIdentityProvider `json:"result"`
	ResultInfo `json:"result_info"`
}

// AccessIdentityProviderResponse is the API response for a single
// Access Identity Provider.
type AccessIdentityProviderResponse struct {
	Response
	Result AccessIdentityProvider `json:"result"`
}

// ListAccessIdentityProviders returns all Access Identity Providers for an
// account.
//
// Account API Reference: https://developers.cloudflare.com/api/operations/access-identity-providers-list-access-identity-providers
// Zone API Reference: https://developers.cloudflare.com/api/operations/zone-level-access-identity-providers-list-access-identity-providers
func (api *API) ListAccessIdentityProviders(ctx context.Context, rc *ResourceContainer, pageOpts PaginationOptions) ([]AccessIdentityProvider, *ResultInfo, error) {
	baseURL := fmt.Sprintf("/%s/%s/access/identity_providers", rc.Level, rc.Identifier)

	autoPaginate := true
	if pageOpts.PerPage >= 1 || pageOpts.Page >= 1 {
		autoPaginate = false
	}

	if pageOpts.PerPage < 1 {
		pageOpts.PerPage = 25
	}

	if pageOpts.Page < 1 {
		pageOpts.Page = 1
	}

	resultInfo := ResultInfo{
		Page:    pageOpts.Page,
		PerPage: pageOpts.PerPage,
	}

	var accessProviders []AccessIdentityProvider
	var r AccessIdentityProvidersListResponse

	for {
		uri := buildURI(baseURL, resultInfo)
		res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
		if err != nil {
			return []AccessIdentityProvider{}, &ResultInfo{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
		}

		err = json.Unmarshal(res, &r)
		if err != nil {
			return []AccessIdentityProvider{}, &ResultInfo{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
		}
		resultInfo = r.ResultInfo.Next()
		if resultInfo.Done() || autoPaginate {
			break
		}
	}

	return accessProviders, &r.ResultInfo, nil
}

// GetAccessIdentityProvider returns a single Access Identity
// Provider for an account.
//
// Account API Reference: https://developers.cloudflare.com/api/operations/access-identity-providers-get-an-access-identity-provider
// Zone API Reference: https://developers.cloudflare.com/api/operations/zone-level-access-identity-providers-get-an-access-identity-provider
func (api *API) GetAccessIdentityProvider(ctx context.Context, rc *ResourceContainer, identityProviderID string) (AccessIdentityProvider, error) {
	uri := fmt.Sprintf(
		"/%s/%s/access/identity_providers/%s",
		rc.Level,
		rc.Identifier,
		identityProviderID,
	)

	res, err := api.makeRequestContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return AccessIdentityProvider{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var accessIdentityProviderResponse AccessIdentityProviderResponse
	err = json.Unmarshal(res, &accessIdentityProviderResponse)
	if err != nil {
		return AccessIdentityProvider{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	return accessIdentityProviderResponse.Result, nil
}

// CreateAccessIdentityProvider creates a new Access Identity Provider.
//
// Account API Reference: https://developers.cloudflare.com/api/operations/access-identity-providers-add-an-access-identity-provider
// Zone API Reference: https://developers.cloudflare.com/api/operations/zone-level-access-identity-providers-add-an-access-identity-provider
func (api *API) CreateAccessIdentityProvider(ctx context.Context, rc *ResourceContainer, identityProviderConfiguration AccessIdentityProvider) (AccessIdentityProvider, error) {
	uri := fmt.Sprintf("/%s/%s/access/identity_providers", rc.Level, rc.Identifier)

	res, err := api.makeRequestContext(ctx, http.MethodPost, uri, identityProviderConfiguration)
	if err != nil {
		return AccessIdentityProvider{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var accessIdentityProviderResponse AccessIdentityProviderResponse
	err = json.Unmarshal(res, &accessIdentityProviderResponse)
	if err != nil {
		return AccessIdentityProvider{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	return accessIdentityProviderResponse.Result, nil
}

// UpdateAccessIdentityProvider updates an existing Access Identity
// Provider.
//
// Account API Reference: https://developers.cloudflare.com/api/operations/access-identity-providers-update-an-access-identity-provider
// Zone API Reference: https://developers.cloudflare.com/api/operations/zone-level-access-identity-providers-update-an-access-identity-provider
func (api *API) UpdateAccessIdentityProvider(ctx context.Context, rc *ResourceContainer, identityProviderUUID string, identityProviderConfiguration AccessIdentityProvider) (AccessIdentityProvider, error) {
	uri := fmt.Sprintf(
		"/%s/%s/access/identity_providers/%s",
		rc.Level,
		rc.Identifier,
		identityProviderUUID,
	)

	res, err := api.makeRequestContext(ctx, http.MethodPut, uri, identityProviderConfiguration)
	if err != nil {
		return AccessIdentityProvider{}, err
	}

	var accessIdentityProviderResponse AccessIdentityProviderResponse
	err = json.Unmarshal(res, &accessIdentityProviderResponse)
	if err != nil {
		return AccessIdentityProvider{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	return accessIdentityProviderResponse.Result, nil
}

// DeleteAccessIdentityProvider deletes an Access Identity Provider.
//
// Account API Reference: https://developers.cloudflare.com/api/operations/access-identity-providers-delete-an-access-identity-provider
// Zone API Reference: https://developers.cloudflare.com/api/operations/zone-level-access-identity-providers-delete-an-access-identity-provider
func (api *API) DeleteAccessIdentityProvider(ctx context.Context, rc *ResourceContainer, identityProviderUUID string) (AccessIdentityProvider, error) {
	uri := fmt.Sprintf(
		"/%s/%s/access/identity_providers/%s",
		rc.Level,
		rc.Identifier,
		identityProviderUUID,
	)

	res, err := api.makeRequestContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return AccessIdentityProvider{}, fmt.Errorf("%s: %w", errMakeRequestError, err)
	}

	var accessIdentityProviderResponse AccessIdentityProviderResponse
	err = json.Unmarshal(res, &accessIdentityProviderResponse)
	if err != nil {
		return AccessIdentityProvider{}, fmt.Errorf("%s: %w", errUnmarshalError, err)
	}

	return accessIdentityProviderResponse.Result, nil
}
