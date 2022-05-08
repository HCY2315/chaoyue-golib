package microservice

import "fmt"

//Identity 唯一界定一个微服务实例，没有命名成ID是为了不与ServiceID混淆
type Identity struct {
	Group
	Region
	ServerType
	ID
	Version
}

func (i Identity) String() string {
	return fmt.Sprintf("%s_%s_%s_%s", i.ServerType, i.Region.Id, i.ID, i.Version)
}

func (i Identity) EqualWith(group string, region string, id string) bool {
	return i.Group == Group(group) && i.Region.Id == region && i.ID == ID(id)
}

type ServerType string

func (st ServerType) String() string {
	return string(st)
}

type ID string

type Version string

type Group string

const (
	DefaultGroup Group = "default"
)

type RegionType string

const (
	RegionTypeCenter RegionType = "center_domain"
	RegionTypeLeaf   RegionType = "leaf_domain"
)

type Region struct {
	RegionType RegionType
	Id         string
}
