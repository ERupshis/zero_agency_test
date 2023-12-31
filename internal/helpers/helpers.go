package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/erupshis/zero_agency_test/internal/logger"
)

// ExecuteWithLogError support method for defer functions call which should return error.
func ExecuteWithLogError(callback func() error, log logger.BaseLogger) {
	if err := callback(); err != nil {
		log.Info("callback execution finished with error: %v", err)
	}
}

// InterfaceToString simple converter any interface into string.
func InterfaceToString(i interface{}) string {
	return fmt.Sprintf("%v", i)
}

func UnmarshalData(body []byte, dst interface{}) error {
	return json.Unmarshal(body, dst)
}

func NormalizeInt64(val *int64) {
	if *val <= 0 {
		*val = 0
	}
}

func RemoveDuplicatesInt(nums []int64) []int64 {
	uniqueMap := make(map[int64]bool)
	var result []int64

	for _, num := range nums {
		if !uniqueMap[num] {
			result = append(result, num)
			uniqueMap[num] = true
		}
	}

	return result
}
