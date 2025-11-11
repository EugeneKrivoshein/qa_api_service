package tests

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/EugeneKrivoshein/qa_api_service/models"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const baseURL = "http://localhost:8080"

func TestCreateAndGetQuestion(t *testing.T) {
	body := map[string]string{"text": "Какой язык самый быстрый?"}
	bodyJSON, _ := json.Marshal(body)

	resp, err := http.Post(baseURL+"/questions", "application/json", bytes.NewBuffer(bodyJSON))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	defer resp.Body.Close()

	var created map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&created)
	id := int(created["id"].(float64))
	assert.NotZero(t, id)

	//получить список всех вопросов
	resp2, err := http.Get(baseURL + "/questions")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	defer resp2.Body.Close()

	var list []map[string]interface{}
	json.NewDecoder(resp2.Body).Decode(&list)
	assert.GreaterOrEqual(t, len(list), 1)

	//получить конкретный вопрос
	resp3, err := http.Get(fmt.Sprintf("%s/questions/%d", baseURL, id))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp3.StatusCode)

	defer resp3.Body.Close()

	var q map[string]interface{}
	json.NewDecoder(resp3.Body).Decode(&q)
	assert.Equal(t, created["text"], q["question"].(map[string]interface{})["text"])
}

func TestAddAnswerAndGetIt(t *testing.T) {
	clearTables(t)

	// Создаем вопрос
	qID := createQuestion(t, "Что ты думаешь о Go?")

	// Добавляем ответ к этому вопросу
	body := map[string]interface{}{
		"user_id": "8d290db8-5b9b-4b39-8a5f-9fd8a1c8c123",
		"text":    "Go — отличный выбор для микросервисов!",
	}
	bodyJSON, _ := json.Marshal(body)

	resp, err := http.Post(
		fmt.Sprintf("%s/questions/%d/answers", baseURL, qID),
		"application/json",
		bytes.NewBuffer(bodyJSON),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
	defer resp.Body.Close()

	// Проверяем, что вопрос теперь возвращает ответ
	resp2, err := http.Get(fmt.Sprintf("%s/questions/%d", baseURL, qID))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp2.StatusCode)
	defer resp2.Body.Close()

	var data map[string]any
	err = json.NewDecoder(resp2.Body).Decode(&data)
	require.NoError(t, err)

	// Смотрим на ответы
	answers, ok := data["answers"].([]any)
	require.True(t, ok)
	require.GreaterOrEqual(t, len(answers), 1, "Ожидается хотя бы один ответ")
}

// вспомогательная функция для создания вопроса
func createQuestion(t *testing.T, text string) uint {
	body := []byte(`{"text":"` + text + `"}`)
	resp, err := http.Post(baseURL+"/questions", "application/json", bytes.NewBuffer(body))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var q models.Question
	err = json.NewDecoder(resp.Body).Decode(&q)
	require.NoError(t, err)
	resp.Body.Close()

	require.NotZero(t, q.ID)
	return q.ID
}

// вспомогательная функция для создания ответа
func createAnswer(t *testing.T, questionID uint, userID, text string) uint {
	body := map[string]interface{}{
		"user_id": userID,
		"text":    text,
	}
	bodyJSON, _ := json.Marshal(body)

	resp, err := http.Post(
		fmt.Sprintf("%s/questions/%d/answers", baseURL, int(questionID)),
		"application/json",
		bytes.NewBuffer(bodyJSON),
	)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)

	var a models.Answer
	err = json.NewDecoder(resp.Body).Decode(&a)
	require.NoError(t, err)
	resp.Body.Close()

	require.Equal(t, questionID, a.QuestionID)
	return a.ID
}

func TestCascadeDeleteQuestion(t *testing.T) {
	clearTables(t)

	//создаем вопрос
	qID := createQuestion(t, "Удалится ли вместе с ответами?")

	//создаем несколько ответов на этот вопрос
	a1 := createAnswer(t, qID, "550e8400-e29b-41d4-a716-446655440000", "Первый ответ")
	a2 := createAnswer(t, qID, "660e8400-e29b-41d4-a716-446655440000", "Второй ответ")

	require.NotZero(t, a1)
	require.NotZero(t, a2)

	// Удаляем вопрос
	req, _ := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/questions/%d", baseURL, qID), nil)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	resp.Body.Close()

	// Пробуем получить вопрос — должен быть 404
	resp2, err := http.Get(fmt.Sprintf("%s/questions/%d", baseURL, qID))
	require.NoError(t, err)
	require.Equal(t, http.StatusNotFound, resp2.StatusCode)
	resp2.Body.Close()

	// Проверяем, что ответы тоже удалены через базу
	db, err := sql.Open("postgres", "host=localhost user=postgres password=postgres dbname=qa_service port=5432 sslmode=disable")
	require.NoError(t, err)
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM answers WHERE question_id=$1", qID).Scan(&count)
	require.NoError(t, err)
	require.Equal(t, 0, count, "После удаления вопроса ответы должны удаляться каскадно")
}

func clearTables(t *testing.T) {
	db, err := sql.Open("postgres", "host=localhost user=postgres password=postgres dbname=qa_service port=5432 sslmode=disable")
	require.NoError(t, err)
	defer db.Close()

	_, err = db.Exec(`TRUNCATE TABLE answers, questions RESTART IDENTITY CASCADE;`)
	require.NoError(t, err)
}
