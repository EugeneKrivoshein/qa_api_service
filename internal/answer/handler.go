package answer

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

// CreateAnswer godoc
// @Summary Create answer
// @Description Добавляет ответ к вопросу
// @Tags answers
// @Accept json
// @Produce json
// @Param id path int true "Question ID"
// @Param answer body models.Answer true "Answer body"
// @Success 200 {object} models.Answer
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /questions/{id}/answers [post]
func (h *Handler) CreateAnswer(w http.ResponseWriter, r *http.Request) {
	questionIDStr := mux.Vars(r)["id"]
	questionID, _ := strconv.Atoi(questionIDStr)

	var a models.Answer
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		logrus.Warnf("Invalid request body for CreateAnswer: %v", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if a.Text == "" || a.UserID == "" {
		http.Error(w, "user_id and text required", http.StatusBadRequest)
		return
	}

	if err := h.Service.Create(&a, questionID); err != nil {
		logrus.Warnf("Failed to create answer for question ID=%d: %v", questionID, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	logrus.Infof("Answer created for question ID=%d, userID=%s", questionID, a.UserID)
	response.WriteJSON(w, a)
}

// GetAnswerByID godoc
// @Summary Get answer by ID
// @Description Получает ответ по ID
// @Tags answers
// @Produce json
// @Param id path int true "Answer ID"
// @Success 200 {object} models.Answer
// @Failure 404 {object} map[string]string
// @Router /answers/{id} [get]
func (h *Handler) GetAnswerByID(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	a, err := h.Service.GetByID(id)
	if err != nil {
		logrus.Warnf("Answer not found: ID=%d", id)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	response.WriteJSON(w, a)
	logrus.Infof("Fetched answer ID=%d", id)
}

// DeleteAnswer godoc
// @Summary Delete answer
// @Description Удаляет ответ
// @Tags answers
// @Param id path int true "Answer ID"
// @Success 204 {string} string "No Content"
// @Failure 404 {object} map[string]string
// @Router /answers/{id} [delete]
func (h *Handler) DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, _ := strconv.Atoi(idStr)

	if err := h.Service.Delete(id); err != nil {
		logrus.Warnf("Failed to delete answer ID=%d: %v", id, err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	logrus.Infof("Deleted answer ID=%d", id)
	w.WriteHeader(http.StatusNoContent)
}
