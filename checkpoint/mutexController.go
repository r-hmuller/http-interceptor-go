package checkpoint

import (
	"errors"
	"sync"
)

var mutexVariables = make(map[string]string)
var mutex = &sync.RWMutex{}

func UpdateVariable(key string, value string) {
	mutex.Lock()
	mutexVariables[key] = value
	mutex.Unlock()
}

func ReadVariable(key string) (string, error) {
	mutex.RLock()
	result, exists := mutexVariables[key]
	mutex.RUnlock()
	if exists {
		return result, nil
	}
	return "", errors.New("key not present")
}
