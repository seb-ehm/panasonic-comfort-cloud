package comfortcloud

type Device struct {
	ID                 string `json:"id"`
	Name               string `json:"name"`
	Group              string `json:"group"`
	Model              string `json:"model"`
	DeviceGuid         string `json:"deviceGuid"`
	DeviceHashGuid     string `json:"deviceHashGuid"`
	DeviceName         string `json:"deviceName"`
	DeviceModuleNumber string `json:"deviceModuleNumber"`
}
