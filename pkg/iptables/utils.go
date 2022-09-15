package iptables

import (
	"sync"

	"github.com/coreos/go-iptables/iptables"
)

var once sync.Once

type single *iptables.IPTables

var singleInstance single

// getIPTablesInstance create a singletone instance for iptables.New()
func getIPTablesInstance() (*iptables.IPTables, error) {
	var err error
	if singleInstance == nil {
		once.Do(func() {
			singleInstance, err = iptables.New()
		})
	}

	return singleInstance, err
}
