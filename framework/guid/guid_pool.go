package guid

import (
	"gitlab.musadisca-games.com/wangxw/musae/framework/logger"
	"sync"
	"time"
)

const (
	StepSecs = 60
	StepMax  = 1000000 * StepSecs
	StepMin  = 100
)

type GUID_TYPE string

const (
	GUID_ACCOUNT  GUID_TYPE = "GUID_ACCOUNT"  // account id
	GUID_PLAYER   GUID_TYPE = "GUID_PLAYER"   // 玩家id,仅用来生成PlayerID
	GUID_EQUIP    GUID_TYPE = "GUID_EQUIP"    // 装备
	GUID_BUILDING GUID_TYPE = "GUID_BUILDING" // 建筑
	GUID_MAIL     GUID_TYPE = "GUID_MAIL"     // 邮件id
	GUID_EVENT    GUID_TYPE = "GUID_EVENT"    // 事件id
	GUID_TOPIC    GUID_TYPE = "GUID_TOPIC"    // 主题订阅id
	GUID_LOG      GUID_TYPE = "GUID_LOG"      // 日志id
	GUID_ALLIANCE GUID_TYPE = "GUID_ALLIANCE" // 联盟id
	GUID_CALLINFO GUID_TYPE = "GUID_CALLINFO" // 通话ID
)

type FDBNext = func(name string, delta uint64) (uint64, error)

type GUIDGen struct {
	name          string
	id            uint64
	idEnd         uint64
	lastStep      float64
	lastBatchTime int64
	dbNext        FDBNext
	mt            sync.Mutex
}

func (g *GUIDGen) Next() (uint64, error) {
	g.mt.Lock()
	defer g.mt.Unlock()
	if g.id >= g.idEnd {
		now := time.Now().UnixNano()
		lastStep := g.lastStep
		if g.lastBatchTime != 0 {
			nowStep := g.lastStep * float64(StepSecs*time.Second) / float64(now-g.lastBatchTime)
			if nowStep > StepMax {
				nowStep = StepMax
			} else if nowStep < 0 {
				nowStep = 0
			}
			lastStep = nowStep
		}
		if lastStep < StepMin {
			lastStep = StepMin
		}
		lastStepInt := uint64(lastStep)
		id, err := g.dbNext(g.name, lastStepInt)
		for err != nil {
			time.Sleep(200 * time.Millisecond)
			id, err = g.dbNext(g.name, lastStepInt)
			logger.Errorf("error incr db guid:[%s], err:[%v]", g.name, err)
		}
		g.idEnd = id
		g.id = g.idEnd - lastStepInt
		g.lastStep = lastStep
		g.lastBatchTime = now
	}
	g.id++
	return g.id, nil
}

type GUIDPool struct {
	pools  *sync.Map
	dbNext FDBNext
}

func NewGUIDPool(fn FDBNext) *GUIDPool {
	pool := &GUIDPool{}
	pool.dbNext = fn
	pool.pools = &sync.Map{}
	return pool
}

func (g *GUIDPool) EndByPoolSafe(name GUID_TYPE) uint64 {
	v, ok := g.pools.Load(name)
	if !ok {
		v, _ = g.pools.LoadOrStore(name, &GUIDGen{
			name:   string(name),
			dbNext: g.dbNext,
		})
	}
	return v.(*GUIDGen).idEnd
}

func (g *GUIDPool) NextByPool(name GUID_TYPE) (uint64, error) {
	v, ok := g.pools.Load(name)
	if !ok {
		v, _ = g.pools.LoadOrStore(name, &GUIDGen{
			name:   string(name),
			dbNext: g.dbNext,
		})
	}

	return v.(*GUIDGen).Next()
}
