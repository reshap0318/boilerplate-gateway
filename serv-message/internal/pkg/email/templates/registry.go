package templates

// TemplateDef describes a named email template: the params it requires and
// how to render its subject + content (still unwrapped — caller applies
// Layout).
type TemplateDef struct {
	RequiredParams []string
	Render         func(params map[string]string, appName string) (subject, content string)
}

// Registry maps a template name (as sent by the caller) to its definition.
// Add new templates here.
var Registry = map[string]TemplateDef{
	"reset_password": {
		RequiredParams: []string{"token", "reset_url"},
		Render: func(params map[string]string, appName string) (string, string) {
			return "Reset Password Request", ResetPasswordContent(params["token"], params["reset_url"])
		},
	},
	"verify_email": {
		RequiredParams: []string{"verify_url"},
		Render: func(params map[string]string, appName string) (string, string) {
			return "Verifikasi Email Anda — " + appName, VerifyEmailContent(params["verify_url"], appName)
		},
	},
	"twofa_code": {
		RequiredParams: []string{"code"},
		Render: func(params map[string]string, appName string) (string, string) {
			return "Kode Verifikasi Login — " + appName, TwoFACodeContent(params["code"])
		},
	},
}
