package domain

import "errors"

type ErrorCode string

const (
	CodeEmailRequired                        ErrorCode = "EMAIL_REQUIRED"
	CodeWeakPassword                         ErrorCode = "WEAK_PASSWORD"
	CodeNotMatchPassword                     ErrorCode = "NOT_MATCH_PASSWORD"
	CodeEmailAlreadyExists                   ErrorCode = "EMAIL_ALREADY_EXISTS"
	CodeUserNotFound                         ErrorCode = "USER_NOT_FOUND"
	CodeInvalidCredentials                   ErrorCode = "INVALID_CREDENTIALS"
	CodeUserBanned                           ErrorCode = "USER_BANNED"
	CodeUserInactive                         ErrorCode = "USER_INACTIVE"
	CodeInvalidToken                         ErrorCode = "INVALID_TOKEN"
	CodeInvalidOTP                           ErrorCode = "INVALID_OTP"
	CodeOTPExpired                           ErrorCode = "OTP_EXPIRED"
	CodeEmptyUsername                        ErrorCode = "EMPTY_USERNAME"
	CodeEmptyEmail                           ErrorCode = "EMPTY_EMAIL"
	CodeEmptyPhone                           ErrorCode = "EMPTY_PHONE"
	CodeUsernameExists                       ErrorCode = "USERNAME_EXISTS"
	CodePhoneAlreadyExists                   ErrorCode = "PHONE_ALREADY_EXISTS"
	CodeEmptyAddressLine                     ErrorCode = "EMPTY_ADDRESS_LINE"
	CodeEmptyUserID                          ErrorCode = "EMPTY_USER_ID"
	CodeEmptyCity                            ErrorCode = "EMPTY_CITY"
	CodeSameEmail                            ErrorCode = "SAME_EMAIL"
	CodeSamePhone                            ErrorCode = "SAME_PHONE"
	CodeEmailNotSet                          ErrorCode = "EMAIL_NOT_SET"
	CodeUnauthorizedUser                     ErrorCode = "UNAUTHORIZED_USER"
	CodeInvalidFullName                      ErrorCode = "INVALID_FULL_NAME"
	CodeInvalidGender                        ErrorCode = "INVALID_GENDER"
	CodeInvalidDob                           ErrorCode = "INVALID_DOB"
	CodeNoFieldsToUpdate                     ErrorCode = "NO_FIELDS_TO_UPDATE"
	CodeAddressNotFound                      ErrorCode = "ADDRESS_NOT_FOUND"
	CodeVerifiedEmailCannotBeUpdatedDirectly ErrorCode = "VERIFIED_EMAIL_CANNOT_BE_UPDATED_DIRECTLY"
	CodeCurrentEmailNotVerified              ErrorCode = "CURRENT_EMAIL_NOT_VERIFIED"
	CodeCurrentPhoneNotVerified              ErrorCode = "CURRENT_PHONE_NOT_VERIFIED"
	CodeInvalidOperation                     ErrorCode = "INVALID_OPERATION"
	CodeTokenAlreadyUsed                     ErrorCode = "TOKEN_ALREADY_USED"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

func (e *AppError) Error() string { return e.Message }
func (e *AppError) Unwrap() error { return e.Err }

func NewAppError(code ErrorCode, err error) *AppError {
	return &AppError{Code: code, Message: err.Error(), Err: err}
}

var (
	ErrEmailRequired                        = errors.New("email is required")
	ErrWeakPassword                         = errors.New("password must be at least 8 characters")
	ErrNotMatchPassword                     = errors.New("confirmed password is not match")
	ErrEmailAlreadyExists                   = errors.New("email already registered")
	ErrUserNotFound                         = errors.New("user not found")
	ErrInvalidCredentials                   = errors.New("invalid credentials")
	ErrUserBanned                           = errors.New("user account is banned")
	ErrUserInactive                         = errors.New("user account is inactive")
	ErrInvalidToken                         = errors.New("invalid or expired token")
	ErrInvalidOTP                           = errors.New("invalid OTP code")
	ErrOTPExpired                           = errors.New("OTP code expired")
	ErrEmptyUsername                        = errors.New("username cannot be empty")
	ErrEmptyEmail                           = errors.New("email cannot be empty")
	ErrEmptyPhone                           = errors.New("phone cannot be empty")
	ErrUsernameExists                       = errors.New("username already exists")
	ErrPhoneAlreadyExists                   = errors.New("phone number already registered")
	ErrEmptyAddressLine                     = errors.New("address line cannot be empty")
	ErrEmptyUserID                          = errors.New("userID cannot be empty")
	ErrEmptyCity                            = errors.New("city cannot be empty")
	ErrSameEmail                            = errors.New("new email must be different from current email")
	ErrSamePhone                            = errors.New("new phone must be different from current phone")
	ErrEmailNotSet                          = errors.New("no email address configured for this user")
	ErrUnauthorizedUser                     = errors.New("user role not authorized")
	ErrInvalidFullName                      = errors.New("invalid full name")
	ErrInvalidGender                        = errors.New("invalid gender")
	ErrInvalidDob                           = errors.New("invalid date of birth")
	ErrNoFieldsToUpdate                     = errors.New("at least 1 field is required to update")
	ErrAddressNotFound                      = errors.New("address not found")
	ErrVerifiedEmailCannotBeUpdatedDirectly = errors.New("verified email cannot be updated directly")
	ErrCurrentEmailNotVerified              = errors.New("current email must be verified first")
	ErrCurrentPhoneNotVerified              = errors.New("current phone must be verified first")
	ErrInvalidOperation                     = errors.New("invalid operation or purpose")
	ErrTokenAlreadyUsed                     = errors.New("token has already been used")
)

var sentinelToCodeMap = map[error]ErrorCode{
	ErrEmailRequired:                        CodeEmailRequired,
	ErrWeakPassword:                         CodeWeakPassword,
	ErrNotMatchPassword:                     CodeNotMatchPassword,
	ErrEmailAlreadyExists:                   CodeEmailAlreadyExists,
	ErrUserNotFound:                         CodeUserNotFound,
	ErrInvalidCredentials:                   CodeInvalidCredentials,
	ErrUserBanned:                           CodeUserBanned,
	ErrUserInactive:                         CodeUserInactive,
	ErrInvalidToken:                         CodeInvalidToken,
	ErrInvalidOTP:                           CodeInvalidOTP,
	ErrOTPExpired:                           CodeOTPExpired,
	ErrEmptyUsername:                        CodeEmptyUsername,
	ErrEmptyEmail:                           CodeEmptyEmail,
	ErrEmptyPhone:                           CodeEmptyPhone,
	ErrUsernameExists:                       CodeUsernameExists,
	ErrPhoneAlreadyExists:                   CodePhoneAlreadyExists,
	ErrEmptyAddressLine:                     CodeEmptyAddressLine,
	ErrEmptyUserID:                          CodeEmptyUserID,
	ErrEmptyCity:                            CodeEmptyCity,
	ErrSameEmail:                            CodeSameEmail,
	ErrSamePhone:                            CodeSamePhone,
	ErrEmailNotSet:                          CodeEmailNotSet,
	ErrUnauthorizedUser:                     CodeUnauthorizedUser,
	ErrInvalidFullName:                      CodeInvalidFullName,
	ErrInvalidGender:                        CodeInvalidGender,
	ErrInvalidDob:                           CodeInvalidDob,
	ErrNoFieldsToUpdate:                     CodeNoFieldsToUpdate,
	ErrAddressNotFound:                      CodeAddressNotFound,
	ErrVerifiedEmailCannotBeUpdatedDirectly: CodeVerifiedEmailCannotBeUpdatedDirectly,
	ErrCurrentEmailNotVerified:              CodeCurrentEmailNotVerified,
	ErrCurrentPhoneNotVerified:              CodeCurrentPhoneNotVerified,
	ErrInvalidOperation:                     CodeInvalidOperation,
	ErrTokenAlreadyUsed:                     CodeTokenAlreadyUsed,
}

func ToAppError(err error) *AppError {
	if err == nil {
		return nil
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	for sentinel, code := range sentinelToCodeMap {
		if errors.Is(err, sentinel) {
			return NewAppError(code, sentinel)
		}
	}
	return nil
}
