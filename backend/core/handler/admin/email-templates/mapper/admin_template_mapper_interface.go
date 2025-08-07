package mapper

import (
	"database/sql"
	db "email-marketing-service/internal/db/sqlc"
	"github.com/google/uuid"
)

type TemplateAdapter db.Template
type TemplateByIDRowAdapter db.GetTemplateByIDRow
type ListTemplatesByTypeRowAdapter db.ListTemplatesByTypeRow

func (t TemplateAdapter) GetID() uuid.UUID                   { return t.ID }
func (t TemplateAdapter) GetUserID() uuid.NullUUID           { return t.UserID }
func (t TemplateAdapter) GetCompanyID() uuid.NullUUID        { return t.CompanyID }
func (t TemplateAdapter) GetTemplateName() string            { return t.TemplateName }
func (t TemplateAdapter) GetType() string                    { return t.Type }
func (t TemplateAdapter) GetSenderName() sql.NullString      { return t.SenderName }
func (t TemplateAdapter) GetFromEmail() sql.NullString       { return t.FromEmail }
func (t TemplateAdapter) GetSubject() sql.NullString         { return t.Subject }
func (t TemplateAdapter) GetEmailHtml() sql.NullString       { return t.EmailHtml }
func (t TemplateAdapter) GetTags() sql.NullString            { return t.Tags }
func (t TemplateAdapter) GetDescription() sql.NullString     { return t.Description }
func (t TemplateAdapter) GetImageUrl() sql.NullString        { return t.ImageUrl }
func (t TemplateAdapter) GetEditorType() sql.NullString      { return t.EditorType }
func (t TemplateAdapter) GetIsEditable() sql.NullBool        { return t.IsEditable }
func (t TemplateAdapter) GetIsPublished() sql.NullBool       { return t.IsPublished }
func (t TemplateAdapter) GetIsPublicTemplate() sql.NullBool  { return t.IsPublicTemplate }
func (t TemplateAdapter) GetIsGalleryTemplate() sql.NullBool { return t.IsGalleryTemplate }
func (t TemplateAdapter) GetIsActive() sql.NullBool          { return t.IsActive }
func (t TemplateAdapter) GetEmailDesign() interface{}        { return t.EmailDesign }
func (t TemplateAdapter) GetCreatedAt() sql.NullTime         { return t.CreatedAt }
func (t TemplateAdapter) GetUpdatedAt() sql.NullTime         { return t.UpdatedAt }
func (t TemplateAdapter) GetDeletedAt() sql.NullTime         { return t.DeletedAt }

// Missing methods for TemplateAdapter - return empty values since db.Template doesn't have user/company info
func (t TemplateAdapter) GetUserFullname() sql.NullString { return sql.NullString{} }
func (t TemplateAdapter) GetUserEmail() sql.NullString    { return sql.NullString{} }
func (t TemplateAdapter) GetUserPicture() sql.NullString  { return sql.NullString{} }
func (t TemplateAdapter) GetCompanyName() sql.NullString  { return sql.NullString{} }

func (t TemplateByIDRowAdapter) GetID() uuid.UUID                   { return t.ID }
func (t TemplateByIDRowAdapter) GetUserID() uuid.NullUUID           { return t.UserID }
func (t TemplateByIDRowAdapter) GetCompanyID() uuid.NullUUID        { return t.CompanyID }
func (t TemplateByIDRowAdapter) GetTemplateName() string            { return t.TemplateName }
func (t TemplateByIDRowAdapter) GetType() string                    { return t.Type }
func (t TemplateByIDRowAdapter) GetSenderName() sql.NullString      { return t.SenderName }
func (t TemplateByIDRowAdapter) GetFromEmail() sql.NullString       { return t.FromEmail }
func (t TemplateByIDRowAdapter) GetSubject() sql.NullString         { return t.Subject }
func (t TemplateByIDRowAdapter) GetEmailHtml() sql.NullString       { return t.EmailHtml }
func (t TemplateByIDRowAdapter) GetTags() sql.NullString            { return t.Tags }
func (t TemplateByIDRowAdapter) GetDescription() sql.NullString     { return t.Description }
func (t TemplateByIDRowAdapter) GetImageUrl() sql.NullString        { return t.ImageUrl }
func (t TemplateByIDRowAdapter) GetEditorType() sql.NullString      { return t.EditorType }
func (t TemplateByIDRowAdapter) GetIsEditable() sql.NullBool        { return t.IsEditable }
func (t TemplateByIDRowAdapter) GetIsPublished() sql.NullBool       { return t.IsPublished }
func (t TemplateByIDRowAdapter) GetIsPublicTemplate() sql.NullBool  { return t.IsPublicTemplate }
func (t TemplateByIDRowAdapter) GetIsGalleryTemplate() sql.NullBool { return t.IsGalleryTemplate }
func (t TemplateByIDRowAdapter) GetIsActive() sql.NullBool          { return t.IsActive }
func (t TemplateByIDRowAdapter) GetEmailDesign() interface{}        { return t.EmailDesign }
func (t TemplateByIDRowAdapter) GetCreatedAt() sql.NullTime         { return t.CreatedAt }
func (t TemplateByIDRowAdapter) GetUpdatedAt() sql.NullTime         { return t.UpdatedAt }
func (t TemplateByIDRowAdapter) GetDeletedAt() sql.NullTime         { return t.DeletedAt }
func (t TemplateByIDRowAdapter) GetUserFullname() sql.NullString    { return t.UserFullname }
func (t TemplateByIDRowAdapter) GetUserEmail() sql.NullString       { return t.UserEmail }
func (t TemplateByIDRowAdapter) GetUserPicture() sql.NullString     { return t.UserPicture }
func (t TemplateByIDRowAdapter) GetCompanyName() sql.NullString     { return t.CompanyName }

