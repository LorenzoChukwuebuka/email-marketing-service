package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/services"
	"github.com/golang-jwt/jwt"
	"net/http"

	"github.com/gorilla/mux"
)

type SupportTicketController struct {
	TicketService *services.SupportTicketService
}

func NewSupportTicketController(ticketService *services.SupportTicketService) *SupportTicketController {
	return &SupportTicketController{
		TicketService: ticketService,
	}
}

func (c *SupportTicketController) CreateTicket(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		response.ErrorResponse(w, "Invalid claims")
		return
	}

	userID, ok := claims["userId"].(string)
	if !ok {
		response.ErrorResponse(w, "Invalid user ID")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.ErrorResponse(w, "Unable to parse form")
		return
	}

	req := &dto.CreateSupportTicketRequest{
		Subject:     r.FormValue("subject"),
		Description: r.FormValue("description"),
		Priority:    dto.Priority(r.FormValue("priority")),
		Message:     r.FormValue("message"),
	}

	file, header, err := r.FormFile("file")
	if err == nil {
		req.File = header
		defer file.Close()
	}

	if err := dto.ValidateCreateSupportTicketRequest(req); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	res, err := c.TicketService.CreateSupportTicket(userID, req)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusCreated, res)
}

func (c *SupportTicketController) ReplyTicket(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		response.ErrorResponse(w, "Invalid claims")
		return
	}

	userID, ok := claims["userId"].(string)
	if !ok {
		response.ErrorResponse(w, "Invalid user ID")
		return
	}

	vars := mux.Vars(r)
	ticketId := vars["ticketId"]

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		response.ErrorResponse(w, "Unable to parse form")
		return
	}

	req := &dto.ReplyTicketRequest{
		Message: r.FormValue("message"),
	}

	file, header, err := r.FormFile("file")
	if err == nil {
		req.File = header
		defer file.Close()
	}

	if err := dto.ValidateReplyTicketRequest(req); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	res, err := c.TicketService.ReplyToTicket(ticketId, userID, req)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, res)
}

func (c *SupportTicketController) GetTicketsByUserID(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		response.ErrorResponse(w, "Invalid claims")
		return
	}

	userID, ok := claims["userId"].(string)
	if !ok {
		response.ErrorResponse(w, "Invalid user ID")
		return
	}

	res, err := c.TicketService.GetTicketsByUserID(userID)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, res)
}

func (c *SupportTicketController) GetSingleTicket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ticketId := vars["ticketId"]

	res, err := c.TicketService.GetTicketWithDetails(ticketId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, res)
}
