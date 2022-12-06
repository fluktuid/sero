package util

const (
	StatusUp Status = iota
	StatusDownscaling
	StatusDown
	StatusUpscaling
)

type Status byte
