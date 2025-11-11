package question

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/EugeneKrivoshein/qa_api_service/models"
	"github.com/EugeneKrivoshein/qa_api_service/pkg/response"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Service *Service
}

func (h *Handler) GetQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := h.Service.GetAll()
	if err != nil {
		logrus.Errorf("Failed to get questions: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response.WriteJSON(w, questions)
	logrus.Infof("Fetched %d questions", len(questions))
}

func (h *Handler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var q models.Question
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		logrus.Warnf("Invalid request body for CreateQuestion: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if q.Text == "" {
		http.Error(w, "text field required", http.StatusBadRequest)
		return
	}

	if err := h.Service.Create(&q); err != nil {
		logrus.Errorf("Failed to create question: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logrus.Infof("Question created: ID=%d", q.ID)
	response.WriteJSON(w, q)
}

func (h *Handler) GetQuestionByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	data, err := h.Service.GetByID(id)
	if err != nil {
		logrus.Warnf("Question not found: ID=%d", id)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	response.WriteJSON(w, data)
	logrus.Infof("Fetched question ID=%d with %d answers", id, len(data["answers"].([]models.Answer)))
}

func (h *Handler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	if err := h.Service.Delete(id); err != nil {
		logrus.Warnf("Failed to delete question ID=%d: %v", id, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	logrus.Infof("Deleted question ID=%d and its answers", id)
	w.WriteHeader(http.StatusNoContent)
}