func (t ListTemplatesByTypeRowAdapter) GetID() uuid.UUID                  { return t.ID }
func (t ListTemplatesByTypeRowAdapter) GetUserID() uuid.NullUUID          { return t.UserID }
func (t ListTemplatesByTypeRowAdapter) GetCompanyID() uuid.NullUUID       { return t.CompanyID }
func (t ListTemplatesByTypeRowAdapter) GetTemplateName() string           { return t.TemplateName }
func (t ListTemplatesByTypeRowAdapter) GetType() string                   { return t.Type }
func (t ListTemplatesByTypeRowAdapter) GetSenderName() sql.NullString     { return t.SenderName }
func (t ListTemplatesByTypeRowAdapter) GetFromEmail() sql.NullString      { return t.FromEmail }
func (t ListTemplatesByTypeRowAdapter) GetSubject() sql.NullString        { return t.Subject }
func (t ListTemplatesByTypeRowAdapter) GetEmailHtml() sql.NullString      { return t.EmailHtml }
func (t ListTemplatesByTypeRowAdapter) GetTags() sql.NullString           { return t.Tags }
func (t ListTemplatesByTypeRowAdapter) GetDescription() sql.NullString    { return t.Description }
func (t ListTemplatesByTypeRowAdapter) GetImageUrl() sql.NullString       { return t.ImageUrl }
func (t ListTemplatesByTypeRowAdapter) GetEditorType() sql.NullString     { return t.EditorType }
func (t ListTemplatesByTypeRowAdapter) GetIsEditable() sql.NullBool       { return t.IsEditable }
func (t ListTemplatesByTypeRowAdapter) GetIsPublished() sql.NullBool      { return t.IsPublished }
func (t ListTemplatesByTypeRowAdapter) GetIsPublicTemplate() sql.NullBool { return t.IsPublicTemplate }
func (t ListTemplatesByTypeRowAdapter) GetIsGalleryTemplate() sql.NullBool {
	return t.IsGalleryTemplate
}
func (t ListTemplatesByTypeRowAdapter) GetIsActive() sql.NullBool       { return t.IsActive }
func (t ListTemplatesByTypeRowAdapter) GetEmailDesign() interface{}     { return t.EmailDesign }
func (t ListTemplatesByTypeRowAdapter) GetCreatedAt() sql.NullTime      { return t.CreatedAt }
func (t ListTemplatesByTypeRowAdapter) GetUpdatedAt() sql.NullTime      { return t.UpdatedAt }
func (t ListTemplatesByTypeRowAdapter) GetDeletedAt() sql.NullTime      { return t.DeletedAt }
func (t ListTemplatesByTypeRowAdapter) GetUserFullname() sql.NullString { return t.UserFullname }
func (t ListTemplatesByTypeRowAdapter) GetUserEmail() sql.NullString    { return t.UserEmail }
func (t ListTemplatesByTypeRowAdapter) GetUserPicture() sql.NullString  { return t.UserPicture }
func (t ListTemplatesByTypeRowAdapter) GetCompanyName() sql.NullString  { return t.CompanyName }
