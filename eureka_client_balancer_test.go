package go_eureka_client

import (
	"testing"
)

func TestNewRoundRobinInstanceBalancer(t *testing.T) {
	bal := NewRoundRobinInstanceBalancer()
	i, e := bal.GetInstance("TEST")
	if i != nil || e == nil {
		t.Fatal()
	}
	if "app have no instance" != e.Error() {
		t.Fatal()
	}

	ia := NewInstance("TEST", "hosta", "1.1.1.1", 1000, 1001)
	ia.Status = STATUS_UP
	ib := NewInstance("TEST", "hostb", "1.1.1.2", 2000, 2001)
	ib.Status = STATUS_UP
	ic := NewInstance("TEST", "hostc", "1.1.1.3", 3000, 3001)
	ic.Status = STATUS_UP
	bal.UpdateInstanceList("TEST", ia, ib, ic)

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostb" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostc" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hosta" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostb" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostc" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hosta" != i.HostName {
		t.Fatal()
	}
}

func TestNewRoundRobinInstanceBalancer2(t *testing.T) {
	bal := NewRoundRobinInstanceBalancer()
	i, e := bal.GetInstance("TEST")
	if i != nil || e == nil {
		t.Fatal()
	}
	if "app have no instance" != e.Error() {
		t.Fatal()
	}

	ia := NewInstance("TEST", "hosta", "1.1.1.1", 1000, 1001)
	ia.Status = STATUS_DOWN
	ib := NewInstance("TEST", "hostb", "1.1.1.2", 2000, 2001)
	ib.Status = STATUS_UP
	ic := NewInstance("TEST", "hostc", "1.1.1.3", 3000, 3001)
	ic.Status = STATUS_UP
	bal.UpdateInstanceList("TEST", ia, ib, ic)

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostb" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostc" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostb" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostc" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostb" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostc" != i.HostName {
		t.Fatal()
	}
}

func TestNewRoundRobinInstanceBalancer3(t *testing.T) {
	bal := NewRoundRobinInstanceBalancer()
	i, e := bal.GetInstance("TEST")
	if i != nil || e == nil {
		t.Fatal()
	}
	if "app have no instance" != e.Error() {
		t.Fatal()
	}

	ia := NewInstance("TEST", "hosta", "1.1.1.1", 1000, 1001)
	ia.Status = STATUS_DOWN
	ib := NewInstance("TEST", "hostb", "1.1.1.2", 2000, 2001)
	ib.Status = STATUS_DOWN
	ic := NewInstance("TEST", "hostc", "1.1.1.3", 3000, 3001)
	ic.Status = STATUS_UP
	bal.UpdateInstanceList("TEST", ia, ib, ic)

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostc" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostc" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostc" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostc" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostc" != i.HostName {
		t.Fatal()
	}

	i, e = bal.GetInstance("TEST")
	if i == nil || e != nil || "hostc" != i.HostName {
		t.Fatal()
	}
}

func TestNewRoundRobinInstanceBalancer4(t *testing.T) {
	bal := NewRoundRobinInstanceBalancer()
	i, e := bal.GetInstance("TEST")
	if i != nil || e == nil {
		t.Fatal()
	}
	if "app have no instance" != e.Error() {
		t.Fatal()
	}

	ia := NewInstance("TEST", "hosta", "1.1.1.1", 1000, 1001)
	ia.Status = STATUS_DOWN
	ib := NewInstance("TEST", "hostb", "1.1.1.2", 2000, 2001)
	ib.Status = STATUS_DOWN
	ic := NewInstance("TEST", "hostc", "1.1.1.3", 3000, 3001)
	ic.Status = STATUS_DOWN
	bal.UpdateInstanceList("TEST", ia, ib, ic)

	i, e = bal.GetInstance("TEST")
	if i != nil || e == nil {
		t.Fatal()
	}
	if "no instance is up" != e.Error() {
		t.Fatal()
	}
}
