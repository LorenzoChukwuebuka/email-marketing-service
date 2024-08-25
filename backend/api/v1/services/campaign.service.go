package services

import (
	"bytes"
	"email-marketing-service/api/v1/dto"
	smtpfactory "email-marketing-service/api/v1/factory/smtpFactory"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"fmt"
	"golang.org/x/net/html"
	"net/url"
	"strings"
	"sync"
	"time"
)

type CampaignService struct {
	CampaignRepo     *repository.CampaignRepository
	ContactRepo      *repository.ContactRepository
	TemplateRepo     *repository.TemplateRepository
	MailUsageRepo    *repository.MailUsageRepository
	SubscriptionRepo *repository.SubscriptionRepository
	UserRepo         *repository.UserRepository
}

func NewCampaignService(campaignRepo *repository.CampaignRepository, contactRepo *repository.ContactRepository,
	templateRepo *repository.TemplateRepository, mailusageRepo *repository.MailUsageRepository,
	subscriptionRepo *repository.SubscriptionRepository, userRepo *repository.UserRepository) *CampaignService {
	return &CampaignService{
		CampaignRepo:     campaignRepo,
		ContactRepo:      contactRepo,
		TemplateRepo:     templateRepo,
		MailUsageRepo:    mailusageRepo,
		UserRepo:         userRepo,
		SubscriptionRepo: subscriptionRepo,
	}
}

const BATCH_SIZE = 20

func (s *CampaignService) CreateCampaign(d *dto.CampaignDTO) (map[string]interface{}, error) {
	if err := utils.ValidateData(d); err != nil {
		return nil, fmt.Errorf("invalid  data: %w", err)
	}

	campaignModel := &model.Campaign{Name: d.Name, UserId: d.UserId, Status: model.CampaignStatus(dto.Draft), SenderFromName: d.SenderFromName}

	campaignExist, err := s.CampaignRepo.CampaignExists(campaignModel)

	if err != nil {
		return nil, err
	}

	if campaignExist {
		return nil, fmt.Errorf("campaign already exists")
	}

	saveCampaign, err := s.CampaignRepo.CreateCampaign(campaignModel)

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"campaignId": saveCampaign,
	}, nil
}

func (s *CampaignService) GetAllCampaigns(userId string, page int, pageSize int) (repository.PaginatedResult, error) {
	paginationParams := repository.PaginationParams{Page: page, PageSize: pageSize}
	campaignRepo, err := s.CampaignRepo.GetAllCampaigns(userId, paginationParams)
	if err != nil {
		return repository.PaginatedResult{}, err
	}

	if campaignRepo.TotalCount == 0 {
		return repository.PaginatedResult{}, nil
	}
	return campaignRepo, nil
}

func (s *CampaignService) GetScheduledCampaigns(userId string, page int, pageSize int) (repository.PaginatedResult, error) {
	paginationParams := repository.PaginationParams{Page: page, PageSize: pageSize}
	campaignRepo, err := s.CampaignRepo.GetScheduledCampaigns(userId, paginationParams)
	if err != nil {
		return repository.PaginatedResult{}, err
	}

	if campaignRepo.TotalCount == 0 {
		return repository.PaginatedResult{}, nil
	}
	return campaignRepo, nil
}

func (s *CampaignService) GetSingleCampaign(userId string, campaignId string) (*model.CampaignResponse, error) {

	campaignRepo, err := s.CampaignRepo.GetSingleCampaign(userId, campaignId)

	if err != nil {
		return nil, err
	}

	if campaignRepo == nil {
		return nil, nil
	}

	return campaignRepo, nil
}

func (s *CampaignService) UpdateCampaign(d *dto.CampaignDTO) error {
	var template *model.TemplateResponse
	var err error

	if d.TemplateId != nil {
		template, err = s.TemplateRepo.GetSingleTemplate(*d.TemplateId)
		if err != nil {
			return err
		}
		if template == nil {
			return fmt.Errorf("template with id %s not found", *d.TemplateId)
		}
	}

	campaignModel := &model.Campaign{
		UUID:           d.UUID,
		Name:           d.Name,
		Subject:        d.Subject,
		PreviewText:    d.PreviewText,
		UserId:         d.UserId,
		SenderFromName: d.SenderFromName,
		RecipientInfo:  d.RecipientInfo,
		IsPublished:    d.IsPublished,
		Status:         model.CampaignStatus(d.Status),
		TrackType:      model.Track,
		IsArchived:     d.IsArchived,
		SentAt:         d.SentAt,
		ScheduledAt:    d.ScheduledAt,
		HasCustomLogo:  d.HasCustomLogo,
	}

	if template != nil {
		campaignModel.TemplateId = &template.ID
	}

	if err := s.CampaignRepo.UpdateCampaign(campaignModel); err != nil {
		return err
	}

	return nil
}

