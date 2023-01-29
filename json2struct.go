package json2struct

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
)

/**
- Convert json to struct structure as much as possible.
- String and number are ignored.
*/

func Version() string {
	return "0.1.0"
}

// NewJson2Struct Create unmarshal object
//
// Attr-OnlyUseStructTag: The value of OnlyUseStructTag defaults to true. If there is no 'Tag' value, the Field will not be parsed.
func NewJson2Struct() *Json2Struct {
	m := new(Json2Struct)
	m.OnlyUseStructTag = true
	return m
}

type Json2Struct struct {
	OnlyUseStructTag bool
}

// UnmarshalJsonArray Default function of unmarshal json data of array, like "[{},{},{}]"
//
// If you need to adjust the details, the way to call NewJson2Struct
func UnmarshalJsonArray(jsonData []byte, sliceObj interface{}) (interface{}, error) {
	m := NewJson2Struct()
	return m.UnmarshalJsonArray(jsonData, sliceObj)
}

// UnmarshalJsonArray Unmarshal json data of array, like "[{},{},{}]"
func (m *Json2Struct) UnmarshalJsonArray(jsonData []byte, sliceObj interface{}) (interface{}, error) {
	var dataArr []interface{}
	d := json.NewDecoder(bytes.NewReader(jsonData))
	d.UseNumber()
	err := d.Decode(&dataArr)
	if err != nil {
		return nil, err
	}

	sliceObjType := reflect.TypeOf(sliceObj)
	elemSlice := reflect.MakeSlice(sliceObjType, 0, len(dataArr))
	err = m.decodeArray(sliceObjType.Elem(), &elemSlice, dataArr)
	if err != nil {
		return nil, err
	}
	return elemSlice.Interface(), err
}

// UnmarshalJsonMap Default function of unmarshal json data of array, like "{{},{},{}}"
//
// If you need to adjust the details, the way to call NewJson2Struct
func UnmarshalJsonMap(jsonData []byte, mapObj interface{}) (interface{}, error) {
	m := NewJson2Struct()
	return m.UnmarshalJsonMap(jsonData, mapObj)
}

// UnmarshalJsonMap Unmarshal json data of array, like "{{},{},{}}"
func (m *Json2Struct) UnmarshalJsonMap(jsonData []byte, mapObj interface{}) (interface{}, error) {
	var dataMap map[string]interface{}
	d := json.NewDecoder(bytes.NewReader(jsonData))
	d.UseNumber()
	err := d.Decode(&dataMap)
	if err != nil {
		return nil, err
	}

	mapObjType := reflect.TypeOf(mapObj)
	elemMap := reflect.MakeMap(reflect.MapOf(mapObjType.Key(), mapObjType.Elem()))
	err = m.decodeMap(mapObjType.Key(), mapObjType.Elem(), elemMap, dataMap)
	if err != nil {
		return nil, err
	}
	return elemMap.Interface(), err
}

func (m *Json2Struct) decodeJson(sType reflect.Type, sVal reflect.Value, data interface{}) error {
	var err error
	if dataMap, ok := data.(map[string]interface{}); ok {
		if sType.Kind() == reflect.Struct {
			err = m.decodeStruct(sType, sVal, dataMap)
		} else if sType.Kind() == reflect.Map {
			err = m.decodeMap(sType.Key(), sType.Elem(), sVal, dataMap)
		} else {
			return errors.New(fmt.Sprintf("The json data of map type cannot be converted into %s type value", sType.Kind()))
		}
	} else if dataArray, ok := data.([]interface{}); ok {
		if sType.Kind() == reflect.Slice {
			err = m.decodeArray(sType, &sVal, dataArray)
		} else {
			return errors.New(fmt.Sprintf("The json data of array type cannot be converted into %s type value", sType.Kind()))
		}
	} else {
		err = setValue(sVal, data)
	}

	return err
}

