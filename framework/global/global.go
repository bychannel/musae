package global

var AppID string    // 服务器标识
var HostName string // 系统主机名
var IsCloud bool    // 是否云环境
var IsDev bool      // 是否开发环境
var SID int64       // 服务ID
var Gateway string  // 网关地址
var TcpAddr string  // 长链接地址
var StartTime int64 // 服务器启动时间
var GateServices []string
var TotalPlayerCount int32   // 总用户在线
var UserActorCount int32     // UserActor数量
var RoomActorCount int32     // UserActor数量
var AllianceActorCount int32 // UserActor数量
var RdsCfgCenterHost string  // redis配置中心 addr+port
var RdsCfgCenterPass string  // redis配置中心 passwd
var RdsCfgNameSpace string   // redis配置中心 namespace
var RdsCfgGroup string       // redis配置中心 group

const (
	//server appid
	GUIDE_SVC  = "guide"
	LOGIN_SVC  = "login"
	GATE_SVC   = "gate"
	LOBBY_SVC  = "lobby"
	CENTER_SVC = "center"
	ACTOR_SVC  = "actor"
	//MAIL_SVC   = "mail"
	IDIP_SVC   = "idip"
	BILL_SVC   = "bill"
	BATTLE_SVC = "battle"

	Platform_Android = "android"
	Platform_IOS     = "ios"

	PlayerCountType = "PlayerCount" // 总在线人数
	CenterActorID   = "centeractor:0"

	UserActorType     = "UserActor"     // UserActor 类型
	RoomActorType     = "RoomActor"     // RoomActor 类型
	AllianceActorType = "AllianceActor" // AllianceActor 类型
	CenterActorType   = "CenterActor"   // CenterActor 类型
	MailActorType     = "MailActor"     // MailActor 类型
)

const (
	TOKEN_LIFE_TIME    = 15 //token有效时长
	SVC_INVOKE_TIMEOUT = 5  //服务调用超时
	DB_INVOKE_TIMEOUT  = 5  //DB调用超时
)
