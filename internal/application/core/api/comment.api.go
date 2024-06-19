package api

import (
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/db_ports"
)

type CommentService struct {
	cr db_ports.CommentDBPort
}

func NewCommentService(cr db_ports.CommentDBPort) *CommentService {
	return &CommentService{
		cr: cr,
	}
}

func (cS *CommentService) CreateComment(newComment domain.Comment) (*domain.Comment, error) {
	comment, err := cS.cr.CreateCommentRepo(newComment)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (cS CommentService) GetCommentById(id string) (*domain.Comment, error) {
	cmt, err := cS.cr.GetCommentByIdRepo(id)
	if err != nil {
		return nil, err
	}
	return cmt, nil
}
func (cS CommentService) GetComments(tutId string) ([]*domain.Comment, error) {

	cmts, err := cS.cr.GetCommentsRepo(tutId)

	if err != nil {
		return nil, err
	}
	return cmts, nil
}
