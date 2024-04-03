package rules

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/alegrey91/fwdctl/pkg/iptables"
	"gopkg.in/yaml.v2"
)

func NewRuleSet() *RuleSet {
	return &RuleSet{
		Rules: make(map[string]iptables.Rule),
	}
}

// NewRuleSet return the struct that contains informations about rules
func NewRuleSetFromFile(file io.Reader) (*RuleSet, error) {
	// Read rules from file
	rulesFile, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// Retrieve rules to fill RuleSet
	rules := supportRuleSet{}
	err = yaml.Unmarshal(rulesFile, &rules)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling content: %v", err)
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

type RuleSetDiff struct {
	ToRemove []*iptables.Rule
	ToAdd    []*iptables.Rule
}

// Diff method returns a *RuleSetDiff struct.
// It contains a list of Rule(s) to be added / remove
// in order to achieve the new RuleSet state.
func Diff(oldRS, newRS *RuleSet) *RuleSetDiff {
	ruleSetDiff := &RuleSetDiff{}
	// loop over old rules set, to find rules to be removed
	for hash := range oldRS.Rules {
		// if key in oldRules is not present in rs,
		// then the old rule must be removed
		if _, ok := newRS.Rules[hash]; !ok {
			rule := oldRS.Rules[hash]
			ruleSetDiff.ToRemove = append(ruleSetDiff.ToRemove, &rule)
		}
	}
	// loop over new rules set, to find rules to be added
	for hash := range newRS.Rules {
		// if key in rs in not present in oldRs,
		// then the new rule must be added
		if _, ok := oldRS.Rules[hash]; !ok {
			rule := newRS.Rules[hash]
			ruleSetDiff.ToAdd = append(ruleSetDiff.ToAdd, &rule)
		}
	}
	return ruleSetDiff
}
