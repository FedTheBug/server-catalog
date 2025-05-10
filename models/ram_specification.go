package models

type RamSpec struct {
	ID   uint
	Type string
}

func (rs *RamSpec) TableName() string {
	return "ram_spec"
}
