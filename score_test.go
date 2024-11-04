package score

import (
	"reflect"
	"testing"
)

func TestCalcPacketUpdateTiming(t *testing.T) {
	type args struct {
		byteHistories [][]int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "packet update timing",
			args: args{
				[][]int{
					{1, 0, 0, 0, 1, 0, 1, 0, 0},
					{1, 1, 0, 0, 1, 0, 1, 0, 1},
					{1, 1, 1, 0, 1, 1, 1, 0, 1},
				},
			},
			want: []int{1, 1, 1, 0, 1, 1, 1, 0, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcPacketUpdateTiming(tt.args.byteHistories); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalcPacketUpdateTiming() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakePackets(t *testing.T) {
	type args struct {
		variables []Mdtpvariable
	}
	tests := []struct {
		name string
		args args
		want []Packet
	}{
		{
			name: "make packets",
			args: args{
				variables: []Mdtpvariable{
					{DataHistory: [][]int{
						{1, 0, 0, 0, 1, 0, 1, 0, 0},
						{1, 1, 0, 0, 1, 0, 1, 0, 1},
						{1, 1, 1, 0, 1, 1, 1, 0, 1},
					},
					},
					{
						DataHistory: [][]int{
							{1, 0, 0, 0, 1, 0, 1, 0, 0},
							{1, 1, 0, 0, 1, 0, 1, 0, 1},
							{1, 1, 1, 0, 1, 1, 1, 0, 1},
						},
					},
					{
						DataHistory: [][]int{
							{1, 0, 0, 0, 1, 0, 1, 0, 0},
							{1, 1, 0, 0, 1, 0, 1, 0, 1},
							{1, 1, 1, 0, 1, 1, 1, 0, 1},
						},
					},
				},
			},
			want: []Packet{
				{DataHistory: [][]int{
					{1, 0, 0, 0, 1, 0, 1, 0, 0},
					{1, 1, 0, 0, 1, 0, 1, 0, 1},
					{1, 1, 1, 0, 1, 1, 1, 0, 1},
					{1, 0, 0, 0, 1, 0, 1, 0, 0},
					{1, 1, 0, 0, 1, 0, 1, 0, 1},
					{1, 1, 1, 0, 1, 1, 1, 0, 1},
					{1, 0, 0, 0, 1, 0, 1, 0, 0}},
				},
				{DataHistory: [][]int{
					{1, 1, 0, 0, 1, 0, 1, 0, 1},
					{1, 1, 1, 0, 1, 1, 1, 0, 1}},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakePackets(tt.args.variables); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakePackets() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalcSimultaneousPacketUpdateScore(t *testing.T) {
	type args struct {
		packets []Packet
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "packet score",
			args: args{
				packets: []Packet{
					{DataHistory: [][]int{
						{1, 0, 0, 0, 1, 0, 1, 0, 0},
						{1, 1, 0, 0, 1, 0, 1, 0, 1},
						{1, 1, 1, 0, 1, 1, 1, 0, 1},
						{1, 0, 0, 0, 1, 0, 1, 0, 0},
						{1, 1, 0, 0, 1, 0, 1, 0, 1},
						{1, 1, 1, 0, 1, 1, 1, 0, 1},
						{1, 0, 0, 0, 1, 0, 1, 0, 0}},
					},
					{DataHistory: [][]int{
						{1, 1, 0, 0, 1, 0, 1, 0, 1},
						{1, 1, 1, 0, 1, 1, 1, 0, 1}},
					},
				},
			},
			want: 2*2 + 2*2 + 2*2 + 0 + 2*2 + 2*2 + 2*2 + 0 + 2*2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcSimultaneousPacketUpdateScore(tt.args.packets); got != tt.want {
				t.Errorf("CalcSimultaneousPacketUpdateScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCalcMdtpvariableScore(t *testing.T) {
	type args struct {
		mdtpv Mdtpvariable
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "calc variable score",
			args: args{
				Mdtpvariable{
					Name: "test",
					DataHistory: [][]int{
						{1, 0, 0, 0, 1, 0, 1, 0, 0},
						{1, 1, 0, 0, 1, 0, 1, 0, 1},
						{1, 1, 1, 0, 1, 1, 1, 0, 1},
						{1, 0, 0, 0, 1, 0, 1, 0, 0},
						{1, 1, 0, 0, 1, 0, 1, 0, 1},
					},
				},
			},
			want: 23 + 5*5 + 3*3 + 1 + 5*5 + 1 + 5*5 + 3*3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcMdtpvariableScore(tt.args.mdtpv); got != tt.want {
				t.Errorf("CalcMdtpvariableScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortByScoreDesc(t *testing.T) {
	type args struct {
		mdtps []Mdtpvariable
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "sort",
			args: args{
				mdtps: []Mdtpvariable{
					{
						Score: 1,
					},
					{
						Score: 0,
					},
					{
						Score: 3,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortByScoreDesc(tt.args.mdtps)
		})
	}
}