func (s *CampaignService) AddOrEditCampaignGroup(d *dto.CampaignGroupDTO) error {
	if err := utils.ValidateData(d); err != nil {
		return fmt.Errorf("invalid  data: %w", err)
	}

	getCampaign, err := s.CampaignRepo.GetSingleCampaign(d.UserId, d.CampaignId)
	if err != nil {
		return err
	}

	getContactGroup, err := s.ContactRepo.GetASingleGroup(d.UserId, d.GroupId)
	if err != nil {
		return err
	}

	cgpModel := &model.CampaignGroup{CampaignId: getCampaign.ID, GroupId: getContactGroup.ID}

	if err := s.CampaignRepo.AddOrEditCampaignGroup(cgpModel); err != nil {
		return err
	}

	return nil
}

func (s *CampaignService) DeleteCampaign(campaignId string, userId string) error {
	if err := s.CampaignRepo.DeleteCampaign(campaignId, userId); err != nil {
		return err
	}
	return nil
}

func (s *CampaignService) SendCampaign(d *dto.SendCampaignDTO) error {
	getGroup, err := s.CampaignRepo.GetSingleCampaign(d.UserId, d.CampaignId)
	if err != nil {
		return err
	}

	var groupIds []int
	for _, group := range getGroup.CampaignGroups {
		groupIds = append(groupIds, int(group.ID))
	}

	var contacts []string
	for _, id := range groupIds {
		getContactsFromGroup, err := s.ContactRepo.GetGroupById(d.UserId, id)
		if err != nil {
			return err
		}
		for _, contact := range getContactsFromGroup.Contacts {
			contacts = append(contacts, contact.Email)
		}
	}

	userModel := &model.User{UUID: d.UserId}
	userId, err := s.UserRepo.FindUserById(userModel)
	if err != nil {
		return fmt.Errorf("error fetching userId: %w", err)
	}

	subscription, err := s.SubscriptionRepo.GetUserCurrentRunningSubscription(userId.ID)
	if err != nil {
		return fmt.Errorf("error fetching subscription record: %w", err)
	}

	mailUsageRecord, err := s.MailUsageRepo.GetCurrentMailUsageRecord(int(subscription.ID))
	if err != nil {
		return fmt.Errorf("error fetching or creating mail usage record: %w", err)
	}

	if mailUsageRecord.RemainingMails == 0 {
		return fmt.Errorf("you have exceeded your plan limit")
	}

	// Update campaign status to "Sending"
	err = s.updateCampaignStatus(d.CampaignId, "Sending", d.UserId)
	if err != nil {
		return err
	}

	// Create a channel to receive errors from goroutines
	errChan := make(chan error, len(contacts))

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	// Create a mutex to synchronize mail usage updates
	var mu sync.Mutex

	// Send emails in batches using goroutines
	for i := 0; i < len(contacts); i += BATCH_SIZE {
		end := i + BATCH_SIZE
		if end > len(contacts) {
			end = len(contacts)
		}
		batch := contacts[i:end]

		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()

			subject := "No Subject"
			if getGroup.Subject != nil {
				subject = *getGroup.Subject
			}

			previewText := ""
			if getGroup.PreviewText != nil {
				previewText = *getGroup.PreviewText
			}

			err := s.sendEmailBatch(getGroup.Template.EmailHtml, d.CampaignId, batch, subject, previewText, mailUsageRecord, &mu, userId.Email, *getGroup.SenderFromName)
			if err != nil {
				errChan <- err
			}
		}(batch)

		// small delay between batches to avoid overwhelming the email server
		time.Sleep(5 * time.Second)
	}

	// Start a goroutine to close the error channel when all workers are done
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Check for any errors from the goroutines
	for err := range errChan {
		// If there's an error, update campaign status to "Failed"
		s.updateCampaignStatus(d.CampaignId, "Failed", d.UserId)
		return err
	}

	// Update campaign status to "Sent"
	err = s.updateCampaignStatus(d.CampaignId, "Sent", d.UserId)
	if err != nil {
		return err
	}

	return nil
}

