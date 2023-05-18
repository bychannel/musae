package metrics

// metric labels
const (
	Redis  = "Redis"  // Redis
	Mongo  = "Mongo"  // Mongo
	Invoke = "Invoke" // Invoke
	PubSub = "PubSub" // PubSub

	Delay = "Delay" // Delay

	TagDelay      = "Delay"
	TagConcurrent = "Concurrent"
	TagFailed     = "Failed"
)

// metric name
const (
	//blocking
	RedisBlock  BlockingType = "Redis"  // Redis Delay
	MongoBlock  BlockingType = "Mongo"  // Mongo Delay
	InvokeBlock BlockingType = "Invoke" // Invoke Delay

	// log
	WarnCount  GaugeType = "WarnCount"  // warn count
	ErrorCount GaugeType = "ErrorCount" // Error count
	FatalCount GaugeType = "FatalCount" // Fatal count

	// gauge
	RedisRCount GaugeType = "RedisRCount" // Redis read count
	RedisWCount GaugeType = "RedisWCount" // Redis write count
	MongoRCount GaugeType = "MongoRCount" // Mongo read count
	MongoWCount GaugeType = "MongoWCount" // Mongo write count
	RedisRErr   GaugeType = "RedisRErr"   // Redis read Err count
	RedisWErr   GaugeType = "RedisWErr"   // Redis write Err count
	MongoRErr   GaugeType = "MongoRErr"   // Mongo read Err count
	MongoWErr   GaugeType = "MongoWErr"   // Mongo write Err count

	InvokePubCount GaugeType = "InvokePubCount" // InvokePub count
	InvokeSubCount GaugeType = "InvokeSubCount" // InvokeSub count
	MsgPubCount    GaugeType = "MsgPubCount"    // MsgPubCount count
	MsgSubCount    GaugeType = "MsgSubCount"    // MsgSubCount count

	RedisWSize    GaugeType = "RedisWSize"    // Redis write size
	MongoWSize    GaugeType = "MongoWSize"    // Mongo write size
	GateConnCount GaugeType = "GateConnCount" // PubSub count

	// histogram
	SrvInvokeDelayHist  HistogramType = "SrvInvokeDelayHist"  // srv invoke delay histogram
	UserInvokeDelayHist HistogramType = "UserInvokeDelayHist" // user invoke delay histogram
	RoomInvokeDelayHist HistogramType = "RoomInvokeDelayHist" // user invoke delay histogram
	ChatInvokeDelayHist HistogramType = "ChatInvokeDelayHist" // user invoke delay histogram
	RedisRDelayHist     HistogramType = "RedisRDelayHist"     // redis r delay histogram
	RedisWDelayHist     HistogramType = "RedisWDelayHist"     // redis w delay histogram
	MongoRDelayHist     HistogramType = "MongoRDelayHist"     // redis r delay histogram
	MongoWDelayHist     HistogramType = "MongoWDelayHist"     // redis r delay histogram

	GuideDelayHist HistogramType = "GuideDelayHist" // guide delay histogram

	LoginDelayHist HistogramType = "LoginDelayHist" // login delay histogram

	EnterDelayHist HistogramType = "EnterDelayHist" // endter delay histogram

	GateDelayHist HistogramType = "GateDelayHist" // gate delay histogram

	// 具体服务
	GuideSucceedCount GaugeType = "GuideSucceedCount" // 请求Guide成功的用户数量
	GuideFailedCount  GaugeType = "GuideFailedCount"  // 请求Guide失败的用户数量
	LoginSucceedCount GaugeType = "LoginSucceedCount" // 登录Login成功的用户数量
	LoginFailedCount  GaugeType = "LoginFailedCount"  // 登录Login失败的用户数量
	EnterSucceedCount GaugeType = "EnterSucceedCount" // 登录成功成功的用户数量
	EnterFailedCount  GaugeType = "EnterFailedCount"  // 登录录失败的用户数量
	EnterDropCount    GaugeType = "EnterDropCount"    // 登录录丢去的用户数量
	GateAuthCount     GaugeType = "GateAuthCount"     // 网关Auth验证成功用户数量
	GateAuthFailCount GaugeType = "GateAuthFailCount" // 网关Auth验证失败用户数量
	UserActorCount    GaugeType = "UserActorCount"    // 用户模型在数据中的数量
	UserCount         GaugeType = "UserCount"         // 用户数量
	UserConn          GaugeType = "UserConn"          // 网关 用户连接数量
	PendingUserCount  GaugeType = "PendingUserCount"  // 网关 挂起的用户数量
	QueueUserCount    GaugeType = "QueueUserCount"    // 网关 登录队列用户数量

	GateMsgCount    GaugeType = "GateMsgCount"    // 网关消息计数
	GateUpMsgSize   GaugeType = "GateUpMsgSize"   // 网关上行消息大小
	GateDownMsgSize GaugeType = "GateDownMsgSize" // 网关下行消息大小

	DropUpMsgCount   GaugeType = "DropUpMsgCount"   // 网关上行消息丢弃数
	DropDownMsgCount GaugeType = "DropDownMsgCount" // 网关下行消息丢弃数

	DDosCount      GaugeType = "DDosCount"      // 被攻击计数统计
	ReplayReqCount GaugeType = "ReplayReqCount" // 防重放次数

	PanicCount GaugeType = "PanicCount" // 恐慌日志输出次数
	ErrCount   GaugeType = "ErrCount"   // 错误日志发生次数

)

// RegisterMetrics
//
//	@Description: 重复注册会报错
func RegisterMetrics() {

	RegisterBlockMetric(RedisBlock, Redis)
	RegisterBlockMetric(MongoBlock, Mongo)
	RegisterBlockMetric(InvokeBlock, Invoke)

	RegisterGauge(WarnCount, true)
	RegisterGauge(ErrorCount, true)
	RegisterGauge(FatalCount, true)

	RegisterGauge(RedisRCount, true)
	RegisterGauge(RedisWCount, true)
	RegisterGauge(MongoRCount, true)
	RegisterGauge(MongoWCount, true)
	RegisterGauge(RedisRErr, false)
	RegisterGauge(RedisWErr, false)
	RegisterGauge(MongoRErr, false)
	RegisterGauge(MongoWErr, false)

	RegisterGauge(MsgPubCount, false)
	RegisterGauge(MsgSubCount, false)
	RegisterGauge(InvokePubCount, true)
	RegisterGauge(InvokeSubCount, true)
	RegisterGauge(RedisWSize, true)
	RegisterGauge(MongoWSize, true)

	RegisterGauge(PanicCount, false)
	RegisterGauge(ErrCount, false)

	RegisterHistogram(SrvInvokeDelayHist, nil, Invoke)
	RegisterHistogram(UserInvokeDelayHist, nil, Invoke)
	RegisterHistogram(RoomInvokeDelayHist, nil, Invoke)
	RegisterHistogram(ChatInvokeDelayHist, nil, Invoke)
	RegisterHistogram(RedisRDelayHist, nil, Redis)
	RegisterHistogram(RedisWDelayHist, nil, Redis)
	RegisterHistogram(MongoRDelayHist, nil, Mongo)
	RegisterHistogram(MongoWDelayHist, nil, Mongo)
}
