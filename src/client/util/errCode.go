package util

const (
	ErrNull                 = 0  //no error
	ErrGenHttpReqFailed     = 1  //generateRequest,failed
	ErrHttpReqFailed        = 2  // http request ,failed
	ErrNoAccessToken        = 3  //no access token
	ErrUnmarshalFailed      = 4  //unmarsh ,failed
	ErrNoGenShareData       = 5  //no generate share data
	ErrReadHttpBodyFailed   = 6  // read http body ,failed
	ErrLackUserNameOrPasswd = 7  // lack user name or password
	ErrUserLoginFailed      = 8  //user login ,failed
	ErrUserLogoutFailed     = 9  //user logout ,failed
	ErrOffline              = 10 //user is off line
	ErrModifyPasswdFailed   = 11 //modify passwd failed
	ErrMarshalFailed        = 12 //marsh ,failed
	ErrInvalidParam         = 13 //invalid parameter
	ErrUpdateFailed         = 14 //update  failed
	ErrShowFailed           = 15 //show  failed
	ErrDeleteFailed         = 16 //delete  failed
	ErrCreateFailed         = 17 //create failed
	ErrListFailed           = 18 //list failed
	ErrOnlineCheckout       = 19 //online checkout, failed
	ErrGbkToUtf8Failed      = 20 //covert gbk to utf8 failed
	ErrUtf8ToGbkFailed      = 21 //covert utf8 to gbk failed
	ErrVoucherAuditFailed   = 22 //voucher audit failed
)
