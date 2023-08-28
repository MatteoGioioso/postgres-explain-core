package pkg

import (
	"fmt"
)

var operationsMap = map[string]Operation{
	SEQUENTIAL_SCAN: {
		RelationName: RELATION_NAME,
		Filter:       FILTER,
		getWorkers:   getScanWorkers,
	},
	INDEX_SCAN: {
		RelationName: RELATION_NAME,
		Index:        INDEX_NAME,
		Filter:       FILTER,
		Condition:    INDEX_CONDITION,
		getWorkers:   getScanWorkers,
		getSpecificProperties: func(node Node) []Property {
			props := make([]Property, 0)

			if node[HEAP_FETCHES] != nil {
				props = append(props, Property{
					ID:          "heap_fetches",
					Name:        "Heap fetches",
					Type:        "float",
					ValueFloat:  ConvertToFloat64(node[HEAP_FETCHES]),
					ValueString: "",
					Skip:        false,
					Kind:        Quantity,
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
		getWorkers:   getScanWorkers,
		getSpecificProperties: func(node Node) []Property {
			props := make([]Property, 0)

			if node[HEAP_FETCHES] != nil {
				props = append(props, Property{
					ID:          "heap_fetches",
					Name:        "Heap fetches",
					Type:        "float",
					ValueFloat:  ConvertToFloat64(node[HEAP_FETCHES]),
					ValueString: "",
					Skip:        false,
					Kind:        Quantity,
				})
			}

			return props
		},
	},
	SORT: {
		Key:        SORT_KEY,
		getWorkers: getSortWorkers,
		getSpecificProperties: func(node Node) []Property {
			props := make([]Property, 0)

			if node[SORT_METHOD] != nil {
				props = append(props, Property{
					ID:          "sort_method",
					Name:        "Sort method",
					Type:        "string",
					ValueFloat:  0,
					ValueString: node[SORT_METHOD].(string),
					Skip:        false,
					Kind:        "",
				})
			}

			if node[SORT_SPACE_TYPE] != nil {
				property := Property{
					ID:          "sort_space_type",
					Name:        node[SORT_SPACE_TYPE].(string),
					Type:        "float",
					ValueString: "",
					Skip:        false,
					Kind:        DiskSize,
				}

				property.ValueFloat = ConvertToFloat64(node[SORT_SPACE_USED])

				props = append(props, property)
			}

			return props
		},
	},
	INCREMENTAL_SORT: {
		Key:        SORT_KEY,
		getWorkers: getSortWorkers,
		getSpecificProperties: func(node Node) []Property {
			props := make([]Property, 0)

			if node[PRESORTED_KEY] != nil {
				props = append(props, Property{
					ID:          "pre_sorted_key",
					Name:        PRESORTED_KEY,
					Type:        "string",
					ValueFloat:  0,
					ValueString: convertPropToString(node[PRESORTED_KEY]),
					Skip:        false,
					Kind:        "",
				})
			}

			// TODO add fullsort groups and finish with all the properties
			if node[PRE_SORTED_GROUPS] != nil {
				sortedGroups := node[PRE_SORTED_GROUPS].(map[string]interface{})

				value := fmt.Sprintf(
					"%v, Sort method: %v",
					sortedGroups["Group Count"],
					sortedGroups["Sort Methods Used"].([]interface{})[0],
				)

				props = append(props, Property{
					ID:          "pre_sorted_group",
					Name:        PRE_SORTED_GROUPS,
					Type:        "string",
					ValueFloat:  0,
					ValueString: value,
					Skip:        false,
					Kind:        "",
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
	"Finalize " + GROUP_AGGREGATE: {
		Key: GROUP_KEY,
	},
	"Partial " + GROUP_AGGREGATE: {
		Key: GROUP_KEY,
	},
	HASH_AGGREGATE: {
		Key: GROUP_KEY,
		getSpecificProperties: func(node Node) []Property {
			props := make([]Property, 0)

			if node[BATCHES] != nil {
				props = append(props, Property{
					ID:          "batches",
					Name:        "Batches",
					Type:        "float",
					ValueFloat:  ConvertToFloat64(node[BATCHES]),
					ValueString: "",
					Skip:        false,
					Kind:        Quantity,
				})
			}
			return props
		},
	},
	BITMAP_HEAP_SCAN: {
		RelationName: RELATION_NAME,
		Condition:    "Recheck Cond",
		Filter:       FILTER,
		getSpecificProperties: func(node Node) []Property {
			props := make([]Property, 0)
			if node[HEAP_BLOCKS] != nil {
				props = append(props, Property{
					ID:          "heap_blocks",
					Name:        "Heap Blocks",
					Type:        "string",
					ValueFloat:  0,
					ValueString: node[HEAP_BLOCKS].(string),
					Skip:        false,
					Kind:        "",
				})
			}

			return props
		},
	},
	BITMAP_INDEX_SCAN: {
		Index:     INDEX_NAME,
		Condition: INDEX_CONDITION,
	},
	NESTED_LOOP_JOIN: {
		Filter: JOIN_FILTER,
	},
	NESTED_LOOP: {
		Filter: JOIN_FILTER,
	},
	NESTED_LOOP_SEMI_JOIN: {
		Filter: JOIN_FILTER,
	},
	HASH_JOIN: {
		Filter:    JOIN_FILTER,
		Condition: HASH_CONDITION_PROP,
	},
	MERGE_JOIN: {
		Filter:    JOIN_FILTER,
		Condition: "Merge Cond",
	},
	"Default": {
		RelationName: RELATION_NAME,
		Index:        INDEX_NAME,
		Filter:       FILTER,
		Key:          GROUP_KEY,
	},
}

var getScanWorkers = func(node Node) [][]Property {
	props := make([][]Property, 0)
	if node[WORKERS] == nil {
		return props
	}

	for _, worker := range node[WORKERS].([]interface{}) {
		w := worker.(map[string]interface{})
		work := make([]Property, 0)

		work = append(work, Property{
			ID:          "worker_number",
			Name:        "Worker Number",
			Type:        "float",
			ValueFloat:  w["Worker Number"].(float64),
			ValueString: "",
			Skip:        false,
			Kind:        "",
		})
		work = append(work, Property{
			ID:          "actual_loops",
			Name:        ACTUAL_LOOPS,
			Type:        "float",
			ValueFloat:  ConvertToFloat64(w[ACTUAL_LOOPS]),
			ValueString: "",
			Skip:        false,
			Kind:        "",
		})
		work = append(work, Property{
			ID:          "actual_rows",
			Name:        ACTUAL_ROWS,
			Type:        "float",
			ValueFloat:  ConvertToFloat64(w[ACTUAL_ROWS]),
			ValueString: "",
			Skip:        false,
			Kind:        "",
		})
		work = append(work, Property{
			ID:          "actual_total_time",
			Name:        ACTUAL_TOTAL_TIME,
			Type:        "float",
			ValueFloat:  ConvertToFloat64(w[ACTUAL_TOTAL_TIME]),
			ValueString: "",
			Skip:        false,
			Kind:        Timing,
		})

		props = append(props, work)
	}

	return props
}

var getSortWorkers = func(node Node) [][]Property {
	props := make([][]Property, 0)
	if node[WORKERS] == nil {
		return props
	}

	for _, worker := range node[WORKERS].([]interface{}) {
		w := worker.(map[string]interface{})
		work := make([]Property, 0)

		work = append(work, Property{
			ID:          "worker_number",
			Name:        "Worker Number",
			Type:        "float",
			ValueFloat:  w["Worker Number"].(float64),
			ValueString: "",
			Skip:        false,
			Kind:        "",
		})
		work = append(work, Property{
			ID:          "sort_method",
			Name:        SORT_METHOD,
			Type:        "string",
			ValueFloat:  0,
			ValueString: w[SORT_METHOD].(string),
			Skip:        false,
			Kind:        "",
		})
		work = append(work, Property{
			ID:          "sort_space_used",
			Name:        SORT_SPACE_USED,
			Type:        "float",
			ValueFloat:  ConvertToFloat64(w[SORT_SPACE_USED]),
			ValueString: "",
			Skip:        false,
			Kind:        DiskSize,
		})
		work = append(work, Property{
			ID:          "sort_space_type",
			Name:        SORT_SPACE_TYPE,
			Type:        "string",
			ValueFloat:  0,
			ValueString: w[SORT_SPACE_TYPE].(string),
			Skip:        false,
			Kind:        "",
		})

		props = append(props, work)
	}

	return props
}

var filtersMap = map[string]string{
	HASH_JOIN:        ROWS_REMOVED_BY_JOIN_FILTER,
	NESTED_LOOP_JOIN: ROWS_REMOVED_BY_JOIN_FILTER,
	MERGE_JOIN:       ROWS_REMOVED_BY_JOIN_FILTER,
}

const (
	Timing   = Kind("timing")
	Quantity = Kind("quantity")
	DiskSize = Kind("disk_size")
	Blocks   = Kind("blocks")
)
