package permission

var PlanFeatureMap = map[string]map[string]bool{
	"free": {
		"upload_csv":    false,
		"max_file_size": true,
		"send_campaign": false,
	},
	"basic": {
		"upload_csv":    true,
		"max_file_size": true,
		"send_campaign": true,
	},
	"premium": {
		"upload_csv":    true,
		"max_file_size": true,
		"send_campaign": true,
		"custom_domain": true,
	},
	"custom": {
		"upload_csv":    true,
		"max_file_size": true,
		"send_campaign": true,
		"custom_domain": true,
		"dedicated_ip":  true,
	},
	"professional": {
		"upload_csv":    true,
		"max_file_size": true,
		"send_campaign": true,
		"custom_domain": true,
		"dedicated_ip":  true,
	},
}
