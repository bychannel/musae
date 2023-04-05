package global

var AppID string       // 服务器标识
var SysUserName string // 系统用户名
var IsCloud bool       // 是否云环境
var SID int64          // 服务ID

const (
	TOKEN_LIFE_TIME    = 15 //token有效时长
	SVC_INVOKE_TIMEOUT = 5  //服务调用超时
	DB_INVOKE_TIMEOUT  = 5  //DB调用超时

	UserActorType = "UserActor" // UserActor 类型
)
