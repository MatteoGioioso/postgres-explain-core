package pkg

import "strconv"

var operationsMap = map[string]Operation{
	SEQUENTIAL_SCAN: {
		RelationName: RELATION_NAME,
		Filter:       FILTER,
	},
	INDEX_SCAN: {
		RelationName: RELATION_NAME,
		Index:        INDEX_NAME,
		Filter:       FILTER,
		Condition:    INDEX_CONDITION,
	},
	INDEX_ONLY_SCAN: {
		RelationName: RELATION_NAME,
		Index:        INDEX_NAME,
		Filter:       FILTER,
		Condition:    INDEX_CONDITION,
	},
	HASH_JOIN: {
		Filter:    "Join Filter",
		Condition: HASH_CONDITION_PROP,
	},
	SORT: {
		Key: SORT_KEY,
		getSpecificProperties: func(node Node) []Property {
			props := make([]Property, 0)

			if node[SORT_METHOD] != nil {
				props = append(props, Property{
					ID:          "sort_method",
					Name:        "Sorty method",
					Type:        "string",
					ValueFloat:  0,
					ValueString: node[SORT_METHOD].(string),
					Skip:        false,
					Kind:        "",
				})
			}

			if node[SORT_SPACE_TYPE] != nil {
				strVal := node[SORT_SPACE_USED].(string)
				float, err := strconv.ParseFloat(strVal, 64)
				if err != nil {
					panic(err)
				}

				props = append(props, Property{
					ID:          "sort_space_type",
					Name:        node[SORT_SPACE_TYPE].(string),
					Type:        "float",
					ValueFloat:  float,
					ValueString: "",
					Skip:        false,
					Kind:        disk_size,
				})
			}

			return props
		},
	},
	CTE_SCAN: {
		RelationName: CTE_NAME,
		Filter:       FILTER,
	},
	FUNCTION_SCAN: {
		RelationName: FUNCTION_NAME,
	},
	GROUP_AGGREGATE: {
		Key: GROUP_KEY,
	},
	BITMAP_HEAP_SCAN: {
		RelationName: RELATION_NAME,
		Condition:    "Recheck Cond",
	},
	"Default": {
		RelationName: RELATION_NAME,
		Index:        INDEX_NAME,
		Filter:       FILTER,
	},
}

var filtersMap = map[string]string{
	HASH_JOIN: ROWS_REMOVED_BY_JOIN_FILTER,
}
