package comfortcloud

import "fmt"

const (
	AppClientId      = "Xmy6xIYIitMxngjB2rHvlm6HSDNnaMJx"
	Auth0Client      = "eyJuYW1lIjoiQXV0aDAuQW5kcm9pZCIsImVudiI6eyJhbmRyb2lkIjoiMzAifSwidmVyc2lvbiI6IjIuOS4zIn0="
	RedirectUri      = "panasonic-iot-cfc://authglb.digital.panasonic.com/android/com.panasonic.ACCsmart/callback"
	BasePathAuth     = "https://authglb.digital.panasonic.com"
	BasePathAcc      = "https://accsmart.panasonic.com"
	XAppVersion      = "1.22.0"
	OAuthScopes      = "openid offline_access comfortcloud.control a2w.control"
	OAuthAudienceURL = "https://digital.panasonic.com/%s/api/v1/"
)

var OAuthAudience = fmt.Sprintf(OAuthAudienceURL, AppClientId)

type Power int

const (
	PowerOff Power = iota
	PowerOn
)

func (p Power) String() string {
	switch p {
	case PowerOff:
		return "Off"
	case PowerOn:
		return "On"
	default:
		return "Unknown"
	}
}

type OperationMode int

const (
	OperationModeAuto OperationMode = iota
	OperationModeDry
	OperationModeCool
	OperationModeHeat
	OperationModeFan
)

type AirSwingUD int

const (
	AirSwingUDAuto AirSwingUD = iota - 1
	AirSwingUDUp
	_
	AirSwingUDUpMid
	AirSwingUDMid
	AirSwingUDDownMid
	AirSwingUDDown
	AirSwingUDSwing
)

type AirSwingLR int

const (
	AirSwingLRAuto AirSwingLR = iota - 1
	AirSwingLRLeft
	AirSwingLRMid
	_
	AirSwingLRRightMid
	AirSwingLRRight
)

type EcoMode int

const (
	EcoModeAuto EcoMode = iota
	EcoModePowerful
	EcoModeQuiet
)

type AirSwingAutoMode int

const (
	AirSwingAutoModeDisabled AirSwingAutoMode = iota
	AirSwingAutoModeBoth
	AirSwingAutoModeAirSwingUD
	AirSwingAutoModeAirSwingLR
)

type FanSpeed int

const (
	FanSpeedAuto FanSpeed = iota
	FanSpeedLow
	FanSpeedLowMid
	FanSpeedMid
	FanSpeedHighMid
	FanSpeedHigh
)

type DataMode int

const (
	DataModeDay DataMode = iota
	DataModeWeek
	DataModeMonth
	_
	DataModeYear
)

var DataModeMap = map[string]DataMode{
	"Day":   DataModeDay,
	"Week":  DataModeWeek,
	"Month": DataModeMonth,
	"Year":  DataModeYear,
}

type NanoeMode int

const (
	NanoeModeUnavailable NanoeMode = iota
	NanoeModeOff
	NanoeModeOn
	NanoeModeModeG
	NanoeModeAll
)
