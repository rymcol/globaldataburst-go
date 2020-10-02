package gdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	baseURI = "https://send.globaldataburst.com/api"
)

// Messenger holds the base resources we need to operate the API
type Messenger interface {
	Send(Message) (err error)
	Repeat(RepeatMessage) (err error)
	Lookup(startDate time.Time, endDate time.Time) ([]Lookup, error)
	Live() (err error)
}

type messenger struct {
	log    *logrus.Logger
	client *http.Client
	apiKey string
}

// NewMessenger creates a new base messenger
func NewMessenger(apiKey string, log *logrus.Logger) (Messenger, error) {
	m := &messenger{
		log:    log,
		client: newClient(),
		apiKey: apiKey,
	}
	_ = Messenger(m)
	return m, nil
}

// Send a constructed message
func (m *messenger) Send(msg Message) error {

	uri := baseURI + "/textmessages"
	if msg.IMEI != nil {
		uri += "/"
		uri += *msg.IMEI
	}

	uri += "?message="
	uri += msg.Message

	if msg.ServiceMask != nil {
		uri += "&serviceMask="
		uri += string(*msg.ServiceMask)
	}

	uri += "&sendDirect="
	uri += strconv.FormatBool(msg.SendDirect)

	json, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	req.Header.Add("X-GDB-APIKEY", m.apiKey)
	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("Non-200 Response")
	}
	return nil
}

// Repeat sends a repeating message
func (m *messenger) Repeat(msg RepeatMessage) error {
	uri := baseURI + "/repeatmessages"

	json, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(json))
	if err != nil {
		return err
	}
	req.Header.Add("X-GDB-APIKEY", m.apiKey)
	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("Non-200 Response")
	}
	return nil
}

// Live returns an error if the service is down
func (m *messenger) Live() error {
	uri := baseURI + "/BurstVpnStatus"

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	req.Header.Add("X-GDB-APIKEY", m.apiKey)
	resp, err := m.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("Non-200 Response")
	}
	return nil
}

// Lookup returns message lookup items between two dates
func (m *messenger) Lookup(startDate time.Time, endDate time.Time) ([]Lookup, error) {
	uri := baseURI + "/textmessages"

	uri += "?startDate="
	uri += startDate.Format(time.RFC3339)

	uri += "&endDate="
	uri += endDate.Format(time.RFC3339)

	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("X-GDB-APIKEY", m.apiKey)
	resp, err := m.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New("Non-200 Response")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []Lookup
	err = json.Unmarshal(body, &results)
	if err != nil {
		return nil, err
	}
	return results, nil
}
