package score

import "sort"

type Mdtpvariable struct {
	Name        string
	DataHistory [][]int
	Score       int
	// [i][j] i番目のバイトがでj回目の試行で更新されたかされなかったのかの配列
}

// Scoreの降順でソートするためのカスタム関数
func sortByScoreDesc(mdtps []Mdtpvariable) {
	sort.Slice(mdtps, func(i, j int) bool {
		return mdtps[i].Score > mdtps[j].Score // Scoreの降順
	})
}

type Packet struct {
	DataHistory [][]int
	// [i][j] i番目のバイトがでj回目の試行で更新されたかされなかったのかの配列
}

func calcSimultaneousScore(datahistory [][]int) int {
	// historyの長さが固定の[][k]intであることが必須
	score := 0
	cycle := len(datahistory[0])
	num_of_byte := len(datahistory)
	for j := 0; j < cycle; j += 1 {
		simultaneous_packet_num := 0
		for i := 0; i < num_of_byte; i += 1 {
			if datahistory[i][j] == 1 {
				simultaneous_packet_num++
			}
		}
		score += simultaneous_packet_num * simultaneous_packet_num
	}
	return score
}

func calcFrequencyScore(datahistory [][]int) int {
	score := 0
	for _, bytehistory := range datahistory {
		for _, updated := range bytehistory {
			if updated == 1 {
				score += 1
			}
		}
	}
	return score
}

func MakeInitialDataOrder(mdtpvs []Mdtpvariable) []Mdtpvariable {
	for idx, mdtpv := range mdtpvs {
		mdtpvs[idx].Score = CalcMdtpvariableScore(mdtpv)
	}

	sortByScoreDesc(mdtpvs)
	return mdtpvs
}

func CalcMdtpvariableScore(mdtpv Mdtpvariable) int {
	score := 0
	// calc update frequency
	score += calcFrequencyScore(mdtpv.DataHistory)
	// calc update simultaneous score
	score += calcSimultaneousScore(mdtpv.DataHistory)

	return score
}

func CalcPacketUpdateTiming(byteHistories [][]int) []int {
	packetHistory := make([]int, len(byteHistories[1]))
	for _, byteHistory := range byteHistories {
		for idx, past := range byteHistory {
			if past == 1 {
				packetHistory[idx] = 1
			}
		}
	}
	return packetHistory
}

func MakePackets(variables []Mdtpvariable) []Packet {
	const packetsize = 7
	packets := []Packet{}
	data_line := [][]int{}
	for _, variable := range variables {
		data_line = append(data_line, variable.DataHistory...)
	}

	for i := 0; i < len(data_line); i += packetsize {
		end := i + packetsize
		if end > len(data_line) {
			end = len(data_line)
		}
		packets = append(packets, Packet{DataHistory: data_line[i:end]})
	}

	return packets
}

func CalcSimultaneousPacketUpdateScore(packets []Packet) int {
	// 各タイミングでの更新パケット数の平方和をスコアとする
	score := 0
	// i番目のパケットがでj回目の試行で更新されたかされなかったのかの配列
	packetUpdateTimings := [][]int{}
	for _, packet := range packets {
		packetUpdateTimings = append(packetUpdateTimings, CalcPacketUpdateTiming(packet.DataHistory))
	}
	cycle := len(packetUpdateTimings[0])
	num_of_packet := len(packetUpdateTimings)
	for j := 0; j < cycle; j += 1 {
		simultaneous_packet_num := 0
		for i := 0; i < num_of_packet; i += 1 {
			if packetUpdateTimings[i][j] == 1 {
				simultaneous_packet_num++
			}
		}
		score += simultaneous_packet_num * simultaneous_packet_num
	}

	return score
}
