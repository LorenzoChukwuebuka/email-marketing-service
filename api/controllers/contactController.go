package controllers

import "net/http"

type ContactController struct{}

func NewContactController() *ContactController {
	return &ContactController{}
}

func (c *ContactController) AddContact(w http.ResponseWriter, r *http.Request) {

}

func (c *ContactController) GetAllContacts(w http.ResponseWriter, r *http.Request) {

}

func (c *ContactController) UpdateContact(w http.ResponseWriter, r *http.Request) {

}

func (c *ContactController) UploadContactsViaCSV(w http.ResponseWriter, r *http.Request) {

}


func (c *ContactController) DeleteContact(w http.ResponseWriter, r *http.Request){
	
}
