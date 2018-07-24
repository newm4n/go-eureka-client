package go_eureka_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	STATUS_UP             = "UP"
	STATUS_DOWN           = "DOWN"
	STATUS_STARTING       = "STARTING"
	STATUS_OUT_OF_SERVICE = "OUT_OF_SERVICE"
	STATUS_UNKNOWN        = "UNKNOWN"

	DATACENTER_OWN    = "MyOwn"
	DATACENTER_AMAZON = "Amazon"
)

type EurekaResponse struct {
	Application  Application  `json:"application,omitempty"`
	Applications Applications `json:"applications,omitempty"`
	Instance     Instance     `json:"applications,omitempty"`
}

type EurekaRequest struct {
	Instance *Instance `json:"instance"`
}

type Applications struct {
	VersionDelta string        `json:"versions__delta,omitempty"`
	AppsHashCode string        `json:"apps__hashcode,omitempty"`
	Application  []Application `json:"application"`
}

type Application struct {
	Name      string     `json:"name,omitempty"`
	Instances []Instance `json:"instance,omitempty"`
	Indice    int        `json:"-"`
}

func (app *Application) NextInstance() *Instance {
	current := app.Indice
	for {
		current++
		if current >= len(app.Instances) {
			current = 0
		}
		if app.Instances[current].Status == STATUS_UP {
			app.Indice = current
			return &app.Instances[current]
		} else if current == app.Indice {
			return nil
		}
	}
}

type Instance struct {
	mutex                         sync.Mutex        `json:"-"`
	client                        http.Client       `json:"-"`
	discoveryServer               *DiscoveryServer  `json:"-"`
	InstanceId                    string            `json:"instanceId,omitempty"`
	HostName                      string            `json:"hostName"`
	App                           string            `json:"app"`
	IpAddr                        string            `json:"ipAddr"`
	Status                        string            `json:"status"`
	OverriddenStatus              string            `json:"overriddenStatus,omitempty"`
	Port                          PortInfo          `json:"port,omitempty"`
	SecurePort                    PortInfo          `json:"securePort,omitempty"`
	CountryId                     int               `json:"countryId,omitempty"`
	DataCenterInfo                DataCenterInfo    `json:"dataCenterInfo,omitempty"`
	LeaseInfo                     LeaseInfo         `json:"leaseInfo,omitempty"`
	HomePageUrl                   string            `json:"homePageUrl,omitempty"`
	StatusPageUrl                 string            `json:"statusPageUrl,omitempty"`
	HealthCheckUrl                string            `json:"healthCheckUrl,omitempty"`
	VipAddress                    string            `json:"vipAddress,omitempty"`
	SecureVipAddress              string            `json:"secureVipAddress,omitempty"`
	IsCoordinatingDiscoveryServer string            `json:"isCoordinatingDiscoveryServer,omitempty"`
	LastUpdatedTimestamp          string            `json:"lastUpdatedTimestamp,omitempty"`
	LastDirtyTimestamp            string            `json:"lastDirtyTimestamp,omitempty"`
	Metadata                      map[string]string `json:"metadata,omitempty"`
	ActionType                    string            `json:"actionType,omitempty"`
}

func NewInstance(app, hostname, ip string, port, securePort int) *Instance {
	i := &Instance{
		InstanceId: fmt.Sprintf("%s:%s:%d", ip, app, port),
		HostName:   hostname,
		App:        strings.ToUpper(app),
		IpAddr:     ip,
		Port: PortInfo{
			PortNo:  port,
			Enabled: "true",
		},
		SecurePort: PortInfo{
			PortNo:  securePort,
			Enabled: "true",
		},
		Status: STATUS_STARTING,
		DataCenterInfo: DataCenterInfo{
			Name:  DATACENTER_OWN,
			Class: "com.netflix.appinfo.InstanceInfo$DefaultDataCenterInfo",
		},
		VipAddress:                    app,
		SecureVipAddress:              app,
		IsCoordinatingDiscoveryServer: "false",
	}
	if port == 0 {
		i.Port.Enabled = "false"
	}
	if securePort == 0 {
		i.SecurePort.Enabled = "false"
	}
	return i
}

func (i *Instance) Register(ds *DiscoveryServer) error {
	i.discoveryServer = ds
	if err := i.sendStatus(ds, STATUS_STARTING); err != nil {
		i.discoveryServer = nil
		return err
	}
	go i.emitHeartBeat()
	return nil
}

func (i *Instance) NowUp() error {
	return i.sendStatus(i.discoveryServer, STATUS_UP)
}

func (i *Instance) NowPause() error {
	return i.sendStatus(i.discoveryServer, STATUS_OUT_OF_SERVICE)
}

