package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"github.com/golang-jwt/jwt"
	"net/http"
)

type ContactController struct {
	ContactService   *services.ContactService
	UserRepo         *repository.UserRepository
	SubscriptionRepo *repository.SubscriptionRepository
}

func NewContactController(contactsvc *services.ContactService, userRepo *repository.UserRepository, subscriptionRepo *repository.SubscriptionRepository) *ContactController {
	return &ContactController{
		ContactService:   contactsvc,
		UserRepo:         userRepo,
		SubscriptionRepo: subscriptionRepo,
	}
}

func (c *ContactController) CreateContact(w http.ResponseWriter, r *http.Request) {
	var reqdata dto.ContactDTO

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	utils.DecodeRequestBody(r, &reqdata)

	reqdata.UserId = userId

	result, err := c.ContactService.CreateContact(&reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)

}

func (c *ContactController) UploadContactViaCSV(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	//get the users id in int

	userModel := &model.User{UUID: userId}

	user, err := c.UserRepo.FindUserById(userModel)

	if err != nil {
		response.ErrorResponse(w, err)
		return
	}

	sub, err := c.SubscriptionRepo.GetUserCurrentRunningSubscription(user.ID)

	if err != nil {
		response.ErrorResponse(w, err)
	}

	var fileSizeLimit int64

	switch sub.Plan.PlanName {
	case "free":
		fileSizeLimit = 5 << 20 // 5 MB
	case "basic":
		fileSizeLimit = 20 << 20 // 20 MB
	case "premium":
		fileSizeLimit = 50 << 20 // 50 MB
	default:
		fileSizeLimit = 10 << 20 // 10 MB default
	}

	err = r.ParseMultipartForm(fileSizeLimit)
	if err != nil {
		http.Error(w, "File too large or unable to parse form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("contacts_csv")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check file size
	if header.Size > fileSizeLimit {
		http.Error(w, "File size exceeds the limit for your subscription", http.StatusBadRequest)
		return
	}

	err = c.ContactService.UploadContactViaCSV(file, header.Filename, userId)
	if err != nil {
		http.Error(w, "Error uploading CSV: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.SuccessResponse(w, 200, "contacts uploaded successfully")

}

func (c *ContactController) GetAllContacts(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	result, err := c.ContactService.GetAllContacts(userId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
	}

	response.SuccessResponse(w, 200, result)

}

func (c *ContactController) AddContactToGroup(w http.ResponseWriter, r *http.Request) {}

func (c *ContactController) DeleteContact(w http.ResponseWriter, r *http.Request) {}

func (c *ContactController) RemoveContactFromGroup(w http.ResponseWriter, r *http.Request) {}

func (c *ContactController) UpdateContactGroup(w http.ResponseWriter, r *http.Request) {}
