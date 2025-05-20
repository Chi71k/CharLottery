package natswrap

import "github.com/nats-io/nats.go"

type NatsClient struct {
    conn *nats.Conn
}

func NewNatsClient(url string) (*NatsClient, error) {
    nc, err := nats.Connect(url)
    if err != nil {
        return nil, err
    }
    return &NatsClient{conn: nc}, nil
}

func (n *NatsClient) Publish(subject string, data []byte) error {
    return n.conn.Publish(subject, data)
}

func (n *NatsClient) Subscribe(subject string, handler func(msg []byte)) error {
    _, err := n.conn.Subscribe(subject, func(m *nats.Msg) {
        handler(m.Data)
    })
    return err
}