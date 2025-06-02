package ec2

import (
	"testing"
)

func TestSDConfig_GetLabels(t *testing.T) {
	tests := []struct {
		name    string
		cfg     SDConfig
		wantErr bool
	}{
		{
			name: "ValidProfile",
			cfg: SDConfig{
				Region:  "eu-west-1",
				Profile: "account-one",
			},
			wantErr: false,
		},
		{
			name: "ProfileWithAccessKey",
			cfg: SDConfig{
				Region:    "eu-west-1",
				Profile:   "account-one",
				AccessKey: "key",
			},
			wantErr: true,
		},
		{
			name: "MissingRegion",
			cfg: SDConfig{
				Profile: "account-one",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := tt.cfg.GetLabels("")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLabels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