func (s *CampaignService) sendEmailBatch(templateHtml string, campaignId string, recipients []string, subject string, previewText string, mailUsageRecord *model.MailUsageResponseModel, mu *sync.Mutex, from string, fromName string) error {
	for _, recipient := range recipients {
		// Modify the email template with tracking info
		modifiedTemplate, err := s.addTrackingToTemplate(templateHtml, campaignId, recipient)
		if err != nil {
			return fmt.Errorf("error adding tracking to template for recipient %s: %w", recipient, err)
		}

		// Send the email
		err = s.sendEmail(recipient, modifiedTemplate, subject, previewText, from, fromName)
		if err != nil {
			return fmt.Errorf("error sending email to recipient %s: %w", recipient, err)
		}

		// Lock the mutex before updating mail usage
		mu.Lock()

		// Update the mail usage after sending the email
		mailUsageRecord.MailsSent++

		mailUsageRecord.RemainingMails = mailUsageRecord.LimitAmount - mailUsageRecord.MailsSent

		print(mailUsageRecord.RemainingMails)

		// Prepare the updated mail usage record
		updateMailUsage := &model.MailUsage{
			UUID:           mailUsageRecord.UUID,
			MailsSent:      mailUsageRecord.MailsSent,
			RemainingMails: mailUsageRecord.RemainingMails,
		}

		// Update the mail usage record in the repository
		if err := s.MailUsageRepo.UpdateMailUsageRecord(updateMailUsage); err != nil {
			mu.Unlock() // Unlock before returning
			return fmt.Errorf("error updating mail usage for recipient %s: %w", recipient, err)
		}

		// Unlock the mutex after updating
		mu.Unlock()
	}

	return nil
}

func (s *CampaignService) sendEmail(recipient string, emailContent string, subject string, previewText string, from string, fromName string) error {

	validEmail := utils.IsValidEmail(recipient)

	if !validEmail {
		return nil
	}

	receiver := dto.Recipient{Email: recipient}
	sender := &dto.Sender{Email: from, Name: &fromName}
	request := &dto.EmailRequest{
		Sender:      *sender,
		To:          receiver,
		Subject:     subject,
		HtmlContent: &emailContent,
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

func (s *CampaignService) updateCampaignStatus(campaignId string, status string, userId string) error {

	htime := time.Now().UTC()

	campaignModel := &model.Campaign{UUID: campaignId, UserId: userId, Status: model.CampaignStatus(status), SentAt: &htime, IsPublished: true}

	if err := s.CampaignRepo.UpdateCampaign(campaignModel); err != nil {
		return err
	}

	return nil
}

func (s *CampaignService) addTrackingToTemplate(template string, campaignId string, recipientEmail string) (string, error) {
	if template == "" {
		return "", fmt.Errorf("empty template provided")
	}

	// Create tracking pixel and unsubscribe link
	trackingPixel := fmt.Sprintf(`<img src="https://yourserver.com/track/open/%s?email=%s" alt="" width="1" height="1" style="display:none;" />`,
		campaignId, url.QueryEscape(recipientEmail))
	unsubscribeLink := fmt.Sprintf(`<div style="margin-top: 20px; font-size: 12px; color: #666666; text-align: center;">
        <a href="https://yourserver.com/unsubscribe?email=%s&campaign=%s" style="color: #666666; text-decoration: underline;">Unsubscribe</a>
    </div>`, url.QueryEscape(recipientEmail), url.QueryEscape(campaignId))

	// Inject tracking pixel and unsubscribe link at the end of the body or document
	if strings.Contains(template, "</body>") {
		template = strings.Replace(template, "</body>", trackingPixel+unsubscribeLink+"</body>", 1)
	} else {
		template += trackingPixel + unsubscribeLink
	}

	// Modify links for click tracking
	doc, err := html.Parse(strings.NewReader(template))
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML template: %w", err)
	}

	var modifyLinks func(*html.Node)
	modifyLinks = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for i, a := range n.Attr {
				if a.Key == "href" && !strings.Contains(a.Val, "unsubscribe") {
					originalURL := a.Val
					trackingURL := fmt.Sprintf("https://yourserver.com/track/click/%s?url=%s", campaignId, url.QueryEscape(originalURL))
					n.Attr[i].Val = trackingURL
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			modifyLinks(c)
		}
	}

	modifyLinks(doc)

	var buf bytes.Buffer
	if err := html.Render(&buf, doc); err != nil {
		return "", fmt.Errorf("failed to render modified HTML: %w", err)
	}

	return buf.String(), nil
}
