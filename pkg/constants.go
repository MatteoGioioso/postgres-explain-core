package pkg

var operationsMap = map[string]Operation{
	SEQUENTIAL_SCAN: {
		RelationName: RELATION_NAME,
		Filter:       FILTER,
		getWorkers:   getGenericWorkers,
	},
	INDEX_SCAN: {
		RelationName: RELATION_NAME,
		Index:        INDEX_NAME,
		Filter:       FILTER,
		Condition:    INDEX_CONDITION,
		getWorkers:   getGenericWorkers,
		getSpecificProperties: func(node Node) []Property {
			props := make([]Property, 0)

			if node[HEAP_FETCHES] != nil {
				props = append(props, Property{
					ID:         "heap_fetches",
					Name:       "Heap fetches",
					Type:       "float",
					ValueFloat: ConvertToFloat64(node[HEAP_FETCHES]),
					Kind:       Quantity,
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
		getWorkers:   getGenericWorkers,
		getSpecificProperties: func(node Node) []Property {
			props := make([]Property, 0)

			if node[HEAP_FETCHES] != nil {
				props = append(props, Property{
					ID:         "heap_fetches",
					Name:       "Heap fetches",
					Type:       "float",
					ValueFloat: ConvertToFloat64(node[HEAP_FETCHES]),
					Kind:       Quantity,
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
			props = getSortProperties(node, props)
			return props
		},
	},
	INCREMENTAL_SORT: {
		Key:        SORT_KEY,
		getWorkers: getSortWorkers,
		getSpecificProperties: func(node Node) []Property {
			props := make([]Property, 0)
			props = getSortProperties(node, props)

			if node[PRESORTED_KEY] != nil {
				props = append(props, Property{
					ID:          "pre_sorted_key",
					Name:        PRESORTED_KEY,
					Type:        "string",
					ValueString: ConvertScopeToString(node[PRESORTED_KEY]),
				})
			}

			if node[PRE_SORTED_GROUPS] != nil {
				sortedGroups := node[PRE_SORTED_GROUPS].(map[string]interface{})

				props = append(props, Property{
					ID:         "pre_sort_group_count",
					Name:       "Pre Sort Group Count",
					Type:       "float",
					ValueFloat: ConvertToFloat64(sortedGroups["Group Count"]),
					Kind:       Quantity,
				})

				if sortedGroups["Sort Space Memory"] != nil {
					sortSpaceMemory := sortedGroups["Sort Space Memory"].(map[string]interface{})
					props = append(props, Property{
						ID:         "average_pre_sort_space_used",
						Name:       "Pre Sort Average Space Used",
						Type:       "float",
						ValueFloat: ConvertToFloat64(sortSpaceMemory["Average Sort Space Used"]),
						Kind:       DiskSize,
					})
					props = append(props, Property{
						ID:         "peak_pre_sort_space_used",
						Name:       "Pre Sort Peak Space Used",
						Type:       "float",
						ValueFloat: ConvertToFloat64(sortSpaceMemory["Peak Sort Space Used"]),
						Kind:       DiskSize,
					})
				}
			}

			if node[FULL_SORT_GROUPS] != nil {
				sortedGroups := node[FULL_SORT_GROUPS].(map[string]interface{})

				props = append(props, Property{
					ID:         "full_group_count",
					Name:       "Full Sort Group Count",
					Type:       "float",
					ValueFloat: ConvertToFloat64(sortedGroups["Group Count"]),
					Kind:       Quantity,
				})

				if sortedGroups["Sort Space Memory"] != nil {
					sortSpaceMemory := sortedGroups["Sort Space Memory"].(map[string]interface{})
					props = append(props, Property{
						ID:         "average_full_sort_space_used",
						Name:       "Full Average Sort Space Used",
						Type:       "float",
						ValueFloat: ConvertToFloat64(sortSpaceMemory["Average Sort Space Used"]),
						Kind:       DiskSize,
					})
					props = append(props, Property{
						ID:         "peak_full_sort_space_used",
						Name:       "Full Sort Peak Space Used",
						Type:       "float",
						ValueFloat: ConvertToFloat64(sortSpaceMemory["Peak Sort Space Used"]),
						Kind:       DiskSize,
					})
				}
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
		Key:        GROUP_KEY,
		getWorkers: getHashWorkers,
	},
	"Partial " + HASH_AGGREGATE: {
		Key:        GROUP_KEY,
		getWorkers: getHashWorkers,
	},
	HASH: {
		getSpecificProperties: func(node Node) []Property {
			props := make([]Property, 0)
			props = append(props, hashBucketsAndBatches(node, props)...)

			return props
		},
		getWorkers: getGenericWorkers,
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
					ValueString: node[HEAP_BLOCKS].(string),
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
		Filter:     JOIN_FILTER,
		getWorkers: getGenericWorkers,
	},
	NESTED_LOOP: {
		Filter:     JOIN_FILTER,
		getWorkers: getGenericWorkers,
	},
	NESTED_LOOP_SEMI_JOIN: {
		Filter: JOIN_FILTER,
	},
	HASH_JOIN: {
		Filter:     JOIN_FILTER,
		Condition:  HASH_CONDITION_PROP,
		getWorkers: getGenericWorkers,
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
		getWorkers:   getGenericWorkers,
	},
}

var getGenericWorkers = func(node Node) [][]Property {
	props := make([][]Property, 0)
	if node[WORKERS] == nil {
		return props
	}

	for _, worker := range node[WORKERS].([]interface{}) {
		w := worker.(map[string]interface{})
		work := make([]Property, 0)

		work = append(work, Property{
			ID:         "worker_number",
			Name:       "Worker Number",
			Type:       "float",
			ValueFloat: w["Worker Number"].(float64),
		})
		props = append(props, getGenericWorkerProperties(w, work))
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
			ID:         "worker_number",
			Name:       "Worker Number",
			Type:       "float",
			ValueFloat: w["Worker Number"].(float64),
		})
		props = append(props, getSortProperties(w, getGenericWorkerProperties(w, work)))
	}

	return props
}

var getHashWorkers = func(node Node) [][]Property {
	workers := make([][]Property, 0)

	if node[WORKERS] == nil {
		return workers
	}

	for _, worker := range node[WORKERS].([]interface{}) {
		w := worker.(map[string]interface{})
		work := make([]Property, 0)
		work = append(work, Property{
			ID:         "worker_number",
			Name:       "Worker Number",
			Type:       "float",
			ValueFloat: w["Worker Number"].(float64),
		})

		workers = append(workers, hashBucketsAndBatches(w, getGenericWorkerProperties(w, work)))
	}

	return workers
}

func getGenericWorkerProperties(w Node, work []Property) []Property {
	work = append(work, Property{
		ID:         "actual_loops",
		Name:       ACTUAL_LOOPS,
		Type:       "float",
		ValueFloat: ConvertToFloat64(w[ACTUAL_LOOPS]),
		Kind:       Quantity,
	})
	work = append(work, Property{
		ID:         "actual_rows",
		Name:       ACTUAL_ROWS,
		Type:       "float",
		ValueFloat: ConvertToFloat64(w[ACTUAL_ROWS]),
		Kind:       Quantity,
	})
	work = append(work, Property{
		ID:         "actual_total_time",
		Name:       ACTUAL_TOTAL_TIME,
		Type:       "float",
		ValueFloat: ConvertToFloat64(w[ACTUAL_TOTAL_TIME]),
		Kind:       Timing,
	})

	return work
}

func hashBucketsAndBatches(node Node, props []Property) []Property {
	if node["Memory Usage"] != nil {
		props = append(props, Property{
			ID:         "memory_usage",
			Name:       "Memory Usage",
			Type:       "float",
			ValueFloat: ConvertToFloat64(node["Memory Usage"]),
			Kind:       DiskSize,
		})
	}

	if node["Disk Usage"] != nil {
		props = append(props, Property{
			ID:         "disk_usage",
			Name:       "Disk Usage",
			Type:       "float",
			ValueFloat: ConvertToFloat64(node["Disk Usage"]),
			Kind:       DiskSize,
		})
	}

	if node[BATCHES] != nil {
		props = append(props, Property{
			ID:         "batches",
			Name:       BATCHES,
			Type:       "float",
			ValueFloat: ConvertToFloat64(node[BATCHES]),
			Kind:       Quantity,
		})
	}

	if node[BATCHES+" Originally"] != nil {
		props = append(props, Property{
			ID:         "batches_originally",
			Name:       BATCHES + " Originally",
			Type:       "float",
			ValueFloat: ConvertToFloat64(node[BATCHES+" Originally"]),
			Kind:       Quantity,
		})
	}

	if node["Buckets"] != nil {
		props = append(props, Property{
			ID:         "buckets",
			Name:       "Buckets",
			Type:       "float",
			ValueFloat: ConvertToFloat64(node["Buckets"]),
			Kind:       Quantity,
		})
	}

	if node["Buckets"+" Originally"] != nil {
		props = append(props, Property{
			ID:         "buckets_originally",
			Name:       "Buckets" + " Originally",
			Type:       "float",
			ValueFloat: ConvertToFloat64(node["Buckets Originally"]),
			Kind:       Quantity,
		})
	}

	return props
}

func getSortProperties(node Node, props []Property) []Property {
	if node[SORT_METHOD] != nil {
		props = append(props, Property{
			ID:          "sort_method",
			Name:        "Sort method",
			Type:        "string",
			ValueString: node[SORT_METHOD].(string),
		})
	}

	if node[SORT_SPACE_TYPE] != nil {
		property := Property{
			ID:          "sort_space_type",
			Name:        SORT_SPACE_TYPE,
			Type:        "string",
			ValueString: node[SORT_SPACE_TYPE].(string),
		}
		props = append(props, property)
	}

	if node[SORT_SPACE_USED] != nil {
		property := Property{
			ID:         "sort_space_used",
			Name:       SORT_SPACE_USED,
			Type:       "float",
			ValueFloat: ConvertToFloat64(node[SORT_SPACE_USED]),
			Kind:       DiskSize,
		}
		props = append(props, property)
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
