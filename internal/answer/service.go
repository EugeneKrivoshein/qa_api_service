package answer

import (
	"errors"

	"github.com/EugeneKrivoshein/qa_api_service/internal/repo"
	"github.com/EugeneKrivoshein/qa_api_service/models"
)

type Service struct {
	answers   *repo.AnswerRepo
	questions *repo.QuestionRepo
}

func NewService(ar *repo.AnswerRepo, qr *repo.QuestionRepo) *Service {
	return &Service{
		answers:   ar,
		questions: qr,
	}
}

func (s *Service) Create(a *models.Answer, qid int) error {
	q, err := s.questions.GetByID(qid)
	if err != nil || q == nil {
		return errors.New("question not found")
	}

	a.QuestionID = uint(qid)
	return s.answers.Create(a)
}
func (s *Service) GetByID(id int) (*models.Answer, error) {
	return s.answers.GetByID(id)
}

func (s *Service) Delete(id int) error {
	return s.answers.Delete(id)
}
