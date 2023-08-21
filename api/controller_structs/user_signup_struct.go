package controllerstructs

type UserSignUp struct {
	FirstName  string `json:"firstname" validate:"required"`
	MiddleName string `json:"middlename"`
	LastName   string `json:"lastname" validate:"required"`
	Email      string `json:"email" validate:"required"`
}
