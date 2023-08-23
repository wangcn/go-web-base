package util

import (
	"bytes"
	"encoding/binary"
	"sync"
	"time"

	"mybase/pkg/xserver"
)

// A Node struct holds the basic information needed for a snowflake generator
// node
type Node struct {
	mu   sync.Mutex
	time int64
	node uint32
	step uint32
}

// An ID is a custom type used for a snowflake ID.  This is used so we can
// attach methods onto the ID.
type ID int64

// NewNode returns a new snowflake node that can be used to generate snowflake
// IDs
func NewNode() *Node {

	return &Node{
		time: 0,
		node: xserver.InetAtoN(xserver.GetInternalIp()),
		step: 0,
	}
}

// Generate creates and returns a unique snowflake ID
func (n *Node) Generate() []byte {
	n.mu.Lock()
	defer n.mu.Unlock()

	now := time.Now().UnixNano() / 1000000

	diff := n.time - now
	// 因ntp导致时间回退时，需要等待
	if diff > 0 {
		time.Sleep(time.Duration(diff) * time.Millisecond)
	}
	if n.time == now {
		n.step++

		if n.step == 0 {
			for now <= n.time {
				time.Sleep(1 * time.Millisecond)
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		n.step = 0
	}

	n.time = now

	bytesBuffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(bytesBuffer, binary.BigEndian, now)
	_ = binary.Write(bytesBuffer, binary.BigEndian, n.node)
	_ = binary.Write(bytesBuffer, binary.BigEndian, n.step)

	return bytesBuffer.Bytes()
}
