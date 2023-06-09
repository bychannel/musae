package baseactor

import (
	"github.com/pkg/errors"
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"gitlab.musadisca-games.com/wangxw/musae/framework/safe"

	"gitlab.musadisca-games.com/wangxw/musae/framework/service"

	"google.golang.org/protobuf/proto"
)

type IBaseHandler interface {
	// Init 初始化数据块
	Init() error

	// LoadDBData 从DB中加载数据
	LoadDBData(mongoDbName service.MongoDbType, dbKey string, dbData proto.Message) error

	// SetDBData 将数据块绑定到handler中
	SetDBData(dbData proto.Message) error

	// DBTable 获取handler管理数据块的相关信息
	// @return 	1.数据块保存的mongo库名称
	//			2.数据块的key
	//			3.handler中的数据块
	DBTable() (service.MongoDbType, string, proto.Message)

	// EnterGame 玩家登陆游戏时执行
	EnterGame() error

	// DailyRefresh 玩家登陆游戏和每天5点都会执行
	DailyRefresh() error

	// IsDirty actor生命周期内有过修改数据,则标记
	IsDirty() bool

	// OfflineSync2DB actor生命周期结束时同步数据到db
	//OfflineSync2DB() error

	//IsMongoDirty() bool
	//SetMongoDirty()
	//CleanMongoDirty()

	// IsRedisDirty 标记修改过的数据是否同步到redis缓存中
	IsRedisDirty() bool
	SetRedisDirty()
	CleanRedisDirty()
	SetSupportMini()
	IsSupportMini() bool

	// SaveDB 持久化数据
	SaveDB(...bool) error

	// Cache2Redis 缓存数据到redis
	Cache2Redis(...bool) error

	//GetPlayerDataKV() (service.MongoDbType, string, proto.Message)
	//SetPlayerDataBound() error

	//InitPlayerData() error
}

func NewBaseHandler(actor IBaseActor, handler string) *BaseHandler {
	h := &BaseHandler{
		actor:       actor,
		ActorLogger: &ActorLogger{a: actor, h: handler},
	}
	return h
}

type BaseHandler struct {
	*ActorLogger
	actor        IBaseActor
	isDirty      bool // actor生命周期内有过修改数据,则标记
	isRedisDirty bool // 标记修改过的数据是否同步到redis缓存中
	//isMongoDirty bool // 废弃
	isSupportMini bool // 是否支持模拟模式

	ChildHandler IBaseHandler // 具体的handler
}

func (h *BaseHandler) IsDirty() bool {
	return h.isDirty
}

//func (h *BaseHandler) IsMongoDirty() bool {
//	return h.isMongoDirty
//}
//
//func (h *BaseHandler) SetMongoDirty() {
//	h.isMongoDirty = true
//}
//
//func (h *BaseHandler) CleanMongoDirty() {
//	h.isMongoDirty = false
//}

func (h *BaseHandler) IsRedisDirty() bool {
	return h.isRedisDirty
}

func (h *BaseHandler) SetRedisDirty() {
	h.isRedisDirty = true
	h.isDirty = true
}

func (h *BaseHandler) CleanRedisDirty() {
	h.isRedisDirty = false
}

func (h *BaseHandler) SetSupportMini() {
	h.isSupportMini = true
}

func (h *BaseHandler) IsSupportMini() bool {
	return h.isSupportMini
}

// 加载玩家数据
func loadPlayerData[T proto.Message](actor IBaseActor, mongoDbName service.MongoDbType, dbKey string, dbData T) (T, error) {
	if safe.IsNull(dbData) {
		dbData = dbData.ProtoReflect().New().Interface().(T)
	}

	_, err := actor.GetCache(mongoDbName, dbKey, dbData)
	if err != nil {
		var zero T // 返回默认值 - nil
		return zero, err
	}

	return dbData, nil
}

