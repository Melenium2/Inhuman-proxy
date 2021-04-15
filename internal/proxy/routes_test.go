package proxy_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Melenium2/inhuman-reverse-proxy/internal/proxy"
	"github.com/Melenium2/inhuman-reverse-proxy/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewProxy(t *testing.T) {
	var tt = []struct {
		name         string
		body         map[string]interface{}
		setupMock    func() *mocks.ProxyStorage
		expectedCode int
	}{
		{
			name: "should return success message",
			body: map[string]interface{}{
				"code":    "ru",
				"address": []string{"1", "2", "3"},
			},
			setupMock: func() *mocks.ProxyStorage {
				mockStorage := mocks.ProxyStorage{}
				mockStorage.
					On("Set", mock.IsType(&fasthttp.RequestCtx{}), "ru", "1", "2", "3").
					Return(nil).
					Once()

				return &mockStorage
			},
			expectedCode: 200,
		},
		{
			name: "should return error if can not parse body",
			body: map[string]interface{}{
				"123123213": "1231231",
			},
			setupMock: func() *mocks.ProxyStorage {
				return &mocks.ProxyStorage{}
			},
			expectedCode: 400,
		},
		{
			name: "should return error if error occurs in set method",
			body: map[string]interface{}{
				"code":    "ru",
				"address": []string{"1", "2", "3"},
			},
			setupMock: func() *mocks.ProxyStorage {
				mockStorage := mocks.ProxyStorage{}
				mockStorage.
					On("Set", mock.IsType(&fasthttp.RequestCtx{}), "ru", "1", "2", "3").
					Return(errors.New("error :(")).
					Once()

				return &mockStorage
			},
			expectedCode: 400,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			storeMock := test.setupMock()
			app := fiber.New()
			app.Get("/test", proxy.NewProxy(storeMock))

			byts, _ := json.Marshal(test.body)

			req := httptest.NewRequest(http.MethodGet, "/test", bytes.NewReader(byts))
			resp, err := app.Test(req, 1000*60)
			require.NoError(t, err)
			assert.Equal(t, test.expectedCode, resp.StatusCode)
		})
	}
}
