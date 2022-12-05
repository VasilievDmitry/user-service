package pkg

const (
	ServiceName = "lot.user.v1"

	ErrorInternalError           = "internal_server_error"
	ErrorUserNotFound            = "user_not_found"
	ErrorWalletNotFound          = "wallet_not_found"
	ErrorAuthenticationNotFound  = "authentication_not_found"
	ErrorLoginAlreadyConfirmed   = "login_already_confirmed"
	ErrorLoginAlreadyExists      = "login_already_exists"
	ErrorRecoveryCodeInvalid     = "recovery_code_invalid"
	ErrorConfirmationCodeInvalid = "confirmation_code_invalid"
	ErrorInvalidPassword         = "invalid_password"
	ErrorTokenOwnerInvalid       = "token_owner_invalid"
	ErrorWalletUnsupportedType   = "wallet_unsupported_type"

	WalletTypePhantom = "phantom"
)

var (
	WalletList = map[string]bool{
		WalletTypePhantom: true,
	}
)

func IsSupportedWalletType(t string) bool {
	if _, ok := WalletList[t]; !ok {
		return false
	}

	return true
}
