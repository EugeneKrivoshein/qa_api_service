package api

import (
	"net/http"

	"github.com/EugeneKrivoshein/qa_api_service/internal/answer"
	"github.com/EugeneKrivoshein/qa_api_service/internal/question"
	"github.com/gorilla/mux"
)

func NewRouter(qHandler *question.Handler, aHandler *answer.Handler) *mux.Router {
	r := mux.NewRouter()

	//question
	r.HandleFunc("/questions", qHandler.GetQuestions).Methods(http.MethodGet)
	r.HandleFunc("/questions", qHandler.CreateQuestion).Methods(http.MethodPost)
	r.HandleFunc("/questions/{id:[0-9]+}", qHandler.GetQuestionByID).Methods(http.MethodGet)
	r.HandleFunc("/questions/{id:[0-9]+}", qHandler.DeleteQuestion).Methods(http.MethodDelete)

	//answer
	r.HandleFunc("/questions/{id:[0-9]+}/answers", aHandler.CreateAnswer).Methods(http.MethodPost)
	r.HandleFunc("/answers/{id:[0-9]+}", aHandler.GetAnswerByID).Methods(http.MethodGet)
	r.HandleFunc("/answers/{id:[0-9]+}", aHandler.DeleteAnswer).Methods(http.MethodDelete)

	return r
}
