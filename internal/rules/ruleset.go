package rules

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/alegrey91/fwdctl/pkg/iptables"
	"gopkg.in/yaml.v2"
)

func NewRuleSet() *RuleSet {
	return &RuleSet{
		Rules: make(map[string]iptables.Rule),
	}
}

// NewRuleSet return the struct that contains informations about rules
func NewRuleSetFromFile(path string) (*RuleSet, error) {
	// Read rules from file
	rulesFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Retrieve rules to fill RuleSet
	rules := supportRuleSet{}
	err = yaml.Unmarshal(rulesFile, &rules)
	if err != nil {
		return nil, err
	}

	// Fill RuleSet with rules taken from file
	rs := NewRuleSet()
	for _, rule := range rules.Rules {
		ruleHash := hash(rule)
		rs.Rules[ruleHash] = rule
	}
	return rs, nil
}

func hash(rule iptables.Rule) string {
	md5.New()
	ruleString := fmt.Sprintf("%s%s%d%s%d",
		rule.Iface,
		rule.Proto,
		rule.Dport,
		rule.Saddr,
		rule.Sport,
	)
	hash := md5.Sum([]byte(ruleString))
	return hex.EncodeToString(hash[:])
}

func (rs *RuleSet) GetHash(rule iptables.Rule) string {
	return hash(rule)
}

func (rs *RuleSet) Add(rule iptables.Rule) {
	ruleHash := hash(rule)
	rs.Rules[ruleHash] = rule
}

func (rs *RuleSet) Remove(ruleHash string) {
	delete(rs.Rules, ruleHash)
}

// ApplyChanges add and remove rules based on the differences
// between the old and current rules set.
// Return an error in case of fail, nil otherwise.
func (rs *RuleSet) ApplyChanges(oldRS *RuleSet) error {
	// loop over old rules set, to find rules to be removed
	for hash := range oldRS.Rules {
		// if key in oldRules is not present in rs,
		// then the old rule must be removed
		if _, ok := rs.Rules[hash]; !ok {
			err := iptables.DeleteForwardByRule(
				oldRS.Rules[hash].Iface,
				oldRS.Rules[hash].Proto,
				oldRS.Rules[hash].Dport,
				oldRS.Rules[hash].Saddr,
				oldRS.Rules[hash].Sport,
			)
			if err != nil {
				return fmt.Errorf("%v", err)
			}
		}
	}
	// loop over new rules set, to find rules to be added
	for hash := range rs.Rules {
		// if key in rs in not present in oldRs,
		// then the new rule must be added
		if _, ok := oldRS.Rules[hash]; !ok {
			err := iptables.CreateForward(
				rs.Rules[hash].Iface,
				rs.Rules[hash].Proto,
				rs.Rules[hash].Dport,
				rs.Rules[hash].Saddr,
				rs.Rules[hash].Sport,
			)
			if err != nil {
				return fmt.Errorf("%v", err)
			}
		}
	}
	return nil
}
