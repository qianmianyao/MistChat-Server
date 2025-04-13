package encryption

import (
	"fmt"
	"testing"
)

func Test_generateUID(t *testing.T) {
	type args struct {
		username string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "valid username",
			args:    args{username: "testuser"},
			want:    "u_Na3_Z2bCxlY-jO7e",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateUID("r_")
			if (err != nil) != tt.wantErr {
				t.Errorf("generateUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.want {
				t.Errorf("generateUID() got = %v, want %v", got, tt.want)
			}
			fmt.Println(got)
		})
	}
}
