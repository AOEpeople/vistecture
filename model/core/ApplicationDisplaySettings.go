package core

type ApplicationDisplaySettings struct {
	Rotate      bool   `json:"rotate" yaml:"rotate"`
	BorderColor string `json:"bordercolor,omitempty" yaml:"bordercolor,omitempty"`
}
