package yandex

type StateType struct {
	Instance string       `json:"instance"`
	Value    *interface{} `json:"value,omitempty"`
}

type Action struct {
	Type  string    `json:"type"`
	State StateType `json:"state"`
}
