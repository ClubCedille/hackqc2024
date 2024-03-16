package notifications

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

type Notification struct {
	APIUsername string `json:"api_username"`
	APIPassword string `json:"api_password"`
	Message     string `json:"message"`
	Dst         string `json:"dst"`
	Did         string `json:"did"`
}

// Send notification to all phone numbers
func SendNotification(message string, recipients []string) {
	// TODO: Implement notification sending
	fmt.Println("Sending Notification, message: ", message, " recipients: ", recipients)

	for _, recipient := range recipients {
		fmt.Println("Sending notification to: ", recipient)

		go SendNotificationToPhoneNumber(message, recipient)
	}
}

func SendNotificationToPhoneNumber(message string, recipient string) {

	api_username := os.Getenv("VOIPMS_API_USERNAME")
	api_password := os.Getenv("VOIPMS_API_PASSWORD")
	api_did := os.Getenv("VOIPMS_API_DID")

	params := url.Values{
		"method":       []string{"sendSMS"},
		"message":      []string{message},
		"dst":          []string{recipient},
		"did":          []string{api_did},
		"api_username": []string{api_username},
		"api_password": []string{api_password},
	}

	u := &url.URL{
		Scheme:   "https",
		Host:     "voip.ms",
		Path:     "/api/v1/rest.php",
		RawQuery: params.Encode(),
	}
	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		fmt.Println("Error creating request: ", err)
	}

	client := &http.Client{}
	result, err := client.Do(req)

	if err != nil || result.StatusCode != 200 {
		fmt.Println("Error sending notification to phone number: ", err)
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		fmt.Println("Error reading response body: ", err)
	}

	fmt.Println("Result sending notification to phone number: ", string(body))
}
