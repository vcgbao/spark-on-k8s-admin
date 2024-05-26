package models

type StatusCode int

const (
	SUCCESS = 0
	ERROR   = 1
)

type Response[T any] struct {
	StatusCode StatusCode `json:"statusCode"`
	Message    string     `json:"message"`
	Data       T          `json:"data"`
}
