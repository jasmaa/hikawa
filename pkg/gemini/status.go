package gemini

const (
	STATUS_INPUT                       = 10
	STATUS_SENSITIVE_INPUT             = 11
	STATUS_SUCCESS                     = 20
	STATUS_REDIRECT_TEMPORARY          = 30
	STATUS_REDIRECT_PERMANENT          = 31
	STATUS_TEMPORARY_FAILURE           = 40
	STATUS_SERVER_UNAVAILABLE          = 41
	STATUS_CGI_ERROR                   = 42
	STATUS_PROXY_ERROR                 = 43
	STATUS_SLOW_DOWN                   = 44
	STATUS_PERMANENT_FAILURE           = 50
	STATUS_NOT_FOUND                   = 51
	STATUS_GONE                        = 52
	STATUS_PROXY_REQUEST_REFUSED       = 53
	STATUS_BAD_REQUEST                 = 59
	STATUS_CLIENT_CERTIFICATE_REQUIRED = 60
	STATUS_CERTIFICATE_NOT_AUTHORISED  = 61
	STATUS_CERTIFICATE_NOT_VALID       = 62
)

func statusCodeToMessage(code int) string {
	switch code {
	case 10:
		return "INPUT"
	case 11:
		return "SENSITIVE INPUT"
	case 20:
		return "SUCCESS"
	case 30:
		return "REDIRECT - TEMPORARY"
	case 31:
		return "REDIRECT - PERMANENT"
	case 40:
		return "TEMPORARY FAILURE"
	case 41:
		return "SERVER UNAVAILABLE"
	case 42:
		return "CGI ERROR"
	case 43:
		return "PROXY ERROR"
	case 44:
		return "SLOW DOWN"
	case 50:
		return "PERMANENT FAILURE"
	case 51:
		return "NOT FOUND"
	case 52:
		return "GONE"
	case 53:
		return "PROXY REQUEST REFUSED"
	case 59:
		return "BAD REQUEST"
	case 60:
		return "CLIENT CERTIFICATE REQUIRED"
	case 61:
		return "CERTIFICATE NOT AUTHORISED"
	case 62:
		return "CERTIFICATE NOT VALID"
	default:
		return "Invalid code"
	}
}
