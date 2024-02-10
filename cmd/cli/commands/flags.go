package commands

type InputFlags struct {
	ConfigPath  string  `json:"config_path"`
	Verbose     bool    `json:"verbose"`
	SitHeight   float64 `json:"sit_height"`
	StandHeight float64 `json:"stand_height"`
	Position    float64 `json:"position"`
}
