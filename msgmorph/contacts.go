package msgmorph

import (
	"context"
	"fmt"
	"net/http"
)

// ContactsResource provides methods to manage contacts in MsgMorph.
//
// Contacts represent users in your application who can receive feedback requests.
// Use this resource to create, list, get, update, and delete contacts.
//
// Example usage:
//
//	// Create a contact
//	contact, err := client.Contacts.Create(ctx, msgmorph.CreateContactInput{
//	    ExternalID: "user-123",
//	    Email:      "user@example.com",
//	    ProjectID:  "proj-456",
//	})
//
//	// List contacts
//	contacts, err := client.Contacts.List(ctx, msgmorph.ListContactsParams{
//	    ProjectID: "proj-456",
//	})
//
//	// Get a contact
//	contact, err := client.Contacts.Get(ctx, "contact-id")
//
//	// Update a contact
//	updated, err := client.Contacts.Update(ctx, "contact-id", msgmorph.UpdateContactInput{
//	    Name: "New Name",
//	})
//
//	// Delete a contact
//	err = client.Contacts.Delete(ctx, "contact-id")
type ContactsResource struct {
	client *Client
}

// Create creates a new contact in MsgMorph.
//
// The ExternalID field should be your system's user ID. This is used to prevent
// duplicate contacts and to link contacts to users in your system.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - input: Contact creation parameters
//
// Returns the created Contact or an error.
//
// Example:
//
//	contact, err := client.Contacts.Create(ctx, msgmorph.CreateContactInput{
//	    ExternalID: "user-123",       // Required: Your system's user ID
//	    Email:      "alice@example.com", // Required: User's email
//	    Name:       "Alice Smith",    // Optional: User's display name
//	    ProjectID:  os.Getenv("MSGMORPH_PROJECT_ID"), // Required: MsgMorph project ID
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Created contact: %s\n", contact.ID)
//
// Errors:
//   - ErrValidationError: If required fields are missing
//   - ErrAlreadyExists: If a contact with the same externalId already exists
//   - ErrUnauthorized: If the API key is invalid
func (r *ContactsResource) Create(ctx context.Context, input CreateContactInput) (*Contact, error) {
	var contact Contact
	err := r.client.request(ctx, http.MethodPost, "/api/v1/contacts", input, &contact)
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

// List retrieves all contacts for a project.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - params: Query parameters for filtering contacts
//
// Returns a slice of Contact objects or an error.
//
// Example:
//
//	contacts, err := client.Contacts.List(ctx, msgmorph.ListContactsParams{
//	    ProjectID: os.Getenv("MSGMORPH_PROJECT_ID"),
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	for _, c := range contacts {
//	    fmt.Printf("Contact: %s (%s)\n", c.Name, c.Email)
//	}
//
// Errors:
//   - ErrValidationError: If projectId is missing
//   - ErrUnauthorized: If the API key is invalid
func (r *ContactsResource) List(ctx context.Context, params ListContactsParams) ([]Contact, error) {
	path := fmt.Sprintf("/api/v1/contacts?projectId=%s", params.ProjectID)

	var contacts []Contact
	err := r.client.request(ctx, http.MethodGet, path, nil, &contacts)
	if err != nil {
		return nil, err
	}
	return contacts, nil
}

// Get retrieves a single contact by ID.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - id: The contact's unique ID in MsgMorph
//
// Returns the Contact or an error.
//
// Example:
//
//	contact, err := client.Contacts.Get(ctx, "cnt_abc123")
//	if err != nil {
//	    var msgErr *msgmorph.Error
//	    if errors.As(err, &msgErr) && msgErr.IsNotFound() {
//	        fmt.Println("Contact not found")
//	        return
//	    }
//	    log.Fatal(err)
//	}
//	fmt.Printf("Contact: %s (%s)\n", contact.Name, contact.Email)
//
// Errors:
//   - ErrNotFound: If the contact doesn't exist
//   - ErrUnauthorized: If the API key is invalid
func (r *ContactsResource) Get(ctx context.Context, id string) (*Contact, error) {
	path := fmt.Sprintf("/api/v1/contacts/%s", id)

	var contact Contact
	err := r.client.request(ctx, http.MethodGet, path, nil, &contact)
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

// Update modifies an existing contact.
//
// Only the fields provided in the input will be updated.
// All fields in UpdateContactInput are optional.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - id: The contact's unique ID in MsgMorph
//   - input: Fields to update
//
// Returns the updated Contact or an error.
//
// Example:
//
//	updated, err := client.Contacts.Update(ctx, "cnt_abc123", msgmorph.UpdateContactInput{
//	    Email: "newemail@example.com",
//	    Name:  "Alice Johnson",
//	})
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Updated contact: %s\n", updated.Name)
//
// Errors:
//   - ErrNotFound: If the contact doesn't exist
//   - ErrValidationError: If the input is invalid
//   - ErrUnauthorized: If the API key is invalid
func (r *ContactsResource) Update(ctx context.Context, id string, input UpdateContactInput) (*Contact, error) {
	path := fmt.Sprintf("/api/v1/contacts/%s", id)

	var contact Contact
	err := r.client.request(ctx, http.MethodPatch, path, input, &contact)
	if err != nil {
		return nil, err
	}
	return &contact, nil
}

// Delete removes a contact.
//
// This operation is permanent and cannot be undone.
//
// Parameters:
//   - ctx: Context for cancellation and timeout control
//   - id: The contact's unique ID in MsgMorph
//
// Returns nil on success or an error.
//
// Example:
//
//	err := client.Contacts.Delete(ctx, "cnt_abc123")
//	if err != nil {
//	    var msgErr *msgmorph.Error
//	    if errors.As(err, &msgErr) && msgErr.IsNotFound() {
//	        fmt.Println("Contact already deleted or doesn't exist")
//	        return
//	    }
//	    log.Fatal(err)
//	}
//	fmt.Println("Contact deleted successfully")
//
// Errors:
//   - ErrNotFound: If the contact doesn't exist
//   - ErrUnauthorized: If the API key is invalid
func (r *ContactsResource) Delete(ctx context.Context, id string) error {
	path := fmt.Sprintf("/api/v1/contacts/%s", id)
	return r.client.request(ctx, http.MethodDelete, path, nil, nil)
}
