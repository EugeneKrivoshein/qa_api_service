package repo

import (
	"github.com/EugeneKrivoshein/qa_api_service/models"
	"gorm.io/gorm"
)

type QuestionRepo struct {
	DB *gorm.DB
}

func NewQuestionRepo(db *gorm.DB) *QuestionRepo {
	return &QuestionRepo{DB: db}
}

func (r *QuestionRepo) GetAll() ([]models.Question, error) {
	var questions []models.Question
	err := r.DB.Order("created_at desc").Find(&questions).Error
	return questions, err
}

func (r *QuestionRepo) GetByID(id int) (*models.Question, error) {
	var q models.Question
	err := r.DB.First(&q, id).Error
	if err != nil {
		return nil, err
	}
	return &q, nil
}

func (r *QuestionRepo) Create(q *models.Question) error {
	return r.DB.Create(q).Error
}

func (r *QuestionRepo) Delete(id int) error {
	return r.DB.Delete(&models.Question{}, id).Error
}
