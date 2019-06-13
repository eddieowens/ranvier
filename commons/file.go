package commons

import "os"

type Filer interface {
	Exists(path string) bool
}

type filer struct {
}

func NewFiler() Filer {
	return &filer{}
}

func (f *filer) Exists(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
