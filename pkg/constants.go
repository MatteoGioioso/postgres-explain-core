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
		Key:    SORT_KEY,
		Method: SORT_METHOD,
		getSpecificProperties: func(node Node) map[string]interface{} {
			props := make(map[string]interface{})
			if node[SORT_SPACE_TYPE] != nil {
				props[SORT_SPACE_TYPE] = node[SORT_SPACE_TYPE].(string)
			}

			if node[SORT_SPACE_USED] != nil {
				props[SORT_SPACE_USED] = node[SORT_SPACE_USED].(string)
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
