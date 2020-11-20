package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jrieck1991/golang-di-example/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// testCase represents a single test for our http server
type testCase struct {
	name               string
	mockStore          *storage.MockClient
	path               string
	expectedStatusCode int
}

func TestServer(t *testing.T) {

	// create test cases for our server
	// table driven testing helps keep things streamlined
	testCases := []testCase{
		{
			name:               "index 200",
			path:               "/",
			mockStore:          &storage.MockClient{},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "get 200",
			path: "/get?key=foo",
			// inject a mock which returns a value
			mockStore: &storage.MockClient{
				DoGet: func(key string) (string, error) {
					return "bar", nil
				},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "get value not found",
			path: "/get?key=foo",
			// inject a mock which returns an error
			mockStore: &storage.MockClient{
				DoGet: func(key string) (string, error) {
					return "", errors.New("item not found")
				},
			},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "no route",
			path:               "/foo",
			mockStore:          &storage.MockClient{},
			expectedStatusCode: http.StatusNotFound,
		},
	}

	// run all test cases
	for _, tc := range testCases {

		fmt.Printf("running test %s\n", tc.name)

		// form request
		req, err := http.NewRequest(http.MethodGet, tc.path, nil)
		require.Nil(t, err)

		// init our server storage client dependency injection of mock
		// this way we don't actually need our datastore to test the different execution paths of our http server
		server := New(WithStorageClient(tc.mockStore))

		// record response and serve http request
		recorder := httptest.NewRecorder()
		server.ServeHTTP(recorder, req)

		// get recorded http response
		res := recorder.Result()

		// verify http status code
		// we use assert here so we can run all the test cases to see all results if there's a failure
		assert.Equal(t, tc.expectedStatusCode, res.StatusCode)
	}
}
