package images

import (
	"reflect"
	"testing"
	"gopkg.in/h2non/bimg.v1"
)

func TestParseSize(t *testing.T) {
	type args struct {
		width     string
		height    string
		keepRatio bool
		oldSize   bimg.ImageSize
	}
	tests := []struct {
		name    string
		args    args
		want    bimg.ImageSize
		wantErr bool
	}{
		{
			name: "-/-",
			args: args{
				keepRatio: true,
				oldSize: bimg.ImageSize{
				},
			},
			want:    bimg.ImageSize{},
			wantErr: true,
		},
		{
			name: "px/px",
			args: args{
				width:     "500px",
				height:    "500px",
				keepRatio: true,
				oldSize: bimg.ImageSize{
					Width:  3064,
					Height: 2128,
				},
			},
			want:    bimg.ImageSize{Width: 500, Height: 347},
			wantErr: false,
		},
		{
			name: "px/-",
			args: args{
				width: "500px",
				oldSize: bimg.ImageSize{
					Width:  3064,
					Height: 2128,
				},
			},
			want:    bimg.ImageSize{Width: 500, Height: 347},
			wantErr: false,
		},
		{
			name: "-/px",
			args: args{
				height: "500px",
				oldSize: bimg.ImageSize{
					Width:  3064,
					Height: 2128,
				},
			},
			want:    bimg.ImageSize{Width: 720, Height: 500},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSize(tt.args.width, tt.args.height, tt.args.keepRatio, tt.args.oldSize)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSize() error = %#v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseSize() = %#v, want %v", got, tt.want)
			}
		})
	}
}
