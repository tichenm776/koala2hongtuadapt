package util

import (
	"github.com/alecthomas/log4go"
	"github.com/robfig/cron"
	"sync"
	"time"
	"zhiyuan/zyutil"
)

var G_map_SubjectsLocal = Demo{
	Data: make(map[string]interface{},0),
	Lock: &sync.Mutex{},
	RLock:&sync.RWMutex{},
}
var G_map = Demo{
	Data: make(map[string]interface{},0),
	Lock: &sync.Mutex{},
	RLock:&sync.RWMutex{},
}

var G_map_GroupsLocal = Demo{
	Data: make(map[string]interface{},0),
	Lock: &sync.Mutex{},
	RLock:&sync.RWMutex{},
}

type Demo struct {
	Data map[string]interface{}
	Lock *sync.Mutex//golang struct 在赋值时是浅拷贝,会导致lock的不是同一份锁 所以需要加上指针操作
	RLock *sync.RWMutex//golang struct 在赋值时是浅拷贝,会导致lock的不是同一份锁 所以需要加上指针操作
}




func (d *Demo) Get(k string) interface{}{
	//d.Lock.Lock()
	//defer d.Lock.Unlock()
	d.RLock.RLock()
	defer d.RLock.RUnlock()
	return d.Data[k]
}

func (d *Demo) Set(k string,v interface{}) {
	//d.Lock.Lock()
	//defer d.Lock.Unlock()
	d.RLock.Lock()
	defer d.RLock.Unlock()
	d.Data[k]=v
	//log4go.Info("sety value",k,v)
}
func (d *Demo) Getmap()(map[string]interface{}) {
	defer zyutil.Recover()
	d.RLock.RLock()
	defer d.RLock.RUnlock()
	//d.Lock.Lock()
	//defer d.Lock.Unlock()
	copy_map := d.Data
	return copy_map
}
func (d *Demo) Clean() {
	d.RLock.Lock()
	defer d.RLock.Unlock()
	d.Data = nil //释放内存
	d.Data = make(map[string]interface{},0)
}
func (d *Demo) Delete(k string) {
	if d.Get(k) != nil{
		d.RLock.Lock()
		defer d.RLock.Unlock()
		delete(d.Data,k)
	}
}
func (d *Demo) Deletechild(k,k2 string) {
	if d.Get(k) != nil{
		d.RLock.Lock()
		defer d.RLock.Unlock()
		delete(d.Get(k).(map[string]interface{}),k2)
		//delete(d.Data,k)
	}
}

func (d *Demo) CronDelete(duration int64) {
	cronTarget := cron.New(cron.WithSeconds())
	spec := "0 */5 * * * ?"
	cronTarget.AddFunc(spec, func() {
		nowtime := time.Now().Unix()
		log4go.Info("cron start check value")
		G_data := d.Getmap()
		log4go.Info("G_data is",G_data)
		for k,_ :=range G_data{
			log4go.Info("key is ",k)
			for k1,v1 := range G_data[k].(map[string]interface{}){
				log4go.Info("v1",v1)
				if (nowtime - v1.(int64)) > 600{
					//delete(G_data[k].(map[string]interface{}),k1)
					G_map.Deletechild(k,k1)
				}
			}
			//G_map.Set(k,G_data)
		}
	})
	cronTarget.Start()
}