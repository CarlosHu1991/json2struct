package json2struct

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type A struct {
	I int32  `json:"i"`
	N string `json:"n"`
}

type Hero struct {
	Id         int32         `json:"id"`
	Name       string        `json:"name"`
	TestArr    []int         `json:"testArr"`
	TestMap    map[int]int   `json:"testMap"`
	TestObjArr []*A          `json:"testObjArr"`
	TestObjMap map[string]*A `json:"testObjMap"`
}

func loadJSON(file string) ([]byte, error) {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return f, nil
}

//may change path
const path = "E:\\work\\sgyy\\server\\src\\github.com\\json2struct\\"

func TestArray() {
	jsonData, err := loadJSON(path + "testArray.json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("jsonString: %s\n", string(jsonData))
	json2struct := NewJson2Struct()
	ret, err := json2struct.UnmarshalJsonArray(jsonData, make([]*Hero, 0))
	if err != nil {
		fmt.Printf("UnmarshalJsonArray failed! error:%s", err)
		return
	}
	objArr := ret.([]*Hero)
	fmt.Printf("json arrry: %v\n", objArr)

	bts, _ := json.Marshal(objArr)
	fmt.Printf("jsonString: %s\n", string(bts))
}

func TestMap() {
	jsonData, err := loadJSON(path + "testMap.json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("jsonString: %s\n", string(jsonData))
	json2struct := NewJson2Struct()
	ret, err := json2struct.UnmarshalJsonMap(jsonData, make(map[int]*Hero))
	if err != nil {
		fmt.Printf("UnmarshalJsonArray failed! error:%s", err)
		return
	}
	objMap := ret.(map[int]*Hero)
	fmt.Printf("json map: %v\n", objMap)

	bts, _ := json.Marshal(objMap)
	fmt.Printf("jsonString: %s\n", string(bts))
}
