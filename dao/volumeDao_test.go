package dao

import (
	"encoding/json"
	"fmt"
	"hcc/cello/lib/logger"
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

func TestReadVolumeAll(t *testing.T) {
	celloParams := make(map[string]interface{})
	celloParams["row"] = 999999
	celloParams["page"] = 1
	var volumes []model.Volume
	qwe, err := ReadVolumeAll(celloParams)

	if err != nil {
		fmt.Println("test err")
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
