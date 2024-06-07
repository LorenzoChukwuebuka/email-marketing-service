package observers

import (
	"email-marketing-service/api/dto"
	"email-marketing-service/api/repository"
	"fmt"
)

// DatabaseObserver stores the email status to the database.
type CreateEmailStatusObserver struct {
	EmailRepo *repository.SMTPWebHookRepository
}

func NewCreateEmailStatusObserver(emailRepo *repository.SMTPWebHookRepository) *CreateEmailStatusObserver {
	return &CreateEmailStatusObserver{
		EmailRepo: emailRepo,
	}
}

// Notify handles the event by storing the email status to the database.
func (db *CreateEmailStatusObserver) Notify(event Event) {
	err := StoreEmailStatus(event.EmailRequest, event.Type)
	if err != nil {
		fmt.Printf("Error storing email status: %v\n", err)
	} else {
		fmt.Printf("Email status stored successfully: %s\n", event.Message)
	}

}

// StoreEmailStatus is a placeholder function to store email status in the database.
func StoreEmailStatus(emailRequest *dto.EmailRequest, status string) error {
	// Implement the actual database storage logic here.
	// For example, you can use an ORM like GORM to store the data.
	// db.Create(&EmailStatus{Email: emailRequest.To, Status: status})

	// This is a placeholder implementation.
	fmt.Printf("Storing email status for %s: %s\n", emailRequest.To, status)
	return nil
}
