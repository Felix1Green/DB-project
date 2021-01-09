package delivery

import (
	"encoding/json"
	"github.com/Felix1Green/DB-project/internal/pkg/models"
	"github.com/Felix1Green/DB-project/internal/pkg/post"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"strings"
)

type PostDelivery struct {
	UseCase post.UseCase
}

func NewPostDelivery(usecase post.UseCase) *PostDelivery{
	return &PostDelivery{
		UseCase: usecase,
	}
}


func getRelatedParameter(r *http.Request) (bool, bool, bool){
	related := r.URL.Query().Get(post.RelatedPathName)
	relatedArr := strings.Split(related, ",")
	user, forum,thread := false,false,false
	for _, val := range relatedArr{
		if val == "user"{
			user = true
			continue
		}else if val == "forum"{
			forum = true
			continue
		}else if val == "thread"{
			thread = true
			continue
		}
	}
	return user, forum, thread
}

func (t *PostDelivery) GetPostDetails(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet{
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	postID, castErr := strconv.Atoi(mux.Vars(r)[post.PathPostName])
	if castErr != nil{
		w.WriteHeader(400)
		return
	}
	user, forum, thread := getRelatedParameter(r)
	resp, err := t.UseCase.GetPostDetails(uint64(postID), user, thread, forum)
	if err != nil{
		w.WriteHeader(models.ErrorsStatusCodes[err])
		outputBuf, _ := json.Marshal(models.ErrorMessage{
			Message: err.Error(),
		})
		w.Write(outputBuf)
		return
	}

	outputBuf, _ := json.Marshal(resp)
	w.Write(outputBuf)
}

func (t *PostDelivery) UpdatePostMessage(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		w.WriteHeader(405)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	postID, castErr := strconv.Atoi(mux.Vars(r)[post.PathPostName])
	if castErr != nil{
		w.WriteHeader(400)
		return
	}
	input := new(models.PostUpdateRequestInput)
	defer r.Body.Close()
	decodeErr := json.NewDecoder(r.Body).Decode(&input)
	if decodeErr != nil{
		w.WriteHeader(400)
		return
	}

	resp, err := t.UseCase.UpdatePost(uint64(postID), input)
	if err != nil{
		w.WriteHeader(models.ErrorsStatusCodes[err])
		outputBuf, _ := json.Marshal(err)
		w.Write(outputBuf)
		return
	}
	output, _ := json.Marshal(resp)
	w.Write(output)
}