func (h *BaseHandler) LoadDBData(mongoDbName service.MongoDbType, dbKey string, dbData proto.Message) error {
	var (
		err error
	)
	if mongoDbName == service.MongoDbType_MongoNil || dbKey == "" {
		return nil
	}

	if dbData, err = loadPlayerData(h.actor, mongoDbName, dbKey, dbData); err != nil {
		if errors.Is(err, service.DB_ERROR_NOT_EXIST) {
			if err = h.ChildHandler.Init(); err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		err = h.ChildHandler.SetDBData(dbData)
		if err != nil {
			return err
		}
	}
	return nil
}

//func (h *BaseHandler) GetPlayerDataKV() (service.MongoDbType, string, proto.Message) {
//	// 具体地派生类类型
//	switch h.ChildHandler.(type) {
//	//case *HeartBeatHandler:
//	//	return
//	case *OfflineDataHandler:
//		return service.MongoDbType_MongoGame, db.KeyOfflineInfo(h.actor.ID()), h.actor.Data.OfflineData
//	case *AccountHandler:
//		return service.MongoDbType_MongoAccount, db.KeyAccountInfo(h.actor.GetUID()), h.actor.Account
//	case *BattleHandler:
//		// do nothing...
//	case *LoginHandler: // 玩家基础数据
//		return service.MongoDbType_MongoGame, db.KeyUserBaseInfo(h.actor.ID()), h.actor.GetUserData()
//	case *PlayerLevelHandler: // 玩家业务数据
//		return service.MongoDbType_MongoGame, db.KeyUserLevelData(h.actor.ID()), h.actor.Data.PlayerLevelData
//	case *BagHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserItems(h.actor.ID()), h.actor.Data.ItemData
//	case *ChapterHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserLevelInfo(h.actor.ID()), h.actor.Data.LevelsData
//	case *CampHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserCamp(h.actor.ID()), h.actor.Data.Camp
//	case *CardHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserCard(h.actor.ID()), h.actor.GetUserCardData()
//	case *CurrencyHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserCurrency(h.actor.ID()), h.actor.Data.Currency
//	case *DutyHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserDutyInfo(h.actor.ID()), h.actor.Data.DutyData
//	case *EquipHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserEquipInfo(h.actor.ID()), h.actor.Data.EquipData
//	case *CampaignHandler:
//		return service.MongoDbType_MongoGame, db.KeyCampaign(h.actor.ID()), h.actor.Data.CampaignInfo
//	case *GmHandler, *GiftHandler:
//		// do nothing...
//	case *HandBookHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserHandBook(h.actor.ID()), h.actor.Data.Handbooks
//	case *MailHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserMail(h.actor.ID()), h.actor.Data.UserMail
//	case *PoolHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserCardPool(h.actor.ID()), h.actor.Data.Pools
//	case *QuestHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserQuestInfo(h.actor.ID()), h.actor.Data.QuestData
//	case *ShopHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserShopInfo(h.actor.ID()), h.actor.Data.ShopData
//	case *TroopHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserCardTroop(h.actor.ID()), h.actor.Data.Troops
//	case *TutorialHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserTutorial(h.actor.ID()), h.actor.Data.Tutorial
//	case *StoryFlagHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserStoryFlag(h.actor.ID()), h.actor.Data.StoryFlagData
//	case *SkinHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserCardSkin(h.actor.ID()), h.actor.Data.SkinData
//	case *SignHandler:
//		return service.MongoDbType_MongoGame, db.KeyUserSign(h.actor.ID()), h.actor.Data.Sign
//	default:
//		logger.Errorf("BaseHandler GetPlayerDataKV 未知的类型, %s", reflect.TypeOf(h.ChildHandler))
//	}
//
//	return "", "", nil
//}

// OfflineSync2DB 玩家下线时同步修改的数据到db
/*func (h *BaseHandler) OfflineSync2DB() error {
	var err error

	if h.IsDirty() {
		//mongoDbName, dbKey, dbData := h.GetPlayerDataKV()
		dbType, dbKey, dbVal := h.ChildHandler.DBTable()
		err = h.actor.SaveMongoDB(dbType, dbKey, dbVal)
		if err != nil {
			return err
		}

		err = h.actor.Cache2Redis(dbType, dbKey, dbVal)
		if err != nil {
			return err
		}

		h.isDirty = false
	}

	return nil
}*/

// SaveDB 持久化到mongo
//@param isImm (true:立即落库; false:延迟落库)
func (h *BaseHandler) SaveDB(isImm ...bool) error {
	err := h.Cache2Redis(isImm...)
	if err != nil {
		logger.Debugf(err.Error())
		return err
	}

	// save 2 mongodb
	if _isImm2Mongo(isImm...) {
		//mongoDbName, dbKey, dbData := h.GetPlayerDataKV()
		dbType, dbKey, dbVal := h.ChildHandler.DBTable()
		err = h.actor.SaveMongoDB(dbType, dbKey, dbVal)
		if err != nil {
			logger.Debugf(err.Error())

			//// 重试提交
			//time.Sleep(time.Millisecond * 100)
			//err = h.actor.SaveMongoDB(dbType, dbKey, dbVal)
			//if err != nil {
			//	err = errors.Wrap(err, "重试SaveDbByKvTable")
			//	logger.Debugf(err.Error())
			//}

			return err
		}
		//} else {
		//	h.SetMongoDirty() // 由派生类发起调用
	}

	return nil
}

// Cache2Redis 缓存到redis
//@param isImm (true:立即落库; false:延迟落库)
func (h *BaseHandler) Cache2Redis(isImm ...bool) error {
	// cache 2 redis
	if _isCommit2Redis(isImm...) {
		//mongoDbName, dbKey, dbData := h.GetPlayerDataKV()
		dbType, dbKey, dbVal := h.ChildHandler.DBTable()
		err := h.actor.Cache2Redis(dbType, h.actor.ID(), dbKey, dbVal)
		if err != nil {
			return err
		}
	} else {
		h.SetRedisDirty() // 由派生类发起调用
	}

	return nil
}

// 解析SaveDB参数
func _isImm2Mongo(isImm ...bool) bool {
	_isImm := false // 默认延迟落库

	if len(isImm) > 0 {
		_isImm = isImm[0] // 有传参就取第一个参数
	}

	return _isImm
}

// 解析SaveDB参数
func _isCommit2Redis(isImm ...bool) bool {
	_isImm := false // 默认延迟落库

	if len(isImm) > 0 {
		_isImm = isImm[0] // 有传参就取第二个参数
	}

	return _isImm
}
