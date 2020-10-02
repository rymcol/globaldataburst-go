package gdb

import "time"

// Message holds a message you want to send
type Message struct {
	Message                string  `json:"-"` // Text of the message, limited to 100 characters. Any messages in excess of 100 characters will be accepted, but truncated to the first 100 characters.
	ServiceMask            *byte   `json:"-"`
	IMEI                   *string `json:"-"`
	SendDirect             bool    `json:"-"` // true is the default, but still needs to be set
	RepeatOptions          RepeatOptions
	DeliveryTime           time.Time
	ExpirationTime         time.Time
	BroadcastCoverageAreas []BroadcastCoverageArea
	GlobalDeliveryAreas    []int
	DeliveryPoints         []DeliveryPoint
}

// RepeatMessage holds a repeatable message for sending
type RepeatMessage struct {
	DeliveryOptions      RepeatOptions
	GdbRepeatMessageType int // See GdbRepeatMessageType Enumeration Below
	IMEI                 string
	TransmissionDelay    int
	RepeatCount          int
	RepeatInterval       int
	RicCode              string
	GroupCode            string
	Payload              []byte
}

// RepeatOptions holds redelivery parameters
type RepeatOptions struct {
	Count    int
	Interval int
}

// GdbRepeatMessageType holds details of message type for repeats
var GdbRepeatMessageType = struct {
	TextMessage          int
	ConfigurationMessage int
	StatusRequest        int
}{
	TextMessage:          1,
	ConfigurationMessage: 2,
	StatusRequest:        3,
}

// BroadcastCoverageArea holds a delivery area
type BroadcastCoverageArea struct {
	Code     string
	Name     string
	GdaCount int
	Gdas     []int
}

// DeliveryPoint holds a lat/long and radius based delivery point
type DeliveryPoint struct {
	Latitude  float64
	Longitude float64
	Radius    int // likely km
}

// Lookup holds the details of historical messages the API returns
type Lookup struct {
	MessageID                  int
	CreatedOn                  time.Time
	MessageText                string
	DeliveryStatus             string // The status of a submitted message to the Iridium Burst system. This field will be empty if no delivery status check was performed from the portal.
	IridiumMessageID           string // A globally unique identifier returned from Iridium (to be used for support purposes)
	RepeatCount                int
	RepeatInterval             int
	BroadcastCoverageAreaCodes string // Broadcast coverage area code (can contain multiple values separated by commas)
	Latitude                   float64
	Longitude                  float64
	Radius                     int
	ServiceMask                int
	MessageCost                float64  // The cost based on the message size and delivery area in USD (This is an estimate. The actual charge for the message will be determined by Iridium)
	Sender                     string   // The sender's user name (email address)
	AccountID                  int      // The Account identifier used for this message
	CustomerNumber             int      // The customer account number
	Services                   []string // A list of services to which this message was targetted
}
