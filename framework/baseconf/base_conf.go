package baseconf

type RedisConf struct {
	Addr            string `json:"addr"`    // ip:port
	AddrDev         string `json:"addrDev"` // dev ip:port
	UserName        string `json:"userName"`
	Password        string `json:"password"`
	DB              int    `json:"db"`
	DialTimeout     int    `json:"dialTimeout"`     // 连接超时
	ReadTimeout     int    `json:"readTimeout"`     // 读超时
	WriteTimeout    int    `json:"writeTimeout"`    // 写超时
	MaxRetries      int    `json:"maxRetries"`      // 最大重试次数
	MinRetryBackoff int    `json:"minRetryBackoff"` // 重试最小backoff
	MaxRetryBackoff int    `json:"maxRetryBackoff"` // 重试最大backoff
}

type ESConf struct {
	Addr     []string `json:"addr"`    // http://ip:port
	AddrDev  string   `json:"addrDev"` // dev ip:port
	UserName string   `json:"userName"`
	Password string   `json:"password"`
	Timeout  int      `json:"timeout"`
}

type BaseConf struct {
	IsDebug                 bool    `json:"debug"`
	IsDBPrefix              bool    `json:"dbPrefix"`                // 存储db key添加ns:grout:前缀正式环境必须false
	Protocol                string  `json:"protocol"`                // app-protocol, http or grpc
	DaprClientRetry         int     `json:"daprClientRetry"`         // dapr client 重连次数
	MusaeDbGetRetryCount    int     `json:"musaeDbGetRetryCount"`    // db读超时重试次数
	MusaeDbSetRetryCount    int     `json:"musaeDbSetRetryCount"`    // db写超时重试次数
	MusaeDbRetryInterval    int     `json:"musaeDbRetryInterval"`    // db读写间隔时间(毫秒)
	AniwarDbGetRetryCount   int     `json:"aniwarDbGetRetryCount"`   // db读超时重试次数
	AniwarDbSetRetryCount   int     `json:"aniwarDbSetRetryCount"`   // db写超时重试次数
	AniwarDbRetryInterval   int     `json:"aniwarDbRetryInterval"`   // db读写间隔时间(毫秒)
	VersionCheck            bool    `json:"versionCheck"`            // 版本检查开关
	Version                 string  `json:"version"`                 // 服务器版本  [渠道].[大版本号].[小版本号]
	VersionAndroid          string  `json:"versionAndroid"`          // 安卓版本 [渠道].[大版本号].[小版本号]
	VersionIOS              string  `json:"versionIOS"`              // IOS版本 [渠道].[大版本号].[小版本号]
	AutoGateIp              bool    `json:"autoGateIp"`              //自动获取本机IP作为GateIP,IsDebug为true是有效
	GateIP                  string  `json:"gateIP"`                  //服务器网关IP
	GatePort                int32   `json:"gatePort"`                //服务器网关端口
	GateMsgMaxSize          int     `json:"gateMsgMaxSize"`          //网关包体大小限制
	SrvMsgMaxSize           int     `json:"srvMsgMaxSize"`           //服务器包体大小限制
	ServerId                string  `json:"serverId"`                //服务器ID
	ServerName              string  `json:"serverName"`              //服务器名称
	ServerHeartbeatTimout   int32   `json:"serverHeartbeatTimout"`   //服务器之间心跳超时，单位秒
	ServerHeartbeatInterval int     `json:"serverHeartbeatInterval"` //服务器之间心跳间隔，单位秒
	AccTokenTTL             int     `json:"accTokenTTL"`             //账号Token有效时长
	RoomTokenTTL            int     `json:"roomTokenTTL"`            //room的Token有效时长
	LogLevel                string  `json:"logLevel"`                //日志等级debug、info、warn、error、fatal
	LogDir                  string  `json:"logDir"`                  //日志输出目录,子目录程序日志[log],埋点日志[dlog],指标日志[mlog]
	LogMaxLen               int     `json:"logMaxLen"`               //单条日志最大长度
	LogCutType              int32   `json:"logCutType"`              //日志切割类型(0=按大小，1=按时间)
	LogRotationTime         int32   `json:"logRotationTime"`         //时间切割间隔（单位：分钟）
	LogMaxSize              int     `json:"logMaxSize"`              //在进行切割之前，日志文件的最大大小（以 MB 为单位）
	LogMaxBackups           int     `json:"logMaxBackups"`           //保留旧文件的最大个数
	LogMaxAges              int     `json:"logMaxAges"`              //保留旧文件的最大天数
	LogCompress             bool    `json:"logCompress"`             //是否压缩 / 归档旧文件
	LogBufSize              int     `json:"logBufSize"`              //最大buf缓存大小（单位kb）
	LogFlushInterval        int     `json:"logFlushInterval"`        //最大flush间隔（单位秒）
	Metric                  bool    `json:"metric"`                  //是否开启指标输出
	LoginReqRate            int32   `json:"loginReqRate"`            //loginReq处理频率每秒
	LoginReqQueue           int32   `json:"loginReqQueue"`           //login最大请求队列
	GatePendingNumLimit     int32   `json:"gatePendingNumLimit"`     //gate排队人数限制
	GateLoginRateLimit      int32   `json:"gateLoginRateLimit"`      //gate登录频率限制
	GateLoginThreadNum      int32   `json:"gateLoginThreadNum"`      //gate登录协成数量
	GateUserNumLimit        int32   `json:"gateUserNumLimit"`        //每个gate上的登录用户上限
	RoleIdCheck             bool    `json:"roleIdCheck"`             //创建角色时检测重复ID
	HeartbeatInterval       int32   `json:"heartbeatInterval"`       //gate心跳检查间隔
	HeartbeatTimout         int32   `json:"heartbeatTimeout"`        //心跳超时，单位秒
	UserCacheTTL            int32   `json:"UserCacheTTL"`            //Gate上user实体换成时长
	Actor2GateType          int32   `json:"actor2GateType"`          //推送actor消息到gate的模式, 1: rpc; 2: gate private topic
	UserActorGCTime         string  `json:"userActorGCTime"`         //UserActor 空闲超时gc, m:分钟, s:秒
	UserActorGCScan         string  `json:"userActorGCScan"`         //UserActor gc扫描间隔, m:分钟, s:秒
	UseEncrypt              int32   `json:"useEncrypt"`              //是否使用加密
	IgnoreEncryptCmdIds     []int32 `json:"ignoreEncryptCmdIds"`     //不做数据加密的协议号
	UseReqIdx               int32   `json:"useReqIdx"`               //是否启用防重放
	OpenCheckBattle         int32   `json:"openCheckBattle"`         //是否开启战斗校验
	ExcelDataZip            int32   `json:"excelDataZip"`            //配置表数据是否使用压缩 1:是, 0:否
	DirtyWords              string  `json:"dirtyWords"`              //屏蔽字库
	RedisLogKey             string  `json:"redisLogKey"`             //redis log key
	FeishuLogRobot          string  `json:"feishuLogRobot"`          //日志推送到飞书聊天群
	FeishuNotifyRobot       string  `json:"feishuNotifyRobot"`       //通知推送到飞书聊天群
	DelayLogLimit           int64   `json:"delayLogLimit"`           //耗时收集日志阈值
	MailActorMin            int32   `json:"mailActorMin"`            //邮件actor最小启用数量
	MailActorPercent        int32   `json:"mailActorPercent"`        //邮件Actor启用数量万分比
	CanVirtualPay           int32   `json:"canVirtualPay"`           //是否支持模拟充值, 1:是, 0:否
	//DefaultEncrypt      string `json:"defaultEncrypt"`      //默认的秘钥
	RedisConf RedisConf `json:"RedisConf"`
	ESConf    ESConf    `json:"ESConf"`
	SpChars   string    `json:"spChars"` // 特殊字符
	CfgKeys   []string  `json:"cfgKeys"` // 配置中心keys
}

var Iconf IConf

type IConf interface {
	BaseConf() *BaseConf
}

func GetBaseConf() *BaseConf {
	if Iconf == nil {
		return nil
	}

	return Iconf.BaseConf()
}
