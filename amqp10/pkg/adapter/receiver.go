package adapter

import (
	"context"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	ceamqp "github.com/cloudevents/sdk-go/pkg/cloudevents/transport/amqp"
	"github.com/knative/eventing-contrib/amqp10/pkg/config"
	"pack.ag/amqp"
)

type Receiver struct {
	conf  *config.Config
	link  *amqp.Receiver
	codec *ceamqp.Codec
}

func (r *Receiver) Receive(ctx context.Context) (event *cloudevents.Event, err error) {
	if m, err := r.link.Receive(ctx); err == nil {

		aMsg := &ceamqp.Message{
			ContentType:           m.Properties.ContentType,
			ApplicationProperties: m.ApplicationProperties,
			Body:                  m.GetData(),
		}

		if ev, err := r.codec.Decode(aMsg); err != nil {
			return nil, err
		} else {
			return ev, nil
		}

	}

	return nil, err
}

func NewReceiver(c *config.Config, src string, opts []amqp.LinkOption) *Receiver {

	r := &Receiver{
		conf:  c,
		codec: &ceamqp.Codec{},
	}

	o := append(opts, amqp.LinkSourceAddress(src))
	if l, err := c.Session.NewReceiver(o...); err == nil {
		r.link = l
	}

	return r
}
