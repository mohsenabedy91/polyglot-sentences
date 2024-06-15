package serviceerror

type ErrorMessage string

var (
	// General
	ServerError        ErrorMessage = "errors.serverError"
	ServiceUnavailable ErrorMessage = "errors.serviceIsUnavailable"
	PermissionDenied   ErrorMessage = "errors.permissionDenied"
	RecordNotFound     ErrorMessage = "errors.recordNotFound"
	Unauthorized       ErrorMessage = "errors.unauthorized"
	IsNotDeletable     ErrorMessage = "errors.isNotDeletable"
	NoRowsEffected     ErrorMessage = "errors.noRowsEffected"

	// User
	UserIsBanned      ErrorMessage = "errors.userIsBanned"
	UserInActive      ErrorMessage = "errors.userInActive"
	UserUnVerified    ErrorMessage = "errors.userUnVerified"
	EmailRegistered   ErrorMessage = "errors.emailRegistered"
	CredentialInvalid ErrorMessage = "errors.credentialInvalid"
	UserLogout        ErrorMessage = "errors.userLogout"
	PasswordIsNull    ErrorMessage = "errors.passwordIsNull"

	// OTP
	InvalidOTP ErrorMessage = "errors.invalidOTP"
	OTPExpired ErrorMessage = "errors.OTPExpired"

	// Token
	InvalidToken ErrorMessage = "errors.invalidToken"
	TokenExpired ErrorMessage = "errors.tokenExpired"

	// Validation
	InvalidRequestBody ErrorMessage = "errors.invalidRequestBody"

	// Role
	RoleExisted ErrorMessage = "errors.roleExisted"
)
