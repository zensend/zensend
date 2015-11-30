package zensend
import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func ValidMessage() *Message {
	return &Message{
		Body:       "This is a test",
		Originator: "Originator",
		Numbers:    []string{"447877878787"}}
}

func StubHttpResponseAndTest(statusCode int, stubbedResponse string, f func(client Client)) *http.Request {
	var lastRequest *http.Request

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		lastRequest = r
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, stubbedResponse)
	}))

	defer server.Close()

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: transport}

	f(Client{APIKey: "API KEY", HTTPClient: httpClient, URL: "http://localhost:8084"})
	return lastRequest
}

func TestSendSMSSuccess(t *testing.T) {
	stubbedResponse := `{"success": {"txguid":"some-guid-123","numbers":1,"smsparts":1,"encoding":"alpha","cost_in_pence":12.34,"new_balance_in_pence":10.2}}`

	request := StubHttpResponseAndTest(200, stubbedResponse, func(client Client) {
		response, error := client.SendSMS(ValidMessage())

		correctResponse := &SendSMSResponse{
			TxGuid:            "some-guid-123",
			Numbers:           1,
			SmsParts:          1,
			Encoding:          "alpha",
			CostInPence:       12.34,
			NewBalanceInPence: 10.2,
		}

		assert.Nil(t, error)
		assert.Equal(t, correctResponse, response)
	})

	assert.Equal(t, "BODY=This+is+a+test&NUMBERS=447877878787&ORIGINATOR=Originator", request.Form.Encode())
}

func TestOperatorLookupSuccess(t *testing.T) {
	stubbedResponse := `{"success": {"mcc":"123","mnc":"456","operator":"o2-uk","cost_in_pence":2.5,"new_balance_in_pence":10.2}}`

	request := StubHttpResponseAndTest(200, stubbedResponse, func(client Client) {
		response, error := client.LookupOperator("441234567890")

		correctResponse := &OperatorLookupResponse{
			MCC:               "123",
			MNC:               "456",
			Operator:          "o2-uk",
			CostInPence:       2.5,
			NewBalanceInPence: 10.2,
		}

		assert.Nil(t, error)
		assert.Equal(t, correctResponse, response)
	})

	assert.Equal(t, "NUMBER=441234567890", request.Form.Encode())
}

func TestOperatorLookupError(t *testing.T) {
	stubbedResponse := `{"failure": {"failcode":"DATA_MISSING","cost_in_pence":2.5,"new_balance_in_pence":10.2}}`

	request := StubHttpResponseAndTest(503, stubbedResponse, func(client Client) {
		response, error := client.LookupOperator("441234567890")

		assert.Nil(t, response)

		zenSendError := error.(ZenSendError)

		assert.Equal(t, 2.5, *zenSendError.CostInPence)
		assert.Equal(t, 10.2, *zenSendError.NewBalanceInPence)
		assert.Equal(t, "DATA_MISSING", zenSendError.FailCode)
		assert.Equal(t, 503, zenSendError.StatusCode)

	})

	assert.Equal(t, "NUMBER=441234567890", request.Form.Encode())
}

func TestSendAll(t *testing.T) {

	stubbedResponse := `{"success": {"txguid":"some-guid-123","numbers":1,"smsparts":1,"encoding":"alpha","cost_in_pence":12.34,"new_balance_in_pence":10.2}}`

	request := StubHttpResponseAndTest(200, stubbedResponse, func(client Client) {
		message := ValidMessage()

		message.Encoding = UCS2
		message.OriginatorType = MSISDN
		message.TimeToLiveInMinutes = 100

		client.SendSMS(message)

	})

	assert.Equal(t, "BODY=This+is+a+test&ENCODING=ucs2&NUMBERS=447877878787&ORIGINATOR=Originator&ORIGINATOR_TYPE=msisdn&TIMETOLIVE=100", request.Form.Encode())
}

func TestSendSMSInvalidNumbers(t *testing.T) {
	client := New("API KEY")

	_, error := client.SendSMS(&Message{
		Body:       "This is a test",
		Originator: "Originator",
		Numbers:    []string{"447877,878787"}})

	assert.NotNil(t, error)
	assert.Equal(t, errors.New("invalid character in number: 447877,878787"), error)
}

func TestSendSMSKnownFailure(t *testing.T) {
	stubbedResponse := `{"failure": {"failcode":"invalid numbers","parameter":"numbers"}}`
	StubHttpResponseAndTest(400, stubbedResponse, func(client Client) {
		_, error := client.SendSMS(&Message{
			Body:       "This is a test",
			Originator: "Originator",
			Numbers:    []string{"4444"}})

		assert.NotNil(t, error)
		assert.Equal(t, ZenSendError{StatusCode: 400, FailCode: "invalid numbers", Parameter: "numbers"}, error)
	})
}

func TestSendSMSUnknownFailure(t *testing.T) {
	stubbedResponse := `{}`
	StubHttpResponseAndTest(500, stubbedResponse, func(client Client) {
		_, error := client.SendSMS(ValidMessage())

		assert.NotNil(t, error)
		assert.Equal(t, ZenSendError{StatusCode: 500, FailCode: "", Parameter: ""}, error)
	})
}

func TestCheckBalanceSuccess(t *testing.T) {
	stubbedResponse := `{"success": {"balance":10.2}}`
	StubHttpResponseAndTest(200, stubbedResponse, func(client Client) {
		response, error := client.CheckBalance()

		assert.Nil(t, error)
		assert.Equal(t, 10.2, response)
	})
}

func TestGetPrices(t *testing.T) {
	stubbedResponse := `{"success": {"prices_in_pence":{"GB":1.23,"US":1.24}}}`
	StubHttpResponseAndTest(200, stubbedResponse, func(client Client) {
		response, error := client.GetPrices()

		assert.Nil(t, error)
		assert.Equal(t, map[string]float64{
			"GB": 1.23,
			"US": 1.24,
		}, response)
	})
}

func TestCheckBalanceKnownFailure(t *testing.T) {
	stubbedResponse := `{"failure": {"failcode":"failcode","parameter":"param"}}`
	StubHttpResponseAndTest(400, stubbedResponse, func(client Client) {
		_, error := client.CheckBalance()

		assert.NotNil(t, error)
		assert.Equal(t, ZenSendError{StatusCode: 400, FailCode: "failcode", Parameter: "param"}, error)
	})
}

func TestSendCheckBalanceFailure(t *testing.T) {
	stubbedResponse := `{}`
	StubHttpResponseAndTest(500, stubbedResponse, func(client Client) {
		_, error := client.CheckBalance()

		assert.NotNil(t, error)
		assert.Equal(t, ZenSendError{StatusCode: 500, FailCode: "", Parameter: ""}, error)
	})
}
