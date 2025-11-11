package repo

import (
	"github.com/EugeneKrivoshein/qa_api_service/models"
	"gorm.io/gorm"
)

type AnswerRepo struct {
	DB *gorm.DB
}

func NewAnswerRepo(db *gorm.DB) *AnswerRepo {
	return &AnswerRepo{DB: db}
}

func (r *AnswerRepo) GetByID(id int) (*models.Answer, error) {
	var a models.Answer
	err := r.DB.First(&a, id).Error
	if err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *AnswerRepo) GetByQuestionID(qid int) ([]models.Answer, error) {
	var answers []models.Answer
	err := r.DB.Where("question_id = ?", qid).Find(&answers).Error
	return answers, err
}

func (r *AnswerRepo) Create(a *models.Answer) error {
	return r.DB.Create(a).Error
}

func (r *AnswerRepo) Delete(id int) error {
	return r.DB.Delete(&models.Answer{}, id).Error
}

func (r *AnswerRepo) DeleteByQuestionID(qid int) error {
	return r.DB.Where("question_id = ?", qid).Delete(&models.Answer{}).Error
}
