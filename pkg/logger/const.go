package logger

func mapToZapParams(keys map[ExtraKey]interface{}) []interface{} {
	params := make([]interface{}, 0, len(keys))
	for k, v := range keys {
		params = append(params, string(k), v)
	}
	return params
}

type Category string
type SubCategory string
type ExtraKey string

const (
	General         Category = "General"
	Internal        Category = "Internal"
	Database        Category = "Database"
	Redis           Category = "Redis"
	Validation      Category = "Validation"
	RequestResponse Category = "RequestResponse"
	Prometheus      Category = "Prometheus"
	JWT             Category = "JWT"
	Authorization   Category = "Authorization"
	Notification    Category = "Notification"
	Twilio          Category = "Twilio"
	Vonage          Category = "Vonage"
	Email           Category = "Email"
	Slack           Category = "Slack"
	Google          Category = "Google"
	Facebook        Category = "Facebook"
	Apple           Category = "Apple"
	Queue           Category = "Queue"
)

const (
	InternalInfo SubCategory = "InternalInfo"

	Startup         SubCategory = "Startup"
	Shutdown        SubCategory = "Shutdown"
	ExternalService SubCategory = "ExternalService"

	API                 SubCategory = "API"
	DefaultRoleNotFound SubCategory = "DefaultRoleNotFound"

	DatabaseConnectionError SubCategory = "DatabaseConnectionError"
	DatabaseQueryError      SubCategory = "DatabaseQueryError"
	DatabaseSelect          SubCategory = "DatabaseSelect"
	DatabaseInsert          SubCategory = "DatabaseInsert"
	DatabaseUpdate          SubCategory = "DatabaseUpdate"
	DatabaseDelete          SubCategory = "DatabaseDelete"
	DatabaseRollback        SubCategory = "DatabaseRollback"
	MigrationUp             SubCategory = "MigrationUp"
	MigrationDown           SubCategory = "MigrationDown"

	RedisRemember SubCategory = "RedisRemember"
	RedisSet      SubCategory = "RedisSet"
	RedisGet      SubCategory = "RedisGet"
	RedisDel      SubCategory = "RedisDel"
	RedisPing     SubCategory = "RedisPing"

	ValidationFailed SubCategory = "ValidationFailed"

	RequestError SubCategory = "RequestError"

	RemoveFile SubCategory = "RemoveFile"

	JWTGenerate SubCategory = "JWTGenerate"

	CheckAccess SubCategory = "CheckAccess"

	NotificationSend SubCategory = "NotificationSend"

	SlackSendMessage SubCategory = "SlackSendMessage"

	TwilioWebhook   SubCategory = "TwilioWebhook"
	TwilioSendSMS   SubCategory = "TwilioSendSMS"
	TwilioCheck     SubCategory = "TwilioCheck"
	TwilioRetrySMS  SubCategory = "TwilioRetrySMS"
	TwilioUpdateSMS SubCategory = "TwilioUpdateSMS"
	VonageWebhook   SubCategory = "VonageWebhook"
	VonageSendSMS   SubCategory = "VonageSendSMS"
	VonageCheck     SubCategory = "VonageCheck"
	VonageRetrySMS  SubCategory = "VonageRetrySMS"
	VonageUpdateSMS SubCategory = "VonageUpdateSMS"
	EmailSend       SubCategory = "EmailSend"

	GoogleLogin   SubCategory = "GoogleLogin"
	FacebookLogin SubCategory = "FacebookLogin"
	AppleLogin    SubCategory = "AppleLogin"

	DataConversion SubCategory = "DataConversion"
)

const (
	AppName          ExtraKey = "AppName"
	LoggerName       ExtraKey = "LoggerName"
	ClientIp         ExtraKey = "ClientIp"
	ListeningAddress ExtraKey = "ListeningAddress"
	Method           ExtraKey = "Method"
	StatusCode       ExtraKey = "StatusCode"
	BodySize         ExtraKey = "BodySize"
	Path             ExtraKey = "Path"
	Latency          ExtraKey = "Latency"
	Body             ExtraKey = "Body"
	ErrorMessages    ExtraKey = "ErrorMessages"
	Headers          ExtraKey = "Headers"
	RequestBody      ExtraKey = "RequestBody"
	ResponseBody     ExtraKey = "ResponseBody"
	ErrorMessage     ExtraKey = "ErrorMessage"

	SelectDBArg ExtraKey = "SelectDBArg"
	InsertDBArg ExtraKey = "InsertDBArg"
)
