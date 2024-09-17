package services

import (
	"bytes"
	"email-marketing-service/api/v1/dto"
	smtpfactory "email-marketing-service/api/v1/factory/smtpFactory"
	"email-marketing-service/api/v1/model"
	"email-marketing-service/api/v1/repository"
	"email-marketing-service/api/v1/utils"
	"encoding/base64"
	"fmt"
	"golang.org/x/net/html"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"
	"unsafe"
)

type CampaignService struct {
	CampaignRepo     *repository.CampaignRepository
	ContactRepo      *repository.ContactRepository
	TemplateRepo     *repository.TemplateRepository
	MailUsageRepo    *repository.MailUsageRepository
	SubscriptionRepo *repository.SubscriptionRepository
	UserRepo         *repository.UserRepository
	DomainRepo       *repository.DomainRepository
}

func NewCampaignService(campaignRepo *repository.CampaignRepository, contactRepo *repository.ContactRepository,
	templateRepo *repository.TemplateRepository, mailusageRepo *repository.MailUsageRepository,
	subscriptionRepo *repository.SubscriptionRepository, userRepo *repository.UserRepository, domainRepo *repository.DomainRepository) *CampaignService {
	return &CampaignService{
		CampaignRepo:     campaignRepo,
		ContactRepo:      contactRepo,
		TemplateRepo:     templateRepo,
		MailUsageRepo:    mailusageRepo,
		UserRepo:         userRepo,
		SubscriptionRepo: subscriptionRepo,
		DomainRepo:       domainRepo,
	}
}

const BATCH_SIZE = 100

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

