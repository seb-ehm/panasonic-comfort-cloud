package comfortcloud

type DeviceOption func(*ParameterOptions)

func WithTemperature(temperature int) DeviceOption {
	return func(o *ParameterOptions) {
		o.TemperatureSet = &temperature
	}
}

func WithPower(power Power) DeviceOption {
	return func(o *ParameterOptions) {
		o.Operate = &power
	}
}
