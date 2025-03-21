// kunique 提供一个基于雪花算法的唯一ID生成器
package kunique

import (
	"sync"
	"time"
)

const (
	epoch          = int64(1700465775306)              // 设置起始时间(时间戳/毫秒)
	timestampBits  = uint(41)                          // 时间戳占用位数
	nodeIDBits     = uint(6)                           // 机器id所占位数
	sequenceBits   = uint(16)                          // 序列所占的位数
	timestampMax   = int64(-1 ^ (-1 << timestampBits)) // 时间戳最大值
	nodeIDMax      = int64(-1 ^ (-1 << nodeIDBits))    // 支持的最大机器id数量
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits))  // 支持的最大序列id数量
	nodeIDShift    = sequenceBits                      // 机器id左移位数
	timestampShift = sequenceBits + nodeIDBits         // 时间戳左移位数
)

var defaultUniqueNode *UniqueNode
var defaultUniqueNodeOnce sync.Once

type UniqueNode struct {
	mu        sync.Mutex
	nodeID    int64 // 机器ID
	sequence  int64 // 序列号
	timestamp int64 // 时间戳 ，毫秒
}

// NewUniqueNode 创建一个新的唯一ID生成节点
//
// 参数说明:
//   - nodeID: 节点ID，范围必须在0到1023之间
//
// 返回值说明:
//   - *UniqueNode: 返回初始化后的唯一ID生成节点
//
// 注意事项:
//   - 如果nodeID超出范围，会触发panic
//   - 每个节点ID对应一个唯一的生成器实例
//   - 建议在系统启动时初始化并保持单例
//
// 示例:
//
//	node := NewUniqueNode(1) // 创建节点ID为1的生成器
func NewUniqueNode(nodeID int64) *UniqueNode {
	if nodeID < 0 || nodeID > nodeIDMax {
		panic("nodeID must be between 0 and 1023")
	}
	return &UniqueNode{nodeID: nodeID}
}

// GenerateUniqueID 生成一个全局唯一的ID
//
// 参数说明:
//   - nodeID: 可选参数，节点ID，范围必须在0到1023之间，默认为1
//
// 返回值说明:
//   - int64: 返回生成的64位唯一ID
//
// 注意事项:
//   - 该函数只会在第一次使用时进行初始化
//   - 该函数是线程安全的，使用sync.Once保证单例初始化
//   - 如果未提供nodeID参数，默认使用节点ID为1
//   - 如果nodeID超出范围，会触发panic
//   - ID结构: 41位时间戳 | 6位节点ID | 16位序列号
//
// 示例:
//
//	id := GenerateUniqueID() // 使用默认节点ID 1 生成唯一ID
func GenerateUniqueID(nodeID ...int64) int64 {
	nodeIDFlag := int64(1)
	if len(nodeID) > 0 {
		nodeIDFlag = nodeID[0]
	}
	defaultUniqueNodeOnce.Do(func() {
		defaultUniqueNode = NewUniqueNode(nodeIDFlag)
	})
	return defaultUniqueNode.Generate()
}

// Generate 生成一个全局唯一的ID
//
// 参数说明:
//   - 无
//
// 返回值说明:
//   - int64: 返回生成的64位唯一ID
//
// 注意事项:
//   - 该方法是线程安全的，使用互斥锁保证并发安全
//   - 如果时间戳超出最大值(41位)，会返回0
//   - 同一毫秒内生成的ID会递增序列号
//   - 当序列号超出最大值(16位)时，会等待到下一毫秒再生成
//   - ID结构: 41位时间戳 | 6位节点ID | 16位序列号
//
// 示例:
//
//	node := NewUniqueNode(1)
//	id := node.Generate() // 生成唯一ID
func (s *UniqueNode) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UnixMilli() // 转毫秒
	if s.timestamp == now {
		// 当同一时间戳（精度：毫秒）下多次生成id会增加序列号
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 如果当前序列超出12bit长度，则需要等待下一毫秒
			// 下一毫秒将使用sequence:0
			for now <= s.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		// 不同时间戳（精度：毫秒）下直接使用序列号：0
		s.sequence = 0
	}
	t := now - epoch
	if t > timestampMax {
		return 0
	}
	s.timestamp = now
	r := t<<timestampShift | (s.nodeID << nodeIDShift) | (s.sequence)

	return r
}
