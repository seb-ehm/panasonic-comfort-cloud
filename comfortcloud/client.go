package comfortcloud

type Client struct {
	auth          *Authentication
	groups        []map[string]interface{}
	devices       []Device
	deviceIndexer map[string]string // Maps device ID to device GUID
	accClientID   string
}

func NewApiClient(auth *Authentication, raw bool) *Client {
	return &Client{
		auth:          auth,
		deviceIndexer: make(map[string]string),
	}
}

func (c Client) Login() error {
	//TODO implement me
	panic("implement me")
}

func (c Client) Logout() error {
	//TODO implement me
	panic("implement me")
}

func (c Client) GetDevices() error {
	//TODO implement me
	panic("implement me")
}

func (c Client) GetDevice() error {
	//TODO implement me
	panic("implement me")
}

func (c Client) GetDeviceState() error {
	//TODO implement me
	panic("implement me")
}

func (c Client) SetDeviceState() error {
	//TODO implement me
	panic("implement me")
}
