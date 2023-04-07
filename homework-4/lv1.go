package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

// 定义 ICMP 报文
type ICMP struct {
	Type        uint8  // 类型
	Code        uint8  // 代码
	Checksum    uint16 // 校验和
	Identifier  uint16 // 标识符
	SequenceNum uint16 // 序列号
}

func CheckSum(data []byte) uint16 {
	var (
		sum    uint32
		length = len(data)
		index  int
	)
	// 首先将每两个字节拆分出来进行累加
	for length > 1 {
		sum += uint32(data[index])<<8 + uint32(data[index+1])
		index += 2
		length -= 2
	}

	// 如果数据长度是奇数个，需要额外处理最后一个字节
	if length > 0 {
		sum += uint32(data[index]) << 8
	}

	// 将累加得到的结果进行一次反码取反
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	return uint16(^sum)
}
func Ping(url string, size int, count int) {
	// 解析主机名为 IP 地址
	addr, err := net.ResolveIPAddr("ip", "cqupt.edu.cn")
	if err != nil {
		fmt.Println("解析主机名失败:", err)
		return
	}

	// 创建 ICMP 连接
	conn, err := net.DialIP("ip4:icmp", nil, addr)
	if err != nil {
		fmt.Println("连接主机失败:", err)
		return
	}
	defer conn.Close()
	// 初始化报文
	icmp := ICMP{
		Type:        8,
		Code:        0,
		Checksum:    0,
		Identifier:  0,
		SequenceNum: 0,
	}
	// 创建一个bytes.Buffer对象，用于保存ICMP报文和数据
	var buffer bytes.Buffer
	// 原始报文
	initialBytes := make([]byte, 20000)
	binary.Write(&buffer, binary.BigEndian, icmp)
	// 截断CMP报文携带的数据部分
	binary.Write(&buffer, binary.BigEndian, initialBytes[0:size])
	b := buffer.Bytes()
	binary.BigEndian.PutUint16(b[2:], CheckSum(b))
	fmt.Printf("\n正在 Ping [%s] 具有 %d 字节的数据:\n", addr.String(), size)

	// 丢包处理
	// ping应答数据
	recvive := make([]byte, 1024)
	// 每个ping应答的延迟时间
	delayTimes := []float64{}
	// 统计丢包的次数，用于计算丢包率
	dropPack := 0
	maxTime := 3000.0
	minTime := 0.0
	avgTime := 0.0
	for i := count; i > 0; i-- {

		// 向目标地址发送二进制报文包,如果发送失败就dropPack++

		if _, err := conn.Write(buffer.Bytes()); err != nil {
			dropPack++
			time.Sleep(time.Second)
			continue
		}
		// 否则记录当前得时间
		t_start := time.Now()
		conn.SetReadDeadline((time.Now().Add(time.Second * 3)))
		len, err := conn.Read(recvive)
		// 减去IP数据报的头部长度20字节，ICMP报文的头部长度8字节
		len = len - 20 - 8
		// 返回失败丢包++
		if err != nil {
			dropPack++
			time.Sleep(time.Second)
			continue
		}
		time_end := time.Now()
		dur := float64(time_end.Sub(t_start).Nanoseconds()) / 1e6
		delayTimes = append(delayTimes, dur)
		if dur < maxTime {
			maxTime = dur
		}
		if dur > minTime {
			minTime = dur
		}
		fmt.Printf("来自 %s 的回复: 字节 = %d byte 时间 = %.3fms\n", addr.String(), len, dur)

	}
	fmt.Printf("%s 的 Ping 统计信息:", addr.String())
	fmt.Printf("    数据包: 已发送 = %d，已接收 = %d，丢失 = %d (%.2f%%丢失)\n", count, count-dropPack, dropPack, float64(dropPack)/float64(count*100))
	if len(delayTimes) == 0 {
		avgTime = 3000.0
	} else {
		sum := 0.0
		for _, n := range delayTimes {
			sum += n
		}
		avgTime = sum / float64(len(delayTimes))
	}
	fmt.Print("往返行程的估计时间(以毫秒为单位):\n")
	fmt.Printf("    最短 = %.2f ms 平均 = %.2fms 最长 = %.2fms\n", minTime, avgTime, maxTime)

}
func main() {
	// ping 重邮，32字节数据，四次
	Ping("cqupt.edu.cn", 32, 4)
}
