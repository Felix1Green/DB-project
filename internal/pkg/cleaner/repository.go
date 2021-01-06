package cleaner

import "github.com/Felix1Green/DB-project/internal/pkg/models"

type Repository interface{
	Status()(*models.Status, error)
	Clear() error
}
