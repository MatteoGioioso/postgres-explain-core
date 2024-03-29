package pkg

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func GetRootNodeFromPlans(plans string) (Node, error) {
	p := Plans{}
	if err := json.Unmarshal([]byte(plans), &p); err != nil {
		return nil, fmt.Errorf("could not unmarshal plan: %v", err)
	}

	return p[0].Plan, nil
}

func getMaxBlocksRead(rootNode Node) float64 {
	sum := 0.0
	if rootNode[SHARED_READ_BLOCKS] != nil {
		sum += rootNode[SHARED_READ_BLOCKS].(float64)
	}
	if rootNode[TEMP_READ_BLOCKS] != nil {
		sum += rootNode[TEMP_READ_BLOCKS].(float64)
	}
	if rootNode[LOCAL_READ_BLOCKS] != nil {
		sum += rootNode[LOCAL_READ_BLOCKS].(float64)
	}

	return sum
}

func getMaxBlocksWritten(rootNode Node) float64 {
	sum := 0.0
	if rootNode[SHARED_WRITTEN_BLOCKS] != nil {
		sum += rootNode[SHARED_WRITTEN_BLOCKS].(float64)
	}
	if rootNode[TEMP_WRITTEN_BLOCKS] != nil {
		sum += rootNode[TEMP_WRITTEN_BLOCKS].(float64)
	}
	if rootNode[LOCAL_WRITTEN_BLOCKS] != nil {
		sum += rootNode[LOCAL_WRITTEN_BLOCKS].(float64)
	}

	return sum
}

func getMaxBlocksHits(rootNode Node) float64 {
	sum := 0.0
	if rootNode[LOCAL_HIT_BLOCKS] != nil {
		sum += rootNode[LOCAL_HIT_BLOCKS].(float64)
	}
	if rootNode[SHARED_HIT_BLOCKS] != nil {
		sum += rootNode[SHARED_HIT_BLOCKS].(float64)
	}
	return sum
}

func IsCTE(node Node) bool {
	return node[PARENT_RELATIONSHIP] == "InitPlan" && strings.HasPrefix(node[SUBPLAN_NAME].(string), "CTE")
}

func isSubPlan(node Node) bool {
	return node[PARENT_RELATIONSHIP] == "SubPlan" && strings.HasPrefix(node[SUBPLAN_NAME].(string), "SubPlan")
}

func ConvertStringToFloat64(val string) float64 {
	float, err := strconv.ParseFloat(val, 64)
	if err != nil {
		panic(err)
	}

	return float
}

func ConvertToFloat64(val interface{}) float64 {
	if val == nil {
		return 0.0
	}
	if isFloat64(val) {
		return val.(float64)
	} else {
		return ConvertStringToFloat64(val.(string))
	}
}

func ConvertScopeToString(prop interface{}) string {
	if prop == nil {
		return ""
	}

	switch r := prop.(type) {
	case string:
		return r
	case []interface{}: // When Sorting we can have an array of sorting keys
		marshal, err := json.MarshalIndent(r, "", "    ")
		if err != nil {
			panic(fmt.Errorf("could not marshal node operation scope into []string: %v", err))
		}
		return string(marshal)
	default:
		return ""
	}
}

func isFloat64(val interface{}) bool {
	typeOf := reflect.TypeOf(val).Kind()
	return typeOf == reflect.Float64
}

func isString(val interface{}) bool {
	typeOf := reflect.TypeOf(val).Kind()
	return typeOf == reflect.String
}

func isBool(val interface{}) bool {
	typeOf := reflect.TypeOf(val).Kind()
	return typeOf == reflect.Bool
}

func getEffectiveBlocksRead(node Node) float64 {
	sum := 0.0
	if node[EXCLUSIVE+LOCAL_READ_BLOCKS] != nil {
		sum += node[EXCLUSIVE+LOCAL_READ_BLOCKS].(float64)
	}
	if node[EXCLUSIVE+TEMP_READ_BLOCKS] != nil {
		sum += node[EXCLUSIVE+TEMP_READ_BLOCKS].(float64)
	}
	if node[EXCLUSIVE+SHARED_READ_BLOCKS] != nil {
		sum += node[EXCLUSIVE+SHARED_READ_BLOCKS].(float64)
	}
	return sum
}

func getEffectiveBlocksWritten(node Node) float64 {
	sum := 0.0
	if node[EXCLUSIVE+LOCAL_WRITTEN_BLOCKS] != nil {
		sum += node[EXCLUSIVE+LOCAL_WRITTEN_BLOCKS].(float64)
	}
	if node[EXCLUSIVE+TEMP_WRITTEN_BLOCKS] != nil {
		sum += node[EXCLUSIVE+TEMP_WRITTEN_BLOCKS].(float64)
	}
	if node[EXCLUSIVE+SHARED_WRITTEN_BLOCKS] != nil {
		sum += node[EXCLUSIVE+SHARED_WRITTEN_BLOCKS].(float64)
	}
	return sum
}

func getEffectiveBlocksHits(node Node) float64 {
	sum := 0.0
	if node[EXCLUSIVE+LOCAL_HIT_BLOCKS] != nil {
		sum += node[EXCLUSIVE+LOCAL_HIT_BLOCKS].(float64)
	}
	if node[EXCLUSIVE+SHARED_HIT_BLOCKS] != nil {
		sum += node[EXCLUSIVE+SHARED_HIT_BLOCKS].(float64)
	}
	return sum
}

func getRowsRemovedByFilter(node Node) float64 {
	op := node[NODE_TYPE].(string)
	removedByFilter := 0.0
	filter, ok := filtersMap[op]
	if ok {
		removedByFilter = node[filter+REVISED].(float64)
	}
	if removedByFilter == 0.0 {
		removedByFilter = node[ROWS_REMOVED_BY_FILTER+REVISED].(float64)
	}

	return removedByFilter
}
