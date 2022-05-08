package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDistinctIntList(t *testing.T) {
	type args struct {
		listList [][]int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			args: struct{ listList [][]int }{listList: [][]int{
				{
					1, 2,
				},
				{
					2, 3,
				},
			}},
			want: []int{
				1, 2, 3,
			},
		},
		{
			args: struct{ listList [][]int }{listList: [][]int{
				{},
				{
					1, 2, 3,
				},
			}},
			want: []int{1, 2, 3},
		},
		{
			args: struct{ listList [][]int }{listList: [][]int{
				{},
				{},
			}},
			want: []int{},
		},
		{
			args: struct{ listList [][]int }{listList: [][]int{
				{
					1,
				},
				{2},
			}},
			want: []int{
				1, 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, DistinctIntList(tt.args.listList...))
		})
	}
}
