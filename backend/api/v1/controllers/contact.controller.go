package controllers

import (
	"email-marketing-service/api/v1/dto"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/services"
	"email-marketing-service/api/v1/utils"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type ContactController struct {
	ContactService   *services.ContactService
	UserRepo         *repository.UserRepository
	SubscriptionRepo *repository.SubscriptionRepository
	//Response *utils.ApiResponse
}

func NewContactController(contactsvc *services.ContactService, userRepo *repository.UserRepository, subscriptionRepo *repository.SubscriptionRepository) *ContactController {
	return &ContactController{
		ContactService:   contactsvc,
		UserRepo:         userRepo,
		SubscriptionRepo: subscriptionRepo,
		//Response : &utils.ApiResponse{},
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

if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
        return
    }

	reqdata.UserId = userId

	result, err := c.ContactService.CreateContact(&reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)
}


//will have to revisit this again

func (c *ContactController) UploadContactViaCSV(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		response.ErrorResponse(w, "Invalid claims")
		return
	}

	userId := claims["userId"].(string)

	userModel := &model.User{UUID: userId}
	user, err := c.UserRepo.FindUserById(userModel)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	sub, err := c.SubscriptionRepo.GetUserCurrentRunningSubscription(user.ID)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	var fileSizeLimit int64
	switch sub.Plan.PlanName {
	case "free":
		fileSizeLimit = 2 << 20 // 2 MB
	case "basic":
		fileSizeLimit = 5 << 20 // 5 MB
	case "premium":
		fileSizeLimit = 10 << 20 // 10 MB
	default:
		fileSizeLimit = 2 << 20 // 2 MB default
	}

	// Set a reasonable limit for the entire form, separate from file size limit
	err = r.ParseMultipartForm(15 << 20) // 32 MB
	if err != nil {
		response.ErrorResponse(w, "Error parsing form")
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("contacts_csv")
	if err != nil {
		response.ErrorResponse(w, "Error retrieving the file")
		return
	}
	defer file.Close()

	// Check file size
	if header.Size > fileSizeLimit {
		response.ErrorResponse(w, "File size exceeds the limit for your subscription")
		return
	}

	err = c.ContactService.UploadContactViaCSV(file, header.Filename, userId)
	if err != nil {
		response.ErrorResponse(w, "Error uploading CSV: "+err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "Contacts uploaded successfully")

}

func (c *ContactController) GetAllContacts(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		response.ErrorResponse(w, "Invalid claims")
		return
	}

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

	userId, ok := claims["userId"].(string)
	if !ok {
		response.ErrorResponse(w, "Invalid user ID")
		return
	}

	result, err := c.ContactService.GetAllContacts(userId, page, pageSize, searchQuery)
	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, result)

}

func (c *ContactController) UpdateContact(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.EditContactDTO

	vars := mux.Vars(r)

	contactId := vars["contactId"]

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
        return
    }

	reqdata.UserId = userId
	reqdata.ContactId = contactId

	if err := c.ContactService.UpdateContact(reqdata); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "contact updated successfully")

}

func (c *ContactController) DeleteContact(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	contactId := vars["contactId"]

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	if err := c.ContactService.DeleteContact(userId, contactId); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "contact deleted successfully")

}

func (c *ContactController) CreateGroup(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.ContactGroupDTO

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
        return
    }

	reqdata.UserId = userId

	result, err := c.ContactService.CreateGroup(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)

}

func (c *ContactController) AddContactToGroup(w http.ResponseWriter, r *http.Request) {

	var reqdata *dto.AddContactsToGroupDTO

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
        return
    }

	reqdata.UserId = userId

	result, err := c.ContactService.AddContactsToGroup(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 201, result)

}

func (c *ContactController) RemoveContactFromGroup(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.AddContactsToGroupDTO

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
        return
    }

	reqdata.UserId = userId

	err := c.ContactService.RemoveContactFromGroup(reqdata)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 201, "contact removed successfully")

}

func (c *ContactController) UpdateContactGroup(w http.ResponseWriter, r *http.Request) {
	var reqdata *dto.ContactGroupDTO

	vars := mux.Vars(r)

	groupId := vars["groupId"]

	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId, ok := claims["userId"].(string)
	if !ok {
		response.ErrorResponse(w, "invalid user id in claims")
		return
	}

if err := utils.DecodeRequestBody(r, &reqdata); err != nil {
		response.ErrorResponse(w, "unable to decode request body")
        return
    }

	reqdata.UserId = userId

	err := c.ContactService.UpdateContactGroup(reqdata, groupId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 201, "contact group updated successfully")
}

func (c *ContactController) DeleteContactGroup(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	vars := mux.Vars(r)

	groupId := vars["groupId"]

	if err := c.ContactService.DeleteContactGroup(userId, groupId); err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, "group deleted successfully")

}

func (c *ContactController) GetAllContactGroups(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

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

	result, err := c.ContactService.GetAllContactGroups(userId, page, pageSize, searchQuery)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)

}

func (c *ContactController) GetASingleGroupWithContacts(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	vars := mux.Vars(r)

	groupId := vars["groupId"]

	userId := claims["userId"].(string)

	result, err := c.ContactService.GetASingleGroupWithContacts(userId, groupId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)

}


func (c *ContactController) GetContactCount(w http.ResponseWriter, r *http.Request){
	claims, ok := r.Context().Value("authclaims").(jwt.MapClaims)
	if !ok {
		http.Error(w, "Invalid claims", http.StatusInternalServerError)
		return
	}

	userId := claims["userId"].(string)

	result,err := c.ContactService.GetContactCount(userId)

	if err != nil {
		response.ErrorResponse(w, err.Error())
		return
	}

	response.SuccessResponse(w, 200, result)



}
