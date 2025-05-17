package adminController

import (
	"email-marketing-service/api/v1/dto"
	adminservice "email-marketing-service/api/v1/services/admin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type AdminSupportTicketController struct {
	AdminSupportService *adminservice.AdminSupportService
}

func NewAdminSupportTicketController(adminSupportService *adminservice.AdminSupportService) *AdminSupportTicketController {
	return &AdminSupportTicketController{
		AdminSupportService: adminSupportService,
	}
}

func (c *AdminSupportTicketController) ReplyTicketRequest(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("adminclaims").(jwt.MapClaims)
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

	formFiles := r.MultipartForm.File["file"] // "files" should match the form field name
	for _, header := range formFiles {
		file, err := header.Open()
		if err != nil {
			response.ErrorResponse(w, "Error opening file")
			return
		}
		defer file.Close()
		req.File = append(req.File, header) // Store file header
	}

	if err := dto.ValidateReplyTicketRequest(req); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	res, err := c.AdminSupportService.ReplyToTicket(ticketId, userID, req)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, res)
}

func (c *AdminSupportTicketController) GetAllTickets(w http.ResponseWriter, r *http.Request) {
	page1 := r.URL.Query().Get("page")
	pageSize1 := r.URL.Query().Get("page_size")
	searchQuery := r.URL.Query().Get("search")

	page, err := strconv.Atoi(page1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page number")
		return
	}

	pageSize, err := strconv.Atoi(pageSize1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page size")
		return
	}

	result, err := c.AdminSupportService.GetAllTickets(searchQuery, page, pageSize)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *AdminSupportTicketController) GetClosedTickets(w http.ResponseWriter, r *http.Request) {
	page1 := r.URL.Query().Get("page")
	pageSize1 := r.URL.Query().Get("page_size")
	searchQuery := r.URL.Query().Get("search")

	page, err := strconv.Atoi(page1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page number")
		return
	}

	pageSize, err := strconv.Atoi(pageSize1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page size")
		return
	}

	result, err := c.AdminSupportService.GetClosedTickets(searchQuery, page, pageSize)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}

func (c *AdminSupportTicketController) GetPendingTickets(w http.ResponseWriter, r *http.Request) {
	page1 := r.URL.Query().Get("page")
	pageSize1 := r.URL.Query().Get("page_size")
	searchQuery := r.URL.Query().Get("search")

	page, err := strconv.Atoi(page1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page number")
		return
	}

	pageSize, err := strconv.Atoi(pageSize1)
	if err != nil {
		response.ErrorResponse(w, "Invalid page size")
		return
	}

	result, err := c.AdminSupportService.GetPendingTickets(searchQuery, page, pageSize)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}
