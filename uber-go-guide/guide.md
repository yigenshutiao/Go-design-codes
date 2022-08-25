


# Uber Go Style Guide


> Uber公司推出的Go语言规范，建议没看过的同学看一遍，内容同步到了我的GitHub上，后续会补充一些case，有兴趣的同学可以点击原文

## Table of Contents

- [介绍](#介绍)
- [指南](#指南)
    - [指向interface的指针](#指向interface的指针)
    - [验证接口合法性](#验证接口合法性)
    - [接收者和接口](#接收者和接口)
    - [Mutexes的零值是有效的](#Mutexes的零值是有效的)
    - [在边界拷贝Slices和Maps](#在边界拷贝Slices和Maps)
    - [使用Defer释放资源](#使用Defer释放资源)
    - [Channel大小应为0或1](#Channel大小应为0或1)
    - [枚举从 开始](#枚举从1开始)
    - [使用time包来处理时间](#使用time包来处理时间)
    - [错误](#错误)
        - [错误类型](#错误类型)
        - [错误包装](#错误包装)
        - [错误命名](#错误命名)
    - [处理断言失败](#处理断言失败)
    - [不要使用Panic](#不要使用Panic)
    - [使用go.uber.org/atomic](#使用gouberorgatomic)
    - [避免可变全局变量](#避免可变全局变量)
    - [避免在公共结构体中内嵌类型](#避免在公共结构体中内嵌类型)
    - [避免使用内置名称](#避免使用内置名称)
    - [避免使用init()](#避免使用init())
    - [优雅退出主函数](#优雅退出主函数)
        - [退出一次](#退出一次)
    - [在序列化结构中使用字段标记](#在序列化结构中使用字段标记)
- [性能](#性能)
    - [优先使用strconv而不是fmt](#优先使用strconv而不是fmt)
    - [避免字符串到字节的转换](#避免字符串到字节的转换)
    - [指定容器容量](#指定容器容量)
        - [指定Map容量](#指定Map容量)
        - [指定切片容量](#指定切片容量)
- [规范](#规范)
    - [避免过长的行](#避免过长的行)
    - [一致性](#一致性)
    - [相似的声明放在一组](#相似的声明放在一组)
    - [import分组](#import分组)
    - [包名](#包名)
    - [函数名](#函数名)
    - [导入别名](#导入别名)
    - [函数分组与顺序](#函数分组与顺序)
    - [减少嵌套](#减少嵌套)
    - [不必要的else](#不必要的else)
    - [顶层变量声明](#顶层变量声明)
    - [对于未导出的顶层常量和变量，使用_作为前缀](#对于未导出的顶层常量和变量，使用_作为前缀)
    - [结构体中的嵌入](#结构体中的嵌入)
    - [本地变量声明](#本地变量声明)
    - [nil是一个有效的slice](#nil是一个有效的slice)
    - [缩小变量作用域](#缩小变量作用域)
    - [避免参数语义不明确](#避免参数语义不明确)
    - [使用原始字符串字面值，避免转义](#使用原始字符串字面值，避免转义)
    - [初始化结构体](#初始化结构体)
        - [使用字段名初始化结构](#使用字段名初始化结构)
        - [省略结构中的零值字段](#省略结构中的零值字段)
        - [空结构体用var声明](#空结构体用var声明)
        - [初始化 Struct 引用](#初始化 Struct 引用)
    - [初始化 Maps](#初始化 Maps)
    - [字符串 string format](#字符串 string format)
    - [命名 Printf 样式的函数](#命名 Printf 样式的函数)
- [编程模式](#编程模式)
    - [表驱动测试](#表驱动测试)
    - [功能选项](#功能选项)
- [Linting](#linting)

## 介绍

风格是管理我们代码的规范。风格这个词有点问题，因为这些规划所涵盖的内容远远超过了源文件的格式化 -- gofmt为我们处理了这些。

本指南的目的是通过详细描述在Uber编写Go代码的注意事项来管理这种复杂性。这些规则的存在是为了保持代码库的可管理性，同时还允许工程师有效地使用Go语言的特性。

本指南最初是由Prashant Varanasi和Simon Newton创建的，是为了让一些同事尽快掌握Go的使用方法。多年来，我们根据其他人的反馈对它进行了修改。

这记录了我们在 Uber 所遵循的 Go 代码中的习惯性约定。其中很多是Go里面的一般准则，而其他的则是根据外部资源进行扩展：

1. [Effective Go](https://golang.org/doc/effective_go.html)
2. [Go Common Mistakes](https://github.com/golang/go/wiki/CommonMistakes)
3. [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)


我们的目标是使代码样本准确地适用于Go的两个最新的次要版本。

所有代码在通过golint和go vet运行时应该是没有错误的。我们建议将您的编辑器设置为：

- 保存时运行 goimports
- 运行 golint 和 go vet 检查错误

你可以在这里找到编辑器支持Go工具的信息：
<https://github.com/golang/go/wiki/IDEsAndTextEditorPlugins>



## 指南

### 指向interface的指针

你几乎不需要一个指向 interface 的指针，interface类型数据应该直接传递，但实际上 interface
底层是一个指针。

interface 类型包括两部分：

1. 一个指向特定类型的指针。可以将其视为 "类型"。
2. 数据指针。如果底层数据是指针，会被直接存储。如果底层数据是值，那会存储这个数据的指针。

如果你想要接口方法修改基础数据，那必须使用指针。



### 验证接口合法性

在编译期验证接口的合法性，需要验证的有：

- 验证导出类型在作为API时是否实现了特定接口
- 实现一个接口的导出和非导出类型是集合的一部分
- 违反接口合理性无法编译通过，通知用户


<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Handler struct {
  // ...
}



func (h *Handler) ServeHTTP(
  w http.ResponseWriter,
  r *http.Request,
) {
  ...
}
```

</td><td>

```go
type Handler struct {
  // ...
}

var _ http.Handler = (*Handler)(nil)

func (h *Handler) ServeHTTP(
  w http.ResponseWriter,
  r *http.Request,
) {
  // ...
}
```

</td></tr>
</tbody></table>

如果`*Handler`没有实现`http.Handler`接口，那么`var _ http.Handler = (*Handler)(nil)`语句在编译期就会报错；

赋值语句的右边部门应是断言类型的零值。对于指针类型(像`*Handler`)、slice和map类型，零值为`nil`,对于结构体类型，零值为空结构体，下面是空结构体的例子。


```go
type LogHandler struct {
  h   http.Handler
  log *zap.Logger
}

// LogHandler{}是空结构体
var _ http.Handler = LogHandler{}

func (h LogHandler) ServeHTTP(
  w http.ResponseWriter,
  r *http.Request,
) {
  // ...
}
```

### 接收者和接口

使用 值类型 接收者的方法既可以通过值调用，也可以通过指针调用。

使用 指针类型 接收者的方法只能通过指针或者 [addressable values](https://golang.org/ref/spec#Method_values)调用。

例如：

```go
type S struct {
  data string
}

func (s S) Read() string {
  return s.data
}

func (s *S) Write(str string) {
  s.data = str
}

sVals := map[int]S{1: {"A"}}

// 值类型可以调用Read()
sVals[1].Read()

// 值类型调用Write方法会报编译错误
//  sVals[1].Write("test")

sPtrs := map[int]*S{1: {"A"}}

// 指针类型 Read 和 Write 方法都可以调用
sPtrs[1].Read()
sPtrs[1].Write("test")
```
同样的，接口可以通过指针调用，即使这个方法的接收者是指类型。

```go
type F interface {
  f()
}

type S1 struct{}

func (s S1) f() {}

type S2 struct{}

func (s *S2) f() {}

s1Val := S1{}
s1Ptr := &S1{}
s2Val := S2{}
s2Ptr := &S2{}

var i F
i = s1Val
i = s1Ptr
i = s2Ptr

// 这个例子编译会报错，因为s2Val是值类型，而S2的方法里接收者是指针类型.
//   i = s2Val
```

Effective Go 有一段关于 [Pointers vs. Values] 的优秀讲解

[Pointers vs. Values]: https://golang.org/doc/effective_go.html#pointers_vs_values

### Mutexes的零值是有效的

`sync.Mutex` 和 `sync.RWMutex` 的零值是有效的，所以不需要实例化一个Mutex的指针。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
mu := new(sync.Mutex)
mu.Lock()
```

</td><td>

```go
var mu sync.Mutex
mu.Lock()
```

</td></tr>
</tbody></table>

如果结构体中包含mutex，在使用结构体的指针时，mutex应该是结构体的非指针字段，也不要把mutex内嵌到结构体中，即使结构体是非导出类型。


<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type SMap struct {
  sync.Mutex

  data map[string]string
}

func NewSMap() *SMap {
  return &SMap{
    data: make(map[string]string),
  }
}

func (m *SMap) Get(k string) string {
  m.Lock()
  defer m.Unlock()

  return m.data[k]
}
```

</td><td>

```go
type SMap struct {
  mu sync.Mutex

  data map[string]string
}

func NewSMap() *SMap {
  return &SMap{
    data: make(map[string]string),
  }
}

func (m *SMap) Get(k string) string {
  m.mu.Lock()
  defer m.mu.Unlock()

  return m.data[k]
}
```

</td></tr>

<tr><td>

隐式嵌入`Mutex`, 其`Lock`和`Unlock`方法是`SMap`公开API中不明确说明的一部分。

</td><td>

mutex和它`SMap`方法的实现细节对调用方屏蔽。


</td></tr>
</tbody></table>

### 在边界拷贝Slices和Maps

slice 和 map 类型包含指向data数据的指针，所以当你需要复制时应格外注意。

#### 接收 Slices 和 Maps

如果在函数调用中传递 map 或 slice, 请记住这个函数可以修改它。

<table>
<thead><tr><th>Bad</th> <th>Good</th></tr></thead>
<tbody>
<tr>
<td>

```go
func (d *Driver) SetTrips(trips []Trip) {
  d.trips = trips
}

trips := ...
d1.SetTrips(trips)

// 这样赋值会影响到 d1.trips
trips[0] = ...
```

</td>
<td>

```go
func (d *Driver) SetTrips(trips []Trip) {
  d.trips = make([]Trip, len(trips))
  copy(d.trips, trips)
}

trips := ...
d1.SetTrips(trips)

// 这样赋值不会影响 d1.trips(因为 SetTrips 内部有copy).
trips[0] = ...
```

</td>
</tr>

</tbody>
</table>

#### 返回 Slices 和 Maps

同样，请注意用户对 maps 或 slices 的修改暴露了内部状态。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Stats struct {
  mu sync.Mutex
  counters map[string]int
}

// Snapshot 返回当前的 stats.
func (s *Stats) Snapshot() map[string]int {
  s.mu.Lock()
  defer s.mu.Unlock()

  return s.counters
}

// snapshot 变量不在受 mutex锁保护，任何对 snapshot 的访问会受数据竟态的影响
snapshot := stats.Snapshot()
```

</td><td>

```go
type Stats struct {
  mu sync.Mutex
  counters map[string]int
}

func (s *Stats) Snapshot() map[string]int {
  s.mu.Lock()
  defer s.mu.Unlock()

  result := make(map[string]int, len(s.counters))
  for k, v := range s.counters {
    result[k] = v
  }
  return result
}

// snapshot 现在只是个copy。
snapshot := stats.Snapshot()
```

</td></tr>
</tbody></table>

### 使用Defer释放资源

在读写文件、使用锁时，使用 defer 释放资源

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
p.Lock()
if p.count < 10 {
  p.Unlock()
  return p.count
}

p.count++
newCount := p.count
p.Unlock()

return newCount

// 当有多分支 return 时，很容易漏写Unlock().
```

</td><td>

```go
p.Lock()
defer p.Unlock()

if p.count < 10 {
  return p.count
}

p.count++
return p.count

// 使用defer在一个地方 Unlock, 代码可读性更好
```

</td></tr>
</tbody></table>

调用 defer 的性能开销非常小，但如果你需要纳秒级别的函数调用，那可能需要避免使用 defer。
使用 defer 带来的可读性 胜过引入其它带来的性能开销。defer 尤其适用于适用于那些不仅是内存
放在的行数较多、逻辑较为复杂的大方法，这些方法中其他代码逻辑的执行成本比 defer 执行成本更大。

### Channel 大小应为 0 或 1

Channels 的大小应该是1或无缓冲的。默认情况下，channels 是无缓冲的，size为0。其他size需经过严格的审查。
考虑 channel 的size 是如何定义的，是什么造成了 channel 在负荷情况下被写满而无法写入，以及无法写入会发生什么。
 

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// 这个 size 对任何操作都够了！
c := make(chan int, 64)
```

</td><td>

```go
// Size 是1
c := make(chan int, 1) // or
// 无缓冲 channel
c := make(chan int)
```

</td></tr>
</tbody></table>

### 枚举类型值从1开始

在Go中声明枚举值的标准方法是使用`const`包`iota`。由于变量的默认值为0，因此枚举类型的值需要从1开始。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Operation int

const (
  Add Operation = iota
  Subtract
  Multiply
)

// Add=0, Subtract=1, Multiply=2
```

</td><td>

```go
type Operation int

const (
  Add Operation = iota + 1
  Subtract
  Multiply
)

// Add=1, Subtract=2, Multiply=3
```

</td></tr>
</tbody></table>

当你需要将0值视为默认行为时，枚举类型从0开始是有意义的。

```go
type LogOutput int

const (
  LogToStdout LogOutput = iota
  LogToFile
  LogToRemote
)

// LogToStdout=0, LogToFile=1, LogToRemote=2
```

### 使用time包来处理时间

时间处理很复杂，关于时间错误预估有以下这些点。

1. 一天有24小时
2. 一小时有60分钟
3. 一周有7天
4. 一年有365天 
5. [其他易错点](https://infiniteundo.com/post/25326999628/falsehoods-programmers-believe-about-time)

举例来说, *1* 表示在一个时间点加上24小时并不一定会产生新的一天。

因此，在处理时间时应始终使用`"time"`包，因为它会用更安全、准确的方式来处理这些不正确的假设。

[`"time"`]: https://golang.org/pkg/time/

#### 用 `time.Time` 表示瞬时时间

需要瞬时时间语义时，使用[`time.Time`] ，在进行比较、增加或减少时间段时，使用`time.Time`包里的方法。

[`time.Time`]: https://golang.org/pkg/time/#Time

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func isActive(now, start, stop int) bool {
  return start <= now && now < stop
}
```

</td><td>

```go
func isActive(now, start, stop time.Time) bool {
  return (start.Before(now) || start.Equal(now)) && now.Before(stop)
}
```

</td></tr>
</tbody></table>

#### 用 `time.Duration` 表示时间段

应使用 [`time.Duration`] 来表示时间段

[`time.Duration`]: https://golang.org/pkg/time/#Duration

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func poll(delay int) {
  for {
    // ...
    time.Sleep(time.Duration(delay) * time.Millisecond)
  }
}

poll(10) // 这里单位是秒还是毫秒
```

</td><td>

```go
func poll(delay time.Duration) { // 使用 time.Duration 表示时间段
  for {
    // ...
    time.Sleep(delay)
  }
}

poll(10*time.Second) // 明确单位
```

</td></tr>
</tbody></table>

回到刚刚的例子，在一个瞬时时间加上24小时，怎么加这个 "24小时" 取决于我们的意图。如果我们想获取
下一天的当前时间，我们应该使用 [`Time.AddDate`]。如果我们想获取比当前时间晚24小时的瞬时时间，
我们应该应该使用 [`Time.Add`]。

[`Time.AddDate`]: https://golang.org/pkg/time/#Time.AddDate
[`Time.Add`]: https://golang.org/pkg/time/#Time.Add

```go
newDay := t.AddDate(0 /* years */, 0 /* months */, 1 /* days */)
maybeNewDay := t.Add(24 * time.Hour)
```

#### 对外交互 使用 `time.Time` 和 `time.Duration` 

在对外交互时尽可能使用 `time.Duration` 和 `time.Time`，例如：

- Command-line 标记: [`flag`] 通过支持 [`time.ParseDuration`] 来支持 `time.Duration`
- JSON: [`encoding/json`] 通过  [`UnmarshalJSON` 方法] 支持把 `time.Time` 解码为 [RFC 3339] 字符串
- SQL: [`database/sql`] 支持把 `DATETIME` 或 `TIMESTAMP` 类型转化为 `time.Time` 
- YAML: [`gopkg.in/yaml.v2`] 支持把 `time.Time` 作为一个 [RFC 3339] 字符串, 通过支持[`time.ParseDuration`] 来支持`time.Duration`。

  [`flag`]: https://golang.org/pkg/flag/
  [`time.ParseDuration`]: https://golang.org/pkg/time/#ParseDuration
  [`encoding/json`]: https://golang.org/pkg/encoding/json/
  [RFC 3339]: https://tools.ietf.org/html/rfc3339
  [`UnmarshalJSON` method]: https://golang.org/pkg/time/#Time.UnmarshalJSON
  [`database/sql`]: https://golang.org/pkg/database/sql/
  [`gopkg.in/yaml.v2`]: https://godoc.org/gopkg.in/yaml.v2


如果交互中不支持使用`time.Duration`，那字段名中应包含单位，类型应为`int`或`float64`。

例如, 由于 `encoding/json` 不支持 `time.Duration` 类型, 因此字段名中应包含时间单位。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// {"interval": 2}
type Config struct {
  Interval int `json:"interval"`
}
```

</td><td>

```go
// {"intervalMillis": 2000}
type Config struct { // Millis 是单位
  IntervalMillis int `json:"intervalMillis"`
}
```

</td></tr>
</tbody></table>

如果在交互中不能使用 `time.Time`，除非有额外约定，否则应该使用 `string` 和 [RFC 3339]定义的
时间戳格式。默认情况下， [`Time.UnmarshalText`] 使用这种格式，并可以通过[`time.RFC3339`]在
`Time.Format` 和 `time.Parse` 中使用。

[`Time.UnmarshalText`]: https://golang.org/pkg/time/#Time.UnmarshalText
[`time.RFC3339`]: https://golang.org/pkg/time/#RFC3339

尽管在实践中这不是什么问题，但是你需要记住`"time"`不能解析闰秒时间戳([8728]),在计算中也不考虑闰秒([15190])。
因此如果你要比较两个瞬时时间，比较结果不会包含这两个瞬时时间可能会出现的闰秒。

[8728]: https://github.com/golang/go/issues/8728
[15190]: https://github.com/golang/go/issues/15190

<!-- TODO: section on String methods for enums -->

### 错误

#### 错误类型

声明错误的选项很少。 在为你的代码选择合适的用例之前，考虑这些事项：

- 调用方需要匹配错误吗，还是调用方需要自己处理错误。如果需要匹配错误，那应该声明顶级
  错误类型或自定义类型 来让 [`errors.Is`] 或 [`errors.As`] 匹配。
- 错误信息是静态字符串吗，还是说错误信息是需要上下文的动态字符串。如果是静态字符串，
  我们可以使用 [`errors.New`]，如果是动态字符串，我们应该使用[`fmt.Errorf`]来
  自定义错误类型。
- 我们是否正在传递由下游返回的新的错误类型，如果是这样，参考[section on error wrapping](#error-wrapping).

[`errors.Is`]: https://golang.org/pkg/errors/#Is
[`errors.As`]: https://golang.org/pkg/errors/#As

| 错误匹配? | 错误信息 | 使用                       |
|-------|------|--------------------------|
| 无     | 静态   | [`errors.New`]           |
| 无     | 动态   | [`fmt.Errorf`]           |
| 有     | 静态   | 用[`errors.New`]声明顶级错误类型  |
| 有     | 动态   | 定制 `error` 类型            |

[`errors.New`]: https://golang.org/pkg/errors/#New
[`fmt.Errorf`]: https://golang.org/pkg/fmt/#Errorf

举例，使用 [`errors.New`] 表示一个静态字符串错误。如果调用方需要匹配并处理这个错误，
就把这个错误声明为变量来支持和 `errors.Is` 匹配。

<table>
<thead><tr><th>无错误匹配</th><th>有错误匹配</th></tr></thead>
<tbody>
<tr><td>

```go
// package foo

func Open() error {
  return errors.New("could not open")
}

// package bar

if err := foo.Open(); err != nil {
  //无法处理错误
  panic("unknown error")
}
```

</td><td>

```go
// package foo

var ErrCouldNotOpen = errors.New("could not open")

func Open() error {
  return ErrCouldNotOpen
}

// package bar

if err := foo.Open(); err != nil {
  if errors.Is(err, foo.ErrCouldNotOpen) {
    // 处理这个错误
  } else {
    panic("unknown error")
  }
}
```

</td></tr>
</tbody></table>

对于动态字符串的错误，如果调用方不需要匹配就是用[`fmt.Errorf`]，需要匹配就搞一个自定义`error`。

<table>
<thead><tr><th>无错误匹配</th><th>有错误匹配</th></tr></thead>
<tbody>
<tr><td>

```go
// package foo

func Open(file string) error {
  return fmt.Errorf("file %q not found", file)
}

// package bar

if err := foo.Open("testfile.txt"); err != nil {
  // Can't handle the error.
  panic("unknown error")
}
```

</td><td>

```go
// package foo

type NotFoundError struct {
  File string
}

func (e *NotFoundError) Error() string {
  return fmt.Sprintf("file %q not found", e.File)
}

func Open(file string) error {
  return &NotFoundError{File: file}
}


// package bar

if err := foo.Open("testfile.txt"); err != nil {
  var notFound *NotFoundError
  if errors.As(err, &notFound) {
    // handle the error
  } else {
    panic("unknown error")
  }
}
```

</td></tr>
</tbody></table>

注意，如果你的包里导出了错误变量或错误类型，那这个错误将变成你包里公共API的一部分。

#### 错误包装

调用函数失败时，有三种选择供你选择：

- 返回原始错误
- 用 `fmt.Errorf` 和 `%w` 包装上下文信息
- 用 `fmt.Errorf` 和 `%v` 包装上下文信息

返回原始错误不会附加上下文信息，这样就保持了原始错误类型和信息。比较适用于lib库类型代码
展示底层错误信息。

如果不是lib库，就需要增加所需的上下文信息，不然就会出现 "connection refused" 这样非常
模糊的错误，理论上应该添加上下文，来得到这样的报错信息："call service foo: connection refused"。

在错误类型上使用 `fmt.Errorf` 来添加上下文信息，根据调用方不同的使用方式，可以选择 `%w` 或 `%v` 动词。

- 如果调用方需要访问底层错误，使用`%w`动词，这是一个用来包装错误的动词，如果你在代码中使用到了它，请注意
  调用方会对此产生依赖，所以当你的包装的错误是用`var`声明的已知类型，需要在你的代码里对其进行测试。
- 使用 `%v` 会混淆你的底层错误类型，调用方将无法进行匹配，如果有匹配需求，应该使用`%w`动词。


当为返回错误增加上下文信息时，避免在上下文中增加像 "failed to" 这样的没啥用的短语，这样没用的短语在错误
堆栈中堆积起来的话，反而不利于你定位bug。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
s, err := store.New()
if err != nil {
    return fmt.Errorf(
        "failed to create new store: %w", err)
}
```

</td><td>

```go
s, err := store.New()
if err != nil {
    return fmt.Errorf(
        "new store: %w", err)
}
```

</td></tr><tr><td>

```
failed to x: failed to y: failed to create new store: the error
```

</td><td>

```
x: y: new store: the error
```

</td></tr>
</tbody></table>

然而当你的错误传给别的系统时，错误信息应该足够清晰。(比如, 错误信息在日志中以 "Failed" 开头) 

其他参考信息： [Don't just check errors, handle them gracefully].

[`"pkg/errors".Cause`]: https://godoc.org/github.com/pkg/errors#Cause
[Don't just check errors, handle them gracefully]: https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully

#### 错误命名

对于全局变量类型，根据是否导出使用  `Err` 或 `err` 前缀。 详情参考：[Prefix Unexported Globals with _](#prefix-unexported-globals-with-_).

```go
var (
  // 以下两个错误是导出类型，所以他们的命名以 Err 作为开头，用户可以使用 errors.Is 来匹配错误类型
  ErrBrokenLink = errors.New("link is broken")
  ErrCouldNotOpen = errors.New("could not open")

  // 这个错误是非导出类型，不会作为我们公共API的一部分，但是你可以在包内使用errors.Is匹配它。
  errNotFound = errors.New("not found")
)
```

对于自定义错误类型，请使用`Error`后缀。

```go
// 这个错误是被导出的，用户可以使用errors.As去匹配它。
type NotFoundError struct {
  File string
}

func (e *NotFoundError) Error() string {
  return fmt.Sprintf("file %q not found", e.File)
}

// 这个错误是非导出类型，不会作为我们公共API的一部分，但是你可以在包内使用errors.As匹配它。 

type resolveError struct {
  Path string
}

func (e *resolveError) Error() string {
  return fmt.Sprintf("resolve %q", e.Path)
}
```

todo： errors.Is和errors.As的区别

### 处理断言失败

在不正确的类型断言上 使用单返回值来处理会导致 panic, 因此请使用 "comma ok" 习俗.

[type assertion]: https://golang.org/ref/spec#Type_assertions

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
t := i.(string)
```

</td><td>

```go
t, ok := i.(string)
if !ok {
  // 在这里优雅处理错误
}
```

</td></tr>
</tbody></table>

<!-- TODO: There are a few situations where the single assignment form is
fine. -->

### 不要使用Panic

在生产环境的业务代码避免使用panic。Panics 是级联问题[cascading failures]的主要来源。
如果发生错误，函数必须返回错误，让调用方决定如何处理这种情况。

[cascading failures]: https://en.wikipedia.org/wiki/Cascading_failure

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func run(args []string) {
  if len(args) == 0 {
    panic("an argument is required")
  }
  // ...
}

func main() {
  run(os.Args[1:])
}
```

</td><td>

```go
func run(args []string) error {
  if len(args) == 0 {
    return errors.New("an argument is required")
  }
  // ...
  return nil
}

func main() {
  if err := run(os.Args[1:]); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
```

</td></tr>
</tbody></table>

Panic/recover 不是错误处理策略。当系统发生像空指针异常这种 不可恢复的 Fatal 异常时，才需要使用Panic。
唯一的意外情况是项目启动：如果程序启动阶段出现问题需要抛出异常。

```go
var _statusTemplate = template.Must(template.New("name").Parse("_statusHTML"))
```

即使在测试中，优先使用`t.Fatal` 或 `t.FailNow` 而不是异常，来确保失败情况被记录。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// func TestFoo(t *testing.T)

f, err := os.CreateTemp("", "test")
if err != nil {
  panic("failed to set up test")
}
```

</td><td>

```go
// func TestFoo(t *testing.T)

f, err := os.CreateTemp("", "test")
if err != nil {
  t.Fatal("failed to set up test")
}
```

</td></tr>
</tbody></table>

<!-- TODO: Explain how to use _test packages. -->

### 使用 go.uber.org/atomic

使用 [sync/atomic] 包的原子操作
Atomic operations with the [sync/atomic] package operate on the raw types
(`int32`, `int64`, etc.) so it is easy to forget to use the atomic operation to
read or modify the variables.

[go.uber.org/atomic] adds type safety to these operations by hiding the
underlying type. Additionally, it includes a convenient `atomic.Bool` type.

[go.uber.org/atomic]: https://godoc.org/go.uber.org/atomic
[sync/atomic]: https://golang.org/pkg/sync/atomic/

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type foo struct {
  running int32  // atomic
}

func (f* foo) start() {
  if atomic.SwapInt32(&f.running, 1) == 1 {
     // already running…
     return
  }
  // start the Foo
}

func (f *foo) isRunning() bool {
  return f.running == 1  // race!
}
```

</td><td>

```go
type foo struct {
  running atomic.Bool
}

func (f *foo) start() {
  if f.running.Swap(true) {
     // already running…
     return
  }
  // start the Foo
}

func (f *foo) isRunning() bool {
  return f.running.Load()
}
```

</td></tr>
</tbody></table>

### 避免使用全局可变对象

Avoid mutating global variables, instead opting for dependency injection.
This applies to function pointers as well as other kinds of values.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// sign.go

var _timeNow = time.Now

func sign(msg string) string {
  now := _timeNow()
  return signWithTime(msg, now)
}
```

</td><td>

```go
// sign.go

type signer struct {
  now func() time.Time
}

func newSigner() *signer {
  return &signer{
    now: time.Now,
  }
}

func (s *signer) Sign(msg string) string {
  now := s.now()
  return signWithTime(msg, now)
}
```
</td></tr>
<tr><td>

```go
// sign_test.go

func TestSign(t *testing.T) {
  oldTimeNow := _timeNow
  _timeNow = func() time.Time {
    return someFixedTime
  }
  defer func() { _timeNow = oldTimeNow }()

  assert.Equal(t, want, sign(give))
}
```

</td><td>

```go
// sign_test.go

func TestSigner(t *testing.T) {
  s := newSigner()
  s.now = func() time.Time {
    return someFixedTime
  }

  assert.Equal(t, want, s.Sign(give))
}
```

</td></tr>
</tbody></table>

### 避免在公共结构体中内嵌类型

嵌入类型会暴露实现细节，无法类型演化，让文档也变得模糊。

假设你用`AbstractList`结构体实现了公共的 list 方法，避免在其他实现中内嵌`AbstractList`类型。
而是应该在其他结构体中显式声明list，并在方法实现中调用list的方法。

```go
type AbstractList struct {}

// Add adds an entity to the list.
func (l *AbstractList) Add(e Entity) {
  // ...
}

// Remove removes an entity from the list.
func (l *AbstractList) Remove(e Entity) {
  // ...
}
```

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// ConcreteList is a list of entities.
type ConcreteList struct {
  *AbstractList
}
```

</td><td>

```go
// ConcreteList is a list of entities.
type ConcreteList struct {
  list *AbstractList
}

// Add adds an entity to the list.
func (l *ConcreteList) Add(e Entity) {
  l.list.Add(e)
}

// Remove removes an entity from the list.
func (l *ConcreteList) Remove(e Entity) {
  l.list.Remove(e)
}
```

</td></tr>
</tbody></table>

Go允许内嵌类型[type embedding]作为组合和继承的折中方案。外部的结构体会获得内嵌类型的隐式拷贝。默认情况下，内嵌类型的方法会嵌入实例的同一方法。

[type embedding]: https://golang.org/doc/effective_go.html#embedding

外部的结构体会获取嵌入类型的同名字段。如果嵌入类型的字段是公开(public)的，那嵌入后也是公开的。
为保证向后兼容性，外部结构体未来每个版本都需要保留嵌入类型。

很少场景需要嵌入类型，虽然嵌入类型很方便，让你避免编写冗长的方法。

即使是用interface而不是结构体来嵌入方法，这是给开发人员带来了一定的灵活性，但是仍然暴露了具体实现列表的抽象细节。


<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// AbstractList is a generalized implementation
// for various kinds of lists of entities.
type AbstractList interface {
  Add(Entity)
  Remove(Entity)
}

// ConcreteList is a list of entities.
type ConcreteList struct {
  AbstractList
}
```

</td><td>

```go
// ConcreteList is a list of entities.
type ConcreteList struct {
  list AbstractList
}

// Add adds an entity to the list.
func (l *ConcreteList) Add(e Entity) {
  l.list.Add(e)
}

// Remove removes an entity from the list.
func (l *ConcreteList) Remove(e Entity) {
  l.list.Remove(e)
}
```

</td></tr>
</tbody></table>

不管是嵌入结构体还是嵌入接口，都会限制类型的演化。

- 若嵌入接口，当你增加一个方法是一种破坏性改变
- 若嵌入结构体，当你删除一个方法是一种破坏性改变
- 删除嵌入类型是一种破坏性改变
- 即使用满足接口约束的类型去替换嵌入类型，也是一种破坏性改变

尽管编写内嵌类型已实现的方法是乏味的。但是这些工作隐藏了实现细节，留下了更多更改的机会，
并消除了在文档中发现完整List接口的间接方法。

### 避免使用内建命名

Go语言的spec中列举了一些内建命名，在你的Go程序中应该避免使用预声明的标识符；

根据上下文的不同，用预声明标识符命名变量可能会在当前作用域下覆盖官方标识符，让你的代码变得难以理解。
最好的情况下，编译器会直接报错，最糟糕的情况下，这样的代码会引入难以排查的bug。

[language specification]: https://golang.org/ref/spec
[predeclared identifiers]: https://golang.org/ref/spec#Predeclared_identifiers

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
var error string
// `error` 覆盖了内建的error

// or

func handleErrorMessage(error string) {
    // `error` 覆盖了内建的error
}
```

</td><td>

```go
var errorMessage string
// `error` 指向内置的 error 

// or

func handleErrorMessage(msg string) {
    // `error` 指向内置的 error
}
```

</td></tr>
<tr><td>

```go
type Foo struct {
    // While these fields technically don't
    // constitute shadowing, grepping for
    // `error` or `string` strings is now
    // ambiguous.
    error  error
    string string
}

func (f Foo) Error() error {
    // `error` and `f.error` are
    // visually similar
    return f.error
}

func (f Foo) String() string {
    // `string` and `f.string` are
    // visually similar
    return f.string
}
```

</td><td>

```go
type Foo struct {
    // `error` and `string` strings are
    // now unambiguous.
    err error
    str string
}

func (f Foo) Error() error {
    return f.err
}

func (f Foo) String() string {
    return f.str
}
```

</td></tr>
</tbody></table>


注意当你使用预声明标识符时编译器不会报错，但是像 `go vet` 这样的工具会告诉你标识符被覆盖的情况。


### 避免使用init()

尽可能避免使用`init()`。如果实在依赖 `init()`，可以使用以下方式：

1. 不管程序环境或调用方式如何，初始化要完全确定。
2. 避免依赖其他`init()`函数的顺序或者产生的结果。虽然`init()`顺序是明确的，但是代码可以更改。
   `init()`函数之间的关系会让代码变得易错和脆弱。
3. 避免读写全局变量、环境变量，比如机器信息、环境变量、工作目录，程序的参数和输入等等。
4. 避免 I/O 操作，比如文件系统，网络和系统调用。

如果代码不能满足这些需求，那可能属于帮助代码，需要作为`main()`函数的一部分进行调用(或者封装初始化逻辑，让main
函数去调用)。 需要注意的是，被其他模块依赖的代码应该完全指定初始化顺序的确定性，而不是依赖"初始化魔法"。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Foo struct {
    // ...
}

var _defaultFoo Foo

func init() {
    _defaultFoo = Foo{
        // ...
    }
}
```

</td><td>

```go
var _defaultFoo = Foo{
    // ...
}

// or, better, for testability:

var _defaultFoo = defaultFoo()

func defaultFoo() Foo {
    return Foo{
        // ...
    }
}
```

</td></tr>
<tr><td>

```go
type Config struct {
    // ...
}

var _config Config

func init() {
    // Bad: based on current directory
    cwd, _ := os.Getwd()

    // Bad: I/O
    raw, _ := os.ReadFile(
        path.Join(cwd, "config", "config.yaml"),
    )

    yaml.Unmarshal(raw, &_config)
}
```

</td><td>

```go
type Config struct {
    // ...
}

func loadConfig() Config {
    cwd, err := os.Getwd()
    // handle err

    raw, err := os.ReadFile(
        path.Join(cwd, "config", "config.yaml"),
    )
    // handle err

    var config Config
    yaml.Unmarshal(raw, &config)

    return config
}
```

</td></tr>
</tbody></table>

但是在某些情况下，`init()`函数可能更具优势：

- 单个赋值语句中无法表示的复杂表达式
- 插件钩子，比如 `database/sql`，编码信息注册表等
- 对 [Google Cloud Functions] 和其他形式确定性预计算的优化

  [Google Cloud Functions]: https://cloud.google.com/functions/docs/bestpractices/tips#use_global_variables_to_reuse_objects_in_future_invocations

### 优雅退出主函数

 Go程序使用[`os.Exit`] 或 [`log.Fatal*`]来立即退出。(Panic 不是优雅的程序退出方式，可以参考 [don't panic](#dont-panic))

[`os.Exit`]: https://golang.org/pkg/os/#Exit
[`log.Fatal*`]: https://golang.org/pkg/log/#Fatal

应该只在`main()`函数里调用`os.Exit` 或 `log.Fatal*`函数。其他函数应该返回错误来表示失败，在`main`中进行退出。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func main() {
  body := readFile(path)
  fmt.Println(body)
}

func readFile(path string) string {
  f, err := os.Open(path)
  if err != nil {
    log.Fatal(err)
  }

  b, err := io.ReadAll(f)
  if err != nil {
    log.Fatal(err)
  }

  return string(b)
}
```

</td><td>

```go
func main() {
  body, err := readFile(path)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println(body)
}

func readFile(path string) (string, error) {
  f, err := os.Open(path)
  if err != nil {
    return "", err
  }

  b, err := io.ReadAll(f)
  if err != nil {
    return "", err
  }

  return string(b), nil
}
```

</td></tr>
</tbody></table>

程序中多个函数都能退出的话会有一些问题：

- 不明显的控制流：多个函数都能退出的话，找出程序的控制流会变得困难。
- 测试困难：如果一个函数让程序退出，那它也会让测试退出。这样会让函数难以测试。而且可能会让`go text`
  无法测试其他函数。
- 跳过清理：当一个函数退出程序时，会跳过已经进入`defer`队列的函数调用。这样会增加跳过清理任务的风险。

#### 一次性退出

有条件的情况下，`main()`函数中最好只调用`os.Exit` 或 `log.Fatal` 一次。如果有多种错误情况会停止
程序的执行，将这些错误放在一个独立的函数中，并返回错误，`main()`中处理错误并退出。

把所有的关键逻辑放在一个独立的可测试的函数中，会让你的`main()`函数变得简短。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
package main

func main() {
  args := os.Args[1:]
  if len(args) != 1 {
    log.Fatal("missing file")
  }
  name := args[0]

  f, err := os.Open(name)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  // If we call log.Fatal after this line,
  // f.Close will not be called.

  b, err := io.ReadAll(f)
  if err != nil {
    log.Fatal(err)
  }

  // ...
}
```

</td><td>

```go
package main

func main() {
  if err := run(); err != nil {
    log.Fatal(err)
  }
}

func run() error {
  args := os.Args[1:]
  if len(args) != 1 {
    return errors.New("missing file")
  }
  name := args[0]

  f, err := os.Open(name)
  if err != nil {
    return err
  }
  defer f.Close()

  b, err := io.ReadAll(f)
  if err != nil {
    return err
  }

  // ...
}
```

</td></tr>
</tbody></table>

### 在序列化结构体中使用字段标签。

要编码成JSON、YAML或其他支持tag格式的结构体字段应该用指定对应项tag标签进行注释。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Stock struct {
  Price int
  Name  string
}

bytes, err := json.Marshal(Stock{
  Price: 137,
  Name:  "UBER",
})
```

</td><td>

```go
type Stock struct {
  Price int    `json:"price"`
  Name  string `json:"name"`
  // Safe to rename Name to Symbol.
}

bytes, err := json.Marshal(Stock{
  Price: 137,
  Name:  "UBER",
})
```

</td></tr>
</tbody></table>

Rationale:
结构体的序列化方式是不同系统通信的契约。修改结构体的结构和字段会破坏这个契约。在结构体中声明tag
可以防止重构结构体中意外违反约定。

## 性能

性能方面的指导准则只适用于高频调用场景。

### 使用 strconv 而不是 fmt

当需要原始类型和字符串互相转化时，`strconv`比`fmt`性能更好。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
for i := 0; i < b.N; i++ {
  s := fmt.Sprint(rand.Int())
}
```

</td><td>

```go
for i := 0; i < b.N; i++ {
  s := strconv.Itoa(rand.Int())
}
```

</td></tr>
<tr><td>

```
BenchmarkFmtSprint-4    143 ns/op    2 allocs/op
```

</td><td>

```
BenchmarkStrconv-4    64.2 ns/op    1 allocs/op
```

</td></tr>
</tbody></table>

### 避免 字符串到字节的转化

不要在for循环中创建[]byte类型，应该在for循环开始前把[]byte数据准备好。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
for i := 0; i < b.N; i++ {
  w.Write([]byte("Hello world"))
}
```

</td><td>

```go
data := []byte("Hello world")
for i := 0; i < b.N; i++ {
  w.Write(data)
}
```

</td></tr>
<tr><td>

```
BenchmarkBad-4   50000000   22.2 ns/op
```

</td><td>

```
BenchmarkGood-4  500000000   3.25 ns/op
```

</td></tr>
</tbody></table>

### 预先指定容器类型的容量

尽可能指定容器类型变量的容量来预先分配容器类型所需的内存大小。这样可以预防由于后续由分配元素
（由于拷贝或重新指定容器大小）而导致的内存分配。

#### 指定Map容量

如果有可能，用make来初始化map类型，并指定map的大小。

```go
make(map[T1]T2, hint)
```

使用 `make()` 初始化map时，提供一个容量来执行size，这样会减少后续将给map添加元素时引起的内存分配。

注意，和 slice 不同，给map指定容量不意味着抢占式内存分配完成，而是会用于预估的哈希表内部 buckets。
因此，当你给 map 添加元素，或者给 map 指定值时，仍有可能发生内存分配。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
m := make(map[string]os.FileInfo)

files, _ := os.ReadDir("./files")
for _, f := range files {
    m[f.Name()] = f
}
```

</td><td>

```go

files, _ := os.ReadDir("./files")

m := make(map[string]os.DirEntry, len(files))
for _, f := range files {
    m[f.Name()] = f
}
```

</td></tr>
<tr><td>

`m` 没有指定内存大小，因此在运行期间可能会有更多的内存分配。

</td><td>

`m` 指定了内存大小，因此在运行期间可能会有较少的内存分配。


</td></tr>
</tbody></table>

#### 指定Slice容量

如果有可能的话，在使用`make()`初始化slice的时候提供容量大小，尤其是后面需要 append 操作时。

```go
make([]T, length, capacity)
```

和 map 不同，slice的容量不是一个提示：编译器会根据 `make()` 提供的容量信息申请足够的内存，
这意味着后续的 `append()` 操作不会申请内存（除非slice的长度和容量相等，这样的话后续添加元素
会申请内存来调整 slice 的大小）。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
for n := 0; n < b.N; n++ {
  data := make([]int, 0)
  for k := 0; k < size; k++{
    data = append(data, k)
  }
}
```

</td><td>

```go
for n := 0; n < b.N; n++ {
  data := make([]int, 0, size)
  for k := 0; k < size; k++{
    data = append(data, k)
  }
}
```

</td></tr>
<tr><td>

```
BenchmarkBad-4    100000000    2.48s
```

</td><td>

```
BenchmarkGood-4   100000000    0.21s
```

</td></tr>
</tbody></table>

## 规范

### 避免代码过长

避免由于长度过长而需要水平滚动或者太需要转动头部的代码。

我们建议一行代码长度为 **99个字符**，如果代码超过了这个限制就应该换行。但是这也不是绝对的，代码
可以超过这个限制。

### 一致性

本文的一些指导准则可以被客观评估，其他准则可以根据实际情况进行选择。

但是最重要是，在你的代码中要**保持一致**。

一致性的代码更易于维护，更容易合理化，需要的认知开销较少；当新的管理出现时或 bug 被修复后也更易于
迁移和更新。

与之相反，如果单个代码库中有多种冲突的风格，会让维护成本升高、不确定性增高、认知不协调，这些问题会
导致开发效率降低，code review 困难，且容易产生 bug。

当你在代码库中实施标准时，建议最低在包层面进行修改：在子包层面进行应用违反了上述约定，因为在一种代码
里面引入了多种风格。

### 相似声明放一组

Go语言支持组引用。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
import "a"
import "b"
```

</td><td>

```go
import (
  "a"
  "b"
)
```

</td></tr>
</tbody></table>

组声明同样适用于常量、变量和类型声明。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go

const a = 1
const b = 2



var a = 1
var b = 2



type Area float64
type Volume float64
```

</td><td>

```go
const (
  a = 1
  b = 2
)

var (
  a = 1
  b = 2
)

type (
  Area float64
  Volume float64
)
```

</td></tr>
</tbody></table>

注意只把相关的变量声明到一个组里，不想管的声明放在多个组里。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Operation int

const (
  Add Operation = iota + 1
  Subtract
  Multiply
  EnvVar = "MY_ENV"
)
```

</td><td>

```go
type Operation int

const (
  Add Operation = iota + 1
  Subtract
  Multiply
)

const EnvVar = "MY_ENV"
```

</td></tr>
</tbody></table>

组声明不限制在哪使用。比如，你可以在函数中使用组声明。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func f() string {
  red := color.New(0xff0000)
  green := color.New(0x00ff00)
  blue := color.New(0x0000ff)

  // ...
}
```

</td><td>

```go
func f() string {
  var (
    red   = color.New(0xff0000)
    green = color.New(0x00ff00)
    blue  = color.New(0x0000ff)
  )

  // ...
}
```

</td></tr>
</tbody></table>

例外：对于变量声明，尤其是函数中的变量声明，不管他们之间是否有关系，都应该被放在一个组内。 

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func (c *client) request() {
  caller := c.name
  format := "json"
  timeout := 5*time.Second
  var err error

  // ...
}
```

</td><td>

```go
func (c *client) request() {
  var (
    caller  = c.name
    format  = "json"
    timeout = 5*time.Second
    err error
  )

  // ...
}
```

</td></tr>
</tbody></table>

### 包导入顺序

包中应该有两种导入顺序：

- 标准库
- 其他库

默认情况下，应该使用 goimports 的导入顺序。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
import (
  "fmt"
  "os"
  "go.uber.org/atomic"
  "golang.org/x/sync/errgroup"
)
```

</td><td>

```go
import (
  "fmt"
  "os"

  "go.uber.org/atomic"
  "golang.org/x/sync/errgroup"
)
```

</td></tr>
</tbody></table>

### 包命名

当命名包时，应按照下面原则命名：

- 全小写字母。无大写字母或下划线。
- 大多数导入包的情况下，不需要对包重新命名。
- 简短而简洁，因为当你使用包名时你都需要完成输入包名称。
- 不要使用复数。比如：命名为 `net/url`, 而不是 `net/urls`。
- 不要使用"common", "util", "shared", 或 "lib"。这些包含有信息太少了。

可以参考 [Package Names] 和 [Style guideline for Go packages].

[Package Names]: https://blog.golang.org/package-names
[Style guideline for Go packages]: https://rakyll.org/style-packages/

### 函数命名

我们遵守 Go 社区 [MixedCaps for function names] 约定。一种其他情况是使用测试函数。测试函数
命名可以包含下划线以便于相关测试函数进行分组。比如：`TestMyFunction_WhatIsBeingTested`。

[MixedCaps for function names]: https://golang.org/doc/effective_go.html#mixed-caps

### 导入别名

如果包名称和导入路径最后一个元素不匹配，就需要使用导入别名。

```go
import (
  "net/http"

  client "example.com/client-go"
  trace "example.com/trace/v2"
)
```

其他情况下，除非几个包之间有导入冲突，否则应该避免使用导入别名。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
import (
  "fmt"
  "os"


  nettrace "golang.net/x/trace"
)
```

</td><td>

```go
import (
  "fmt"
  "os"
  "runtime/trace"

  nettrace "golang.net/x/trace"
)
```

</td></tr>
</tbody></table>

### 函数分组和排序

- 函数应该按照大概调用顺序排序。
- 一个文件中的函数应该按照接收者分组。

因此，导入的函数时，应该放在 `struct`, `const`, `var` 的下面。

像 `newXYZ()`/`NewXYZ()` 这样的函数可能会出现在类型定义下、接收者的其他方法之上。

由于函数是按照接收者进行分组的，普通的工具函数应该放在文件末尾。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func (s *something) Cost() {
  return calcCost(s.weights)
}

type something struct{ ... }

func calcCost(n []int) int {...}

func (s *something) Stop() {...}

func newSomething() *something {
    return &something{}
}
```

</td><td>

```go
type something struct{ ... }

func newSomething() *something {
    return &something{}
}

func (s *something) Cost() {
  return calcCost(s.weights)
}

func (s *something) Stop() {...}

func calcCost(n []int) int {...}
```

</td></tr>
</tbody></table>

### 减少嵌套

代码应通过尽早处理错误/特殊情况尽早处理/循环中使用 continue 等手段，来减少嵌套代码过多问题。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
for _, v := range data {
  if v.F1 == 1 {
    v = process(v)
    if err := v.Call(); err == nil {
      v.Send()
    } else {
      return err
    }
  } else {
    log.Printf("Invalid v: %v", v)
  }
}
```

</td><td>

```go
for _, v := range data {
  if v.F1 != 1 {
    log.Printf("Invalid v: %v", v)
    continue
  }

  v = process(v)
  if err := v.Call(); err != nil {
    return err
  }
  v.Send()
}
```

</td></tr>
</tbody></table>

### 没用的else

如啊a变量在两个if分支中都进行赋值操作，则可以被替换为只在一个if分支中声明。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
var a int
if b {
  a = 100
} else {
  a = 10
}
```

</td><td>

```go
a := 10
if b {
  a = 100
}
```

</td></tr>
</tbody></table>

### 顶层变量声明

在顶层使用`var`来声明变量。不要指定类型，除非它和表达式的类型不同。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
var _s string = F()

func F() string { return "A" }
```

</td><td>

```go
var _s = F()
// Since F already states that it returns a string, we don't need to specify
// the type again.

func F() string { return "A" }
```

</td></tr>
</tbody></table>

如果表达式的类型和所需类型不一样，需要指定类型。

```go
type myError struct{}

func (myError) Error() string { return "error" }

func F() myError { return myError{} }

var _e error = F()
// F returns an object of type myError but we want error.
```

### 非导出变量使用_前缀

对于非导出类型的变量，在用`var`s and `const`声明时加上`_`前缀，来表示他们是全局符号。

原因：顶层声明的变量作用域一般是包范围。用一个常见的名字可能会导致在其他包中被意外修改。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// foo.go

const (
  defaultPort = 8080
  defaultUser = "user"
)

// bar.go

func Bar() {
  defaultPort := 9090
  ...
  fmt.Println("Default port", defaultPort)

  // We will not see a compile error if the first line of
  // Bar() is deleted.
}
```

</td><td>

```go
// foo.go

const (
  _defaultPort = 8080
  _defaultUser = "user"
)
```

</td></tr>
</tbody></table>

**异常**：非导出的错误类型一般使用不带 `_` 的 `err` 前缀。参考
[Error Naming](#error-naming).


### 结构体内嵌类型

嵌入类型应该放在结构体的最上面，应该和结构体的常规字段用一个空行分隔开。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Client struct {
  version int
  http.Client
}
```

</td><td>

```go
type Client struct {
  http.Client

  version int
}
```

</td></tr>
</tbody></table>

内嵌类型会带来足够的好处，比如在语义上会增加或增强功能。但应该在对用户没有影响的情况下使用内嵌。
(参考: [避免在公共结构中嵌入类型]).

例外：即使是未导出类型，Mutex 也不应该被内嵌。参考：: [Mutex的零值是有效的].

[避免在公共结构中嵌入类型]: #avoid-embedding-types-in-public-structs
[Mutex的零值是有效的]: #zero-value-mutexes-are-valid

这些情况下**避免内嵌**:

- 单纯为了便利和美观。
- 让外部类型构造起来或使用起来更困难。
- 影响了外部的零值。如果外部类型的零值是有用的，嵌入类型应该也有一个有用的零值。
- 作为嵌入类型的副作用，公开外部类型的不相关函数或字段。
- 公开非导出类型。
- 影响外部类型的复制语义。
- 影响外部类型的API或类型语义。
- 影响内部类型的非规范形式。
- 公开外部类型的详细实现信息。
- 允许用户观察和控制内部类型。
- 通过包装的形式改变了内部函数的行为，这种包装的方式会给用户造成意外观感。

简单概括，使用嵌入类型时要明确目的。一个不错的方式是："这些嵌入的字段/方法是否需要被直接添加到外部
类型"，如果答案是"一些"或者"No"，不要使用内嵌类型，而是使用命名字段。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type A struct {
    // Bad: A.Lock() and A.Unlock() are
    //      now available, provide no
    //      functional benefit, and allow
    //      users to control details about
    //      the internals of A.
    sync.Mutex
}
```

</td><td>

```go
type countingWriteCloser struct {
    // Good: Write() is provided at this
    //       outer layer for a specific
    //       purpose, and delegates work
    //       to the inner type's Write().
    io.WriteCloser

    count int
}

func (w *countingWriteCloser) Write(bs []byte) (int, error) {
    w.count += len(bs)
    return w.WriteCloser.Write(bs)
}
```

</td></tr>
<tr><td>

```go
type Book struct {
    // Bad: pointer changes zero value usefulness
    io.ReadWriter

    // other fields
}

// later

var b Book
b.Read(...)  // panic: nil pointer
b.String()   // panic: nil pointer
b.Write(...) // panic: nil pointer
```

</td><td>

```go
type Book struct {
    // Good: has useful zero value
    bytes.Buffer

    // other fields
}

// later

var b Book
b.Read(...)  // ok
b.String()   // ok
b.Write(...) // ok
```

</td></tr>
<tr><td>

```go
type Client struct {
    sync.Mutex
    sync.WaitGroup
    bytes.Buffer
    url.URL
}
```

</td><td>

```go
type Client struct {
    mtx sync.Mutex
    wg  sync.WaitGroup
    buf bytes.Buffer
    url url.URL
}
```

</td></tr>
</tbody></table>

### 本地变量声明

如果将变量声明为某个值，应该使用短变量命名方式：`:=`。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
var s = "foo"
```

</td><td>

```go
s := "foo"
```

</td></tr>
</tbody></table>

然而，有些情况下用 `var` 会让声明语句更加清晰，比如[声明空slice].

[声明空slice]: https://github.com/golang/go/wiki/CodeReviewComments#declaring-empty-slices

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func f(list []int) {
  filtered := []int{}
  for _, v := range list {
    if v > 10 {
      filtered = append(filtered, v)
    }
  }
}
```

</td><td>

```go
func f(list []int) {
  var filtered []int
  for _, v := range list {
    if v > 10 {
      filtered = append(filtered, v)
    }
  }
}
```

</td></tr>
</tbody></table>

### nil值的slice有效

`nil` 代表长度为0的有效slice, 意味着：

- 你不应该声明一个长度为0的空slice，而是用`nil`来代替。 


  <table>
  <thead><tr><th>Bad</th><th>Good</th></tr></thead>
  <tbody>
  <tr><td>

  ```go
  if x == "" {
    return []int{}
  }
  ```

  </td><td>

  ```go
  if x == "" {
    return nil
  }
  ```

  </td></tr>
  </tbody></table>

- 要检查slice是否为空，不应该检查 `nil`, 而是用长度判断`len(s) == 0`。

  <table>
  <thead><tr><th>Bad</th><th>Good</th></tr></thead>
  <tbody>
  <tr><td>

  ```go
  func isEmpty(s []string) bool {
    return s == nil
  }
  ```

  </td><td>

  ```go
  func isEmpty(s []string) bool {
    return len(s) == 0
  }
  ```

  </td></tr>
  </tbody></table>

- 用 `var` 声明的零值slice是有效的，没必要用 `make` 来创建。

  <table>
  <thead><tr><th>Bad</th><th>Good</th></tr></thead>
  <tbody>
  <tr><td>

  ```go
  nums := []int{}
  // or, nums := make([]int)
  
  if add1 {
    nums = append(nums, 1)
  }
  
  if add2 {
    nums = append(nums, 2)
  }
  ```

  </td><td>

  ```go
  var nums []int
  
  if add1 {
    nums = append(nums, 1)
  }
  
  if add2 {
    nums = append(nums, 2)
  }
  ```

  </td></tr>
  </tbody></table>

另外记住，虽然 nil 的slice有效，但是它不等于长度为0的 slice。在一些情况下(比如说序列化)，
这两种slice的表现不同。

### 缩小变量作用域

尽可能减小变量的作用域。如果与 [减少嵌套](#减少嵌套) 冲突，就不要缩小。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
err := os.WriteFile(name, data, 0644)
if err != nil {
 return err
}
```

</td><td>

```go
if err := os.WriteFile(name, data, 0644); err != nil {
 return err
}
```

</td></tr>
</tbody></table>

但是如果作用域是 if 范围之外，不应该减少作用域。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
if data, err := os.ReadFile(name); err == nil {
  err = cfg.Decode(data)
  if err != nil {
    return err
  }

  fmt.Println(cfg)
  return nil
} else {
  return err
}
```

</td><td>

```go
data, err := os.ReadFile(name)
if err != nil {
   return err
}

if err := cfg.Decode(data); err != nil {
  return err
}

fmt.Println(cfg)
return nil
```

</td></tr>
</tbody></table>

### 不面参数语义不明确

在函数中裸传参数值会让代码语义不明确，可以添加 C 风格(`/* ... */`)的注释。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// func printInfo(name string, isLocal, done bool)

printInfo("foo", true, true)
```

</td><td>

```go
// func printInfo(name string, isLocal, done bool)

printInfo("foo", true /* isLocal */, true /* done */)
```

</td></tr>
</tbody></table>

当然，更好的处理方式将上面的 `bool` 换成自定义类型。因为未来可能不仅仅局限于两个bool值(true/false)。

```go
type Region int

const (
  UnknownRegion Region = iota
  Local
)

type Status int

const (
  StatusReady Status = iota + 1
  StatusDone
  // 获取未来会有 StatusInProgress 枚举
)

func printInfo(name string, region Region, status Status)
```

### 字符串中避免转义

Go中支持 [字符串原始值](https://golang.org/ref/spec#raw_string_lit),当需要
转义时，尽量使用 "`" 来包装字符串。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
wantError := "unknown name:\"test\""
```

</td><td>

```go
wantError := `unknown error:"test"`
```

</td></tr>
</tbody></table>

### 初始化结构体

#### 初始化结构体时声明字段名

在你初始化结构时，几乎应该始终指定字段名。目前由[`go vet`]强制执行。

[`go vet`]: https://golang.org/cmd/vet/

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
k := User{"John", "Doe", true}
```

</td><td>

```go
k := User{
    FirstName: "John",
    LastName: "Doe",
    Admin: true,
}
```

</td></tr>
</tbody></table>

例外：当有 3 个或更少的字段时，测试表中的字段名也许可以省略。

```go
tests := []struct{
  op Operation
  want string
}{
  {Add, "add"},
  {Subtract, "subtract"},
}
```

#### 省略结构体中的零值字段

当初始化结构体字段时，除非需要提供一个有意义的上下文，否则需要忽略对零值字段进行赋值。因为Go
会自动给这些零值字段进行填充。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
user := User{
  FirstName: "John",
  LastName: "Doe",
  MiddleName: "",
  Admin: false,
}
```

</td><td>

```go
user := User{
  FirstName: "John",
  LastName: "Doe",
}
```

</td></tr>
</tbody></table>

这种行为让我们忽略了上下文无关的噪音信息。只关注有意义的特殊值。

当零值代表有意义的上下文时需要提供零值。比如在  [表驱动测试](#表驱动测试) 中零值字段
是有意义的。


```go
tests := []struct{
  give string
  want int
}{
  {give: "0", want: 0},
  // ...
}
```

#### 空结构体用var声明

当结构体中所有的字段都为空时，用 `var` 来声明结构体。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
user := User{}
```

</td><td>

```go
var user User
```

</td></tr>
</tbody></table>

这种 零值结构体 和具有非零值字段的结构体有所不同，和 [map初始化] 更相似，
和我们更想用的 [声明空Slices][声明空Slices] 更匹配。

[map初始化]: #map初始化

#### 初始化结构体引用

初始化结构引用时，请使用`&T{}`代替`new(T)`，以使其与结构体初始化一致。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
sval := T{Name: "foo"}

// 非一致的
sptr := new(T)
sptr.Name = "bar"
```

</td><td>

```go
sval := T{Name: "foo"}

sptr := &T{Name: "bar"}
```

</td></tr>
</tbody></table>

### 初始化Map

优先使用make来创建空map，这样使得map的初始化不同于声明，而且你还可以在 make 中添加map的大小提示。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
var (
  // m1 的读写操作都是安全的
  // m2 的写操作会panic
  m1 = map[T1]T2{}
  m2 map[T1]T2
)
```

</td><td>

```go
var (
  // m1 的读写操作都是安全的
  // m2 的写操作会panic
  m1 = make(map[T1]T2)
  m2 map[T1]T2
)
```

</td></tr>
<tr><td>

声明和初始化在形式上相似

</td><td>

声明和初始化在形式上隔离


</td></tr>
</tbody></table>

尽可能在 `make()` 中制定map的初始化容量，可以参考：[Specifying Map Capacity Hints](#specifying-map-capacity-hints)。

另外，如果map初始化的时候需要赋值固定信息，使用 map literals 方式来初始化。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
m := make(map[T1]T2, 3)
m[k1] = v1
m[k2] = v2
m[k3] = v3
```

</td><td>

```go
m := map[T1]T2{
  k1: v1,
  k2: v2,
  k3: v3,
}
```

</td></tr>
</tbody></table>


原则上：在初始化map时增加一组固定的元素，就使用map literals。否则就使用 `make`(如果可以，
尽可能指定map的容量)。


### 在Printf外面格式化字符串

如果你在函数外声明 `Printf` 风格 函数的格式字符串，请将其设置为 `const` 常量。

这有助于go vet对格式字符串执行静态分析。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
msg := "unexpected values %v, %v\n"
fmt.Printf(msg, 1, 2)
```

</td><td>

```go
const msg = "unexpected values %v, %v\n"
fmt.Printf(msg, 1, 2)
```

</td></tr>
</tbody></table>

### 命名Printf样式函数

使用`Printf`函数时，应保证`go vet`可以检测到他的格式化字符串。


这意味着你需要使用预定义的`Printf`函数名称，`go vet`会默认检查这些。更多信息，请参考：[Printf family]

[Printf family]: https://golang.org/cmd/vet/#hdr-Printf_family

如果不能使用预定义的名称，请以 f 结束选择的名称：`Wrapf`，而不是`Wrap`。`go vet`可以要求检查特定的 `Printf`
样式名称，但名称必须以`f`结尾。

```shell
$ go vet -printfuncs=wrapf,statusf
```

参考 [go vet: Printf family check].

[go vet: Printf family check]: https://kuzminva.wordpress.com/2017/11/07/go-vet-printf-family-check/

## Patterns

### Test Tables

当你的测试用例形式上重复时，用 [subtests] 方式编写case会让测试用例看起来更加简洁。

[subtests]: https://blog.golang.org/subtests

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// func TestSplitHostPort(t *testing.T)

host, port, err := net.SplitHostPort("192.0.2.0:8000")
require.NoError(t, err)
assert.Equal(t, "192.0.2.0", host)
assert.Equal(t, "8000", port)

host, port, err = net.SplitHostPort("192.0.2.0:http")
require.NoError(t, err)
assert.Equal(t, "192.0.2.0", host)
assert.Equal(t, "http", port)

host, port, err = net.SplitHostPort(":8000")
require.NoError(t, err)
assert.Equal(t, "", host)
assert.Equal(t, "8000", port)

host, port, err = net.SplitHostPort("1:8")
require.NoError(t, err)
assert.Equal(t, "1", host)
assert.Equal(t, "8", port)
```

</td><td>

```go
// func TestSplitHostPort(t *testing.T)

tests := []struct{
  give     string
  wantHost string
  wantPort string
}{
  {
    give:     "192.0.2.0:8000",
    wantHost: "192.0.2.0",
    wantPort: "8000",
  },
  {
    give:     "192.0.2.0:http",
    wantHost: "192.0.2.0",
    wantPort: "http",
  },
  {
    give:     ":8000",
    wantHost: "",
    wantPort: "8000",
  },
  {
    give:     "1:8",
    wantHost: "1",
    wantPort: "8",
  },
}

for _, tt := range tests {
  t.Run(tt.give, func(t *testing.T) {
    host, port, err := net.SplitHostPort(tt.give)
    require.NoError(t, err)
    assert.Equal(t, tt.wantHost, host)
    assert.Equal(t, tt.wantPort, port)
  })
}
```

</td></tr>
</tbody></table>

显然，如果你用了test table的方式，在拓展测试用例时也会显得更加清晰。

我们遵守这样的准则：搞一个slice类型的struct测试用例，每个测试case叫做`tt`。然后使用`give`和
`want`说明测试用例的输入和输出。


```go
tests := []struct{
  give     string
  wantHost string
  wantPort string
}{
  // ...
}

for _, tt := range tests {
  // ...
}
```

对于并行测试，比如一些特殊的循环 (比如那些生产 goroutine 或 在循环中捕获引用的循环), 需要
注意在循环中明确分配循环变量来确保不会产生闭包。

```go
tests := []struct{
  give string
  // ...
}{
  // ...
}

for _, tt := range tests {
  tt := tt // for t.Parallel
  t.Run(tt.give, func(t *testing.T) {
    t.Parallel()
    // ...
  })
}
```
在上面的例子中，由于循环中使用了 `t.Parallel()`，我们必须在外部循环中声明一个 `tt` 变量。
如果不这么做，大多数测试用例都会收到一个非预期的 `tt`，或是一个在运行期改变的值。

### 函数功能选项API

功能选项是一种模式，你可以声明一个对用户不透明的 `Option` 类型，在一些内部结构中记录信息。
函数接收不定长的参数选项，并根据参数做不同的行为。

对于需要拓展参数的构造方法或是其他需要可选参数的公共API可以考虑这种模式，对于参数在三个及以上
的函数更应该考虑。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// package db

func Open(
  addr string,
  cache bool,
  logger *zap.Logger
) (*Connection, error) {
  // ...
}
```

</td><td>

```go
// package db

type Option interface {
  // ...
}

func WithCache(c bool) Option {
  // ...
}

func WithLogger(log *zap.Logger) Option {
  // ...
}

// Open 创建一个连接
func Open(
  addr string,
  opts ...Option,
) (*Connection, error) {
  // ...
}
```

</td></tr>
<tr><td>

即使用户默认不需要 cache 和 logger，也需要提供这俩参数。

```go
db.Open(addr, db.DefaultCache, zap.NewNop())
db.Open(addr, db.DefaultCache, log)
db.Open(addr, false /* cache */, zap.NewNop())
db.Open(addr, false /* cache */, log)
```

</td><td>

Options 只在需要时才被提供。


```go
db.Open(addr)
db.Open(addr, db.WithLogger(log))
db.Open(addr, db.WithCache(false))
db.Open(
  addr,
  db.WithCache(false),
  db.WithLogger(log),
)
```

</td></tr>
</tbody></table>

我们建议这种模式的实现方式是 提供一个 `Option` 接口，里面有一个非导出类型方法，在一个非
导出类型的 `options` 结构体中记录选项。

```go
type options struct {
  cache  bool
  logger *zap.Logger
}

type Option interface {
  apply(*options)
}

type cacheOption bool

func (c cacheOption) apply(opts *options) {
  opts.cache = bool(c)
}

func WithCache(c bool) Option {
  return cacheOption(c)
}

type loggerOption struct {
  Log *zap.Logger
}

func (l loggerOption) apply(opts *options) {
  opts.logger = l.Log
}

func WithLogger(log *zap.Logger) Option {
  return loggerOption{Log: log}
}

// Open 创建一个连接
func Open(
  addr string,
  opts ...Option,
) (*Connection, error) {
  options := options{
    cache:  defaultCache,
    logger: zap.NewNop(),
  }

  for _, o := range opts {
    o.apply(&options)
  }

  // ...
}
```

还有一种用闭包实现这种方法的模式，但我们认为上面提供的这种模式给作者提供了更高的灵活性，更
易于调试和测试。这种方式可以在测试和mock中进行比较，而闭包方式难以做到。此外，它允许 option
实现其他接口，比如 `fmt.Stringer`，会 string 类型的可读性更高。

还可以参考：

- [Self-referential functions and the design of options]
- [Functional options for friendly APIs]

  [Self-referential functions and the design of options]: https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
  [Functional options for friendly APIs]: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

<!-- TODO: replace this with parameter structs and functional options, when to
use one vs other -->

## Linting

比其他任何 "神圣" linter 工具更重要的是，在你的代码库里使用一致性的 lint 工具。

我们建议最少要使用下面这些 linters 工具吗，因为我们认为这些工具可以帮你捕获最常见的问题，有助于
在没有规定的前提下提高代码质量：

- [errcheck] 确保错误被处理
- [goimports] 格式化代码和管理包引用
- [golint] 指出常见的文本错误
- [govet] 分析代码中的常见错误
- [staticcheck] 各种静态分析检查

  [errcheck]: https://github.com/kisielk/errcheck
  [goimports]: https://godoc.org/golang.org/x/tools/cmd/goimports
  [golint]: https://github.com/golang/lint
  [govet]: https://golang.org/cmd/vet/
  [staticcheck]: https://staticcheck.io/


### Lint Runners

由于优秀的性能表现，我们推荐 [golangci-lint] 作为Go代码的首选 lint 工具。这个仓库有在一个
[.golangci.yml] 例子，里面有配置的 linters 工具和设置。


golangci-lint 有一系列 [various linters] 可供使用。建议将这些 linters 作为基础集合，
我们鼓励团队内部将其他有意义的 linters 工具在他们的项目中进行使用。

[golangci-lint]: https://github.com/golangci/golangci-lint
[.golangci.yml]: https://github.com/uber-go/guide/blob/master/.golangci.yml
[various linters]: https://golangci-lint.run/usage/linters/
