package common

import "time"

const (
	ResourceName string = "hxm.com/gopher"
	DevicePath string = "/etc/gophers"
	DeviceSocket string = "gopher.sock"
	ConnectTimeout time.Duration = 5 * time.Second
)