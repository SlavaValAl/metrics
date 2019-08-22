package models

type Metric struct {
	Url       string `json:"url"`
	EventType string `json:"event_type"`
}

func (m *Metric) IsEmpty() bool {
	if m.EventType == "" && m.Url == "" {
		return true
	}
	return false
}
