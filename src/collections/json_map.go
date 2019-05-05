package collections

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func NewJsonMap(json []byte) JsonMap {
	return JsonMap{
		jsonData: json,
		json:     gjson.ParseBytes(json),
	}
}

type JsonMap struct {
	jsonData []byte
	json     gjson.Result
}

func (j *JsonMap) GetRaw() []byte {
	return j.jsonData
}

func (j *JsonMap) UnmarshalJSON(data []byte) error {
	j.jsonData = data
	j.json = gjson.ParseBytes(data)
	return nil
}

func (j JsonMap) MarshalJSON() ([]byte, error) {
	return j.jsonData, nil
}

func (j *JsonMap) GetAll() interface{} {
	return j.json.Value()
}

func (j *JsonMap) Get(key string) (interface{}, bool) {
	r := j.json.Get(key)
	return r.Value(), r.Exists()
}

func (j *JsonMap) Set(key string, value interface{}) {
	if data, err := sjson.SetBytes(j.jsonData, key, value); err == nil {
		j.overwrite(data)
	}
}

func (j *JsonMap) Overwrite(json []byte) {
	if gjson.ValidBytes(json) {
		j.overwrite(json)
	}
}

func (j *JsonMap) overwrite(json []byte) {
	j.jsonData = json
	j.json = gjson.ParseBytes(json)
}
