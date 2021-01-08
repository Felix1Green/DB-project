package Delivery

import (
	"encoding/json"
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/Felix1Green/DB-project/internal/pkg/thread"
	"github.com/Felix1Green/DB-project/internal/pkg/utils"
	"github.com/gorilla/mux"
	"net/http"
)

type ThreadDelivery struct {
	UseCase thread.UseCase
}

func NewThreadDelivery(uc thread.UseCase) *ThreadDelivery{
	return &ThreadDelivery{
		UseCase: uc,
	}
}

func (t *ThreadDelivery) CreateNewPosts(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		w.WriteHeader(405)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	slug := mux.Vars(r)[thread.PathThreadName]

	defer r.Body.Close()
	input := make([]models.PostCreateRequestInput, 0)
	decodeErr := json.NewDecoder(r.Body).Decode(&input)
	if decodeErr != nil{
		w.WriteHeader(400)
		return
	}

	resp, err := t.UseCase.CreatePosts(slug, &input)
	if err != nil{
		w.WriteHeader(models.ErrorsStatusCodes[err])
		outputBuf, _ := json.Marshal(models.ErrorMessage{Message: err.Error()})
		_, _ = w.Write(outputBuf)
		return
	}

	w.WriteHeader(201)
	outputBuf, _ := json.Marshal(resp)
	_, _ = w.Write(outputBuf)
}

func (t *ThreadDelivery) GetThreadDetails(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	threadSlug := mux.Vars(r)[thread.PathThreadName]

	resp, err :=t.UseCase.GetThreadDetails(threadSlug)
	if err != nil{
		w.WriteHeader(models.ErrorsStatusCodes[err])
		outputBuf, _ := json.Marshal(models.ErrorMessage{Message: err.Error()})
		_, _ = w.Write(outputBuf)
		return
	}

	outputBuf, _ := json.Marshal(resp)
	_, _ = w.Write(outputBuf)
}

func (t *ThreadDelivery) UpdateThreadInfo(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	threadSlug := mux.Vars(r)[thread.PathThreadName]
	input := new(models.ThreadUpdateInput)
	decodeErr := json.NewDecoder(r.Body).Decode(&input)
	if decodeErr != nil{
		w.WriteHeader(400)
		return
	}

	resp, err := t.UseCase.UpdateThreadDetails(threadSlug, input)
	if err != nil{
		w.WriteHeader(models.ErrorsStatusCodes[err])
		outputBuf, _ := json.Marshal(models.ErrorMessage{Message: err.Error()})
		w.Write(outputBuf)
		return
	}
	outputBuf, _ := json.Marshal(resp)
	w.Write(outputBuf)
}

func (t *ThreadDelivery) GetThreadPosts(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodGet{
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	threadSlug := mux.Vars(r)[thread.PathThreadName]

	limit, since, desc := utils.GetLimitSinceDescQueryParams(r)
	sort := r.URL.Query().Get(thread.QuerySortName)
	resp, err := t.UseCase.GetThreadPosts(threadSlug, limit, int64(since), sort, desc)
	if err != nil {
		w.WriteHeader(models.ErrorsStatusCodes[err])
		outputBuf, _ := json.Marshal(models.ErrorMessage{Message: err.Error()})
		w.Write(outputBuf)
		return
	}
	outputBuf, _ := json.Marshal(resp)
	w.Write(outputBuf)
}

func (t *ThreadDelivery) SetThreadVote(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	threadSlug := mux.Vars(r)[thread.PathThreadName]

	input := new(models.ThreadVoteInput)
	defer r.Body.Close()
	decodeErr := json.NewDecoder(r.Body).Decode(&input)
	if decodeErr != nil{
		w.WriteHeader(400)
		return
	}

	resp, err := t.UseCase.SetThreadVote(threadSlug, *input)
	if err != nil{
		w.WriteHeader(models.ErrorsStatusCodes[err])
		outputBuf, _ := json.Marshal(models.ErrorMessage{Message: err.Error()})
		w.Write(outputBuf)
		return
	}

	outputBuf, _ := json.Marshal(resp)
	w.Write(outputBuf)
}