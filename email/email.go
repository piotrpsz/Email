package email

import (
	"fmt"
	"time"

	"Email/shared/tr"
)

type Parcel struct {
	header map[string]string
	body   []byte
}

func New(header map[string]string, body []byte) *Parcel {
	return &Parcel{
		header: header,
		body:   body,
	}
}

func (p *Parcel) Header() map[string]string {
	return p.header
}

func (p *Parcel) DisplayHeader() {
	for k, v := range p.header {
		fmt.Printf("%s: %s\n", k, v)
	}
}

func (p *Parcel) From() string {
	return p.HeaderValueFor("From")
}

func (p *Parcel) Subject() string {
	return p.HeaderValueFor("Subject")
}

func (p *Parcel) Date() time.Time {
	if value := p.HeaderValueFor("Date"); value != "" {
		if t, err := time.Parse(time.RFC1123Z, value); tr.IsOK(err) {
			return t
		}
	}
	return time.Time{}
}

func (p *Parcel) DeliveryDate() time.Time {
	if value := p.HeaderValueFor("Delivery-date"); value != "" {
		if t, err := time.Parse(time.RFC1123Z, value); tr.IsOK(err) {
			return t
		}
	}
	return time.Time{}
}

func (p *Parcel) HeaderValueFor(key string) string {
	if value, ok := p.header[key]; ok {
		return value
	}
	return ""
}
