package deploy

import (
	"github.com/pritunl/pritunl-cloud/state"
)

func Deploy(stat *state.State) (err error) {
	network := NewNetwork(stat)
	err = network.Deploy()
	if err != nil {
		return
	}

	ipset := NewIpset(stat)
	err = ipset.Deploy()
	if err != nil {
		return
	}

	iptables := NewIptables(stat)
	err = iptables.Deploy()
	if err != nil {
		return
	}

	err = ipset.Clean()
	if err != nil {
		return
	}

	disks := NewDisks(stat)
	err = disks.Deploy()
	if err != nil {
		return
	}

	instances := NewInstances(stat)
	err = instances.Deploy()
	if err != nil {
		return
	}

	namespaces := NewNamespace(stat)
	err = namespaces.Deploy()
	if err != nil {
		return
	}

	domains := NewDomains(stat)
	err = domains.Deploy()
	if err != nil {
		return
	}

	return
}