func (m *Json2Struct) decodeStruct(sType reflect.Type, sVal reflect.Value, dataMap map[string]interface{}) error {
	for i := 0; i < sType.NumField(); i++ {
		sTypeField := sType.Field(i)
		sValField := sVal.Field(i)
		tag := sTypeField.Tag.Get("json") //struct must need json tag
		if !m.OnlyUseStructTag && tag == "" {
			tag = sTypeField.Name
		}

		data, has := dataMap[tag]
		if !has {
			continue //no value
		}
		dataType := reflect.TypeOf(data)
		fieldType := sValField.Type()

		if fieldType.Kind() == reflect.Slice && dataType.Kind() == reflect.Slice {
			//[] json
			dataArr := data.([]interface{})
			elemType := sTypeField.Type.Elem() //Call elem() to get the element type of the array
			//create array
			elemSlice := reflect.MakeSlice(reflect.SliceOf(elemType), 0, len(dataArr))
			err := m.decodeArray(elemType, &elemSlice, dataArr)
			if err != nil {
				return err
			}
			sValField.Set(elemSlice)
		} else if fieldType.Kind() == reflect.Map && dataType.Kind() == reflect.Map {
			// {} json
			dataMap, ok := data.(map[string]interface{})
			if !ok {
				return errors.New("There is a map structure in the json data, but it cannot be converted into a map[string]interface ")
			}
			keyType := sTypeField.Type.Key()    //map key's type
			valueType := sTypeField.Type.Elem() //map value's type
			//create map
			elemMap := reflect.MakeMapWithSize(reflect.MapOf(keyType, valueType), len(dataMap))
			err := m.decodeMap(keyType, valueType, elemMap, dataMap)
			if err != nil {
				return err
			}
			sValField.Set(elemMap)
		} else {
			e := m.decodeValue(sValField, data)
			if e != nil {
				return e
			}
		}
	}
	return nil
}

func (m *Json2Struct) decodeMap(keyType reflect.Type, valueType reflect.Type, mapValue reflect.Value, dataMap map[string]interface{}) error {
	for keyData, valueData := range dataMap {
		//Create key, the key must be a normal variable, so use decodeValue() to get key value
		key := reflect.New(keyType).Elem()
		err := m.decodeValue(key, keyData)
		if err != nil {
			return errors.New(fmt.Sprintf("decode map key failed！%s", err.Error()))
		}

		//Create elem object
		elem, canSetElem, canSetElemType := newElem(valueType)
		e := m.decodeJson(canSetElemType, canSetElem, valueData)
		if e != nil {
			return e
		}

		mapValue.SetMapIndex(key, elem)
	}

	return nil
}

func (m *Json2Struct) decodeArray(elemType reflect.Type, arrayValue *reflect.Value, dataArray []interface{}) error {
	for _, d := range dataArray {
		//create elem object
		elem, canSetElem, canSetElemType := newElem(elemType)

		e := m.decodeJson(canSetElemType, canSetElem, d)
		if e != nil {
			return e
		}

		//append
		*arrayValue = reflect.Append(*arrayValue, elem)
	}

	return nil
}

func (m *Json2Struct) decodeValue(sVal reflect.Value, data interface{}) error {
	dataType := reflect.TypeOf(data)
	fieldType := sVal.Type()

	if fieldType == dataType {
		sVal.Set(reflect.ValueOf(data))
	} else {
		if dataType.ConvertibleTo(fieldType) {
			sVal.Set(reflect.ValueOf(data).Convert(fieldType))
		} else {
			err := setValue(sVal, data)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// NewElem Create a variable of type valueType.
//
// @valueElem: Real variable.
//
// @canSetValueElem: Pointers are ignored, variables whose value can be set.
//
// @canSetElemType: The type of canSetValueElem.
//
func newElem(valueType reflect.Type) (valueElem, canSetValueElem reflect.Value, canSetElemType reflect.Type) {
	if valueType.Kind() == reflect.Ptr {
		valueElem = reflect.New(valueType.Elem())
		canSetValueElem = valueElem.Elem()
		canSetElemType = valueType.Elem()
	} else {
		valueElem = reflect.New(valueType).Elem()
		canSetValueElem = valueElem
		canSetElemType = valueType
	}

	return valueElem, canSetValueElem, canSetElemType
}
