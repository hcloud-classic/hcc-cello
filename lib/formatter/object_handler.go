package formatter

import (
	"fmt"
	"hcc/cello/model"
	"strconv"
	"strings"
	"sync"
)

type Pool struct {
	Size     string
	Free     string
	Capacity string
	Health   string
	Name     string
}

type Volpool struct {
	PoolMap map[string]*Pool
	lock    sync.RWMutex
}

type lun struct {
	UUID    string
	Path    string
	Size    int
	UseType string
	Pool    string
	Name    string
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

//PoolObjectMap : For iscsi conf object, it has zpool status
// var VolObjectMap IscsiMap
var PoolObjectMap *Volpool

//VolObjectMap : For iscsi conf object
var VolObjectMap *IscsiMap

//GlobalVolumesDB : DB info
var GlobalVolumesDB []model.Volume

func init() {
	// VolObjectMap.Domain = make(map[string]*Clusterdomain)
	VolObjectMap = NewVolObj()
	PoolObjectMap = NewPoolObj()
}

//NewPoolObj : For iscsi config field, Generate New Object
func NewPoolObj() *Volpool {
	return &Volpool{PoolMap: make(map[string]*Pool)}
}

//NewVolObj : For iscsi config field, Generate New Object
func NewVolObj() *IscsiMap {
	return &IscsiMap{Domain: make(map[string]*Clusterdomain)}
}

//DevPathBuilder : For iscsi configuration, zvol path builder
func DevPathBuilder(volume model.Volume) string {
	if strings.ToUpper(volume.UseType) == "OS" {
		return "/dev/zvol/" + volume.Pool + "/" + strings.ToLower(volume.Filesystem) + "-" + strings.ToUpper(volume.UseType) + "-" + volume.ServerUUID

	} else {
		return "/dev/zvol/" + volume.Pool + "/" + strings.ToLower(volume.Filesystem) + "-" + strings.ToUpper(volume.UseType) + "-" + volume.ServerUUID + "-"

	}
}

//VolNameBuilder : Build Vol name consist of pool and server uuid
func VolNameBuilder(volume model.Volume) string {
	return volume.Pool + "/" + strings.ToLower(volume.Filesystem) + "-" + strings.ToUpper(volume.UseType) + "-" + volume.ServerUUID
}

//PreLoad : vol data
func (m *IscsiMap) PreLoad() {
	for _, args := range GlobalVolumesDB {
		if VolObjectMap.Domain[args.ServerUUID] == nil {
			VolObjectMap.PutDomain(args.ServerUUID)
		}
		if VolObjectMap.Domain[args.ServerUUID] != nil {
			lun, _ := VolObjectMap.GetIscsiLun(args)
			if lun.UUID == "" {
				VolObjectMap.SetIscsiLun(args)
			}
		}
		// args.Pool = config.VolumeConfig.VOLUMEPOOL
	}
}

//PutPool : make serveruuid map
func (m *Volpool) PutPool(pool Pool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	// fmt.Println("PutPool => ", pool.Name, "\n\n", m.PoolMap[pool.Name])
	if pool.Name != "" {
		if m.PoolMap[pool.Name] == nil {

			m.PoolMap[pool.Name] = new(Pool)
		}
		m.PoolMap[pool.Name].Capacity = pool.Capacity
		m.PoolMap[pool.Name].Free = pool.Free
		m.PoolMap[pool.Name].Health = pool.Health
		m.PoolMap[pool.Name].Size = pool.Size
		m.PoolMap[pool.Name].Name = pool.Name
	} else {
		fmt.Println("Pool Obj Can't put in structure")
	}
}

//GetPool : make serveruuid map
func (m *Volpool) GetPool(poolname string) *Pool {
	m.lock.Lock()
	defer m.lock.Unlock()
	var tempPool *Pool
	if poolname != "" && m.PoolMap[poolname] == nil {
		tempPool = m.PoolMap[poolname]
	} else {
		fmt.Println("Pool Obj Can't put in structure")
		return nil
	}
	return tempPool
}

//SetIscsiLun : local var has iscsi config(Target Name , Lun number)
func (m *IscsiMap) SetIscsiLun(volume model.Volume) string {
	m.lock.Lock()
	defer m.lock.Unlock()
	if volume.ServerUUID != "" && volume.UseType != "" {

		var templLun lun
		var lunNuber int
		templLun.UUID = volume.UUID
		templLun.Path = DevPathBuilder(volume)

		templLun.Size = volume.Size
		templLun.UseType = volume.UseType
		templLun.Pool = volume.Pool
		templLun.Name = strings.Split(VolNameBuilder(volume), "/")[1]
		lunNuber = 0

		for range m.Domain[volume.ServerUUID].Lun {
			lunNuber++
		}
		if strings.ToUpper(volume.UseType) == "DATA" {
			templLun.Path += strconv.Itoa(lunNuber)
		}

		m.Domain[volume.ServerUUID].Lun = append(m.Domain[volume.ServerUUID].Lun, templLun)
		return strconv.Itoa(lunNuber)
	}
	return "object handler Setiscsi val err"
}

//To-Do : Remove Volume sequence
//RemoveIscsiLun : local var has iscsi config(Target Name , Lun number)
func (m *IscsiMap) RemoveIscsiLun(volume model.Volume, lunOrder int) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if volume.ServerUUID != "" && volume.Filesystem != "" {
		if m.Domain[volume.ServerUUID] != nil {
			if lunOrder < len(m.Domain[volume.ServerUUID].Lun) {
				if strings.ToUpper(volume.UseType) == "DATA" {
					m.Domain[volume.ServerUUID].Lun = append(m.Domain[volume.ServerUUID].Lun[:lunOrder], m.Domain[volume.ServerUUID].Lun[lunOrder+1:]...)

				}
			}
			// All lun and Domain map delete
			// m.Domain[volume.ServerUUID].Lun = nil
			// delete(m.Domain, volume.ServerUUID)
		}
	} else {
		fmt.Println("object handler put val err")
	}
}

//To Do : Remove Volume sequence
//InsertIscsiLun : local var has iscsi config(Target Name , Lun number)
func (m *IscsiMap) InsertIscsiLun(volume model.Volume, lunOrder int) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if volume.ServerUUID != "" && volume.Filesystem != "" {
		if m.Domain[volume.ServerUUID] != nil {
			if lunOrder < len(m.Domain[volume.ServerUUID].Lun) {
				m.Domain[volume.ServerUUID].Lun = append(m.Domain[volume.ServerUUID].Lun[:lunOrder+1], m.Domain[volume.ServerUUID].Lun[lunOrder:]...)
			}
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
		return nil
	}
	return tempVal
}
func (m *IscsiMap) GetIscsiLun(volume model.Volume) (lun, string) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var tempVal lun
	if volume.ServerUUID != "" && m.Domain[volume.ServerUUID] != nil {
		for i, args := range m.Domain[volume.ServerUUID].Lun {
			if args.UUID == volume.UUID {
				tempVal = args
				return tempVal, strconv.Itoa(i)
			}
		}
	}
	return tempVal, "GetIscsiLun : Can't find lun"
}

//PutDomain : make serveruuid map
func (m *IscsiMap) PutDomain(serveruuid string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if serveruuid != "" && m.Domain[serveruuid] == nil {
		m.Domain[serveruuid] = new(Clusterdomain)
		m.Domain[serveruuid].TargetName = serveruuid
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
