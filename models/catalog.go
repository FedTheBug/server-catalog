package models

type ServerCatalog struct {
	ID       uint    `json:"-" gorm:"primaryKey;autoIncrement;column:id"`
	Model    string  `json:"model" gorm:"type:varchar(128);not null;column:model"`
	RamSize  int     `json:"ram_size" gorm:"not null;column:ram_size"`
	RamType  int     `json:"ram_type" gorm:"not null;column:ram_type;foreignKey:RamType;references:ID"`
	HDDSize  int     `json:"hdd_size" gorm:"not null;column:hdd_size"`
	HDDCount int     `json:"hdd_count" gorm:"not null;column:hdd_count"`
	HDDType  int     `json:"hdd_type" gorm:"not null;column:hdd_type;foreignKey:HDDType;references:ID"`
	Location string  `json:"location" gorm:"type:varchar(128);not null;column:location"`
	Price    float64 `json:"price" gorm:"type:decimal(20,2);unsigned;not null;column:price"`
	Currency int     `json:"currency" gorm:"not null;column:currency;foreignKey:Currency;references:ID"`
}

func (sc *ServerCatalog) TableName() string {
	return "server_catalog"
}
