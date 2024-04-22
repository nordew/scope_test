package v1

import (
	"github.com/nordew/scope_test/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_handlePostRequest(t *testing.T) {
	mockWorkerService := &mocks.WorkerService{}
	mockWorkerService.On("Submit", mock.AnythingOfType("model.Job")).Return()

	handler := NewHandler(mockWorkerService)

	testCases := []struct {
		name             string
		requestBody      string
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "ValidRequest",
			requestBody:      `{"test":"value to process"}`,
			expectedStatus:   http.StatusOK,
			expectedResponse: `{"message": "Success"}`,
		},
		{
			name:             "InvalidRequestBody",
			requestBody:      `invalid JSON`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `{"error": "invalid request body"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("POST", "/api/", strings.NewReader(tc.requestBody))
			assert.NoError(t, err, "error creating request")

			rr := httptest.NewRecorder()

			handler.Init().ServeHTTP(rr, req)

			assert.Equal(t, tc.expectedStatus, rr.Code, "unexpected status code")

			assert.JSONEq(t, tc.expectedResponse, strings.TrimSpace(rr.Body.String()), "unexpected response body")

			if tc.expectedStatus == http.StatusOK {
				mockWorkerService.AssertCalled(t, "Submit", mock.AnythingOfType("model.Job"))
			} else {
				mockWorkerService.AssertNotCalled(t, "Submit")
			}
		})
	}
}
