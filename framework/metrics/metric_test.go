package metrics

import (
	"testing"
	"time"
)

func TestMetric(t *testing.T) {

	//logger.Init("./", "test")
	err := InitMetric("./", "TestMetric", []string{}, nil, func() int64 {
		return time.Now().Unix()
	})
	if err != nil {
		return
	}
	count := 10
	for count > 0 {
		GaugeInc(RedisRCount)
		GaugeInc(MongoRCount)
		GaugeInc(GateConnCount)
		HistogramPut(SrvInvokeDelayHist, 1, Invoke)
		count--

		//request
		stub := StartBlock(&BlockStub{}, RedisBlock, Redis)
		time.Sleep(120 * time.Millisecond)
		stub.End(true)

		stub2 := StartBlock(&BlockStub{}, RedisBlock, Redis)
		time.Sleep(110 * time.Millisecond)
		stub2.End(false)

		stub3 := StartBlock(&BlockStub{}, MongoBlock, Mongo)
		time.Sleep(10 * time.Millisecond)
		stub3.End(true)
		time.Sleep(time.Second * 1)
	}
}
