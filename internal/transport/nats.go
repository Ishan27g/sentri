package transport

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

var opts []nats.Option
var urls = os.Getenv("NATS")

var subjects map[string]*subjectMeta
var lock = sync.RWMutex{}

type SubCb func()

type subjectMeta struct {
	subjectName string
	docId       string
}

func setupConnOptions(opts []nats.Option) []nats.Option {
	totalWait := 5 * time.Second
	reconnectDelay := time.Second

	opts = append(opts, nats.ReconnectWait(reconnectDelay))
	opts = append(opts, nats.MaxReconnects(int(totalWait/reconnectDelay)))
	opts = append(opts, nats.DisconnectErrHandler(func(nc *nats.Conn, err error) {
		fmt.Printf("Disconnected due to:%s, will attempt reconnects for %.0fm", err, totalWait.Minutes())
	}))
	opts = append(opts, nats.ReconnectHandler(func(nc *nats.Conn) {
		fmt.Printf("Reconnected [%s]", nc.ConnectedUrl())
	}))
	opts = append(opts, nats.ClosedHandler(func(nc *nats.Conn) {
		fmt.Printf("Exiting: %v", nc.LastError())
	}))
	return opts
}
func init() {

	subjects = make(map[string]*subjectMeta)
	lock = sync.RWMutex{}
	opts := []nats.Option{nats.Name("-nats-")}
	opts = setupConnOptions(opts)

}
func connect() {
	if nc != nil {
		return
	}
	for {
		<-time.After(1 * time.Second)
		natsC, err := nats.Connect(urls, opts...)
		if err != nil {
			//fmt.Println(err)
		} else {
			nc = natsC
			break
		}
	}
}
func sub(subj string, cb func(msg *nats.Msg)) {
	nc.Subscribe(subj, func(msg *nats.Msg) {
		cb(msg)
	})
	nc.Flush()
	if err := nc.LastError(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("nats-sub for", subj)

}

var nc *nats.Conn = nil

func NatsSubscribe(subj string, cb func(msg *nats.Msg)) {
	if nc == nil {
		connect()
	}
	lock.Lock()
	if subjects[subj] == nil {
		subjects[subj] = &subjectMeta{subjectName: subj, docId: ""}
	}
	lock.Unlock()
	sub(subj, cb)
}
func NatsPublish(subj string, msg []byte) bool {
	if nc == nil {
		connect()
	}
	lock.Lock()
	if subjects[subj] == nil {
		subjects[subj] = &subjectMeta{subjectName: subj, docId: ""}
	}
	lock.Unlock()
	err := nc.Publish(subj, msg)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	err = nc.Flush()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if err = nc.LastError(); err != nil {
		fmt.Println(err)
	} else {
		//fmt.Printf("Published [%s] : '%s'\n", subj, msg)
	}
	return true
}
