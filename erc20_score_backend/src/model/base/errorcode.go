package base

type ErrorCode string

const (
	Success            ErrorCode = "0_success"
	Login_fail         ErrorCode = "11001_用户名或密码错误！"
	Token_missing      ErrorCode = "11002_token丢失"
	Token_invalid      ErrorCode = "11003_无效token"
	Operation_error    ErrorCode = "11100_执行出错"
	Param_error        ErrorCode = "11101_参数出错"
	No_data_permission ErrorCode = "11102_没有操作此数据权限"
	No_data            ErrorCode = "11103_没有操作此数据权限"
)
