package handlers

import (
	"net/http"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/ports/api_ports"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
)

type CommentHandler struct {
	cS api_ports.CommentAPIPort
}

func NewCommentHandler(cS api_ports.CommentAPIPort) *CommentHandler {
	return &CommentHandler{
		cS: cS,
	}
}

func (cH CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {

	err := r.ParseMultipartForm(10 << 20)

	if err != nil {
		http.Error(w, "Could not parse form : "+err.Error(), http.StatusInternalServerError)
		return
	}

	payload, err := utils.GetPayload(r)
	if err != nil {
		http.Error(w, "Could not get payload ", http.StatusInternalServerError)
		return
	}
	userRole := domain.Role(payload["role"])

	if userRole != domain.ClientRole {
		http.Error(w, "Only clients can rate", http.StatusForbidden)
		return
	}

	clientId := payload["id"]

	commentBody := r.PostForm.Get("comment")

	tutorId := r.URL.Query().Get("tutId")

	if tutorId == "" {
		http.Error(w, "tutor id can not be empty", http.StatusBadRequest)
		return
	}

	comment := domain.Comment{
		Body:     commentBody,
		Giver:    clientId,
		Receiver: tutorId,
	}

	cmt, err := cH.cS.CreateComment(comment)

	if err != nil {
		http.Error(w, "Could not create comment : "+err.Error(), http.StatusInternalServerError)
		return
	}

	res := Response{
		Success: true,
		Data:    cmt,
	}

	err = utils.WriteJSON(w, http.StatusOK, res, nil)

	if err != nil {
		http.Error(w, "Could not serialize to json", http.StatusInternalServerError)
		return
	}
}

func (cH CommentHandler) GetCommentById(w http.ResponseWriter, r *http.Request) {
	commentId := r.URL.Query().Get("cmtId")

	if commentId == "" {
		http.Error(w, "Comment Id can not be null", http.StatusBadRequest)
		return
	}

	comment, err := cH.cS.GetCommentById(commentId)

	if err != nil {
		http.Error(w, "Could not get comment by id "+err.Error(), http.StatusInternalServerError)
		return
	}
	res := Response{
		Success: true,
		Data:    comment,
	}
	err = utils.WriteJSON(w, http.StatusOK, res, nil)

	if err != nil {
		http.Error(w, "Could not serialize to json", http.StatusInternalServerError)
		return
	}

}
func (cH CommentHandler) GetComments(w http.ResponseWriter, r *http.Request) {

	tutorId := r.URL.Query().Get("tutId")
	cmts, err := cH.cS.GetComments(tutorId)

	if err != nil {
		http.Error(w, "Could not get comments : "+err.Error(), http.StatusInternalServerError)
		return
	}

	var comments []domain.Comment

	for _, comment := range cmts {
		comments = append(comments, *comment)
	}

	res := Response{
		Success: true,
		Data:    comments,
	}
	err = utils.WriteJSON(w, http.StatusOK, res, nil)

	if err != nil {
		http.Error(w, "Could not serialize to json", http.StatusInternalServerError)
		return
	}
}
