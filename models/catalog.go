package models

// ServerCatalog represents a server in the catalog
// @Description Server catalog information
type ServerCatalog struct {
	ID       uint    `json:"-" gorm:"primaryKey;autoIncrement;column:id" swaggerignore:"true"`
	Model    string  `json:"model" gorm:"type:varchar(128);not null;column:model" example:"HP DL120G7Intel G850"`
	RamSize  int     `json:"ram_size" gorm:"not null;column:ram_size" example:"4"`
	RamType  int     `json:"ram_type" gorm:"not null;column:ram_type;foreignKey:RamType;references:ID" example:"1"`
	HDDSize  int     `json:"hdd_size" gorm:"not null;column:hdd_size" example:"1000"`
	HDDCount int     `json:"hdd_count" gorm:"not null;column:hdd_count" example:"4"`
	HDDType  int     `json:"hdd_type" gorm:"not null;column:hdd_type;foreignKey:HDDType;references:ID" example:"2"`
	Location string  `json:"location" gorm:"type:varchar(128);not null;column:location" example:"AmsterdamAMS-01"`
	Price    float64 `json:"price" gorm:"type:decimal(20,2);unsigned;not null;column:price" example:"39.99"`
	Currency int     `json:"currency" gorm:"not null;column:currency;foreignKey:Currency;references:ID" example:"1"`
}

func (sc *ServerCatalog) TableName() string {
	return "server_catalog"
}
