package formatter

import (
	"fmt"
	"hcc/cello/lib/config"
	"hcc/cello/model"
	"strings"
	"sync"
)

type lun struct {
	Path    string
	Size    int
	UseType string
}

// Clusterdomain is configuration object field.
type Clusterdomain struct {
	TargetName string
	Lun        []lun
}

// IscsiMap is the key-value configuration object about iscsi config data
type IscsiMap struct {
	Domain map[string]*Clusterdomain
	lock   sync.RWMutex
}

//VolObjectMap : For iscsi conf object
// var VolObjectMap IscsiMap
var VolObjectMap *IscsiMap

//GlobalVolumesDB : DB info
var GlobalVolumesDB []model.Volume

func init() {
	// VolObjectMap.Domain = make(map[string]*Clusterdomain)
	VolObjectMap = New()
}

//New : For iscsi config field, Generate New Object
func New() *IscsiMap {
	return &IscsiMap{Domain: make(map[string]*Clusterdomain)}
}

//DevPathBuilder : For iscsi configuration, zvol path builder
func DevPathBuilder(volume model.Volume) string {
	return "/dev/zvol/" + volume.Pool + "/" + strings.ToUpper(volume.Filesystem) + "-" + strings.ToUpper(volume.UseType) + "-" + volume.ServerUUID
}

//VolNameBuilder : Build Vol name consist of pool and server uuid
func VolNameBuilder(volume model.Volume) string {
	return volume.Pool + "/" + strings.ToUpper(volume.Filesystem) + "-" + strings.ToUpper(volume.UseType) + "-" + volume.ServerUUID
}

//PreLoad : vol data
func (m *IscsiMap) PreLoad() {
	for _, args := range GlobalVolumesDB {
		if VolObjectMap.Domain[args.ServerUUID] == nil {
			VolObjectMap.PutDomain(args.ServerUUID)
		}
		if VolObjectMap.Domain[args.ServerUUID] != nil {
			VolObjectMap.SetIscsiLun(args, DevPathBuilder(args))
		}
		args.Pool = config.VolumeConfig.VOLUMEPOOL
	}
}

//SetIscsiLun : local var has iscsi config(Target Name , Lun number)
func (m *IscsiMap) SetIscsiLun(volume model.Volume, volPath string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if volume.ServerUUID != "" && volume.UseType != "" {
		var templLun lun
		templLun.Path = volPath
		templLun.Size = volume.Size
		templLun.UseType = volume.UseType
		m.Domain[volume.ServerUUID].Lun = append(m.Domain[volume.ServerUUID].Lun, templLun)
	} else {
		// fmt.Println("object handler set iscsi error Check filesystem")
		fmt.Println("object handler put val err")
	}
}

//RemoveIscsiLun : local var has iscsi config(Target Name , Lun number)
func (m *IscsiMap) RemoveIscsiLun(volume model.Volume, lunOrder int) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if volume.ServerUUID != "" && volume.Filesystem != "" {
		if m.Domain[volume.ServerUUID] != nil {
			m.Domain[volume.ServerUUID].Lun = append(m.Domain[volume.ServerUUID].Lun[:lunOrder], m.Domain[volume.ServerUUID].Lun[lunOrder+1:]...)
			// All lun and Domain map delete
			// m.Domain[volume.ServerUUID].Lun = nil
			// delete(m.Domain, volume.ServerUUID)
		}
	} else {
		fmt.Println("object handler put val err")
	}
}

//GetIscsiData : make serveruuid map
func (m *IscsiMap) GetIscsiData(serveruuid string) *Clusterdomain {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var tempVal *Clusterdomain
	if serveruuid != "" && m.Domain[serveruuid] != nil {
		tempVal = m.Domain[serveruuid]
	} else {
		fmt.Println("object handler get val err")
	}
	return tempVal
}

// GetIscsiLunNum : Retrun Lun number
func (m *IscsiMap) GetIscsiLunNum(serveruuid string) *Clusterdomain {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var tempVal *Clusterdomain
	if serveruuid != "" && m.Domain[serveruuid] != nil {
		tempVal = m.Domain[serveruuid]
	} else {
		fmt.Println("object handler get val err")
	}
	return tempVal
}

//PutDomain : make serveruuid map
func (m *IscsiMap) PutDomain(serveruuid string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if serveruuid != "" && m.Domain[serveruuid] == nil {
		m.Domain[serveruuid] = new(Clusterdomain)
		m.Domain[serveruuid].TargetName = serveruuid
	} else {
		fmt.Println("object handler put val err > ", serveruuid)
	}
}

//GetDomain : make serveruuid map
func (m *IscsiMap) GetDomain(serveruuid string) *Clusterdomain {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var tempVal *Clusterdomain
	if serveruuid != "" && m.Domain[serveruuid] != nil {
		tempVal = m.Domain[serveruuid]
	} else {
		fmt.Println("object handler get val err")
	}
	return tempVal
}

//RemoveDomain : make serveruuid map
func (m *IscsiMap) RemoveDomain(serveruuid string) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if serveruuid != "" && VolObjectMap.Domain[serveruuid] != nil {
		delete(VolObjectMap.Domain, serveruuid)
	} else {
		fmt.Println("object handler There is not exist map ", serveruuid)
	}
}

//MapSize : return map size
func (m *IscsiMap) MapSize(serveruuid string) int {
	m.lock.RLock()
	defer m.lock.RUnlock()
	return len(m.Domain)
}
