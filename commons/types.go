package commons

import "errors"

type Raw []byte

func (r *Raw) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(r)
}

func (r Raw) MarshalYAML() (interface{}, error) {
	if r == nil {
		return []byte("null"), nil
	}

	return r, nil
}

func (r Raw) MarshalJSON() ([]byte, error) {
	if r == nil {
		return []byte("null"), nil
	}
	return r, nil
}

func (r *Raw) UnmarshalJSON(data []byte) error {
	if r == nil {
		return errors.New("UnmarshalJSON on nil pointer")
	}
	*r = append((*r)[0:0], data...)
	return nil
}
