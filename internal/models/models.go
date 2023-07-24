package models

type RequiredInputs struct {
	Azimuth        string `json:"azimuth"`
	SystemCapacity string `json:"system_capacity"`
	Losses         string `json:"losses"`
	ArrayType      string `json:"array_type"`
	ModuleType     string `json:"module_type"`
	Tilt           string `json:"tilt"`
	Adress         string `json:"address"`
}

type Outputs struct {
	DcMonthly []float64 `json:"dc_monthly"`
	AcMonthly []float64 `json:"ac_monthly"`
	AcAnnual  float64   `json:"ac_annual"`
}

type PowerEstimate struct {
	Inputs  RequiredInputs  `json:"inputs"`
	Outputs Outputs `json:"outputs"`
}

type User struct {
	ID int 
	Username string
	Email string
	Password string
}
