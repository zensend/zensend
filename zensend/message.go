package zensend

import (
	"errors"
	"net/url"
	"strconv"
	"strings"
)

type Encoding string

const (
	GSM  Encoding = "gsm"
	UCS2 Encoding = "ucs2"
)

type OriginatorType string

const (
	ALPHA  OriginatorType = "alpha"
	MSISDN OriginatorType = "msisdn"
)

type Message struct {
	Body                string
	Numbers             []string
	Originator          string
	OriginatorType      OriginatorType
	Encoding            Encoding
	TimeToLiveInMinutes int
}

func (message *Message) toPostParams() (url.Values, error) {
	postParams := url.Values{}

	postParams.Set("BODY", message.Body)
	postParams.Set("ORIGINATOR", message.Originator)

	for i := 0; i < len(message.Numbers); i++ {
		number := message.Numbers[i]
		if strings.Contains(number, ",") {
			return nil, errors.New("invalid character in number: " + number)
		}
	}

	postParams.Set("NUMBERS", strings.Join(message.Numbers, ","))

	if len(message.OriginatorType) > 0 {
		postParams.Set("ORIGINATOR_TYPE", string(message.OriginatorType))
	}

	if len(message.Encoding) > 0 {
		postParams.Set("ENCODING", string(message.Encoding))
	}

	if message.TimeToLiveInMinutes != 0 {
		postParams.Set("TIMETOLIVE", strconv.Itoa(message.TimeToLiveInMinutes))
	}

	return postParams, nil
}
