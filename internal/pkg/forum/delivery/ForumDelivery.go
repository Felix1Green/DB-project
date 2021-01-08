package delivery

import (
	"encoding/json"
	"github.com/Felix1Green/DB-project/internal/pkg/forum"
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/Felix1Green/DB-project/internal/pkg/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type ForumDelivery struct {
	UseCase forum.UseCase
}

func NewForumDelivery(usecase forum.UseCase) *ForumDelivery {
	return &ForumDelivery{
		UseCase: usecase,
	}
}

func (t *ForumDelivery) CreateForum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	input := new(models.ForumRequestInput)
	defer r.Body.Close()
	translationErr := json.NewDecoder(r.Body).Decode(&input)
	if translationErr != nil {
		utils.ServerErrorResponse(models.IncorrectInputParams, models.ErrorsStatusCodes[models.IncorrectInputParams], w)
		return
	}

	response, responseErr := t.UseCase.CreateForum(input)
	if response != nil && responseErr != nil{
		w.WriteHeader(409)
		outputBuf, _ := json.Marshal(response)
		_, _ = w.Write(outputBuf)
		return
	}else if responseErr != nil {
		w.WriteHeader(models.ErrorsStatusCodes[responseErr])
		utils.ServerErrorResponse(responseErr, models.ErrorsStatusCodes[responseErr], w)
		return
	}

	w.WriteHeader(201)
	outputBuf, _ := json.Marshal(response)

	_, _ = w.Write(outputBuf)
}

func (t *ForumDelivery) GetForum(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slug, ok := vars[forum.SlugPathName]
	if !ok {
		w.WriteHeader(models.ErrorsStatusCodes[models.IncorrectInputParams])
		utils.ServerErrorResponse(models.IncorrectInputParams, models.ErrorsStatusCodes[models.IncorrectInputParams], w)
		return
	}

	response, err := t.UseCase.GetForum(slug)
	if err != nil {
		w.WriteHeader(models.ErrorsStatusCodes[err])
		utils.ServerErrorResponse(err, models.ErrorsStatusCodes[err], w)
		return
	}


	outputBuf, _ := json.Marshal(response)
	_, _ = w.Write(outputBuf)
}

func (t *ForumDelivery) CreateForumThread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	slug, ok := vars[forum.SlugPathName]
	if !ok {
		w.WriteHeader(models.ErrorsStatusCodes[models.IncorrectInputParams])
		utils.ServerErrorResponse(models.IncorrectInputParams, models.ErrorsStatusCodes[models.IncorrectInputParams], w)
		return
	}

	input := new(models.ThreadRequestInput)
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(models.ErrorsStatusCodes[models.IncorrectInputParams])
		utils.ServerErrorResponse(models.IncorrectInputParams, models.ErrorsStatusCodes[models.IncorrectInputParams], w)
		return
	}

	response, err := t.UseCase.CreateForumThread(slug, input)
	if response != nil && err != nil{
		w.WriteHeader(409)
		outputBuf, _ := json.Marshal(response)
		_, _ = w.Write(outputBuf)
		return
	}else if err != nil {
		w.WriteHeader(models.ErrorsStatusCodes[err])
		utils.ServerErrorResponse(err, models.ErrorsStatusCodes[err], w)
		return
	}

	w.WriteHeader(201)
	outputBuf, _ := json.Marshal(response)
	_, _ = w.Write(outputBuf)
}

func (t *ForumDelivery) GetForumUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	slug, ok := mux.Vars(r)[forum.SlugPathName]
	if !ok {
		w.WriteHeader(405)
		return
	}

	limit, since, desc := utils.GetLimitSinceDescQueryParams(r)

	response, err := t.UseCase.GetForumUsers(slug, limit, since, desc)
	if err != nil {
		w.WriteHeader(models.ErrorsStatusCodes[err])
		utils.ServerErrorResponse(err, models.ErrorsStatusCodes[err], w)
		return
	}
	outputBuf, _ := json.Marshal(response)
	_, _ = w.Write(outputBuf)
}

func (t *ForumDelivery) GetForumThreads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	slug := mux.Vars(r)[forum.SlugPathName]
	limit, since, desc := utils.LimitIntSinceStringParams(r)

	response, err := t.UseCase.GetForumThreads(slug, limit, since, desc)
	if err != nil {
		w.WriteHeader(models.ErrorsStatusCodes[err])
		utils.ServerErrorResponse(err, models.ErrorsStatusCodes[err], w)
		return
	}

	outputBuf, _ := json.Marshal(response)
	_, _ = w.Write(outputBuf)
}
