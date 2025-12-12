// Package msgmorph provides the official Go SDK for the MsgMorph API.
//
// MsgMorph is a feedback collection platform that helps you gather
// user feedback from your applications.
//
// # Installation
//
//	go get github.com/MHamzaAhmad/msgmorph-go-sdk
//
// # Quick Start
//
//	package main
//
//	import (
//	    "context"
//	    "fmt"
//	    "log"
//	    "os"
//
//	    msgmorph "github.com/MHamzaAhmad/msgmorph-go-sdk/msgmorph"
//	)
//
//	func main() {
//	    // Initialize the client
//	    client := msgmorph.NewClient(
//	        os.Getenv("MSGMORPH_API_KEY"),
//	        os.Getenv("MSGMORPH_ORGANIZATION_ID"),
//	    )
//
//	    // Create a contact
//	    contact, err := client.Contacts.Create(context.Background(), msgmorph.CreateContactInput{
//	        ExternalID: "user-123",
//	        Email:      "user@example.com",
//	        ProjectID:  os.Getenv("MSGMORPH_PROJECT_ID"),
//	    })
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    fmt.Printf("Created contact: %s\n", contact.ID)
//	}
//
// # Error Handling
//
// All methods return errors that can be type-asserted to *msgmorph.Error
// for detailed error information:
//
//	contact, err := client.Contacts.Get(ctx, "invalid-id")
//	if err != nil {
//	    var msgErr *msgmorph.Error
//	    if errors.As(err, &msgErr) {
//	        fmt.Printf("Error Code: %s\n", msgErr.Code)
//	        fmt.Printf("Status: %d\n", msgErr.Status)
//	        fmt.Printf("Hint: %s\n", msgErr.Hint)
//
//	        // Check specific error types
//	        if msgErr.IsNotFound() {
//	            fmt.Println("Contact not found")
//	        }
//	    }
//	}
//
// # Configuration
//
// The client can be customized with options:
//
//	client := msgmorph.NewClient(
//	    apiKey,
//	    orgID,
//	    msgmorph.WithBaseURL("http://localhost:3001"),
//	    msgmorph.WithTimeout(60 * time.Second),
//	)
package msgmorph
