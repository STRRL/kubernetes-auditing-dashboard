package gql

import (
	pointer "k8s.io/utils/pointer"
	"testing"
)

func Test_paginationToSQL(t *testing.T) {
	type args struct {
		page     *int
		pageSize *int
	}
	tests := []struct {
		name       string
		args       args
		wantOffset int
		wantLimit  int
	}{
		{
			name: "both are nil",
			args: args{
				page:     nil,
				pageSize: nil,
			},
			wantOffset: 0,
			wantLimit:  10,
		},
		{
			name: "page is nil",
			args: args{
				page:     nil,
				pageSize: pointer.Int(10),
			},
			wantOffset: 0,
			wantLimit:  10,
		},
		{
			name: "pageSize is nil",
			args: args{
				page:     pointer.Int(1),
				pageSize: nil,
			},
			wantOffset: 10,
			wantLimit:  10,
		},

		{
			name: "both are not nil",
			args: args{
				page:     pointer.Int(1),
				pageSize: pointer.Int(10),
			},
			wantOffset: 10,
			wantLimit:  10,
		},
		{
			name: "first page",
			args: args{
				page:     pointer.Int(0),
				pageSize: pointer.Int(10),
			},
			wantOffset: 0,
			wantLimit:  10,
		},
		{
			name: "page 3 pageSize 20",
			args: args{
				page:     pointer.Int(3),
				pageSize: pointer.Int(20),
			},
			wantOffset: 60,
			wantLimit:  20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOffset, gotLimit := paginationToSQL(tt.args.page, tt.args.pageSize)
			if gotOffset != tt.wantOffset {
				t.Errorf("paginationToSQL() gotOffset = %v, want %v", gotOffset, tt.wantOffset)
			}
			if gotLimit != tt.wantLimit {
				t.Errorf("paginationToSQL() gotLimit = %v, want %v", gotLimit, tt.wantLimit)
			}
		})
	}
}
