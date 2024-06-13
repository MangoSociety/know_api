package repos

import (
	"context"
	"go.mongodb.org/mongo-driver/x/mongo/driver/uuid"
	"know_api/internal_1/models"
	"know_api/internal_1/repository/questions"
	"know_api/pkg/db/mongo"
)

type questionsRepository struct {
	db *mongo.Notes
}

func (q questionsRepository) Create(ctx context.Context, question *models.Question) (*uuid.UUID, error) {
	question.ID, _ = uuid.New()
	_, err := q.db.InsertOne(ctx, question)
	if err != nil {
		return nil, err
	}
	return &question.ID, err
}

//func (q questionsRepository) Update(ctx context.Context, question *models.Question) error {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (q questionsRepository) GetByID(ctx context.Context, questionId uuid.UUID) (*models.Question, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (q questionsRepository) GetByTitle(ctx context.Context, title string) (*models.Question, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (q questionsRepository) GetAllByCategoryID(ctx context.Context, categoryId uuid.UUID) ([]*models.Question, error) {
//	//TODO implement me
//	panic("implement me")
//}

func NewQuestionsMGRepository(db *mongo.Storage) questions.MGRepository {
	return &questionsRepository{
		db: db,
	}
}
