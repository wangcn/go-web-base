package xnsq

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

const errNsqConfigNil = "NsqConfigNil: %s"
const errNsqEmptyPool = "NsqEmptyPool"

// NOTE: must call LoadNsqds first if use Producer (Publish / DeferredPublish)
type XNsq struct {
	Env       string
	producers []*nsq.Producer
	mutex     *sync.RWMutex

	consumers []*nsq.Consumer
}

type NsqNode struct {
	IP      string `json:"broadcast_address"`
	TcpPort int    `json:"tcp_port"`
}

type NsqList struct {
	Producers []NsqNode
}

func NewXNsq(env string) *XNsq {
	return &XNsq{
		Env:       env,
		producers: []*nsq.Producer{},
		mutex:     &sync.RWMutex{},
		consumers: []*nsq.Consumer{},
	}
}

func (x *XNsq) lookupNsqdAddrs() ([]string, error) {
	lookupdHTTPAddrs := viper.GetStringSlice(x.Env + ".consumerAddrs")
	var body []byte
	rand.Shuffle(len(lookupdHTTPAddrs), func(i, j int) { lookupdHTTPAddrs[i], lookupdHTTPAddrs[j] = lookupdHTTPAddrs[j], lookupdHTTPAddrs[i] })
	for _, addr := range lookupdHTTPAddrs {
		resp, err := http.Get(fmt.Sprintf("http://%s/nodes", addr))
		if err != nil {
			fmt.Println("lookup nodes failed at "+addr, "err: ", err)
			continue
		}
		content, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}
		body = content
		break
	}
	if len(body) == 0 {
		errMsg := "failed to get nsq nodes"
		return nil, errors.New(errMsg)
	}
	var v NsqList
	err := json.Unmarshal(body, &v)
	if err != nil {
		errMsg := "failed to parser nsq nodes"
		return nil, errors.New(errMsg)
	}
	addrList := make([]string, 0)
	for _, node := range v.Producers {
		addrList = append(addrList, fmt.Sprintf("%s:%d", node.IP, node.TcpPort))
	}
	return addrList, nil
}

// 对于早期没有开lookup的nsqd服务，直接加载producer，不再通过LoadNsqds函数通过lookup去加载
func (x *XNsq) LoadProducer() {
	producerAddrs := viper.GetStringSlice(x.Env + ".producerAddrs")
	// 添加nsqd
	for _, v := range producerAddrs {
		producer, err := nsq.NewProducer(v, nsq.NewConfig())
		if err != nil {
			// TODO: add error log
			continue
		}
		x.producers = append(x.producers, producer)
	}
}

func (x *XNsq) loadNsqds() error {
	nsqdAddrs, err := x.lookupNsqdAddrs()
	if err != nil {
		return err
	}
	newAddrsMap := make(map[string]struct{})
	for _, nsqdAddr := range nsqdAddrs {
		newAddrsMap[nsqdAddr] = struct{}{}
	}

	x.mutex.Lock()

	oldAddrsMap := make(map[string]struct{})
	removeAddrMap := make(map[string]struct{})
	addAddrMap := make(map[string]struct{})
	// 统计旧nsqd和待移除的nsqd地址
	for _, p := range x.producers {
		v := p.String()
		oldAddrsMap[v] = struct{}{}
		if _, ok := newAddrsMap[v]; !ok {
			removeAddrMap[v] = struct{}{}
		}
	}
	// 统计待添加的nsqd地址
	for k := range newAddrsMap {
		if _, ok := oldAddrsMap[k]; !ok {
			addAddrMap[k] = struct{}{}
		}
	}

	// 移除nsqd
	if len(removeAddrMap) > 0 {
		idx := 0
		for _, p := range x.producers {
			if _, ok := removeAddrMap[p.String()]; !ok {
				x.producers[idx] = p
				idx++
			}
		}
		x.producers = x.producers[:idx]
	}
	// 添加nsqd
	for v := range addAddrMap {
		producer, err := nsq.NewProducer(v, nsq.NewConfig())
		if err != nil {
			// TODO: add error log
			continue
		}
		x.producers = append(x.producers, producer)
	}

	x.mutex.Unlock()
	return nil
}

