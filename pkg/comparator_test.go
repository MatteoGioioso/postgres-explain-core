package pkg

import "testing"

func TestComparator_computeJaccardSimilarityIndex(t *testing.T) {
	type fields struct {
		planPrev      Explained
		planOptimized Explained
	}
	type args struct {
		optimized PlanRow
		prev      PlanRow
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "test",
			fields: fields{
				planPrev:      Explained{},
				planOptimized: Explained{},
			},
			args: args{
				optimized: PlanRow{
					Operation: BITMAP_INDEX_SCAN,
					Scopes: NodeScopes{
						Table: "mytable",
					},
					Level: 2,
				},
				prev: PlanRow{
					Operation: "Parallel " + BITMAP_INDEX_SCAN,
					Scopes: NodeScopes{
						Table: "mytable",
					},
					Level: 2,
				},
			},
			want: 0.8,
		},
		{
			name: "test2",
			fields: fields{
				planPrev:      Explained{},
				planOptimized: Explained{},
			},
			args: args{
				optimized: PlanRow{
					Operation: BITMAP_INDEX_SCAN,
					Scopes: NodeScopes{
						Table: "mytable",
					},
					Level: 2,
				},
				prev: PlanRow{
					Operation: NESTED_LOOP_JOIN,
					Scopes: NodeScopes{
						Table: "mytable",
					},
					Level: 2,
				},
			},
			want: 0.5,
		},
		{
			name: "test2",
			fields: fields{
				planPrev:      Explained{},
				planOptimized: Explained{},
			},
			args: args{
				optimized: PlanRow{
					Operation: NESTED_LOOP_JOIN,
					Scopes: NodeScopes{
						Table: "mytable",
					},
					Level: 2,
				},
				prev: PlanRow{
					Operation: NESTED_LOOP_JOIN,
					Scopes: NodeScopes{
						Table: "sometable",
					},
					Level: 2,
				},
			},
			want: 0.2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := JaccardSimilarity(tt.args.optimized, tt.args.prev); got <= tt.want {
				t.Errorf("computeJaccardSimilarityIndex() = %v <= %v", got, tt.want)
			}
		})
	}
}
