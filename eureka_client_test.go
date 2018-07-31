package go_eureka_client

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestDiscoveryClient_FetchInstance(t *testing.T) {
	ds := NewDiscoveryClient("localhost", 8761, "discovery", "discovery")
	resp, err := ds.FetchApplication("CONFIG")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Instance : %s", resp.Application.Instances[0].String())
}

func TestDiscoveryClient_FetchInstances(t *testing.T) {
	ds := NewDiscoveryClient("localhost", 8761, "discovery", "discovery")
	resp, err := ds.FetchAllApplications()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("Instance : %s", resp.Applications.Application[0].Instances[0].String())
}

func TestInstance_Register(t *testing.T) {
	instance := NewInstance("cemungudh", "cemungudh", "10.10.10.2", 1234, 4321)
	ds := NewDiscoveryClient("localhost", 8761, "discovery", "discovery")
	log.Print("Now Registering")
	if err := instance.Register(ds); err != nil {
		t.Fatalf("Failed register %s", err)
	}
	time.Sleep(4 * time.Second)
	log.Print("Now UP")
	if err := instance.NowUp(); err != nil {
		t.Fatalf("Failed UP %s", err)
	}
	time.Sleep(20 * time.Second)
	log.Print("Now UNREG")
	if err := instance.UnRegister(); err != nil {
		t.Fatalf("Failed UNREG %s", err)
	}
}
