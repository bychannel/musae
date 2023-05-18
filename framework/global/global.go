package global

var AppID string         // 服务器标识
var HostName string      // 系统主机名
var IsCloud bool         // 是否云环境
var SID int64            // 服务ID
var GateWay string       // 网关地址
var UserActorCount int32 // UserActor数量
var RoomActorCount int32 // UserActor数量
var ChatActorCount int32 // UserActor数量

const (
	//server appid
	GUIDE_SVC = "guide"
	LOGIN_SVC = "login"
	GATE_SVC  = "gate"
	LOBBY_SVC = "lobby"
	ACTOR_SVC = "actor"
	//MAIL_SVC   = "mail"
	IDIP_SVC   = "idip"
	BILL_SVC   = "bill"
	BATTLE_SVC = "battle"

	UserActorType   = "UserActor"   // UserActor 类型
	RoomActorType   = "RoomActor"   // RoomActor 类型
	ChatActorType   = "ChatActor"   // ChatActor 类型
	CenterActorType = "CenterActor" // CenterActor 类型
)

const (
	TOKEN_LIFE_TIME    = 15 //token有效时长
	SVC_INVOKE_TIMEOUT = 5  //服务调用超时
	DB_INVOKE_TIMEOUT  = 5  //DB调用超时
)
