package gameconst

import "time"

// service const
const (
	SendBufferSize  = 10
	ReadTimeoutSec  = 6 * time.Second
	WriteTimeoutSec = 3 * time.Second
)
