package comfortcloud

type Response struct {
	UIFlg      bool    `json:"uiFlg"`
	GroupCount int     `json:"groupCount"`
	GroupList  []Group `json:"groupList"`
}

type Group struct {
	GroupID     int      `json:"groupId"`
	GroupName   string   `json:"groupName"`
	PairingList []string `json:"pairingList"` // Assuming pairingList is a list of strings
	DeviceList  []Device `json:"deviceList"`
}

type Device struct {
	DeviceGuid         string     `json:"deviceGuid"`
	DeviceType         string     `json:"deviceType"`
	DeviceName         string     `json:"deviceName"`
	Permission         int        `json:"permission"`
	TemperatureUnit    int        `json:"temperatureUnit"`
	SummerHouse        int        `json:"summerHouse"`
	NanoeStandAlone    bool       `json:"nanoeStandAlone"`
	AutoMode           bool       `json:"autoMode"`
	ModeAvlList        ModeAvl    `json:"modeAvlList"`
	Parameters         Parameters `json:"parameters"`
	DeviceModuleNumber string     `json:"deviceModuleNumber"`
	DeviceHashGuid     string     `json:"deviceHashGuid"`
	ModelVersion       int        `json:"modelVersion"`
	CoordinableFlg     bool       `json:"coordinableFlg"`
}

type ModeAvl struct {
	AutoMode int `json:"autoMode"`
}

type Parameters struct {
	Operate         Power            `json:"operate"`
	OperationMode   OperationMode    `json:"operationMode"`
	TemperatureSet  int              `json:"temperatureSet"`
	FanSpeed        FanSpeed         `json:"fanSpeed"`
	FanAutoMode     AirSwingAutoMode `json:"fanAutoMode"`
	AirSwingLR      AirSwingLR       `json:"airSwingLR"`
	AirSwingUD      AirSwingUD       `json:"airSwingUD"`
	EcoFunctionData int              `json:"ecoFunctionData"`
	EcoMode         EcoMode          `json:"ecoMode"`
	EcoNavi         int              `json:"ecoNavi"`
	Nanoe           NanoeMode        `json:"nanoe"`
	IAuto           int              `json:"iAuto"`
	AirDirection    int              `json:"airDirection"`
	LastSettingMode int              `json:"lastSettingMode"`
}
