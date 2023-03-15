package util

import (
	"errors"
	"net"
	"os"

	// "strconv"
	"sync/atomic"
	"time"
)

type FastID struct {
	//开始时间
	timeStart int64
	//时间位
	timeBits uint
	//序列位
	stepBits uint
	//机器位
	nodeBits uint
	timeMask int64
	stepMask int64
	//机器ID
	machineID     int64
	machineIDMask int64
	lastID        int64
}

func NewFastID(timeBits, nodeBits, stepBits uint, timeStart int64) *FastID {
	nodeId := getMachineID()
	return newFastID(timeBits, nodeBits, stepBits, nodeId, timeStart)
}

func NewFastIDWithNode(timeBits, nodeBits, stepBits uint, timeStart, nodeId int64) *FastID {
	return newFastID(timeBits, nodeBits, stepBits, nodeId, timeStart)
}

//45bit有1000多年，6bit共64个点节，12bit每毫秒4096个，每秒400万
func newFastID(timeBits, nodeBits, stepBits uint, nodeId, timeStart int64) *FastID {
	machineIDMask := ^(int64(-1) << nodeBits)
	return &FastID{
		timeStart:     timeStart,
		timeBits:      timeBits,
		stepBits:      stepBits,
		nodeBits:      nodeBits,
		timeMask:      ^(int64(-1) << timeBits),
		stepMask:      ^(int64(-1) << stepBits),
		machineIDMask: machineIDMask,
		machineID:     nodeId & machineIDMask,
		lastID:        0,
	}
}

//48bits 时间戳，可以表示大概8000多年
//45bits 大概就可以1000多年
func (c *FastID) getCurrentTimestamp() int64 {
	return (time.Now().UnixNano() - c.timeStart) >> 20 & c.timeMask
}

func (c *FastID) NextID() int64 {
	for {
		localLastID := atomic.LoadInt64(&c.lastID)
		seq := c.GetSequence(localLastID)
		lastIDTime := c.GetTime(localLastID)
		now := c.getCurrentTimestamp()
		if now > lastIDTime {
			seq = 0
		} else if seq >= c.stepMask {
			time.Sleep(time.Duration(0xFFFFF - (time.Now().UnixNano() & 0xFFFFF)))
			continue
		} else {
			seq++
		}

		newID := now<<(c.nodeBits+c.stepBits) + seq<<c.nodeBits + c.machineID
		if atomic.CompareAndSwapInt64(&c.lastID, localLastID, newID) {
			return newID
		}
		time.Sleep(time.Duration(20))
	}
}

func (c *FastID) GetSequence(id int64) int64 {
	return (id >> c.nodeBits) & c.stepMask
}

func (c *FastID) GetTime(id int64) int64 {
	return id >> (c.nodeBits + c.stepBits)
}

func getMachineID() int64 {
	if ip, err := getIP(); err == nil {
		return (int64(ip[2]) << 8) + int64(ip[3])
	}
	return 0
}

func getIP() (net.IP, error) {
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, addr := range addrs {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					ip := ipNet.IP.To4()

					if ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168 {
						return ip, nil
					}
				}
			}
		}
	}
	return nil, errors.New("Failed to get ip address")
}
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
