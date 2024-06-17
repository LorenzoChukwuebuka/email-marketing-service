package controllers

import (
	"email-marketing-service/api/dto"
	"email-marketing-service/api/services"

	"github.com/golang-jwt/jwt"
	"net/http"
)

type SupportTicketController struct {
	TicketService services.SupportTicketService
}

func NewTicketController(ticketService *services.SupportTicketService) *SupportTicketController {
	return &SupportTicketController{
		TicketService: *ticketService,
	}
}

func (c *SupportTicketController) CreateTicket(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	// Parse the multipart form data
	err := r.ParseMultipartForm(10 << 20) // limit your max input length!
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	var reqdata dto.SupportTicket
	reqdata.Subject = r.FormValue("subject")
	reqdata.Description = r.FormValue("description")
	reqdata.Status = dto.Status(r.FormValue("status"))
	reqdata.Priority = dto.Priority(r.FormValue("priority"))
	reqdata.Message = r.FormValue("message")

	userId := claims["userId"].(string)
	reqdata.UserID = userId

	// Retrieve file
	file, _, err := r.FormFile("ticket_file")
	if err != nil {
		http.Error(w, "Error Retrieving the File", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileBytes := make([]byte, 10<<20)
	n, err := file.Read(fileBytes)
	if err != nil {
		http.Error(w, "Error Reading the File", http.StatusInternalServerError)
		return
	}

	ticketFile := fileBytes[:n]
	reqdata.TicketFile = &ticketFile

	err = c.TicketService.CreateSupportTicket(&reqdata)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "Your ticket has been successfully created. An agent will reply to you soon.")
}

func (c *SupportTicketController) ReplyTicket(w http.ResponseWriter, r *http.Request) {
	// Implementation for replying to a ticket
}
