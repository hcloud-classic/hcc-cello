package dao

import (
	"encoding/json"
	"fmt"
	"hcc/cello/lib/logger"
	"hcc/cello/lib/mysql"
	"hcc/cello/model"
	"testing"

	"github.com/TylerBrock/colorjson"
)

// func TestloadIscsiDB() {
// 	celloParams := make(map[string]interface{})
// 	celloParams["row"] = 999999
// 	celloParams["page"] = 1
// 	var volumes []model.Volume
// 	volumes, _ = ReadVolumeAll(celloParams)
// 	body, _ := json.Marshal(volumes)
// 	var obj map[string]interface{}
// 	json.Unmarshal([]byte(body), &obj)

// 	// Make a custom formatter with indent set
// 	f := colorjson.NewFormatter()
// 	f.Indent = 4

// 	// Marshall the Colorized JSON
// 	s, _ := f.Marshal(obj)
// 	// fmt.Println(string(s))
// 	logger.Logger.Println("doHcc Action [", string(s), "]")

// }

func Testmain(t *testing.T) error {
	err := mysql.Prepare()
	if err != nil {
		return err
	}
	testReadVolumeAll()
	return nil
}

func testReadVolumeAll() {
	celloParams := make(map[string]interface{})
	celloParams["row"] = 999999
	celloParams["page"] = 1
	var volumes []model.Volume
	qwe, errCode, recvStr := ReadVolumeAll(celloParams)

	if recvStr != "" {
		fmt.Println("test : ", errCode, recvStr)
	}
	volumes = qwe.([]model.Volume)
	body, _ := json.Marshal(volumes)
	var obj map[string]interface{}
	json.Unmarshal([]byte(body), &obj)

	// Make a custom formatter with indent set
	f := colorjson.NewFormatter()
	f.Indent = 4

	// Marshall the Colorized JSON
	s, _ := f.Marshal(obj)
	// fmt.Println(string(s))
	logger.Logger.Println("doHcc Action [", string(s), "]")
}
