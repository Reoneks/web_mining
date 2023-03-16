package settings

import "time"

const (
	WSTicker               = time.Minute
	WaiterTTL              = 30 * time.Minute
	EventsWithoutWaiterTTL = 15 * time.Minute
	ReadDeadline           = 2 * time.Second
)
