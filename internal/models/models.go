package models

type RequiredInputs struct {
	Azimuth        string `json:"azimuth"`         // Azimuth angle (degrees)
	SystemCapacity string `json:"system_capacity"` // Nameplate capacity (KW)
	Losses         string `json:"losses"`          // System losses (percent)
	ArrayType      string `json:"array_type"`      /* Options : 0 - Fixed Open Rack;
	1 - Fixed Roof Mounted; 2 - 1-Axis; 3 - 1-Axis Backtracking; 4 - 2-Axis */
	ModuleType string `json:"module_type"` /* 0 - Standard; 1 - Premium; 2 - Thin Film*/
	Tilt       string `json:"tilt"`        // Tilt angle (degrees)
	Adress     string `json:"address"`     // Adress to use - required if latitude/longitude not specified
}

type OptionalInputs struct {
	Gcr string `json:"gcr"` // Ground coverage ratio. Range 0.01-0.99(default: 0.4)
	DcAcRatio string `json:"dc_ac_ratio"` // Type unsigned float(default: 1.2)
	InvEff string `json:"inv_eff"` //inverter efficiency value. Range 90-99.5(defalut: 96)
	Radius string `json:"radius"` // Radius for searching the closest climate data station(default: 100). Pass radius=0 to use closest
	Dataset string `json:"dataset"` /* Dataset to use (default: nsrdb) 
	Options: nsrdb, tmy2, tmy3, intl */
	Soiling string `json:"soiling"` /*Reduction in incident solar irradiance caused by dust or other seasonal soiling of the module surface .Specify a pipe-delimited array of 12 monthly values. Example: "12|4|45|23|9|99|67|12.54|54|9|0|7.6" */
	Albedo string `json:"albedo"` /*Ground reflectance. A value of 0 would mean that the ground is completely non-reflective, and a value of 1 would mean that it is completely reflective. Specify either a pipe-delimited array of 12 monthly values or a single value to be used for all months.*/
	Bifaciality string `json:"bifaciality"` /*The ratio of rear-side efficiency to front-side efficiency.Typical range 0.65 - 0.9 provided on the bifacial module datasheeet.(default:none)*/
	// Latitude string `json:"latitude"` // Latitude to use - required if adress not specified
	// Longitude string `json:"longitude"` // Longitude to use - required if adress not specified


}

type Outputs struct {
	DcMonthly []float64 `json:"dc_monthly"`
	AcMonthly []float64 `json:"ac_monthly"`
	AcAnnual  float64   `json:"ac_annual"`
}

type PowerEstimate struct {
	Inputs  RequiredInputs `json:"inputs"`
	Outputs Outputs        `json:"outputs"`
}

type User struct {
	ID       uint   `json:"id"`
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Email    string `json:"email" form:"email"`
}
