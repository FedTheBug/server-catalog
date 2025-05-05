package models

type ServerCatalog struct {
	ID       uint
	Model    string
	RamSize  uint
	RamType  uint
	HDDSize  uint
	HDDCount uint
	HDDType  uint
	Location string
	Price    float64
	Currency int
}