// LoadNsqds load balance for nsqds
func (x *XNsq) LoadNsqds() {
	// 服务启动时必须同步调用一次，否则开始发消息时x.producer可能为空
	err := x.loadNsqds()
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			_ = x.loadNsqds()
			time.Sleep(10 * time.Second)
		}
	}()
}

func (x *XNsq) Publish(topic string, body []byte) (err error) {
	var p *nsq.Producer
	x.mutex.RLock()
	defer x.mutex.RUnlock()
	if len(x.producers) == 0 {
		err = errors.New(errNsqEmptyPool)
		return
	}
	p = x.producers[rand.Intn(len(x.producers))]
	err = p.Publish(topic, body)
	if err != nil {
		// try other producer
		for _, producer := range x.producers {
			err = producer.Publish(topic, body)
			if err == nil {
				return
			}
		}
	}
	return
}

func (x *XNsq) PublishContext(ctx context.Context, topic string, body []byte) (err error) {
	var p *nsq.Producer

	x.mutex.RLock()
	defer x.mutex.RUnlock()
	if len(x.producers) == 0 {
		err = errors.New(errNsqEmptyPool)
		return
	}
	p = x.producers[rand.Intn(len(x.producers))]
	err = p.Publish(topic, body)
	if err != nil {
		// try other producer
		for _, producer := range x.producers {
			err = producer.Publish(topic, body)
			if err == nil {
				return
			}
		}
	}
	return
}

func (x *XNsq) DeferredPublish(topic string, delay time.Duration, body []byte) (err error) {
	var p *nsq.Producer
	x.mutex.RLock()
	defer x.mutex.RUnlock()
	if len(x.producers) == 0 {
		err = errors.New(errNsqEmptyPool)
		return
	}
	p = x.producers[rand.Intn(len(x.producers))]
	err = p.DeferredPublish(topic, delay, body)
	if err != nil {
		// try other producer
		for _, producer := range x.producers {
			err = producer.DeferredPublish(topic, delay, body)
			if err == nil {
				return
			}
		}
	}
	return
}

func (x *XNsq) DeferredPublishContext(ctx context.Context, topic string, delay time.Duration, body []byte) (err error) {
	var p *nsq.Producer

	x.mutex.RLock()
	defer x.mutex.RUnlock()
	if len(x.producers) == 0 {
		err = errors.New(errNsqEmptyPool)
		return
	}
	p = x.producers[rand.Intn(len(x.producers))]
	err = p.DeferredPublish(topic, delay, body)
	if err != nil {
		// try other producer
		for _, producer := range x.producers {
			err = producer.DeferredPublish(topic, delay, body)
			if err == nil {
				return
			}
		}
	}
	return
}

func (x *XNsq) StartConsumer(topic, channel string, h nsq.HandlerFunc) (err error) {
	if !viper.IsSet(x.Env) {
		err = errors.New(fmt.Sprintf(errNsqConfigNil, x.Env))
		return
	}
	config := nsq.NewConfig()
	config.DialTimeout = viper.GetDuration(x.Env+".dialTimeout") * time.Millisecond
	config.ReadTimeout = viper.GetDuration(x.Env+".readTimeout") * time.Millisecond
	config.WriteTimeout = viper.GetDuration(x.Env+".writeTimeout") * time.Millisecond
	config.MaxInFlight = viper.GetInt(x.Env + ".maxInFlight")
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	concurrentHandlers := viper.GetInt(x.Env + ".concurrentHandlers")
	consumer.AddConcurrentHandlers(h, concurrentHandlers)
	addrs := viper.GetStringSlice(x.Env + ".comsumerAddrs")
	err = consumer.ConnectToNSQLookupds(addrs)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	x.consumers = append(x.consumers, consumer)
	return
}

func (x *XNsq) StopConsumers() {
	for _, consumer := range x.consumers {
		consumer.Stop()
		<-consumer.StopChan
	}
}
