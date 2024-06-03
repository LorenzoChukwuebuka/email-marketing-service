package controllers

import "net/http"

type TicketController struct {
}

func NewTicketController() *TicketController {
	return &TicketController{}
}

func (c *TicketController) CreateTicket(w http.ResponseWriter, r *http.Request) {

}

func (c *TicketController) ReplyTicket(w http.ResponseWriter, r *http.Request) {

}
