package sync

import (
	"time"

	"github.com/pritunl/pritunl-cloud/constants"
	"github.com/pritunl/pritunl-cloud/database"
	"github.com/pritunl/pritunl-cloud/deploy"
	"github.com/pritunl/pritunl-cloud/firewall"
	"github.com/pritunl/pritunl-cloud/instance"
	"github.com/pritunl/pritunl-cloud/iptables"
	"github.com/pritunl/pritunl-cloud/node"
	"github.com/pritunl/pritunl-cloud/state"
	"github.com/pritunl/pritunl-cloud/vpc"
	"github.com/sirupsen/logrus"
)

func deployState() (err error) {
	stat, err := state.GetState()
	if err != nil {
		return
	}

	err = deploy.Deploy(stat)
	if err != nil {
		return
	}

	return
}

func syncNodeFirewall() {
	db := database.GetDatabase()
	defer db.Close()

	if !node.Self.Firewall {
		iptables.UpdateState(node.Self, []*vpc.Vpc{}, []*instance.Instance{},
			[]string{}, nil, map[string][]*firewall.Rule{})
		return
	}

	for i := 0; i < 2; i++ {
		fires, err := firewall.GetRoles(db, node.Self.NetworkRoles)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("sync: Failed to get node firewall rules")
			return
		}

		ingress := firewall.MergeIngress(fires)

		iptables.UpdateStateRecover(node.Self, []*vpc.Vpc{},
			[]*instance.Instance{}, []string{}, ingress,
			map[string][]*firewall.Rule{})

		break
	}
}

func vmRunner() {
	time.Sleep(1 * time.Second)

	for {
		time.Sleep(1 * time.Second)
		if constants.Shutdown {
			return
		}

		if !node.Self.IsHypervisor() {
			syncNodeFirewall()
			continue
		}

		break
	}

	logrus.WithFields(logrus.Fields{
		"production": constants.Production,
	}).Info("sync: Starting hypervisor")

	for {
		time.Sleep(2500 * time.Millisecond)
		if constants.Shutdown {
			return
		}

		if !node.Self.IsHypervisor() {
			syncNodeFirewall()
			continue
		}

		err := deployState()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
			}).Error("sync: Failed to deploy state")
			continue
		}
	}
}

func initVm() {
	go vmRunner()
}
