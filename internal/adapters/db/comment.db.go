package db

import (
	"github.com/google/uuid"
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"gorm.io/gorm"
)

type CommentRepo struct {
	db *gorm.DB
}

func NewCommentRepo(db *gorm.DB) *CommentRepo {
	return &CommentRepo{
		db: db,
	}
}

type comment_table struct {
	gorm.Model
	Id           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Body         string    `gorm:"not null"`
	Giver        string
	Client_table client_table `gorm:"foreignKey:Giver;references:Id"`
	Receiver     string
	Tutor_table  tutor_table `gorm:"foreignKey:Receiver;references:Id"`
}

func (cr *CommentRepo) CreateCommentRepo(newComment domain.Comment) (*domain.Comment, error) {
	cmt := &comment_table{
		Body:     newComment.Body,
		Giver:    newComment.Giver,
		Receiver: newComment.Receiver,
	}

	if err := cr.db.Create(&cmt).Error; err != nil {
		return nil, err
	}

	newComment.Id = cmt.Id.String()
	return &newComment, nil

}

func (cr *CommentRepo) GetCommentByIdRepo(id string) (*domain.Comment, error) {
	var cmt comment_table
	if err := cr.db.Model(&comment_table{}).Where("id=?", id).First(&cmt).Error; err != nil {
		return nil, err
	}

	return &domain.Comment{
		Id:       cmt.Id.String(),
		Body:     cmt.Body,
		Giver:    cmt.Giver,
		Receiver: cmt.Receiver,
	}, nil
}

func (cr *CommentRepo) GetCommentsRepo(tutId string) ([]*domain.Comment, error) {
	var cmts []*comment_table

	var commentList []*domain.Comment

	if err := cr.db.Order("created_at DESC").
		Where("receiver=?", tutId).
		Find(&cmts).Error; err != nil {
		return nil, err
	}

	for _, comment := range cmts {
		cmt := &domain.Comment{
			Id:       comment.Id.String(),
			Body:     comment.Body,
			Giver:    comment.Giver,
			Receiver: comment.Receiver,
		}
		commentList = append(commentList, cmt)
	}
	return commentList, nil

}
