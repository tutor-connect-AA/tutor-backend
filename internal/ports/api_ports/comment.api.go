package api_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type CommentAPIPort interface {
	CreateComment(newComment domain.Comment) (*domain.Comment, error)
	GetCommentById(id string) (*domain.Comment, error)
	GetComments(tutId string) ([]*domain.Comment, error)
}