func (i *Instance) NowDown() error {
	return i.sendStatus(i.discoveryServer, STATUS_DOWN)
}

func (i *Instance) NowStarting() error {
	return i.sendStatus(i.discoveryServer, STATUS_STARTING)
}

func (i *Instance) NowStall() error {
	return i.sendStatus(i.discoveryServer, STATUS_UNKNOWN)
}

func (i *Instance) UnRegister() error {
	log.Printf("Eureka Shutting Down (%s)", i.InstanceId)
	return i.sendInstanceUpdate(i.discoveryServer, "DELETE")
}

func (i *Instance) emitHeartBeat() {
	for {
		if i.discoveryServer != nil {
			time.Sleep(10 * time.Second)
			i.HeartBeat()
		} else {
			break
		}
	}
}

func (i *Instance) HeartBeat() error {
	log.Printf("Eureka Heartbeating (%s - %s)", i.InstanceId, i.Status)
	return i.sendInstanceUpdate(i.discoveryServer, "PUT")
}

func (i *Instance) sendInstanceUpdate(ds *DiscoveryServer, method string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if ds == nil {
		return errors.New("this instance have not been registered")
	}

	req, err := http.NewRequest(method, fmt.Sprintf("http://%s:%d/eureka/apps/%s/%s", i.discoveryServer.Host, i.discoveryServer.Port, i.App, i.InstanceId), nil)
	if err != nil {
		return err
	}
	resp, err := i.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("unexpected response code %d : %s", resp.StatusCode, string(body)))
	}
	return nil
}

func (i *Instance) sendStatus(ds *DiscoveryServer, status string) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	if ds == nil {
		return errors.New("this instance have not been registered")
	}

	i.Status = status
	log.Printf("Eureka Update (%s - %s)", i.InstanceId, i.Status)
	reqdata := &EurekaRequest{i}

	reqBytes, err := json.Marshal(reqdata)
	if err != nil {
		return nil
	}

	//log.Print(string(reqBytes))

	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:%d/eureka/apps/%s", ds.Host, ds.Port, i.App), bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := i.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New(fmt.Sprintf("unexpected response code %d : %s", resp.StatusCode, string(body)))
	}

	return nil
}

func (i *Instance) String() string {
	return fmt.Sprintf("APP : %s, IP : %s, PORT : %d, STATUS : %s", i.App, i.IpAddr, i.Port.PortNo, i.Status)
}

type PortInfo struct {
	PortNo  int    `json:"$"`
	Enabled string `json:"@enabled"`
}

type DataCenterInfo struct {
	Name  string `json:"name"`
	Class string `json:"@class"`
}

type LeaseInfo struct {
	RenewalIntervalInSecs int   `json:"renewalIntervalInSecs,omitempty"`
	DurationInSecs        int   `json:"durationInSecs,omitempty"`
	RegistrationTimestamp int64 `json:"registrationTimestamp,omitempty"`
	LastRenewalTimestamp  int64 `json:"lastRenewalTimestamp,omitempty"`
	EvictionTimeStamp     int64 `json:"evictionTimeStamp,omitempty"`
	ServiceUpTimestamp    int64 `json:"serviceUpTimestamp,omitempty"`
}

type DiscoveryServer struct {
	Host string
	Port int
}

func NewDiscoveryServer(host string, port int) *DiscoveryServer {
	return &DiscoveryServer{
		Host: host,
		Port: port,
	}
}

func (ds *DiscoveryServer) FetchApplication(name string) (*EurekaResponse, error) {
	url := fmt.Sprintf("http://%s:%d/eureka/apps/%s", ds.Host, ds.Port, name)
	log.Println(url)
	return ds.fetch(url)
}

func (ds *DiscoveryServer) FetchAllApplications() (*EurekaResponse, error) {
	url := fmt.Sprintf("http://%s:%d/eureka/apps", ds.Host, ds.Port)
	log.Println(url)
	return ds.fetch(url)
}

func (ds *DiscoveryServer) FetchInstance(app string, instanceId string) (*EurekaResponse, error) {
	url := fmt.Sprintf("http://%s:%d/eureka/apps/%s/%s", ds.Host, ds.Port, app, instanceId)
	log.Println(url)
	return ds.fetch(url)
}

func (ds *DiscoveryServer) fetch(url string) (*EurekaResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, errors.New(fmt.Sprintf("unexpected response code %d", resp.StatusCode))
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	eurekaResponse := &EurekaResponse{}
	if err = json.Unmarshal(body, eurekaResponse); err != nil {
		return nil, err
	}
	return eurekaResponse, nil
}
