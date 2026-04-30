package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const fonnteURL = "https://api.fonnte.com/send"

type fonntePayload struct {
	Target  string `json:"target"`
	Message string `json:"message"`
}

// SendWA mengirim pesan WhatsApp via Fonnte API
func SendWA(noHP, message string) error {
	payload, err := json.Marshal(fonntePayload{
		Target:  noHP,
		Message: message,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", fonnteURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", os.Getenv("FONNTE_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("fonnte API error: status %d", resp.StatusCode)
	}

	return nil
}
