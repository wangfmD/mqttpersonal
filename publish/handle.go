package publish

import (
	"encoding/json"
	"sync/atomic"
	"time"

	"github.com/antlinker/mqttpersonal/config"

	"gopkg.in/mgo.v2/bson"
)

func NewHandleConnect(clientID string, pub *Publish) *HandleConnect {
	return &HandleConnect{
		clientID: clientID,
		pub:      pub,
	}
}

type HandleConnect struct {
	clientID string
	pub      *Publish
}

func (hc *HandleConnect) ErrorHandle(err error) {
	hc.pub.lg.Errorf("客户端%s发生异常:%s,已断开连接!", hc.clientID, err.Error())
	if !hc.pub.cfg.AutoReconnect {
		hc.pub.clients.Remove(hc.clientID)
	}
}

func (hc *HandleConnect) Subscribe(topicName, message []byte) {
	atomic.AddInt64(&hc.pub.receiveNum, 1)
	if hc.pub.cfg.IsStore {
		var packetInfo config.PacketInfo
		json.Unmarshal(message, &packetInfo)
		err := hc.pub.database.C(config.C_Packet).Update(bson.M{"packetid": packetInfo.ID}, bson.M{"receivetime": time.Now()})
		if err != nil {
			hc.pub.lg.Errorf("Handle subscribe store error:%s", err.Error())
		}
	}
}
