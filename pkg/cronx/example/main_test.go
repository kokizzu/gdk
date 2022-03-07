package main

import (
	"context"
	"testing"

	"github.com/rizalgowandy/gdk/pkg/cronx"
)

func Test_alwaysError_Run(t *testing.T) {
	type args struct {
		in context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Error",
			args:    args{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := alwaysError{}
			if err := a.Run(tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_everyJob_Run(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Success",
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := everyJob{}
			if err := ev.Run(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_subscription_Run(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				in0: cronx.SetJobMetadata(context.Background(), cronx.JobMetadata{
					EntryID:    1,
					Wave:       2,
					TotalWave:  3,
					IsLastWave: true,
				}),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			su := subscription{}
			if err := su.Run(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_sendEmail_Run(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Success",
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := sendEmail{}
			if err := e.Run(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_payBill_Run(t *testing.T) {
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Success",
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := payBill{}
			if err := p.Run(tt.args.in0); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterJobs(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cronx.Default()
			RegisterJobs()
		})
	}
}
