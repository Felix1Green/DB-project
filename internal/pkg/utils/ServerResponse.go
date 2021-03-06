package utils

import (
	"encoding/json"
	"github.com/Felix1Green/DB-project/internal/pkg/forum"
	"net/http"
	"strconv"
)

type ErrorResponse struct{
	Message string
}

func ServerErrorResponse(err error, statusCode int, w http.ResponseWriter){
	outputErr, _ := json.Marshal(ErrorResponse{
		Message: err.Error(),
	})
	_, _ = w.Write(outputErr)
}


func GetLimitSinceDescQueryParams(r *http.Request) (int, int, bool) {
	limit := r.URL.Query().Get(forum.LimitQueryName)
	since := r.URL.Query().Get(forum.SinceQueryName)
	desc := r.URL.Query().Get(forum.DescQueryName)
	parsedLimit, ok := strconv.Atoi(limit)
	if ok != nil {
		parsedLimit = 100
	}
	parsedSince, err := strconv.Atoi(since)
	if err != nil {
		parsedSince = -1
	}
	if desc != "" {
		if desc == "true"{
			return parsedLimit, parsedSince, true
		}
		return parsedLimit, parsedSince, false
	}
	return parsedLimit, parsedSince, false
}

func LimitIntSinceStringParams(r *http.Request)(int, string, bool){
	limit := r.URL.Query().Get(forum.LimitQueryName)
	since := r.URL.Query().Get(forum.SinceQueryName)
	desc := r.URL.Query().Get(forum.DescQueryName)
	parsedLimit, ok := strconv.Atoi(limit)
	if ok != nil {
		parsedLimit = 100
	}

	if desc != "" {
		if desc == "true"{
			return parsedLimit, since, true
		}
		return parsedLimit, since, false
	}
	return parsedLimit, since, false
}