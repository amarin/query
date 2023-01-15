package query_test

import (
	"testing"

	"github.com/amarin/query"
)

func TestJoinCondition_Render(t *testing.T) {
	tests := []struct {
		name   string
		fields query.JoinCondition
		want   string
	}{
		{"join_table", query.JoinFields(query.Field("t1.f1"), query.Field("t2.f2")), "t1.f1=t2.f2"},
		{"join_as", query.JoinFields(query.Field("t1.f1 as s1"), query.Field("t2.f2 as s2")), "t1.f1=t2.f2"},
		{"join_as_left", query.JoinFields(query.Field("t1.f1 as s1"), query.Field("t2.f2")), "t1.f1=t2.f2"},
		{"join_as_right", query.JoinFields(query.Field("t1.f1"), query.Field("t2.f2 as s2")), "t1.f1=t2.f2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Render(); got != tt.want {
				t.Errorf("Render() = %v, want %v", got, tt.want)
			}
		})
	}
}
