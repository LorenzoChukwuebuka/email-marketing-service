package dto

type Admin struct {
	FirstName  *string `json:"firstname"`
	MiddleName *string `json:"middlename"`
	LastName   *string `json:"lastname"`
	Email      string  `json:"email"`
	Password   []byte  `json:"password"`
	Type       string  `json:"type"`
}

type AdminLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type AdminEmailLogDTO struct {
	Subject string `json:"subject" validate:"required"`
	Content string `json:"content" validate:"required"`
}

type SystemsDTO struct {
	Domain string `json:"domain" validate:"required"`
}
