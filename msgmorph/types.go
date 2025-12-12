// Package msgmorph provides a Go SDK for the MsgMorph API.
//
// The MsgMorph SDK allows you to manage contacts and feedback collection
// in your Go applications with a simple, idiomatic API.
package msgmorph

import "time"

// Contact represents a user/contact entity in MsgMorph.
// Contacts are individuals who can receive feedback requests.
type Contact struct {
	// ID is the unique identifier for the contact in MsgMorph.
	ID string `json:"id"`

	// ExternalID is your system's user ID, used to link contacts to your users.
	ExternalID string `json:"externalId"`

	// Email is the contact's email address.
	Email string `json:"email"`

	// Name is the contact's display name. May be nil if not provided.
	Name *string `json:"name"`

	// ProjectID is the MsgMorph project ID this contact belongs to.
	ProjectID string `json:"projectId"`

	// FeedbackSent indicates whether feedback has been sent to this contact.
	FeedbackSent bool `json:"feedbackSent"`

	// FeedbackScheduledAt is the time when feedback is scheduled to be sent.
	// May be nil if no feedback is scheduled.
	FeedbackScheduledAt *time.Time `json:"feedbackScheduledAt"`

	// CreatedAt is the timestamp when the contact was created.
	CreatedAt time.Time `json:"createdAt"`

	// UpdatedAt is the timestamp when the contact was last updated.
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateContactInput contains the parameters for creating a new contact.
type CreateContactInput struct {
	// ExternalID is your system's user ID (required).
	// This is used to prevent duplicate contacts and link them to your users.
	ExternalID string `json:"externalId"`

	// Email is the contact's email address (required).
	Email string `json:"email"`

	// Name is the contact's display name (optional).
	Name string `json:"name,omitempty"`

	// ProjectID is the MsgMorph project ID to associate this contact with (required).
	ProjectID string `json:"projectId"`
}

// UpdateContactInput contains the parameters for updating an existing contact.
// All fields are optional; only provided fields will be updated.
type UpdateContactInput struct {
	// Email is the new email address for the contact.
	Email string `json:"email,omitempty"`

	// Name is the new display name for the contact.
	Name string `json:"name,omitempty"`
}

// ListContactsParams contains the parameters for listing contacts.
type ListContactsParams struct {
	// ProjectID filters contacts by project ID (required).
	ProjectID string `url:"projectId"`
}

// APIResponse is the standard response wrapper from the MsgMorph API.
type APIResponse[T any] struct {
	// Data contains the response payload.
	Data T `json:"data"`

	// Error contains an error message if the request failed.
	Error string `json:"error,omitempty"`
}
