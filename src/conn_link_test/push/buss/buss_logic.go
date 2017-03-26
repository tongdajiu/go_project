package main

import (
	"conn_pool"
	//"fmt"
	//"logger"
	//"sync"
)

type ContentData struct {
	stat      bool
	push_data string
	uuid      int64
}
type BussLogic struct {
	pool          conn_pool.IConnPool
	uid_map_array []uuidData
}

func (this *BussLogic) Initialize() {

}
