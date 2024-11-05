package homeassistant

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"
	"testing"
	"time"
)

func Test485(t *testing.T) {
	var a, b uint8
	a = 200
	b = 57
	ss := a + b
	fmt.Println(ss)

	// 定义TCP服务器的地址和端口
	serverAddress := "192.168.123.129:8899"

	// 尝试连接到TCP服务器
	conn, err := net.Dial("tcp", serverAddress)
	if err != nil {
		log.Fatalf("无法连接到服务器: %instruction", err)
	}
	defer conn.Close()

	fmt.Println("成功连接到服务器")

	// 创建一个缓冲读取器
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriterSize(conn, 8)
	go writerLoop(writer)
	// 定义固定包长度
	const packetLength = 8
	buffer := make([]byte, packetLength)
	// 持续监听服务器的输出
	for {
		// 从连接中读取一行数据
		n, err := reader.Read(buffer)
		if err != nil {
			log.Printf("读取数据时出错: %instruction", err)
			break
		}
		if n == packetLength {
			fmt.Printf("收到数据包: %x\n", buffer)
		} else {
			fmt.Printf("收到不完整的数据包: %x\n", buffer[:n])
		}
	}
}

func writerLoop(writer *bufio.Writer) {
	time.Sleep(2 * time.Second)
	bs := []byte{0x02, 0x06, 0x10, 0x08, 0x01, 0x00, 0x00, 0x00}
	bs[5] = 0x01 | 0x01<<1 | 0x01<<3
	bs[5] = 0x01

	// bs = []byte{0x02, 0x06, 0x10, 0x03, 0x00, 0x1F, 0x00, 0x00}
	crc := CRC16(bs[:6])
	bs[6] = byte(crc & 0xFF)
	bs[7] = byte(crc >> 8)
	fmt.Printf("%04X\n", bs)
	writeBytesData(writer, bs)

	// hex字符串转bytes, eg: 02 06 10 08 01 01 CC AB
	// var hexStr string

	// hexStr = "02 06 10 08 01 01 CC AB"
	// hexStr = "02 06 10 08 00 00 0C FB"
	//writeHexData(writer, hexStr)
	//time.Sleep(1 * time.Second)

	//bs := []byte{0x02, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	// bs = []byte{0x02, 0x06, 0x00, 0x01, 0x80, 0x1E, 0x00, 0x9b}
	//for i := 0; i < 256; i++ {
	//	// bs[2] = byte(i)
	//	bs[0] = byte(i)
	//	bs = newP(bs)
	//	fmt.Printf("%02X\n", bs)
	//	writeBytesData(writer, bs)
	//	time.Sleep(500 * time.Millisecond)
	//}
	// writeBytesData(writer, newP(bs))

	//for i := 0; i < 256; i++ {
	//	bs := []byte{0x3C, 0x03, 0x00, 0x33, 0x00, 0x01, 0x00, 0x00}
	//	bs = []byte{0x3C, 0x03, 0x00, 0x29, 0x00, 0x10, 0x91, 0x23}
	//	bs[0] = byte(i)
	//	crc := CRC16(bs[:6])
	//	bs[6] = byte(crc & 0xFF)
	//	bs[7] = byte(crc >> 8)
	//	fmt.Printf("%04X\n", bs)
	//	writeBytesData(writer, bs)
	//	time.Sleep(500 * time.Millisecond)
	//}
}

func CRC16(data []byte) uint16 {
	var crc uint16 = 0xFFFF
	for _, b := range data {
		crc ^= uint16(b)
		for i := 0; i < 8; i++ {
			if crc&0x0001 != 0 {
				crc = (crc >> 1) ^ 0xA001
			} else {
				crc >>= 1
			}
		}
	}
	return crc
}

func newP(bs []byte) []byte {
	bs[7] = bs[0]
	for i := 1; i < 7; i++ {
		bs[7] ^= bs[i]
	}
	return bs
}

func writeHexData(writer *bufio.Writer, hexStr string) {
	bytes, err := hex.DecodeString(strings.Replace(hexStr, " ", "", -1))
	if err != nil {
		log.Fatalf("无法解码十六进制字符串: %instruction", err)
	}
	_, err = writer.Write(bytes)
	if err != nil {
		log.Fatalf("无法写入数据: %instruction", err)
	}
	_ = writer.Flush()
}

func writeBytesData(writer *bufio.Writer, bytes []byte) {
	_, err := writer.Write(bytes)
	if err != nil {
		log.Fatalf("无法写入数据: %instruction", err)
	}
	_ = writer.Flush()
}
