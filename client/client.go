package client

type Client interface {
	Login() error
	Logout() error
	GetDevices() error
	GetDevice() error
	GetDeviceState() error
	SetDeviceState() error
}
