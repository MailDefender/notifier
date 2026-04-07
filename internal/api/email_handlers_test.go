package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/maxatome/go-testdeep/helpers/tdhttp"
)

const (
	mailhogMaxRetries = 5
	mailhogRetryDelay = 500 * time.Millisecond
)

// getMailhogMessageCount retrieves the total message count from Mailhog with retry logic
func getMailhogMessageCount(t *testing.T) int {
	var resp *http.Response
	var err error

	// Retry logic to handle Mailhog startup delays
	for attempt := 0; attempt < mailhogMaxRetries; attempt++ {
		resp, err = http.Get("http://localhost:8025/api/v2/messages")
		if err == nil {
			break
		}
		if attempt < mailhogMaxRetries-1 {
			time.Sleep(mailhogRetryDelay)
		}
	}

	if err != nil {
		t.Fatalf("failed to connect to Mailhog after %d retries: %v", mailhogMaxRetries, err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode Mailhog response: %v", err)
	}

	total, ok := result["total"].(float64)
	if !ok {
		t.Fatal("could not extract total count from Mailhog response")
	}

	return int(total)
}

func TestSendEmail(t *testing.T) {
	tests := []struct {
		name               string
		endpoint           string
		payload            string
		expectedStatusCode int
		iterations         int
		expectedEmailCount int
	}{
		{
			name:               "send email",
			endpoint:           "/v1/send/email",
			payload:            `{"to": ["user@example.com"],"subject": "Test email","body": "<p>This is a test email.</p>"}`,
			expectedStatusCode: http.StatusOK,
			iterations:         1,
			expectedEmailCount: 1,
		},
		{
			name:               "send 200 emails",
			endpoint:           "/v1/send/email",
			payload:            `{"to": ["user@example.com"],"subject": "Test email","body": "<p>This is a test email.</p>"}`,
			expectedStatusCode: http.StatusOK,
			iterations:         200,
			expectedEmailCount: 200,
		},
		{
			name:               "missing required fields",
			endpoint:           "/v1/send/email",
			payload:            `{"to": ["user@example.com"],"body": "<p>This is a test email.</p>"}`,
			expectedStatusCode: http.StatusBadRequest,
			iterations:         1,
			expectedEmailCount: 0,
		},
		{
			name:               "invalid email address",
			endpoint:           "/v1/send/email",
			payload:            `{"to": ["invalid-email"],"subject": "Test email","body": "<p>This is a test email.</p>"}`,
			expectedStatusCode: http.StatusBadRequest,
			iterations:         1,
			expectedEmailCount: 0,
		},
	}

	tester := tdhttp.NewTestAPI(t, api.Handler())
	for _, test := range tests {
		// Get email count before test
		countBefore := getMailhogMessageCount(t)

		// Run the test case iterations times
		for i := 0; i < test.iterations; i++ {
			tt := tester.Name(test.name).Post(test.endpoint, strings.NewReader(test.payload), "Content-Type", "application/json", "Accept", "application/json")
			tt.CmpStatus(test.expectedStatusCode)
		}

		// Get email count after test
		countAfter := getMailhogMessageCount(t)

		// Assert that the expected number of emails were received
		actualReceived := countAfter - countBefore
		if actualReceived != test.expectedEmailCount {
			t.Errorf("expected %d emails, got %d", test.expectedEmailCount, actualReceived)
		}
	}
}
