package db

import "gorm.io/gorm"

type IDatabaseAdapter interface {
	Any(model interface{}, key, value string) (bool, error)
	Insert(model interface{}) (int, error)
}

type PostgresDbHandler struct {
	DB *gorm.DB
}

// below here is DB Model
type Customer struct {
	gorm.Model
	PhoneNumber string
}

func GetAllModels() []interface{} {
	models := []interface{}{
		&Customer{},
	}

	return models
}
