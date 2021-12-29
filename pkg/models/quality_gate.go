package models

type ListQualityGates struct {
	QualityGate []QualityGate `json:"qualitygates,omitempty"`
}

type QualityGate struct {
	Id        int    `form:"id,omitempty"`
	Name      string `form:"Name,omitempty"`
	IsDefault bool   `form:"isDefault,omitempty"`
	IsBuiltIn bool   `form:"isBuiltIn,omitempty"`
}
