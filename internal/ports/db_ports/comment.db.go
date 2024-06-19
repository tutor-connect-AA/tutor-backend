package db_ports

import "github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"

type CommentDBPort interface {
	CreateCommentRepo(newComment domain.Comment) (*domain.Comment, error)
	GetCommentByIdRepo(id string) (*domain.Comment, error)
	GetCommentsRepo(tutId string) ([]*domain.Comment, error)
}
