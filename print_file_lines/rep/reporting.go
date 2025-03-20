package rep

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func EWL(msg interface{}) error {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		return fmt.Errorf("%v", msg) // Fallback without location
	}

	file = filepath.Base(file) // Only the file name, no directory structure

	// Convert msg to string if it's an error
	var errMsg string
	switch v := msg.(type) {
	case string:
		errMsg = v
	case error:
		errMsg = v.Error()
	default:
		errMsg = fmt.Sprintf("%v", v)
	}

	return fmt.Errorf("%s:%d: %s", file, line, errMsg)
}
