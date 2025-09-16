package types

type Message[T any] struct {
	Type    string `json:"type"`
	Payload T      `json:"payload"`
}
