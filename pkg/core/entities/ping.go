package entities

type PingStatus struct {
	Status string
}

func NewPingStatus(status string) *PingStatus {
	return &PingStatus{
		Status: status,
	}
}
