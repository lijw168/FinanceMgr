package util

// 在该模块定义错误码时，值不能为负值。因为负值是从后台服务传递过来的。
const (
	ErrNull                       = 0  //no error
	ErrGenHttpReqFailed           = 1  //generateRequest,failed
	ErrHttpReqFailed              = 2  // http request ,failed
	ErrNoAccessToken              = 3  //no access token
	ErrUnmarshalFailed            = 4  //unmarsh ,failed
	ErrNoGenShareData             = 5  //no generate share data
	ErrReadHttpBodyFailed         = 6  // read http body ,failed
	ErrLackUserNameOrPasswd       = 7  // lack user name or password
	ErrUserLoginFailed            = 8  //user login ,failed
	ErrUserLogoutFailed           = 9  //user logout ,failed
	ErrOffline                    = 10 //user is off line
	ErrModifyPasswdFailed         = 11 //modify passwd failed
	ErrMarshalFailed              = 12 //marsh ,failed
	ErrInvalidParam               = 13 //invalid parameter
	ErrUpdateFailed               = 14 //update  failed
	ErrShowFailed                 = 15 //show  failed
	ErrDeleteFailed               = 16 //delete  failed
	ErrCreateFailed               = 17 //create failed
	ErrListFailed                 = 18 //list failed
	ErrOnlineCheckout             = 19 //online checkout, failed
	ErrGbkToUtf8Failed            = 20 //covert gbk to utf8 failed
	ErrUtf8ToGbkFailed            = 21 //covert utf8 to gbk failed
	ErrVoucherArrangeFailed       = 22 //voucher arrange failed
	ErrGetLatestVoucherInfoFailed = 23 //get latest voucher information,failed
	ErrGetMaxNumOfMonthFailed     = 24 //get the max numOfMonth,failed
	ErrInitResourceInfoFailed     = 25 //init resource information,failed
	ErrAccSubRefQueryFailed       = 27 //query account subject reference,failed
	ErrBatchAuditVouchersFailed   = 28 // audit vouchers,failed
	ErrCalAccuMoneyFailed         = 29 // calculate accumulative money,failed
	ErrBatchCalAccuMoneyFailed    = 30 // batch calculate accumulative money,failed
	ErrCopyAccSubTemplateFailed   = 31 // copy account subject template,failed
	ErrCalAccPeriodFailed         = 32 // calculate account of period,failed
	ErrAnnualClosing              = 33 // annual closing,failed
	ErrCancelAnnualClosing        = 34 // calcel annual closing,failed
	ErrGetAnnualClosingStatus     = 35 // get annual closing status,failed
)
