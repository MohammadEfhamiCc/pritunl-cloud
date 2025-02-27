package spec

import (
	"net"
	"strconv"
	"strings"

	"github.com/pritunl/pritunl-cloud/errortypes"
)

type Firewall struct {
	Ingress []*Rule `bson:"ingress" json:"ingress"`
}

type Rule struct {
	Protocol  string      `bson:"protocol" json:"protocol"`
	Port      string      `bson:"port" json:"port"`
	SourceIps []string    `bson:"source_ips" json:"source_ips"`
	Sources   []*Refrence `bson:"sources" json:"sources"`
}

func (f *Firewall) Validate() (errData *errortypes.ErrorData, err error) {
	if f.Ingress == nil {
		f.Ingress = []*Rule{}
	}

	for _, rule := range f.Ingress {
		switch rule.Protocol {
		case All:
			rule.Port = ""
			break
		case Icmp:
			rule.Port = ""
			break
		case Tcp, Udp, Multicast, Broadcast:
			ports := strings.Split(rule.Port, "-")

			portInt, e := strconv.Atoi(ports[0])
			if e != nil {
				errData = &errortypes.ErrorData{
					Error:   "invalid_ingress_rule_port",
					Message: "Invalid ingress rule port",
				}
				return
			}

			if portInt < 1 || portInt > 65535 {
				errData = &errortypes.ErrorData{
					Error:   "invalid_ingress_rule_port",
					Message: "Invalid ingress rule port",
				}
				return
			}

			parsedPort := strconv.Itoa(portInt)
			if len(ports) > 1 {
				portInt2, e := strconv.Atoi(ports[1])
				if e != nil {
					errData = &errortypes.ErrorData{
						Error:   "invalid_ingress_rule_port",
						Message: "Invalid ingress rule port",
					}
					return
				}

				if portInt < 1 || portInt > 65535 || portInt2 <= portInt {
					errData = &errortypes.ErrorData{
						Error:   "invalid_ingress_rule_port",
						Message: "Invalid ingress rule port",
					}
					return
				}

				parsedPort += "-" + strconv.Itoa(portInt2)
			}

			rule.Port = parsedPort

			break
		default:
			errData = &errortypes.ErrorData{
				Error:   "invalid_ingress_rule_protocol",
				Message: "Invalid ingress rule protocol",
			}
			return
		}

		if rule.Sources == nil {
			rule.Sources = []*Refrence{}
		}

		if rule.SourceIps == nil {
			rule.SourceIps = []string{}
		}

		for i, sourceIp := range rule.SourceIps {
			if sourceIp == "" {
				errData = &errortypes.ErrorData{
					Error:   "invalid_ingress_rule_source_ip",
					Message: "Empty ingress rule source IP",
				}
				return
			}

			if !strings.Contains(sourceIp, "/") {
				if strings.Contains(sourceIp, ":") {
					sourceIp += "/128"
				} else {
					sourceIp += "/32"
				}
			}

			_, sourceCidr, e := net.ParseCIDR(sourceIp)
			if e != nil {
				errData = &errortypes.ErrorData{
					Error:   "invalid_ingress_rule_source_ip",
					Message: "Invalid ingress rule source IP",
				}
				return
			}

			rule.SourceIps[i] = sourceCidr.String()
		}

		if rule.Protocol == Multicast || rule.Protocol == Broadcast {
			rule.Sources = []*Refrence{}
			rule.SourceIps = []string{}
		}
	}

	return
}

type FirewallYaml struct {
	Name    string                `yaml:"name"`
	Kind    string                `yaml:"kind"`
	Ingress []FirewallYamlIngress `yaml:"ingress"`
}

type FirewallYamlIngress struct {
	Protocol string   `yaml:"protocol"`
	Port     string   `yaml:"port"`
	Source   []string `yaml:"source"`
}
