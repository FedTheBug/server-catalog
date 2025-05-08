package models

type ServerCatalog struct {
	ID       uint
	Model    string
	RamSize  int
	RamType  int
	HDDSize  int
	HDDCount int
	HDDType  int
	Location string
	Price    float64
	Currency int
}

func (sc *ServerCatalog) TableName() string {
	return "server_catalog"
}
