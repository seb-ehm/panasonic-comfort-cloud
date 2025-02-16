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
	Operate           Power            `json:"operate"`
	OperationMode     OperationMode    `json:"operationMode"`
	TemperatureSet    float64          `json:"temperatureSet"`
	FanSpeed          FanSpeed         `json:"fanSpeed"`
	FanAutoMode       AirSwingAutoMode `json:"fanAutoMode"`
	AirSwingLR        AirSwingLR       `json:"airSwingLR"`
	AirSwingUD        AirSwingUD       `json:"airSwingUD"`
	EcoFunctionData   int              `json:"ecoFunctionData"`
	EcoMode           EcoMode          `json:"ecoMode"`
	EcoNavi           int              `json:"ecoNavi"`
	Nanoe             NanoeMode        `json:"nanoe"`
	IAuto             int              `json:"iAuto"`
	AirDirection      int              `json:"airDirection"`
	LastSettingMode   int              `json:"lastSettingMode"`
	InsideCleaning    int              `json:"insideCleaning"`
	Fireplace         int              `json:"fireplace"`
	InsideTemperature int              `json:"insideTemperature"`
	OutTemperature    int              `json:"outTemperature"`
	AirQuality        int              `json:"airQuality"`
}

type ParameterOptions struct {
	Operate           *Power            `json:"operate,omitempty"`
	OperationMode     *OperationMode    `json:"operationMode,omitempty"`
	TemperatureSet    *int              `json:"temperatureSet,omitempty"`
	FanSpeed          *FanSpeed         `json:"fanSpeed,omitempty"`
	FanAutoMode       *AirSwingAutoMode `json:"fanAutoMode,omitempty"`
	AirSwingLR        *AirSwingLR       `json:"airSwingLR,omitempty"`
	AirSwingUD        *AirSwingUD       `json:"airSwingUD,omitempty"`
	EcoFunctionData   *int              `json:"ecoFunctionData,omitempty"`
	EcoMode           *EcoMode          `json:"ecoMode,omitempty"`
	EcoNavi           *int              `json:"ecoNavi,omitempty"`
	Nanoe             *NanoeMode        `json:"nanoe,omitempty"`
	IAuto             *int              `json:"iAuto,omitempty"`
	AirDirection      *int              `json:"airDirection,omitempty"`
	LastSettingMode   *int              `json:"lastSettingMode,omitempty"`
	InsideCleaning    *int              `json:"insideCleaning,omitempty"`
	Fireplace         *int              `json:"fireplace,omitempty"`
	InsideTemperature *int              `json:"insideTemperature,omitempty"`
	OutTemperature    *int              `json:"outTemperature,omitempty"`
	AirQuality        *int              `json:"airQuality,omitempty"`
}

/*full response for the GetDevice API:
type FullGetDeviceAnswer struct {
	Timestamp       int64 `json:"timestamp"`
	Permission      int   `json:"permission"`
	SummerHouse     int   `json:"summerHouse"`
	IAutoX          bool  `json:"iAutoX"`
	Nanoe           bool  `json:"nanoe"`
	NanoeStandAlone bool  `json:"nanoeStandAlone"`
	AutoMode        bool  `json:"autoMode"`
	HeatMode        bool  `json:"heatMode"`
	FanMode         bool  `json:"fanMode"`
	DryMode         bool  `json:"dryMode"`
	CoolMode        bool  `json:"coolMode"`
	EcoNavi         bool  `json:"ecoNavi"`
	PowerfulMode    bool  `json:"powerfulMode"`
	QuietMode       bool  `json:"quietMode"`
	AirSwingLR      bool  `json:"airSwingLR"`
	AutoSwingUD     bool  `json:"autoSwingUD"`
	EcoFunction     int   `json:"ecoFunction"`
	TemperatureUnit int   `json:"temperatureUnit"`
	ModeAvlList     struct {
		AutoMode int `json:"autoMode"`
	} `json:"modeAvlList"`
	NanoeList struct {
		VisualizationShow int `json:"visualizationShow"`
	} `json:"nanoeList"`
	ClothesDrying  bool `json:"clothesDrying"`
	InsideCleaning bool `json:"insideCleaning"`
	Fireplace      bool `json:"fireplace"`
	PairedFlg      bool `json:"pairedFlg"`
	Parameters     struct {
		EcoFunctionData   int `json:"ecoFunctionData"`
		InsideCleaning    int `json:"insideCleaning"`
		Fireplace         int `json:"fireplace"`
		LastSettingMode   int `json:"lastSettingMode"`
		Operate           int `json:"operate"`
		OperationMode     int `json:"operationMode"`
		TemperatureSet    int `json:"temperatureSet"`
		FanSpeed          int `json:"fanSpeed"`
		FanAutoMode       int `json:"fanAutoMode"`
		AirSwingLR        int `json:"airSwingLR"`
		AirSwingUD        int `json:"airSwingUD"`
		EcoMode           int `json:"ecoMode"`
		EcoNavi           int `json:"ecoNavi"`
		Nanoe             int `json:"nanoe"`
		IAuto             int `json:"iAuto"`
		AirDirection      int `json:"airDirection"`
		InsideTemperature int `json:"insideTemperature"`
		OutTemperature    int `json:"outTemperature"`
		AirQuality        int `json:"airQuality"`
	} `json:"parameters"`
	DeviceNanoe int `json:"deviceNanoe"`
}*/
