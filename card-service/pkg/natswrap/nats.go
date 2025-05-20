package natswrap

type Publisher interface {
    Publish(subject string, data []byte) error
}

type Subscriber interface {
    Subscribe(subject string, handler func(msg []byte)) error
}