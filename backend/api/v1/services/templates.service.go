package services

import (
	"email-marketing-service/api/v1/dto"
	smtpfactory "email-marketing-service/api/v1/factory/smtpFactory"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"sync"
)

type TemplateService struct {
	TemplateRepo     *repository.TemplateRepository
	SubscriptionRepo *repository.SubscriptionRepository
	DailyCalcRepo    *repository.DailyMailCalcRepository
	UserRepo         *repository.UserRepository
}

func NewTemplateService(templateRepo *repository.TemplateRepository, subscriptionRepository *repository.SubscriptionRepository,
	dailyCalc *repository.DailyMailCalcRepository, userRepo *repository.UserRepository) *TemplateService {
	return &TemplateService{
		TemplateRepo:     templateRepo,
		SubscriptionRepo: subscriptionRepository,
		DailyCalcRepo:    dailyCalc,
		UserRepo:         userRepo,
	}
}

func (s *TemplateService) CreateTemplate(d *dto.TemplateDTO) (map[string]interface{}, error) {

	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid data: %w", err)
	}

	id := uuid.New().String()

	templateModel := &model.Template{
		UUID:              id,
		UserId:            d.UserId,
		TemplateName:      d.TemplateName,
		SenderName:        d.SenderName,
		FromEmail:         d.FromEmail,
		Subject:           d.Subject,
		Type:              model.TemplateType(d.Type),
		EmailHtml:         d.EmailHtml,
		EmailDesign:       d.EmailDesign,
		IsEditable:        d.IsEditable,
		IsPublished:       d.IsPublished,
		IsPublicTemplate:  d.IsPublicTemplate,
		IsGalleryTemplate: d.IsGalleryTemplate,
		Tags:              d.Tags,
		Description:       d.Description,
		ImageUrl:          d.ImageUrl,
		IsActive:          d.IsActive,
		EditorType:        d.EditorType,
	}

	checkIfTempleExists, err := s.TemplateRepo.CheckMarketingNameExists(templateModel)

	if err != nil {
		return nil, err
	}

	if checkIfTempleExists {
		return nil, fmt.Errorf("template name already exists")
	}

	if err := s.TemplateRepo.CreateAndUpdateTemplate(templateModel); err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"templateId":  id,
		"editor-type": d.EditorType,
		"type":        d.Type,
		"message":     "template created successfully",
	}, nil
}

func (s *TemplateService) GetTransactionalTemplate(userId string, templateId string) (*model.TemplateResponse, error) {
	result, err := s.TemplateRepo.GetTransactionalTemplate(userId, templateId)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return result, nil
}

func (s *TemplateService) GetMarketingTemplate(userId string, templateId string) (*model.TemplateResponse, error) {
	result, err := s.TemplateRepo.GetMarketingTemplate(userId, templateId)

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, nil
	}

	return result, nil
}

func (s *TemplateService) GetAllTransactionalTemplates(userId string) ([]model.TemplateResponse, error) {

	result, err := s.TemplateRepo.GetAllTransactionalTemplates(userId)

	if err != nil {
		return []model.TemplateResponse{}, err
	}

	if len(result) < 1 {
		return []model.TemplateResponse{}, nil
	}

	return result, nil
}

func (s *TemplateService) GetAllMarketingTemplates(userId string) ([]model.TemplateResponse, error) {
	result, err := s.TemplateRepo.GetAllMarketingTemplates(userId)

	if err != nil {
		return []model.TemplateResponse{}, err
	}

	if len(result) < 1 {
		return []model.TemplateResponse{}, nil
	}

	return result, nil
}

func (s *TemplateService) UpdateTemplate(d *dto.TemplateDTO, templateId string) error {

	templateModel := &model.Template{
		UUID:              templateId,
		UserId:            d.UserId,
		TemplateName:      d.TemplateName,
		SenderName:        d.SenderName,
		FromEmail:         d.FromEmail,
		Subject:           d.Subject,
		Type:              model.TemplateType(d.Type),
		EmailHtml:         d.EmailHtml,
		EmailDesign:       d.EmailDesign,
		IsEditable:        d.IsEditable,
		IsPublished:       d.IsPublished,
		IsPublicTemplate:  d.IsPublicTemplate,
		IsGalleryTemplate: d.IsGalleryTemplate,
		Tags:              d.Tags,
		Description:       d.Description,
		ImageUrl:          d.ImageUrl,
		IsActive:          d.IsActive,
		EditorType:        d.EditorType,
	}

	if err := s.TemplateRepo.UpdateTemplate(templateModel); err != nil {
		return err
	}
	return nil
}

func (s *TemplateService) DeleteTemplate(userId string, templateId string) error {
	return nil
}

func (s *TemplateService) SendTestMail(d *dto.SendTestMailDTO) error {
	if err := utils.ValidateData(d); err != nil {
		return fmt.Errorf("invalid data: %w", err)
	}

	emails := strings.Split(d.EmailAddress, ",")

	userModel := &model.User{UUID: d.UserId}
	userId, err := s.UserRepo.FindUserById(userModel)
	if err != nil {
		return fmt.Errorf("error fetching userId: %w", err)
	}

	subscription, err := s.SubscriptionRepo.GetUserCurrentRunningSubscription(userId.ID)
	if err != nil {
		return fmt.Errorf("error fetching subscription record: %w", err)
	}

	mailCalcRepo, err := s.DailyCalcRepo.GetDailyMailRecordForToday(int(subscription.ID))
	if err != nil {
		return fmt.Errorf("error fetching mail calc: %w", err)
	}

	template, err := s.TemplateRepo.GetSingleTemplate(d.TemplateId)
	if err != nil {
		return fmt.Errorf("error fetching template: %w", err)
	}

	var wg sync.WaitGroup
	var errChan = make(chan error, len(emails))

	for _, email := range emails {
		wg.Add(1)
		go func(email string) {

			defer wg.Done()

			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("Error sending batch: %v\n", err)
					errChan <- fmt.Errorf("panic: %v", err)
				}
			}()

			if err := s.proccessEmail(template.EmailHtml, d.Subject, email, userId.Email, userId.FullName); err != nil {
				errChan <- fmt.Errorf("error processing email %s: %w", email, err)
				return
			}

			// Update the mail calculator record
			updateCalcData := &model.DailyMailCalc{
				UUID:           mailCalcRepo.UUID,
				RemainingMails: mailCalcRepo.RemainingMails - 1,
				MailsSent:      mailCalcRepo.MailsSent + 1,
			}

			if err := s.DailyCalcRepo.UpdateDailyMailCalcRepository(updateCalcData); err != nil {
				errChan <- fmt.Errorf("error updating mail calc for email %s: %w", email, err)
			}
		}(email)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *TemplateService) proccessEmail(design string, subject string, email string, from string, fromName string) error {

	valid := utils.IsValidEmail(email)

	if !valid {
		return nil
	}

	sender := &dto.Sender{Email: from, Name: &fromName}
	recipient := dto.Recipient{Email: email}

	request := &dto.EmailRequest{
		Sender:      *sender,
		To:          recipient,
		Subject:     subject,
		HtmlContent: &design,
	}

	mailS, err := smtpfactory.MailFactory(config.MAIL_PROCESSOR)
	if err != nil {
		return fmt.Errorf("failed to create mail factory: %w", err)
	}

	if err := mailS.HandleSendMail(request); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
