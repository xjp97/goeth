package models

type LogError struct {
	ErrorID   uint8
	OrderHash [32]byte
}
