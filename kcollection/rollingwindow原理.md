### 滑动窗口源码执行过程详解（附关键源码）

---

#### **示例场景设定**
- **窗口参数**：
  ```go
  size=3       // 窗口大小
  interval=1s  // 桶时间间隔
  ignoreCurrent=true
  ```
- **操作序列**：
  1. `t=0s`：初始化窗口
  2. `t=0.5s`：Add(100)
  3. `t=1.5s`：Add(200)
  4. `t=2.6s`：Add(300)
  5. `t=3.0s`：Reduce()

---

#### **1. 初始化窗口（t=0s）**
```go
// 源码片段：NewRollingWindow
w := &RollingWindow[T, B]{
    size:     size,
    win:      newWindow(newBucket, size),
    interval: interval,
    lastTime: ktime.Now(), // 初始时间为0s
}

// newWindow 初始化桶数组
func newWindow(newBucket func() B, size int) *window[T, B] {
    buckets := make([]B, size)
    for i := 0; i < size; i++ {
        buckets[i] = newBucket() // 创建3个空桶
    }
    return &window[T, B]{buckets: buckets, size: size}
}
```
- **内存状态**：
  ```
  buckets: [B1(0), B2(0), B3(0)]
  offset: 0 → 指向B1
  lastTime: 0s
  ```

---

#### **2. Add(100) at t=0.5s**
```go
// 源码：Add方法
func (rw *RollingWindow[T, B]) Add(v T) {
    rw.lock.Lock()
    defer rw.lock.Unlock()
    rw.updateOffset() // 更新offset
    rw.win.add(rw.offset, v) // 添加值到当前桶
}

// updateOffset核心逻辑
func (rw *RollingWindow[T, B]) updateOffset() {
    span := rw.span() // 计算时间跨度
    if span <= 0 {
        return
    }
    // 重置过期桶（此处span=0，不执行）
    rw.offset = (offset + span) % rw.size
    rw.lastTime = ... // 保持0.5s
}
```
- **关键执行**：
  - `span = (0.5s - 0s)/1s = 0` → 不更新offset
  - 直接调用 `win.add(0, 100)`
- **状态变化**：
  ```
  buckets: [B1(100), B2(0), B3(0)]
  offset: 0 → 仍指向B1
  lastTime: 0.5s
  ```

---

#### **3. Add(200) at t=1.5s**
```go
// updateOffset核心逻辑
func (rw *RollingWindow[T, B]) updateOffset() {
    span := rw.span() // (1.5s-0.5s)/1s=1
    for i := 0; i < span; i++ {
        // 重置offset+1=1的桶（B2）
        rw.win.resetBucket((0 + i + 1) % 3)
    }
    rw.offset = (0 + 1) % 3 = 1 // 指向B2
    rw.lastTime = 1.5s - 0.5s = 1.0s
}

// resetBucket实现
func (w *window[T, B]) resetBucket(offset int) {
    w.buckets[offset%w.size].Reset() // B2被重置
}
```
- **关键执行**：
  - 重置B2（`Sum=0, Count=0`）
  - 添加200到B2
- **状态变化**：
  ```
  buckets: [B1(100), B2(200), B3(0)]
  offset: 1 → 指向B2
  lastTime: 1.0s
  ```

---

#### **4. Add(300) at t=2.6s**
```go
// span计算
func (rw *RollingWindow[T, B]) span() int {
    offset := int(ktime.Since(rw.lastTime)/rw.interval) 
    // 2.6s-1.0s=1.6s → 1
    return 1
}

// updateOffset重置B3
rw.win.resetBucket((1 + 0 + 1) % 3) // 重置B3
rw.offset = (1 + 1) % 3 = 2 // 指向B3
```
- **状态变化**：
  ```
  buckets: [B1(100), B2(200), B3(300)]
  offset: 2 → 指向B3
  lastTime: 2.0s
  ```

---

#### **5. Reduce() at t=3.0s**
```go
// Reduce核心逻辑
func (rw *RollingWindow[T, B]) Reduce(fn func(b B)) {
    span := rw.span() // (3.0s-2.0s)/1s=1
    diff := 3 - 1 = 2 // 需要遍历2个桶
    offset := (2 + 1 + 1) % 3 = 0 // 从B1开始
    rw.win.reduce(0, 2, fn) // 遍历B1和B2
}

// reduce实现
func (w *window[T, B]) reduce(start, count int, fn func(b B)) {
    for i := 0; i < count; i++ {
        fn(w.buckets[(start+i)%w.size]) // 调用回调函数
    }
}
```
- **遍历结果**：
  - B1: Sum=0（被后续操作重置）
  - B2: Sum=200

---

#### **核心源码图解**
```go
// 滑动窗口核心结构
type RollingWindow[T Number, B BucketInterface[T]] struct {
    lock     sync.RWMutex       // 并发锁
    size     int                // 窗口大小
    win      *window[T, B]      // 环形缓冲区
    interval time.Duration      // 桶时间间隔
    offset   int                // 当前桶指针
    lastTime time.Duration      // 上次更新时间
}

// 时间窗口滑动逻辑
func (rw *RollingWindow[T, B]) updateOffset() {
    span := rw.span()
    for i := 0; i < span; i++ {
        rw.win.resetBucket((offset+i+1) % size) // 重置过期桶
    }
    rw.offset = (offset + span) % size          // 移动指针
}
```

---

#### **设计要点总结**
1. **环形缓冲区**：
   - 使用 `offset` 指针循环写入
   - 自动重置过期桶，避免内存分配
   ```go
   func (w *window[T, B]) resetBucket(offset int) {
       w.buckets[offset%w.size].Reset()
   }
   ```

2. **精确时间计算**：
   - 通过 `lastTime` 和模运算保证时间对齐
   ```go
   rw.lastTime = n - (n-rw.lastTime)%rw.interval
   ```

3. **并发安全**：
   - 使用 `sync.RWMutex` 保护关键操作
   ```go
   func (rw *RollingWindow[T, B]) Add(v T) {
       rw.lock.Lock()
       defer rw.lock.Unlock()
       // ...
   }
   ```

---

通过这种设计，滑动窗口能够在 `O(1)` 时间复杂度内完成数据添加和统计，适用于高频数据采集场景（如API调用次数统计）。时间窗口的滑动通过模运算高效实现，内存占用固定为 `size` 个桶，兼具性能和资源效率。