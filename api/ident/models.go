package ident

import (
	"time"

	uuid "github.com/kthomas/go.uuid"

	"github.com/provideservices/provide-go/api"
)

// Application model which is initially owned by the user who created it
type Application struct {
	api.Model

	NetworkID   uuid.UUID              `json:"network_id,omitempty"`
	UserID      uuid.UUID              `json:"user_id,omitempty"` // this is the user that initially created the app
	Name        *string                `json:"name"`
	Description *string                `json:"description"`
	Status      *string                `json:"status,omitempty"` // this is for enrichment purposes only
	Type        *string                `json:"type"`
	Config      map[string]interface{} `json:"config"`
	Hidden      bool                   `json:"hidden"`
}

// Organization model
type Organization struct {
	api.Model

	Name        *string                `json:"name"`
	UserID      *uuid.UUID             `json:"user_id,omitempty"`
	Description *string                `json:"description"`
	Permissions uint32                 `json:"permissions,omitempty"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// Token represents a bearer JWT
type Token struct {
	api.Model

	Token *string `json:"token,omitempty"`

	// OAuth 2 fields
	AccessToken  *string `json:"access_token,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	TokenType    *string `json:"token_type,omitempty"`
	Scope        *string `json:"scope,omitempty"`

	// Ephemeral JWT header fields and claims; these are here for convenience and
	// are not always populated, even if they exist on the underlying token
	Kid       *string    `json:"kid,omitempty"` // key fingerprint
	Audience  *string    `json:"audience,omitempty"`
	Issuer    *string    `json:"issuer,omitempty"`
	IssuedAt  *time.Time `json:"issued_at,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	NotBefore *time.Time `json:"not_before_at,omitempty"`
	Subject   *string    `json:"subject,omitempty"`

	Permissions uint32                 `json:"permissions,omitempty"`
	Data        map[string]interface{} `json:"data,omitempty"`
}

// User represents a user
type User struct {
	api.Model

	Name                   string                 `json:"name"`
	FirstName              string                 `json:"first_name"`
	LastName               string                 `json:"last_name"`
	Email                  string                 `json:"email"`
	Permissions            uint32                 `json:"permissions,omitempty,omitempty"`
	PrivacyPolicyAgreedAt  *time.Time             `json:"privacy_policy_agreed_at,omitempty"`
	TermsOfServiceAgreedAt *time.Time             `json:"terms_of_service_agreed_at,omitempty"`
	Metadata               map[string]interface{} `json:"metadata,omitempty"`
}
