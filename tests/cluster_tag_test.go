package tests

import (
	"context"
	"testing"

	"github.com/galaxy-future/BridgX/internal/service"
)

func TestGetClustersByTags(t *testing.T) {
	type args struct {
		ctx      context.Context
		tags     map[string]string
		pageSize int
		pageNum  int
	}
	tests := []struct {
		name    string
		args    args
		want1   int64
		wantErr bool
	}{
		{
			"",
			args{ctx: context.Background(),
				tags:     map[string]string{"k1": "v1", "k3": "", "k2": "v2"},
				pageSize: 10, pageNum: 1},
			1, false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := service.GetClustersByTags(tt.args.ctx, tt.args.tags, tt.args.pageSize, tt.args.pageNum)
			t.Logf("got:%v\n", got)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetClustersByTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got1 != tt.want1 {
				t.Errorf("GetClustersByTags() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
