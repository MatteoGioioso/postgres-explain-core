package pkg

import "testing"

func TestPlanEnricher_AnalyzePlan(t *testing.T) {
	type fields struct {
		ctes            map[string]Node
		containsBuffers bool
	}
	type args struct {
		plan string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "test",
			fields: fields{
				ctes:            nil,
				containsBuffers: false,
			},
			args: args{
				plan: `[{"Plan":{"Actual Loops":1,"Actual Rows":531,"Actual Startup Time":26.505,"Actual Total Time":203.707,"Node Type":"Nested Loop","Plans":[{"Actual Loops":1,"Actual Rows":1200,"Actual Startup Time":1.718,"Actual Total Time":63.678,"Node Type":"Nested Loop","Plans":[{"Actual Loops":1,"Actual Rows":1200,"Actual Startup Time":1.707,"Actual Total Time":61.295,"Node Type":"Hash Join","Plans":[{"Actual Loops":1,"Actual Rows":1200,"Actual Startup Time":0.09,"Actual Total Time":58.75,"Node Type":"Seq Scan","Plan Rows":1564,"Rows Removed by Filter":825,"Total Cost":94.44,"Relation Name":"pg_class","Alias":"i_1","Startup Cost":0,"Plan Width":72,"Filter":"((relkind = ANY ('{i,I}'::\"char\"[])) AND (pg_get_indexdef(oid) !~* 'unique'::text))","Shared Hit Blocks":24707,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0},{"Actual Loops":1,"Actual Rows":1420,"Actual Startup Time":1.61,"Actual Total Time":1.614,"Node Type":"Hash","Plans":[{"Actual Loops":1,"Actual Rows":1420,"Actual Startup Time":0.687,"Actual Total Time":1.38,"Node Type":"Hash Join","Plans":[{"Actual Loops":1,"Actual Rows":1563,"Actual Startup Time":0.005,"Actual Total Time":0.186,"Node Type":"Seq Scan","Plan Rows":1701,"Total Cost":78.01,"Relation Name":"pg_index","Alias":"x_1","Startup Cost":0,"Plan Width":8,"Shared Hit Blocks":61,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0},{"Actual Loops":1,"Actual Rows":188,"Actual Startup Time":0.67,"Actual Total Time":0.671,"Node Type":"Hash","Plans":[{"Actual Loops":1,"Actual Rows":188,"Actual Startup Time":0.01,"Actual Total Time":0.635,"Node Type":"Seq Scan","Plan Rows":188,"Rows Removed by Filter":1837,"Total Cost":86.84,"Relation Name":"pg_class","Alias":"c_1","Startup Cost":0,"Plan Width":8,"Filter":"(relkind = ANY ('{r,m,p}'::\"char\"[]))","Shared Hit Blocks":59,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":188,"Total Cost":86.84,"Startup Cost":86.84,"Plan Width":8,"Buckets":1024,"Shared Hit Blocks":59,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":158,"Total Cost":171.68,"Startup Cost":89.19,"Plan Width":8,"Hash Cond":"(x_1.indrelid = c_1.oid)","Shared Hit Blocks":120,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":158,"Total Cost":171.68,"Startup Cost":171.68,"Plan Width":8,"Buckets":2048,"Shared Hit Blocks":120,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":122,"Total Cost":275.18,"Startup Cost":173.66,"Plan Width":68,"Hash Cond":"(i_1.oid = x_1.indexrelid)","Shared Hit Blocks":24827,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0},{"Actual Loops":1200,"Actual Rows":1,"Actual Startup Time":0.001,"Actual Total Time":0.001,"Index Name":"pg_namespace_oid_index","Node Type":"Index Scan","Plan Rows":1,"Total Cost":0.47,"Relation Name":"pg_namespace","Alias":"n_1","Startup Cost":0.28,"Plan Width":68,"Index Cond":"(oid = c_1.relnamespace)","Shared Hit Blocks":3600,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":122,"Total Cost":332.58,"Startup Cost":173.93,"Plan Width":128,"Shared Hit Blocks":28427,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0},{"Actual Loops":1200,"Actual Rows":567,"Actual Startup Time":0.009,"Actual Total Time":0.04,"Node Type":"Materialize","Plans":[{"Actual Loops":1,"Actual Rows":567,"Actual Startup Time":10.136,"Actual Total Time":12.579,"Node Type":"Nested Loop","Plans":[{"Actual Loops":1,"Actual Rows":1280,"Actual Startup Time":10.102,"Actual Total Time":11.206,"Node Type":"Hash Join","Plans":[{"Actual Loops":1,"Actual Rows":1280,"Actual Startup Time":4.975,"Actual Total Time":5.732,"Node Type":"Hash Join","Plans":[{"Actual Loops":1,"Actual Rows":2025,"Actual Startup Time":0.004,"Actual Total Time":0.209,"Node Type":"Seq Scan","Plan Rows":2025,"Total Cost":79.25,"Relation Name":"pg_class","Alias":"i","Startup Cost":0,"Plan Width":68,"Shared Hit Blocks":59,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0},{"Actual Loops":1,"Actual Rows":1280,"Actual Startup Time":4.963,"Actual Total Time":4.969,"Node Type":"Hash","Plans":[{"Actual Loops":1,"Actual Rows":1280,"Actual Startup Time":3.732,"Actual Total Time":4.694,"Node Type":"Hash Join","Plans":[{"Actual Loops":1,"Actual Rows":1539,"Actual Startup Time":0.528,"Actual Total Time":1.191,"Node Type":"Hash Join","Plans":[{"Actual Loops":1,"Actual Rows":1563,"Actual Startup Time":0.004,"Actual Total Time":0.163,"Node Type":"Seq Scan","Plan Rows":1701,"Total Cost":78.01,"Relation Name":"pg_index","Alias":"x","Startup Cost":0,"Plan Width":8,"Shared Hit Blocks":61,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0},{"Actual Loops":1,"Actual Rows":328,"Actual Startup Time":0.517,"Actual Total Time":0.518,"Node Type":"Hash","Plans":[{"Actual Loops":1,"Actual Rows":328,"Actual Startup Time":0.006,"Actual Total Time":0.46,"Node Type":"Seq Scan","Plan Rows":327,"Rows Removed by Filter":1697,"Total Cost":86.84,"Relation Name":"pg_class","Alias":"c","Startup Cost":0,"Plan Width":72,"Filter":"(relkind = ANY ('{r,t,m}'::\"char\"[]))","Shared Hit Blocks":59,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":327,"Total Cost":86.84,"Startup Cost":86.84,"Plan Width":72,"Buckets":1024,"Shared Hit Blocks":59,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":275,"Total Cost":173.42,"Startup Cost":90.93,"Plan Width":80,"Hash Cond":"(x.indrelid = c.oid)","Shared Hit Blocks":120,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0},{"Actual Loops":1,"Actual Rows":989,"Actual Startup Time":3.123,"Actual Total Time":3.124,"Node Type":"Hash","Plans":[{"Actual Loops":1,"Actual Rows":989,"Actual Startup Time":0.01,"Actual Total Time":2.931,"Node Type":"Seq Scan","Plan Rows":998,"Rows Removed by Filter":990,"Total Cost":54.69,"Relation Name":"pg_namespace","Alias":"n","Startup Cost":0,"Plan Width":68,"Filter":"((nspname <> ALL ('{pg_catalog,information_schema}'::name[])) AND (nspname !~ '^pg_toast'::text))","Shared Hit Blocks":25,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":998,"Total Cost":54.69,"Startup Cost":54.69,"Plan Width":68,"Buckets":1024,"Shared Hit Blocks":25,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":139,"Total Cost":241.3,"Startup Cost":158.09,"Plan Width":140,"Hash Cond":"(c.relnamespace = n.oid)","Shared Hit Blocks":145,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":139,"Total Cost":241.3,"Startup Cost":241.3,"Plan Width":140,"Buckets":2048,"Shared Hit Blocks":145,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":139,"Total Cost":331.27,"Startup Cost":243.04,"Plan Width":204,"Hash Cond":"(i.oid = x.indexrelid)","Shared Hit Blocks":204,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0},{"Actual Loops":1,"Actual Rows":119,"Actual Startup Time":5.115,"Actual Total Time":5.123,"Node Type":"Hash","Plans":[{"Actual Loops":1,"Actual Rows":119,"Actual Startup Time":5.026,"Actual Total Time":5.103,"Node Type":"Subquery Scan","Plans":[{"Actual Loops":1,"Actual Rows":119,"Actual Startup Time":5.025,"Actual Total Time":5.086,"Node Type":"HashAggregate","Plans":[{"Actual Loops":1,"Actual Rows":1288,"Actual Startup Time":3.924,"Actual Total Time":4.659,"Node Type":"Hash Join","Plans":[{"Actual Loops":1,"Actual Rows":1551,"Actual Startup Time":0.614,"Actual Total Time":1.119,"Node Type":"Hash Join","Plans":[{"Actual Loops":1,"Actual Rows":1563,"Actual Startup Time":0.013,"Actual Total Time":0.153,"Index Name":"pg_index_indrelid_index","Node Type":"Index Only Scan","Plan Rows":1701,"Total Cost":32.79,"Relation Name":"pg_index","Alias":"i_2","Startup Cost":0.28,"Plan Width":4,"Heap Fetches":"0","Shared Hit Blocks":7,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0},{"Actual Loops":1,"Actual Rows":328,"Actual Startup Time":0.594,"Actual Total Time":0.595,"Node Type":"Hash","Plans":[{"Actual Loops":1,"Actual Rows":328,"Actual Startup Time":0.01,"Actual Total Time":0.537,"Node Type":"Seq Scan","Plan Rows":327,"Rows Removed by Filter":1697,"Total Cost":86.84,"Relation Name":"pg_class","Alias":"c_2","Startup Cost":0,"Plan Width":72,"Filter":"(relkind = ANY ('{r,t,m}'::\"char\"[]))","Shared Hit Blocks":59,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":327,"Total Cost":86.84,"Startup Cost":86.84,"Plan Width":72,"Buckets":1024,"Shared Hit Blocks":59,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":327,"Total Cost":128.2,"Join Type":"Right","Startup Cost":91.21,"Plan Width":72,"Hash Cond":"(i_2.indrelid = c_2.oid)","Shared Hit Blocks":66,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0},{"Actual Loops":1,"Actual Rows":989,"Actual Startup Time":3.245,"Actual Total Time":3.246,"Node Type":"Hash","Plans":[{"Actual Loops":1,"Actual Rows":989,"Actual Startup Time":0.016,"Actual Total Time":3.067,"Node Type":"Seq Scan","Plan Rows":998,"Rows Removed by Filter":990,"Total Cost":54.69,"Relation Name":"pg_namespace","Alias":"n_2","Startup Cost":0,"Plan Width":68,"Filter":"((nspname <> ALL ('{pg_catalog,information_schema}'::name[])) AND (nspname !~ '^pg_toast'::text))","Shared Hit Blocks":25,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":998,"Total Cost":54.69,"Startup Cost":54.69,"Plan Width":68,"Buckets":1024,"Shared Hit Blocks":25,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":165,"Total Cost":196.22,"Startup Cost":158.37,"Plan Width":132,"Hash Cond":"(c_2.relnamespace = n_2.oid)","Shared Hit Blocks":91,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":165,"Total Cost":199.11,"Startup Cost":197.46,"Plan Width":292,"Group Key":"c_2.oid, n_2.nspname, c_2.relname","Batches":1,"Shared Hit Blocks":91,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":165,"Total Cost":200.76,"Startup Cost":197.46,"Plan Width":4,"Shared Hit Blocks":91,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":165,"Total Cost":200.76,"Startup Cost":200.76,"Plan Width":4,"Buckets":1024,"Shared Hit Blocks":91,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":11,"Total Cost":534.72,"Startup Cost":445.86,"Plan Width":208,"Hash Cond":"(c.oid = pg_stat_all_tables.relid)","Shared Hit Blocks":295,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0},{"Actual Loops":1280,"Actual Rows":0,"Actual Startup Time":0.001,"Actual Total Time":0.001,"Index Name":"pg_inherits_relid_seqno_index","Node Type":"Index Only Scan","Plan Rows":1,"Total Cost":0.3,"Relation Name":"pg_inherits","Startup Cost":0.28,"Plan Width":4,"Index Cond":"(inhrelid = x.indrelid)","Heap Fetches":"0","Shared Hit Blocks":2561,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":3,"Total Cost":538.11,"Startup Cost":446.13,"Plan Width":196,"Join Filter":"(c.oid = pg_inherits.inhrelid)","Shared Hit Blocks":2856,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":3,"Total Cost":538.13,"Startup Cost":446.13,"Plan Width":196,"Shared Hit Blocks":2856,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0}],"Plan Rows":1,"Rows Removed by Join Filter":679869,"Total Cost":877.12,"Startup Cost":620.06,"Plan Width":40,"Join Filter":"((i.relname = i_1.relname) AND (n.nspname = n_1.nspname))","Shared Hit Blocks":31283,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0},"Shared Hit Blocks":106,"Shared Read Blocks":0,"Shared Written Blocks":0,"Shared Dirtied Blocks":0,"Planning Time":2.27,"Execution Time":204.007}]`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ps := &PlanEnricher{
				ctes:            tt.fields.ctes,
				containsBuffers: tt.fields.containsBuffers,
			}
			node, err := GetRootNodeFromPlans(tt.args.plan)
			if err != nil {
				t.Fatal(err)
			}
			ps.AnalyzePlan(node)
		})
	}
}