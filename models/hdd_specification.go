package models

type HDDSpec struct {
	ID   uint
	Type string
}

func (hs *HDDSpec) TableName() string {
	return "hdd_spec"
}
