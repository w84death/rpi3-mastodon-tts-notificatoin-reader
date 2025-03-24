package main

import (
    "encoding/json"
    "fmt"
    "io"
    "log"
    "net/http"
    "net/url"
    "os"
    "os/exec"
    "strings"

    "golang.org/x/net/html"
)

// Account represents a Mastodon account
type Account struct {
    Username    string `json:"username"`
    DisplayName string `json:"display_name"`
}

// Status represents a Mastodon status
type Status struct {
    Content string `json:"content"`
}

// Notification represents a Mastodon notification
type Notification struct {
    ID      string   `json:"id"`
    Type    string   `json:"type"`
    Account Account  `json:"account"`
    Status  *Status  `json:"status"` // Pointer because it might be null
}

// extractText extracts plain text from HTML content
func extractText(htmlStr string) string {
    doc, err := html.Parse(strings.NewReader(htmlStr))
    if err != nil {
        return htmlStr // Fallback to original if parsing fails
    }
    var f func(*html.Node) string
    f = func(n *html.Node) string {
        if n.Type == html.TextNode {
            return n.Data
        }
        var text string
        for c := n.FirstChild; c != nil; c = c.NextSibling {
            text += f(c)
        }
        return text
    }
    return strings.TrimSpace(f(doc))
}

func main() {
    // Retrieve Mastodon credentials from environment variables
    instanceURL := os.Getenv("MASTODON_INSTANCE_URL")
    accessToken := os.Getenv("MASTODON_ACCESS_TOKEN")
    if instanceURL == "" || accessToken == "" {
        log.Fatal("Please set MASTODON_INSTANCE_URL and MASTODON_ACCESS_TOKEN environment variables")
    }

    // Read the last processed notification ID from a file
    lastIDBytes, err := os.ReadFile("last_notification_id.txt")
    lastID := ""
    if err == nil {
        lastID = strings.TrimSpace(string(lastIDBytes))
    }

    // Construct the API URL with optional since_id parameter
    u, err := url.Parse(instanceURL + "/api/v1/notifications")
    if err != nil {
        log.Fatalf("Failed to parse URL: %v", err)
    }
    q := u.Query()
    if lastID != "" {
        q.Set("since_id", lastID)
    }
    u.RawQuery = q.Encode()

    // Create and send the HTTP request
    req, err := http.NewRequest("GET", u.String(), nil)
    if err != nil {
        log.Fatalf("Failed to create request: %v", err)
    }
    req.Header.Set("Authorization", "Bearer "+accessToken)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatalf("Failed to fetch notifications: %v", err)
    }
    defer resp.Body.Close()

    // Verify the response status
    if resp.StatusCode != http.StatusOK {
        log.Fatalf("API request failed with status %d", resp.StatusCode)
    }

    // Read and parse the response body
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Fatalf("Failed to read response body: %v", err)
    }
    var notifications []Notification
    err = json.Unmarshal(body, &notifications)
    if err != nil {
        log.Fatalf("Failed to parse JSON: %v", err)
    }

    // Process notifications if there are any
    if len(notifications) > 0 {
        for _, notif := range notifications {
            // Use display name if available, otherwise username
            name := notif.Account.DisplayName
            if name == "" {
                name = notif.Account.Username
            }

            // Construct speech based on notification type
            var speech string
            switch notif.Type {
            case "mention":
                if notif.Status != nil {
                    message := extractText(notif.Status.Content)
                    speech = fmt.Sprintf("%s mentioned you: %s", name, message)
                }
            case "favourite":
                speech = fmt.Sprintf("%s favorited your post", name)
            case "reblog":
                speech = fmt.Sprintf("%s boosted your post", name)
            case "follow":
                speech = fmt.Sprintf("%s started following you", name)
            default:
                continue // Skip unsupported types
            }

            if speech != "" {
				// Use a shell command to pipe the speech text through Piper and aplay
				cmd := exec.Command("sh", "-c", fmt.Sprintf("echo %q | piper --model en_US-danny-low.onnx --output-raw | aplay -r 16000 -f S16_LE -t raw -", speech))
				err := cmd.Run()
				if err != nil {
				log.Printf("Failed to speak notification: %v", err)
				}
			}
        }

        // Update the last notification ID to the newest one
        newLastID := notifications[0].ID
        err = os.WriteFile("last_notification_id.txt", []byte(newLastID), 0644)
        if err != nil {
            log.Printf("Failed to write last ID: %v", err)
        }
    }
}