package serviceerror

type ErrorMessage string

var (
	// General
	ServerError      ErrorMessage = "errors.somethingIsWrong"
	PermissionDenied ErrorMessage = "errors.permissionDenied"
	RecordNotFound   ErrorMessage = "errors.recordNotFound"
	Unauthorized     ErrorMessage = "errors.unauthorized"

	// User
	UserIsBanned      ErrorMessage = "errors.userIsBanned"
	UserInActive      ErrorMessage = "errors.userInActive"
	UserUnVerified    ErrorMessage = "errors.userUnVerified"
	EmailRegistered   ErrorMessage = "errors.emailRegistered"
	CredentialInvalid ErrorMessage = "errors.credentialInvalid"

	// OTP
	InvalidOTP ErrorMessage = "errors.invalidOTP"
	OTPExpired ErrorMessage = "errors.OTPExpired"

	// Token
	InvalidToken ErrorMessage = "errors.invalidToken"
	TokenExpired ErrorMessage = "errors.tokenExpired"

	// Validation
	InvalidRequestBody ErrorMessage = "errors.invalidRequestBody"
)
