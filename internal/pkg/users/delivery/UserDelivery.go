package delivery

import (
	"encoding/json"
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/Felix1Green/DB-project/internal/pkg/users"
	"github.com/Felix1Green/DB-project/internal/pkg/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type UserDelivery struct{
	UseCase users.UseCase
}

func NewUserDelivery(usage users.UseCase) *UserDelivery{
	return &UserDelivery{
		UseCase: usage,
	}
}

func (t *UserDelivery) CreateUser(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		w.WriteHeader(405)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	nickname := vars[users.NickNamePath]
	CreateInput := new(models.UserRequestBody)
	translationErr := json.NewDecoder(r.Body).Decode(&CreateInput)
	if translationErr != nil{
		utils.ServerErrorResponse(translationErr, 400, w)
		return
	}

	response, err := t.UseCase.CreateUser(nickname, CreateInput)
	if response != nil && err != nil{
		w.WriteHeader(409)
		outputBuf, _ := json.Marshal(response)
		_, _ = w.Write(outputBuf)
		return
	}
	if err != nil{
		utils.ServerErrorResponse(err, models.ErrorsStatusCodes[err], w)
		return
	}
	w.WriteHeader(201)
	outputBuf, _ := json.Marshal((*response)[0])
	_, _ = w.Write(outputBuf)
	return
}

func (t *UserDelivery) GetUser(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	nickname := vars[users.NickNamePath]
	response, err := t.UseCase.GetProfile(nickname)
	if err != nil{
		w.WriteHeader(models.ErrorsStatusCodes[err])
		utils.ServerErrorResponse(err, models.ErrorsStatusCodes[err], w)
		return
	}

	outputBuf, _ := json.Marshal(response)
	_,_ =w. Write(outputBuf)
	return
}

func (t *UserDelivery) UpdateUser(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	nickname := vars[users.NickNamePath]
	CreateInput := new(models.UserRequestBody)
	translationErr := json.NewDecoder(r.Body).Decode(&CreateInput)
	if translationErr != nil{
		utils.ServerErrorResponse(translationErr, 400, w)
		return
	}

	response, err := t.UseCase.UpdateProfile(nickname, CreateInput)
	if err != nil{
		w.WriteHeader(models.ErrorsStatusCodes[err])
		utils.ServerErrorResponse(err, models.ErrorsStatusCodes[err], w)
		return
	}

	outputBuf, _ := json.Marshal(response)
	_, _ = w.Write(outputBuf)
	return
}