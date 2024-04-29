package model

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"reflect"
	"time"
)

var tableNamePrefix = "tv2okx_"

var registeredModels = map[string]Interface{}

// Interface model interface
type Interface interface {
	TableName() string
	ShortTableName() string
}

// RegisterModel register model
func RegisterModel(models ...Interface) {
	for _, model := range models {
		if _, exist := registeredModels[model.TableName()]; exist {
			panic(fmt.Errorf("model table name %s conflict", model.TableName()))
		}
		registeredModels[model.TableName()] = model
	}
}

// GetRegisterModels will return the register models
func GetRegisterModels() map[string]Interface {
	return registeredModels
}

// JSONStruct json struct, same with runtime.RawExtension
type JSONStruct map[string]interface{}

// NewJSONStructByString new json struct from string
func NewJSONStructByString(source string) (*JSONStruct, error) {
	if source == "" {
		return nil, nil
	}
	var data JSONStruct
	err := json.Unmarshal([]byte(source), &data)
	if err != nil {
		return nil, fmt.Errorf("parse raw data failure %w", err)
	}
	return &data, nil
}

// NewJSONStructByStruct new json struct from struct object
func NewJSONStructByStruct(object interface{}) (*JSONStruct, error) {
	if object == nil {
		return nil, nil
	}
	var data JSONStruct
	out, err := yaml.Marshal(object)
	if err != nil {
		return nil, fmt.Errorf("marshal object data failure %w", err)
	}
	if err := yaml.Unmarshal(out, &data); err != nil {
		return nil, fmt.Errorf("unmarshal object data failure %w", err)
	}
	return &data, nil
}

// JSON Encoded as a JSON string
func (j *JSONStruct) JSON() string {
	b, err := json.Marshal(j)
	if err != nil {
		log.Logger.Errorf("json marshal failure %s", err.Error())
	}
	return string(b)
}

// Properties return the map
func (j *JSONStruct) Properties() map[string]interface{} {
	return *j
}

// RawExtension Encoded as a RawExtension

// BaseModel common model
type BaseModel struct {
	CreateTime time.Time `json:"createTime"`
	UpdateTime time.Time `json:"updateTime"`
}

// SetCreateTime set create time
func (m *BaseModel) SetCreateTime(time time.Time) {
	m.CreateTime = time
}

// SetUpdateTime set update time
func (m *BaseModel) SetUpdateTime(time time.Time) {
	m.UpdateTime = time
}

func deepCopy(src interface{}) interface{} {
	dst := reflect.New(reflect.TypeOf(src).Elem())

	val := reflect.ValueOf(src).Elem()
	nVal := dst.Elem()
	for i := 0; i < val.NumField(); i++ {
		nvField := nVal.Field(i)
		nvField.Set(val.Field(i))
	}

	return dst.Interface()
}
