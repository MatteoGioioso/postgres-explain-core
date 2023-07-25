package pkg

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
		getSpecificProperties: func(node Node) []Property {
			props := make([]Property, 0)

			if node[HEAP_FETCHES] != nil {
				props = append(props, Property{
					ID:          "heap_fetches",
					Name:        "Heat fetches",
					Type:        "float",
					ValueFloat:  ConvertStringToFloat64(node[HEAP_FETCHES].(string)),
					ValueString: "",
					Skip:        false,
					Kind:        quantity,
				})
			}

			return props
		},
	},
	INDEX_ONLY_SCAN: {
		RelationName: RELATION_NAME,
		Index:        INDEX_NAME,
		Filter:       FILTER,
		Condition:    INDEX_CONDITION,
		getSpecificProperties: func(node Node) []Property {
			props := make([]Property, 0)

			if node[HEAP_FETCHES] != nil {
				props = append(props, Property{
					ID:          "heap_fetches",
					Name:        "Heat fetches",
					Type:        "float",
					ValueFloat:  ConvertStringToFloat64(node[HEAP_FETCHES].(string)),
					ValueString: "",
					Skip:        false,
					Kind:        quantity,
				})
			}

			return props
		},
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

				props = append(props, Property{
					ID:          "sort_space_type",
					Name:        node[SORT_SPACE_TYPE].(string),
					Type:        "float",
					ValueFloat:  ConvertStringToFloat64(strVal),
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
