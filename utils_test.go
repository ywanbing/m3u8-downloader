package m3u8

import (
	"net/url"
	"reflect"
	"testing"
)

func Test_absolutist(t *testing.T) {
	type args struct {
		rawurl string
		u      *url.URL
	}
	tests := []struct {
		name    string
		args    args
		wantUri *url.URL
		wantErr bool
	}{
		{
			name: "absolute",
			args: args{
				rawurl: "//www.baidu.com/123.txt",
				u:      &url.URL{Scheme: "https", Host: "www.baidu.com"},
			},
			wantUri: &url.URL{Scheme: "https", Host: "www.baidu.com", Path: "123.txt"},
			wantErr: false,
		},
		{
			name: "absolute2",
			args: args{
				rawurl: "123.txt",
				u:      &url.URL{Scheme: "https", Host: "www.baidu.com"},
			},
			wantUri: &url.URL{Scheme: "https", Host: "www.baidu.com", Path: "123.txt"},
			wantErr: false,
		},
		{
			name: "absolute3",
			args: args{
				rawurl: "/123.txt",
				u:      &url.URL{Scheme: "https", Host: "www.baidu.com"},
			},
			wantUri: &url.URL{Scheme: "https", Host: "www.baidu.com", Path: "123.txt"},
			wantErr: false,
		},
		{
			name: "absolute4",
			args: args{
				rawurl: "https://google.com/123.txt",
				u:      &url.URL{Scheme: "https", Host: "www.baidu.com"},
			},
			wantUri: &url.URL{Scheme: "https", Host: "google.com", Path: "123.txt"},
			wantErr: false,
		},
		{
			name: "absolute5",
			args: args{
				rawurl: "temp/123.txt",
				u:      &url.URL{Scheme: "https", Host: "www.baidu.com", Path: "temp"},
			},
			wantUri: &url.URL{Scheme: "https", Host: "www.baidu.com", Path: "temp/123.txt"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUri, err := absolutist(tt.args.rawurl, tt.args.u)
			if (err != nil) != tt.wantErr {
				t.Errorf("absolutist() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotUri.String(), tt.wantUri.String()) {
				t.Errorf("absolutist() gotUri = %v, want %v", gotUri, tt.wantUri)
			}
		})
	}
}
