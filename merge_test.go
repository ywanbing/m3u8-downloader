package m3u8

import (
	"testing"
)

func TestMergeTsFileListToSingleMp4(t *testing.T) {
	type args struct {
		req *MergeTsFileReq
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				req: &MergeTsFileReq{
					Dir: "./temp/动画",
					TsFileList: []string{
						"0.ts",
						"1.ts",
						"2.ts",
						"3.ts",
						"4.ts",
						"5.ts",
						"6.ts",
						"7.ts",
						"8.ts",
						"9.ts",
						"10.ts",
						"11.ts",
						"12.ts",
						"13.ts",
						"14.ts",
						"15.ts",
						"16.ts",
						"17.ts",
						"18.ts",
						"19.ts",
						"20.ts",
						"21.ts",
						"22.ts",
						"23.ts",
						"24.ts",
						"25.ts",
						"26.ts",
						"27.ts",
						"28.ts",
						"29.ts",
						"30.ts",
						"31.ts",
						"32.ts",
						"33.ts",
						"34.ts",
						"35.ts",
						"36.ts",
						"37.ts",
						"38.ts",
						"39.ts",
						"40.ts",
						"41.ts",
						"42.ts",
						"43.ts",
						"44.ts",
						"45.ts",
						"46.ts",
						"47.ts",
						"48.ts",
						"49.ts",
						"50.ts",
						"51.ts",
						"52.ts",
						"53.ts",
						"54.ts",
						"55.ts",
						"56.ts",
						"57.ts",
						"58.ts",
						"59.ts",
						"60.ts",
						"61.ts",
						"62.ts",
						"63.ts",
						"64.ts",
						"65.ts",
						"66.ts",
						"67.ts",
						"68.ts",
						"69.ts",
					},
					OutputMp4: "动画.mp4",
					AfterDel:  false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MergeTsFileListToSingleMp4(tt.args.req); (err != nil) != tt.wantErr {
				t.Errorf("MergeTsFileListToSingleMp4() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