func (s *CampaignService) GetAllCampaigns(userId string, page int, pageSize int, searchQuery string) (repository.PaginatedResult, error) {
	paginationParams := repository.PaginationParams{Page: page, PageSize: pageSize}
	campaignRepo, err := s.CampaignRepo.GetAllCampaigns(userId, searchQuery, paginationParams)
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
		Sender:         d.Sender,
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

func (s *CampaignService) SendCampaign(d *dto.SendCampaignDTO, isScheduled bool) error {
	campaignG, err := s.CampaignRepo.GetSingleCampaign(d.UserId, d.CampaignId)
	if err != nil {
		return err
	}

	if campaignG.SentAt != nil {
		return fmt.Errorf("you have sent this campaign")
	}

	// Check if the campaign is scheduled and not due yet
	if isScheduled && campaignG.ScheduledAt != nil {
		scheduledTime, err := time.Parse(time.RFC3339, *campaignG.ScheduledAt)
		if err != nil {
			return fmt.Errorf("invalid scheduled time format: %w", err)
		}

		if scheduledTime.After(time.Now()) {
			return nil // Not due yet, exit without sending
		}
	}

	var groupIds []int
	for _, group := range campaignG.CampaignGroups {
		groupIds = append(groupIds, int(group.GroupId))
	}

	var contacts []string
	for _, id := range groupIds {
		getContactsFromGroup, err := s.ContactRepo.GetGroupById(d.UserId, id)
		if err != nil {
			return err
		}
		for _, contact := range getContactsFromGroup.Contacts {

			if !contact.IsSubscribed {
				contacts = append(contacts, contact.Email)
			}

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
			if campaignG.Subject != nil {
				subject = *campaignG.Subject
			}

			previewText := ""
			if campaignG.PreviewText != nil {
				previewText = *campaignG.PreviewText
			}

			err := s.sendEmailBatch(campaignG.Template.EmailHtml,
				d.CampaignId,
				batch, subject,
				previewText,
				mailUsageRecord,
				&mu,
				*campaignG.Sender,
				*campaignG.SenderFromName,
				d.UserId)

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

func (s *CampaignService) sendEmailBatch(templateHtml string, campaignId string, recipients []string, subject string, previewText string, mailUsageRecord *model.MailUsageResponseModel, mu *sync.Mutex, from string, fromName string, userId string) error {
	for _, recipient := range recipients {
		// Modify the email template with tracking info
		modifiedTemplate, err := s.addTrackingToTemplate(templateHtml, campaignId, recipient)
		if err != nil {
			return fmt.Errorf("error adding tracking to template for recipient %s: %w", recipient, err)
		}

		// Send the email
		err = s.sendEmail(recipient, modifiedTemplate, subject, previewText, from, fromName, userId)
		if err != nil {
			return fmt.Errorf("error sending email to recipient %s: %w", recipient, err)
		}

		// Create the EmailCampaignResult for tracking
		err = s.createEmailCampaignResult(campaignId, recipient)
		if err != nil {
			return fmt.Errorf("error creating email campaign result for recipient %s: %w", recipient, err)
		}

		// Lock the mutex before updating mail usage
		mu.Lock()

		// Update the mail usage after sending the email
		mailUsageRecord.MailsSent++

		mailUsageRecord.RemainingMails = mailUsageRecord.LimitAmount - mailUsageRecord.MailsSent

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

func (s *CampaignService) sendEmail(recipient string, emailContent string, subject string, previewText string, from string, fromName string, userId string) error {

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
		PreviewText: &previewText,
	}

	// mailS, err := smtpfactory.MailFactory(config.MAIL_PROCESSOR)

	// if err != nil {
	// 	return fmt.Errorf("failed to create mail factory: %w", err)
	// }

	// if err := mailS.HandleSendMail(request); err != nil {
	// 	return fmt.Errorf("failed to send email: %w", err)
	// }

	// return nil

	// Convert the email content to bytes
	emailBytes := []byte(emailContent)

	// Get the domain from the sender's email
	parts := strings.Split(from, "@")
	if len(parts) != 2 {
		return fmt.Errorf("invalid sender email format")
	}
	senderDomain := parts[1]

	// Fetch the domain information
	domain, err := s.DomainRepo.FindDomain(userId, senderDomain)
	if err != nil {
		if err == repository.ErrDomainNotFound {
			// If the domain is not found, proceed without signing
			return s.sendEmailWithSMTP(request)
		}
		return fmt.Errorf("failed to fetch domain: %w", err)
	}

	if domain != nil && domain.Verified {

		// Check if the sender's domain matches or is a subdomain of the DKIM signing domain
		if !strings.HasSuffix(senderDomain, domain.Domain) {
			return fmt.Errorf("sender domain %s does not align with DKIM signing domain %s", senderDomain, domain.Domain)
		}
		senderModel := &model.Sender{
			UUID:  domain.UUID,
			Email: from,
			Name:  fromName,
		}

		signedBody, err := s.signEmail(senderModel, emailBytes)
		if err != nil {
			// Log the error, but continue with unsigned email
			log.Printf("failed to sign email: %v", err)
		} else {
			emailBytes = signedBody
			request.HtmlContent = (*string)(unsafe.Pointer(&emailBytes))
		}
	}

	return s.sendEmailWithSMTP(request)
}

func (s *CampaignService) signEmail(sender *model.Sender, emailBody []byte) ([]byte, error) {
	// Fetch the domain associated with the sender
	domain, err := s.DomainRepo.GetDomain(sender.UUID)
	if err != nil || !domain.Verified {
		return nil, fmt.Errorf("domain not found or not verified")
	}

	// Load the DKIM private key
	privateKey, err := base64.StdEncoding.DecodeString(domain.DKIMPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decode DKIM private key: %v", err)
	}

	// DKIM signing process
	signedEmail, err := utils.SignEmail(&emailBody, domain.Domain, domain.DKIMSelector, string(privateKey))
	if err != nil {
		return nil, fmt.Errorf("failed to sign email: %v", err)
	}

	return signedEmail, nil
}

func (s *CampaignService) sendEmailWithSMTP(request *dto.EmailRequest) error {
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

func (s *CampaignService) createEmailCampaignResult(campaignId, recipientEmail string) error {
	emailCampaignResult := &model.EmailCampaignResult{
		CampaignID:     campaignId,
		Version:        "1",
		RecipientEmail: recipientEmail,
		SentAt:         time.Now(),
	}

	if err := s.CampaignRepo.CreateEmailCampaignResult(emailCampaignResult); err != nil {
		return fmt.Errorf("error creating email campaign result: %w", err)
	}

	return nil
}

func (s *CampaignService) addTrackingToTemplate(template string, campaignId string, recipientEmail string) (string, error) {
	if template == "" {
		return "", fmt.Errorf("empty template provided")
	}

	// Create tracking pixel and unsubscribe link
	trackingPixel := fmt.Sprintf(`<img src="%s/campaigns/track/open/%s?email=%s" alt="" width="1" height="1" style="display:none;" />`, config.SERVER_URL,
		campaignId, url.QueryEscape(recipientEmail))
	unsubscribeLink := fmt.Sprintf(`<div style="margin-top: 20px; font-size: 12px; color: #666666; text-align: center;">
        <a href="%s/campaigns/unsubscribe?email=%s&campaign=%s" target="_blank" style="color: #666666; text-decoration: underline;">Unsubscribe</a>
    </div>`, config.SERVER_URL, url.QueryEscape(recipientEmail), url.QueryEscape(campaignId))

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
					trackingURL := fmt.Sprintf("%s/campaigns/track/click/%s?email=%s&url=%s", config.SERVER_URL, campaignId, recipientEmail, url.QueryEscape(originalURL))
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

func (s *CampaignService) TrackOpenCampaignEmails(campaignId string, email string, deviceType string, ipAddress string) error {

	htime := time.Now().UTC()

	// Retrieve the existing email campaign result
	existingEmailResult, err := s.CampaignRepo.GetEmailCampaignResult(campaignId, email)
	if err != nil {
		return err
	}

	// Increment the OpenCount by 1
	openCount := 1

	if existingEmailResult != nil {
		openCount = existingEmailResult.OpenCount + 1
	}

	emailResultModel := &model.EmailCampaignResult{CampaignID: campaignId, RecipientEmail: email, OpenedAt: &htime, DeviceType: deviceType, Location: ipAddress, OpenCount: openCount}

	if err := s.CampaignRepo.UpdateEmailCampaignResult(emailResultModel); err != nil {
		return err
	}

	return nil
}

func (s *CampaignService) UnsubscribeFromCampaign(campaignId string, email string) error {

	tx := s.CampaignRepo.DB.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	htime := time.Now().UTC()

	emailResultModel := &model.EmailCampaignResult{CampaignID: campaignId, RecipientEmail: email, UnsubscribeAt: &htime}

	if err := s.CampaignRepo.UpdateEmailCampaignResult(emailResultModel); err != nil {
		tx.Rollback()
		return err
	}

	//update the subscription status on the contacts repo

	if err := s.ContactRepo.UpdateSubscriptionStatus(email); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (s *CampaignService) TrackClickedCampaignsEmails(campaignId string, email string) error {
	htime := time.Now().UTC()

	existingEmailResult, err := s.CampaignRepo.GetEmailCampaignResult(campaignId, email)
	if err != nil {
		return err
	}

	clickCount := 1

	if existingEmailResult != nil {
		clickCount = existingEmailResult.OpenCount + 1
	}

	emailResultModel := &model.EmailCampaignResult{CampaignID: campaignId, RecipientEmail: email, ClickCount: clickCount, ClickedAt: &htime}

	if err := s.CampaignRepo.UpdateEmailCampaignResult(emailResultModel); err != nil {
		return err
	}

	return nil
}

func (s *CampaignService) GetAllRecipientsForACampaign(campaignId string) (*[]model.EmailCampaignResultResponse, error) {
	result, err := s.CampaignRepo.GetAllRecipientsForACampaign(campaignId)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CampaignService) GetEmailResultStats(campaignId string) (map[string]interface{}, error) {
	result, err := s.CampaignRepo.GetEmailResultStats(campaignId)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *CampaignService) GetUserCampaignStats(userId string) (map[string]int64, error) {
	userStats, err := s.CampaignRepo.GetUserCampaignStats(userId)

	if err != nil {
		return nil, err
	}

	return userStats, nil
}

func (s *CampaignService) GetUserCampaignsStats(userId string) ([]map[string]interface{}, error) {

	userStats, err := s.CampaignRepo.GetAllCampaignStatsByUser(userId)

	if err != nil {
		return nil, err
	}

	return userStats, nil
}

//############################################ JOBS ##################################################

func (s *CampaignService) SendScheduledCampaigns() error {
	// Fetch all scheduled campaigns that are due
	scheduledCampaigns, err := s.CampaignRepo.GetDueScheduledCampaigns()
	if err != nil {
		return fmt.Errorf("error fetching due scheduled campaigns: %w", err)
	}

	for _, campaign := range scheduledCampaigns {
		sendDTO := &dto.SendCampaignDTO{
			UserId:     campaign.UserId,
			CampaignId: campaign.UUID,
		}

		err := s.SendCampaign(sendDTO, true)
		if err != nil {
			// Log the error but continue with other campaigns
			log.Printf("Error sending scheduled campaign %s: %v", campaign.UUID, err)
		}
	}

	return nil
}
