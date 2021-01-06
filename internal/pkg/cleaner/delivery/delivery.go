package delivery

import (
	"encoding/json"
	"github.com/Felix1Green/DB-project/internal/pkg/cleaner"
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"net/http"
)

type CleanerDelivery struct {
	Repository cleaner.Repository
}

func NewCleanerDelivery(repository cleaner.Repository) *CleanerDelivery{
	return &CleanerDelivery{
		Repository: repository,
	}
}

func (t *CleanerDelivery) GetStatus(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		w.WriteHeader(405)
		return
	}

	resp, err := t.Repository.Status()
	if err != nil{
		w.WriteHeader(models.ErrorsStatusCodes[err])
		outputBuf, _ := json.Marshal(models.ErrorMessage{Message: err.Error()})
		w.Write(outputBuf)
		return
	}

	outputBuf, _ := json.Marshal(resp)
	w.Write(outputBuf)
}

func (t *CleanerDelivery) ClearDB(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		w.WriteHeader(405)
		return
	}

	err := t.Repository.Clear()
	if err != nil{
		w.WriteHeader(models.ErrorsStatusCodes[err])
		outputBuf, _ := json.Marshal(models.ErrorMessage{Message: err.Error()})
		w.Write(outputBuf)
		return
	}
}