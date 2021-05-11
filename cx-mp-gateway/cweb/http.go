package cweb

const (

	HttpHeaderStringDomainId		= "Domain-Id"
	HttpHeaderStringAppId 	 		= "App-Id"

	HttpHeaderStringCaller			= "Caller"

	HttpHeaderStringSecretKey		= "Secret-Key"

	HttpHeaderStringAccessToken  	= "Access-Token"
	HttpHeaderStringUid 			= "Uid"


	HttpHeaderStringServerName		= "Server-Name"
	HttpHeaderStringCommandId		= "Command-Id"
	HttpHeaderStringCommandVersion	= "Command-Version"
)

type CallerType string

const (
	// 普通用户请求本服务
	CallerTypeGeneralUsers = "GeneralUsers"
	// 其他服务请求本服务
	CallerTypeServices = "Services"
)

type RequestHeader struct {
	// 主域 ID
	DomainID			string
	//
	AppID 				string
	// 调用方
	Caller 				CallerType
	// Service 请求需要的请求头
	SecretKey			string
	// User 请求需要的请求头
	AccessToken			string
	Uid                 string
	// 服务名
	ServerName          string
	// 方法名
	MethodName			string
	// 版本号
	CommandVersion		string
}