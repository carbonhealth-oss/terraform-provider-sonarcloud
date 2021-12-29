package models

type GetComponentResponse struct {
	Key             string `form:"key,omitempty"`
	Organization    string `form:"organization,omitempty"`
	Id              string `form:"id,omitempty"`
	Name            string `form:"Name,omitempty"`
	AutoscanEnabled bool   `form:"autoscanEnabled,omitempty"`
	Visibility      string `json:"visibility,omitempty"`
	QualityGate     struct {
		Key       int    `form:"key,omitempty"`
		Name      string `form:"Name,omitempty"`
		IsDefault bool   `form:"isDefault,omitempty"`
	} `json:"qualityGate,omitempty"`
}
