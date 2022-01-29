package ipset

import (
	"errors"
	"strings"
	"sync"

	"github.com/nadoo/ipset"

	"github.com/nadoo/glider/rule"
)

// Manager struct.
type Manager struct {
	domainSet sync.Map
}

func addToSet(s, item string) error {
	if strings.IndexByte(item, '.') == -1 {
		return ipset.Add(s+"6", item)
	}
	return ipset.Add(s, item)
}

// NewManager returns a Manager
func NewManager(rules []*rule.Config) (*Manager, error) {
	if err := ipset.Init(); err != nil {
		return nil, err
	}

	m := &Manager{}
	sets := make(map[string]struct{})

	for _, r := range rules {
		if r.IPSet == "" {
			continue
		}

		if _, ok := sets[r.IPSet]; !ok {
			sets[r.IPSet] = struct{}{}
			ipset.Create(r.IPSet)
			ipset.Flush(r.IPSet)
			ipset.Create(r.IPSet+"6", ipset.OptIPv6())
			ipset.Flush(r.IPSet + "6")
		}

		for _, domain := range r.Domain {
			m.domainSet.Store(domain, r.IPSet)
		}
		for _, ip := range r.IP {
			addToSet(r.IPSet, ip)
		}
		for _, cidr := range r.CIDR {
			addToSet(r.IPSet, cidr)
		}
	}

	return m, nil
}

// AddDomainIP implements the dns AnswerHandler function, used to update ipset according to domainSet rule.
func (m *Manager) AddDomainIP(domain, ip string) error {
	if domain == "" || ip == "" {
		return errors.New("please specify the domain and ip address")
	}

	domain = strings.ToLower(domain)
	for i := len(domain); i != -1; {
		i = strings.LastIndexByte(domain[:i], '.')
		if setName, ok := m.domainSet.Load(domain[i+1:]); ok {
			addToSet(setName.(string), ip)
		}
	}
	return nil
}
