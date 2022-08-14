package iptables

import (
	"fmt"
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

// getRuleById retrieve iptable rule using Id
func getRuleById(ruleId int) (string, error) {
	ipt, err := getIPTablesInstance()
	if err != nil {
		return "", err
	}

	ruleList, err := ipt.ListWithCounters(fwdTable, fwdChain)
	if err != nil {
		return "", err
	}

	for id, rule := range ruleList {
		if id == ruleId {
			return rule, nil
		}
	}
	return "", fmt.Errorf("no rule found with id: %d", ruleId)
}
