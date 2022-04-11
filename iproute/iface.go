package iproute

import (
	"encoding/json"

	"github.com/dropbox/godropbox/errors"
	"github.com/pritunl/pritunl-cloud/errortypes"
	"github.com/pritunl/pritunl-cloud/utils"
)

type Iface struct {
	Name  string `json:"ifname"`
	State string `json:"operstate"`
}

func IfaceGetAll(namespace string) (ifaces []*Iface, err error) {
	ifaces = []*Iface{}

	var output string
	if namespace != "" {
		output, err = utils.ExecOutputLogged(
			nil,
			"ip", "netns", "exec", namespace,
			"ip", "--json", "--brief",
			"link", "show",
		)
	} else {
		output, err = utils.ExecOutputLogged(
			nil,
			"ip", "--json", "--brief",
			"link", "show",
		)
	}
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(output), &ifaces)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "iproute: Failed to prase ifaces"),
		}
		return
	}

	return
}

func IfaceGetBridges(namespace string) (ifaces []*Iface, err error) {
	ifaces = []*Iface{}

	var output string
	if namespace != "" {
		output, err = utils.ExecOutputLogged(
			nil,
			"ip", "netns", "exec", namespace,
			"ip", "--json", "--brief",
			"link", "show",
			"type", "bridge",
		)
	} else {
		output, err = utils.ExecOutputLogged(
			nil,
			"ip", "--json", "--brief",
			"link", "show",
			"type", "bridge",
		)
	}
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(output), &ifaces)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "iproute: Failed to prase bridges"),
		}
		return
	}

	return
}

func IfaceGetBridgeIfaces(namespace, bridge string) (
	ifaces []*Iface, err error) {

	ifaces = []*Iface{}

	var output string
	if namespace != "" {
		output, err = utils.ExecCombinedOutputLogged(
			[]string{
				"does not exist",
			},
			"ip", "netns", "exec", namespace,
			"ip", "--json", "--brief",
			"link", "show",
			"master", bridge,
		)
	} else {
		output, err = utils.ExecCombinedOutputLogged(
			[]string{
				"does not exist",
			},
			"ip", "--json", "--brief",
			"link", "show",
			"master", bridge,
		)
	}
	if err != nil {
		return
	}

	if output == "" {
		return
	}

	err = json.Unmarshal([]byte(output), &ifaces)
	if err != nil {
		err = &errortypes.ParseError{
			errors.Wrap(err, "iproute: Failed to prase bridges"),
		}
		return
	}

	return
}
