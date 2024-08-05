package observers

import (
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"fmt"
)

// DatabaseObserver stores the email status to the database.
type CreateEmailStatusObserver struct {
	EmailRepo *repository.MailStatusRepository
}

func NewCreateEmailStatusObserver(emailRepo *repository.MailStatusRepository) *CreateEmailStatusObserver {
	return &CreateEmailStatusObserver{
		EmailRepo: emailRepo,
	}
}

// Notify handles the event by storing the email status to the database.
func (db *CreateEmailStatusObserver) Notify(event Event) {
	err := db.StoreEmailStatus(event.Payload, event.Type)
	if err != nil {
		fmt.Printf("Error storing email status: %v\n", err)
	} else {
		fmt.Printf("Email status stored successfully: %s\n", event.Message)
	}

}

// StoreEmailStatus is a placeholder function to store email status in the database.
func (db *CreateEmailStatusObserver) StoreEmailStatus(payload interface{}, status string) error {

	data := &model.SentEmails{}
	err := db.EmailRepo.CreateStatus(data)
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
