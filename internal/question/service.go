package question

import (
	"github.com/EugeneKrivoshein/qa_api_service/internal/repo"
	"github.com/EugeneKrivoshein/qa_api_service/models"
)

type Service struct {
	questions *repo.QuestionRepo
	answers   *repo.AnswerRepo
}

func NewService(qr *repo.QuestionRepo, ar *repo.AnswerRepo) *Service {
	return &Service{
		questions: qr,
		answers:   ar,
	}
}

func (s *Service) GetAll() ([]models.Question, error) {
	return s.questions.GetAll()
}

func (s *Service) GetByID(id int) (map[string]any, error) {
	q, err := s.questions.GetByID(id)
	if err != nil {
		return nil, err
	}
	answer, err := s.answers.GetByID(id)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"question": q,
		"answer":   answer,
	}, nil
}

func (s *Service) Create(q *models.Question) error {
	return s.questions.Create(q)
}

func (s *Service) Delete(id int) error {
	if err := s.answers.DeleteByQuestionID(id); err != nil {
		return err
	}
	return s.questions.Delete(id)
}
