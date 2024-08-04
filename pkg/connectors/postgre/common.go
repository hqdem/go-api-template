package postgre

type Connector struct {
	// TODO: real postgre conn
}

func NewConnector() *Connector {
	return &Connector{}
}

func (c *Connector) PingDB() string {
	return "pong"
}
