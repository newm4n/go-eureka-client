package go_eureka_client

import (
	"errors"
	"log"
	"sync"
)

type ApplicationInstanceBalancer interface {
	GetInstance(AppName string) *Instance
}

type ApplicationInstanceSet struct {
	mutex       sync.Mutex
	InstanceMap map[string][]*Instance
}

func (set *ApplicationInstanceSet) UpdateInstanceList(appName string, instances ...*Instance) {
	set.mutex.Lock()
	defer set.mutex.Unlock()
	delete(set.InstanceMap, appName)
	arr := make([]*Instance, 0)
	for _, i := range instances {
		arr = append(arr, i)
	}
	set.InstanceMap[appName] = arr
}

type RoundRobinInstanceBalancer struct {
	ApplicationInstanceSet
	IndexMap map[string]int
}

func NewRoundRobinInstanceBalancer() *RoundRobinInstanceBalancer {
	ret := &RoundRobinInstanceBalancer{
		ApplicationInstanceSet: ApplicationInstanceSet{
			InstanceMap: make(map[string][]*Instance),
		},
		IndexMap: make(map[string]int),
	}
	return ret
}

func (balancer *RoundRobinInstanceBalancer) GetInstance(AppName string) (*Instance, error) {
	balancer.mutex.Lock()
	defer balancer.mutex.Unlock()
	if len(balancer.InstanceMap[AppName]) == 0 {
		log.Printf("WARN : App %s has no instance", AppName)
		return nil, errors.New("app have no instance")
	}
	index := balancer.IndexMap[AppName]
	count := 0
	for {
		count++
		index++
		if count > len(balancer.InstanceMap[AppName]) {
			break
		}
		if index >= len(balancer.InstanceMap[AppName]) {
			//log.Printf("RESET %s -- %d >= %d", AppName, index, len(balancer.InstanceMap[AppName]))
			index = 0
		}
		if balancer.InstanceMap[AppName][index].Status == STATUS_UP {
			//log.Printf("%s -> %d", AppName, index)
			balancer.IndexMap[AppName] = index
			return balancer.InstanceMap[AppName][index], nil
		}
	}
	log.Printf("WARN : App %s. no instance is up", AppName)
	return nil, errors.New("no instance is up")
}
