package vpc

import (
	"net"

	"github.com/pritunl/mongo-go-driver/bson/primitive"
	"github.com/pritunl/pritunl-cloud/utils"
)

type VpcIp struct {
	Id       primitive.ObjectID `bson:"_id,omitempty"`
	Vpc      primitive.ObjectID `bson:"vpc"`
	Subnet   primitive.ObjectID `bson:"subnet"`
	Ip       int64              `bson:"ip"`
	Instance primitive.ObjectID `bson:"instance"`
}

func (i *VpcIp) GetIp() net.IP {
	return utils.Int2IpAddress(i.Ip * 2)
}

func (i *VpcIp) GetIps() (net.IP, net.IP) {
	return utils.IpIndex2Ip(i.Ip)
}
