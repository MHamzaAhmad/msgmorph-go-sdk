# MsgMorph Go SDK

Official Go SDK for the [MsgMorph](https://msgmorph.com) API.

## Installation

```bash
go get github.com/MHamzaAhmad/msgmorph-go-sdk
```

## Requirements

- Go 1.21 or later

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    msgmorph "github.com/MHamzaAhmad/msgmorph-go-sdk/msgmorph"
)

func main() {
    // Initialize the client
    client := msgmorph.NewClient(
        os.Getenv("MSGMORPH_API_KEY"),
        os.Getenv("MSGMORPH_ORGANIZATION_ID"),
    )

    // Create a contact
    contact, err := client.Contacts.Create(context.Background(), msgmorph.CreateContactInput{
        ExternalID: "user-123",
        Email:      "user@example.com",
        ProjectID:  os.Getenv("MSGMORPH_PROJECT_ID"),
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created contact: %s\n", contact.ID)
}
```

## Getting Your Credentials

1. Go to your [MsgMorph Dashboard](https://msgmorph.com/dashboard)
2. Navigate to **Settings** → **General** to get your **Organization ID** and **Project ID**
3. Go to **Account Settings** → **API Keys** to create an API key

## Configuration

### Basic Initialization

```go
client := msgmorph.NewClient(
    os.Getenv("MSGMORPH_API_KEY"),
    os.Getenv("MSGMORPH_ORGANIZATION_ID"),
)
```

### Configuration Options

```go
// With custom base URL (for local development)
client := msgmorph.NewClient(
    apiKey,
    orgID,
    msgmorph.WithBaseURL("http://localhost:3001"),
)

// With custom timeout
client := msgmorph.NewClient(
    apiKey,
    orgID,
    msgmorph.WithTimeout(60 * time.Second),
)

// With custom HTTP client
httpClient := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns: 10,
    },
}
client := msgmorph.NewClient(
    apiKey,
    orgID,
    msgmorph.WithHTTPClient(httpClient),
)
```

## Usage

### Contacts

#### Create a Contact

```go
contact, err := client.Contacts.Create(ctx, msgmorph.CreateContactInput{
    ExternalID: "user-123",       // Required: Your system's user ID
    Email:      "alice@example.com", // Required: User's email
    Name:       "Alice Smith",    // Optional: User's display name
    ProjectID:  projectID,        // Required: MsgMorph project ID
})
```

#### List Contacts

```go
contacts, err := client.Contacts.List(ctx, msgmorph.ListContactsParams{
    ProjectID: projectID,
})
for _, c := range contacts {
    fmt.Printf("Contact: %s (%s)\n", *c.Name, c.Email)
}
```

#### Get a Contact

```go
contact, err := client.Contacts.Get(ctx, "cnt_abc123")
```

#### Update a Contact

```go
updated, err := client.Contacts.Update(ctx, "cnt_abc123", msgmorph.UpdateContactInput{
    Email: "newemail@example.com",
    Name:  "Alice Johnson",
})
```

#### Delete a Contact

```go
err := client.Contacts.Delete(ctx, "cnt_abc123")
```

## Error Handling

All methods return errors that can be type-asserted to `*msgmorph.Error`:

```go
contact, err := client.Contacts.Get(ctx, "invalid-id")
if err != nil {
    var msgErr *msgmorph.Error
    if errors.As(err, &msgErr) {
        fmt.Printf("Error: %s\n", msgErr.Message)
        fmt.Printf("Code: %s\n", msgErr.Code)
        fmt.Printf("Status: %d\n", msgErr.Status)
        fmt.Printf("Hint: %s\n", msgErr.Hint)

        // Check specific error types
        if msgErr.IsNotFound() {
            fmt.Println("Contact not found")
        }
        if msgErr.IsUnauthorized() {
            fmt.Println("Invalid API key")
        }
        if msgErr.IsValidationError() {
            fmt.Println("Invalid input")
        }
    }
}
```

### Error Codes

| Code                      | Description                        |
| ------------------------- | ---------------------------------- |
| `INVALID_API_KEY`         | Invalid or missing API key         |
| `INVALID_ORGANIZATION_ID` | Invalid or missing organization ID |
| `VALIDATION_ERROR`        | Invalid request data               |
| `UNAUTHORIZED`            | Authentication failed              |
| `FORBIDDEN`               | Access denied                      |
| `NOT_FOUND`               | Resource not found                 |
| `CONFLICT`                | Resource conflict                  |
| `ALREADY_EXISTS`          | Resource already exists            |
| `INTERNAL_ERROR`          | Server error                       |
| `NETWORK_ERROR`           | Network connectivity issue         |
| `TIMEOUT`                 | Request timeout                    |

## Environment Variables

We recommend using environment variables for credentials:

```bash
export MSGMORPH_API_KEY=msgmorph_xxxxxxxxxxxxx
export MSGMORPH_ORGANIZATION_ID=org_xxxxxxxxxxxxx
export MSGMORPH_PROJECT_ID=xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

## License

MIT License - see [LICENSE](LICENSE) for details.

## Support

- 📚 [Documentation](https://docs.msgmorph.com)
- 📧 [Support](mailto:support@msgmorph.com)
- 🐛 [Issues](https://github.com/MHamzaAhmad/msgmorph/issues)
