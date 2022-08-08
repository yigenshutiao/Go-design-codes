


# Uber Go Style Guide

## Table of Contents

- [介绍](#介绍)
- [指南](#指南)
    - [指针转换为interface](#指针转换为interface)
    - [验证接口合法性](#验证接口合法性)
    - [接收者和接口](#接收者和接口)
    - [Mutexes的零值是有效的](#Mutexes的零值是有效的)
    - [在边界拷贝Slices和Maps](#在边界拷贝Slices和Maps)
    - [使用Defer释放资源](#使用Defer释放资源)
    - [Channel大小应为0或1](#Channel大小应为0或1)
    - [Start Enums at One](#start-enums-at-one)
    - [使用time包来处理时间](#使用time包来处理时间)
    - [错误](#错误)
        - [错误类型](#错误类型)
        - [错误包装](#错误包装)
        - [错误命名](#错误命名)
    - [处理断言失败](#处理断言失败)
    - [不要使用Panic](#不要使用Panic)
    - [Use go.uber.org/atomic](#use-gouberorgatomic)
    - [Avoid Mutable Globals](#avoid-mutable-globals)
    - [Avoid Embedding Types in Public Structs](#avoid-embedding-types-in-public-structs)
    - [Avoid Using Built-In Names](#avoid-using-built-in-names)
    - [Avoid `init()`](#avoid-init)
    - [Exit in Main](#exit-in-main)
        - [Exit Once](#exit-once)
    - [Use field tags in marshaled structs](#use-field-tags-in-marshaled-structs)
- [Performance](#performance)
    - [Prefer strconv over fmt](#prefer-strconv-over-fmt)
    - [Avoid string-to-byte conversion](#avoid-string-to-byte-conversion)
    - [Prefer Specifying Container Capacity](#prefer-specifying-container-capacity)
        - [Specifying Map Capacity Hints](#specifying-map-capacity-hints)
        - [Specifying Slice Capacity](#specifying-slice-capacity)
- [Style](#style)
    - [Avoid overly long lines](#avoid-overly-long-lines)
    - [Be Consistent](#be-consistent)
    - [Group Similar Declarations](#group-similar-declarations)
    - [Import Group Ordering](#import-group-ordering)
    - [Package Names](#package-names)
    - [Function Names](#function-names)
    - [Import Aliasing](#import-aliasing)
    - [Function Grouping and Ordering](#function-grouping-and-ordering)
    - [Reduce Nesting](#reduce-nesting)
    - [Unnecessary Else](#unnecessary-else)
    - [Top-level Variable Declarations](#top-level-variable-declarations)
    - [Prefix Unexported Globals with _](#prefix-unexported-globals-with-_)
    - [Embedding in Structs](#embedding-in-structs)
    - [Local Variable Declarations](#local-variable-declarations)
    - [nil is a valid slice](#nil-is-a-valid-slice)
    - [Reduce Scope of Variables](#reduce-scope-of-variables)
    - [Avoid Naked Parameters](#avoid-naked-parameters)
    - [Use Raw String Literals to Avoid Escaping](#use-raw-string-literals-to-avoid-escaping)
    - [Initializing Structs](#initializing-structs)
        - [Use Field Names to Initialize Structs](#use-field-names-to-initialize-structs)
        - [Omit Zero Value Fields in Structs](#omit-zero-value-fields-in-structs)
        - [Use `var` for Zero Value Structs](#use-var-for-zero-value-structs)
        - [Initializing Struct References](#initializing-struct-references)
    - [Initializing Maps](#initializing-maps)
    - [Format Strings outside Printf](#format-strings-outside-printf)
    - [Naming Printf-style Functions](#naming-printf-style-functions)
- [Patterns](#patterns)
    - [Test Tables](#test-tables)
    - [Functional Options](#functional-options)
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

- Run goimports on save

- Run golint and go vet to check for errors

你可以在这里找到编辑器支持Go工具的信息：
<https://github.com/golang/go/wiki/IDEsAndTextEditorPlugins>






## 指南

### 指针转换为interface



You almost never need a pointer to an interface. You should be passing
interfaces as values—the underlying data can still be a pointer.

An interface is two fields:

1. A pointer to some type-specific information. You can think of this as
   "type."
2. Data pointer. If the data stored is a pointer, it’s stored directly. If
   the data stored is a value, then a pointer to the value is stored.

If you want interface methods to modify the underlying data, you must use a
pointer.



### 验证接口合法性

在编译期验证接口的合法性，需要验证的有：

- 验证导出类型在作为API时是否实现了特定接口
- Exported or unexported types that are part of a collection of types
  implementing the same interface
- Other cases where violating an interface would break users

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

These embedded types leak implementation details, inhibit type evolution, and
obscure documentation.

Assuming you have implemented a variety of list types using a shared
`AbstractList`, avoid embedding the `AbstractList` in your concrete list
implementations.
Instead, hand-write only the methods to your concrete list that will delegate
to the abstract list.

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

Go allows [type embedding] as a compromise between inheritance and composition.
The outer type gets implicit copies of the embedded type's methods.
These methods, by default, delegate to the same method of the embedded
instance.

[type embedding]: https://golang.org/doc/effective_go.html#embedding

The struct also gains a field by the same name as the type.
So, if the embedded type is public, the field is public.
To maintain backward compatibility, every future version of the outer type must
keep the embedded type.

An embedded type is rarely necessary.
It is a convenience that helps you avoid writing tedious delegate methods.

Even embedding a compatible AbstractList *interface*, instead of the struct,
would offer the developer more flexibility to change in the future, but still
leak the detail that the concrete lists use an abstract implementation.

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

Either with an embedded struct or an embedded interface, the embedded type
places limits on the evolution of the type.

- Adding methods to an embedded interface is a breaking change.
- Removing methods from an embedded struct is a breaking change.
- Removing the embedded type is a breaking change.
- Replacing the embedded type, even with an alternative that satisfies the same
  interface, is a breaking change.

Although writing these delegate methods is tedious, the additional effort hides
an implementation detail, leaves more opportunities for change, and also
eliminates indirection for discovering the full List interface in
documentation.

### Avoid Using Built-In Names

The Go [language specification] outlines several built-in,
[predeclared identifiers] that should not be used as names within Go programs.

Depending on context, reusing these identifiers as names will either shadow
the original within the current lexical scope (and any nested scopes) or make
affected code confusing. In the best case, the compiler will complain; in the
worst case, such code may introduce latent, hard-to-grep bugs.

[language specification]: https://golang.org/ref/spec
[predeclared identifiers]: https://golang.org/ref/spec#Predeclared_identifiers

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
var error string
// `error` shadows the builtin

// or

func handleErrorMessage(error string) {
    // `error` shadows the builtin
}
```

</td><td>

```go
var errorMessage string
// `error` refers to the builtin

// or

func handleErrorMessage(msg string) {
    // `error` refers to the builtin
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


Note that the compiler will not generate errors when using predeclared
identifiers, but tools such as `go vet` should correctly point out these and
other cases of shadowing.

### 避免使用 `init()`

Avoid `init()` where possible. When `init()` is unavoidable or desirable, code
should attempt to:

1. Be completely deterministic, regardless of program environment or invocation.
2. Avoid depending on the ordering or side-effects of other `init()` functions.
   While `init()` ordering is well-known, code can change, and thus
   relationships between `init()` functions can make code brittle and
   error-prone.
3. Avoid accessing or manipulating global or environment state, such as machine
   information, environment variables, working directory, program
   arguments/inputs, etc.
4. Avoid I/O, including both filesystem, network, and system calls.

Code that cannot satisfy these requirements likely belongs as a helper to be
called as part of `main()` (or elsewhere in a program's lifecycle), or be
written as part of `main()` itself. In particular, libraries that are intended
to be used by other programs should take special care to be completely
deterministic and not perform "init magic".

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

Considering the above, some situations in which `init()` may be preferable or
necessary might include:

- Complex expressions that cannot be represented as single assignments.
- Pluggable hooks, such as `database/sql` dialects, encoding type registries, etc.
- Optimizations to [Google Cloud Functions] and other forms of deterministic
  precomputation.

  [Google Cloud Functions]: https://cloud.google.com/functions/docs/bestpractices/tips#use_global_variables_to_reuse_objects_in_future_invocations

### Exit in Main

Go programs use [`os.Exit`] or [`log.Fatal*`] to exit immediately. (Panicking
is not a good way to exit programs, please [don't panic](#dont-panic).)

[`os.Exit`]: https://golang.org/pkg/os/#Exit
[`log.Fatal*`]: https://golang.org/pkg/log/#Fatal

Call one of `os.Exit` or `log.Fatal*` **only in `main()`**. All other
functions should return errors to signal failure.

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

Rationale: Programs with multiple functions that exit present a few issues:

- Non-obvious control flow: Any function can exit the program so it becomes
  difficult to reason about the control flow.
- Difficult to test: A function that exits the program will also exit the test
  calling it. This makes the function difficult to test and introduces risk of
  skipping other tests that have not yet been run by `go test`.
- Skipped cleanup: When a function exits the program, it skips function calls
  enqueued with `defer` statements. This adds risk of skipping important
  cleanup tasks.

#### Exit Once

If possible, prefer to call `os.Exit` or `log.Fatal` **at most once** in your
`main()`. If there are multiple error scenarios that halt program execution,
put that logic under a separate function and return errors from it.

This has the effect of shortening your `main()` function and putting all key
business logic into a separate, testable function.

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

### Use field tags in marshaled structs

Any struct field that is marshaled into JSON, YAML,
or other formats that support tag-based field naming
should be annotated with the relevant tag.

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
The serialized form of the structure is a contract between different systems.
Changes to the structure of the serialized form--including field names--break
this contract. Specifying field names inside tags makes the contract explicit,
and it guards against accidentally breaking the contract by refactoring or
renaming fields.

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

Specify container capacity where possible in order to allocate memory for the
container up front. This minimizes subsequent allocations (by copying and
resizing of the container) as elements are added.

#### Specifying Map Capacity Hints

Where possible, provide capacity hints when initializing
maps with `make()`.

```go
make(map[T1]T2, hint)
```

Providing a capacity hint to `make()` tries to right-size the
map at initialization time, which reduces the need for growing
the map and allocations as elements are added to the map.

Note that, unlike slices, map capacity hints do not guarantee complete,
preemptive allocation, but are used to approximate the number of hashmap buckets
required. Consequently, allocations may still occur when adding elements to the
map, even up to the specified capacity.

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

`m` is created without a size hint; there may be more
allocations at assignment time.

</td><td>

`m` is created with a size hint; there may be fewer
allocations at assignment time.

</td></tr>
</tbody></table>

#### Specifying Slice Capacity

Where possible, provide capacity hints when initializing slices with `make()`,
particularly when appending.

```go
make([]T, length, capacity)
```

Unlike maps, slice capacity is not a hint: the compiler will allocate enough
memory for the capacity of the slice as provided to `make()`, which means that
subsequent `append()` operations will incur zero allocations (until the length
of the slice matches the capacity, after which any appends will require a resize
to hold additional elements).

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

## Style

### Avoid overly long lines

Avoid lines of code that require readers to scroll horizontally
or turn their heads too much.

We recommend a soft line length limit of **99 characters**.
Authors should aim to wrap lines before hitting this limit,
but it is not a hard limit.
Code is allowed to exceed this limit.

### Be Consistent

Some of the guidelines outlined in this document can be evaluated objectively;
others are situational, contextual, or subjective.

Above all else, **be consistent**.

Consistent code is easier to maintain, is easier to rationalize, requires less
cognitive overhead, and is easier to migrate or update as new conventions emerge
or classes of bugs are fixed.

Conversely, having multiple disparate or conflicting styles within a single
codebase causes maintenance overhead, uncertainty, and cognitive dissonance,
all of which can directly contribute to lower velocity, painful code reviews,
and bugs.

When applying these guidelines to a codebase, it is recommended that changes
are made at a package (or larger) level: application at a sub-package level
violates the above concern by introducing multiple styles into the same code.

### Group Similar Declarations

Go supports grouping similar declarations.

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

This also applies to constants, variables, and type declarations.

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

Only group related declarations. Do not group declarations that are unrelated.

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

Groups are not limited in where they can be used. For example, you can use them
inside of functions.

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

Exception: Variable declarations, particularly inside functions, should be
grouped together if declared adjacent to other variables. Do this for variables
declared together even if they are unrelated.

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

### Import Group Ordering

There should be two import groups:

- Standard library
- Everything else

This is the grouping applied by goimports by default.

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

### Package Names

When naming packages, choose a name that is:

- All lower-case. No capitals or underscores.
- Does not need to be renamed using named imports at most call sites.
- Short and succinct. Remember that the name is identified in full at every call
  site.
- Not plural. For example, `net/url`, not `net/urls`.
- Not "common", "util", "shared", or "lib". These are bad, uninformative names.

See also [Package Names] and [Style guideline for Go packages].

[Package Names]: https://blog.golang.org/package-names
[Style guideline for Go packages]: https://rakyll.org/style-packages/

### Function Names

We follow the Go community's convention of using [MixedCaps for function
names]. An exception is made for test functions, which may contain underscores
for the purpose of grouping related test cases, e.g.,
`TestMyFunction_WhatIsBeingTested`.

[MixedCaps for function names]: https://golang.org/doc/effective_go.html#mixed-caps

### Import Aliasing

Import aliasing must be used if the package name does not match the last
element of the import path.

```go
import (
  "net/http"

  client "example.com/client-go"
  trace "example.com/trace/v2"
)
```

In all other scenarios, import aliases should be avoided unless there is a
direct conflict between imports.

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

### Function Grouping and Ordering

- Functions should be sorted in rough call order.
- Functions in a file should be grouped by receiver.

Therefore, exported functions should appear first in a file, after
`struct`, `const`, `var` definitions.

A `newXYZ()`/`NewXYZ()` may appear after the type is defined, but before the
rest of the methods on the receiver.

Since functions are grouped by receiver, plain utility functions should appear
towards the end of the file.

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

### Reduce Nesting

Code should reduce nesting where possible by handling error cases/special
conditions first and returning early or continuing the loop. Reduce the amount
of code that is nested multiple levels.

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

### Unnecessary Else

If a variable is set in both branches of an if, it can be replaced with a
single if.

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

### Top-level Variable Declarations

At the top level, use the standard `var` keyword. Do not specify the type,
unless it is not the same type as the expression.

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

Specify the type if the type of the expression does not match the desired type
exactly.

```go
type myError struct{}

func (myError) Error() string { return "error" }

func F() myError { return myError{} }

var _e error = F()
// F returns an object of type myError but we want error.
```

### Prefix Unexported Globals with _

Prefix unexported top-level `var`s and `const`s with `_` to make it clear when
they are used that they are global symbols.

Rationale: Top-level variables and constants have a package scope. Using a
generic name makes it easy to accidentally use the wrong value in a different
file.

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

**Exception**: Unexported error values may use the prefix `err` without the underscore.
See [Error Naming](#error-naming).

### Embedding in Structs

Embedded types should be at the top of the field list of a
struct, and there must be an empty line separating embedded fields from regular
fields.

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

Embedding should provide tangible benefit, like adding or augmenting
functionality in a semantically-appropriate way. It should do this with zero
adverse user-facing effects (see also: [Avoid Embedding Types in Public Structs]).

Exception: Mutexes should not be embedded, even on unexported types. See also: [Zero-value Mutexes are Valid].

[Avoid Embedding Types in Public Structs]: #avoid-embedding-types-in-public-structs
[Zero-value Mutexes are Valid]: #zero-value-mutexes-are-valid

Embedding **should not**:

- Be purely cosmetic or convenience-oriented.
- Make outer types more difficult to construct or use.
- Affect outer types' zero values. If the outer type has a useful zero value, it
  should still have a useful zero value after embedding the inner type.
- Expose unrelated functions or fields from the outer type as a side-effect of
  embedding the inner type.
- Expose unexported types.
- Affect outer types' copy semantics.
- Change the outer type's API or type semantics.
- Embed a non-canonical form of the inner type.
- Expose implementation details of the outer type.
- Allow users to observe or control type internals.
- Change the general behavior of inner functions through wrapping in a way that
  would reasonably surprise users.

Simply put, embed consciously and intentionally. A good litmus test is, "would
all of these exported inner methods/fields be added directly to the outer type";
if the answer is "some" or "no", don't embed the inner type - use a field
instead.

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

### Local Variable Declarations

Short variable declarations (`:=`) should be used if a variable is being set to
some value explicitly.

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

However, there are cases where the default value is clearer when the `var`
keyword is used. [Declaring Empty Slices], for example.

[Declaring Empty Slices]: https://github.com/golang/go/wiki/CodeReviewComments#declaring-empty-slices

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

### nil is a valid slice

`nil` is a valid slice of length 0. This means that,

- You should not return a slice of length zero explicitly. Return `nil`
  instead.

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

- To check if a slice is empty, always use `len(s) == 0`. Do not check for
  `nil`.

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

- The zero value (a slice declared with `var`) is usable immediately without
  `make()`.

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

Remember that, while it is a valid slice, a nil slice is not equivalent to an
allocated slice of length 0 - one is nil and the other is not - and the two may
be treated differently in different situations (such as serialization).

### Reduce Scope of Variables

Where possible, reduce scope of variables. Do not reduce the scope if it
conflicts with [Reduce Nesting](#reduce-nesting).

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

If you need a result of a function call outside of the if, then you should not
try to reduce the scope.

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

### Avoid Naked Parameters

Naked parameters in function calls can hurt readability. Add C-style comments
(`/* ... */`) for parameter names when their meaning is not obvious.

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

Better yet, replace naked `bool` types with custom types for more readable and
type-safe code. This allows more than just two states (true/false) for that
parameter in the future.

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
  // Maybe we will have a StatusInProgress in the future.
)

func printInfo(name string, region Region, status Status)
```

### Use Raw String Literals to Avoid Escaping

Go supports [raw string literals](https://golang.org/ref/spec#raw_string_lit),
which can span multiple lines and include quotes. Use these to avoid
hand-escaped strings which are much harder to read.

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

### Initializing Structs

#### Use Field Names to Initialize Structs

You should almost always specify field names when initializing structs. This is
now enforced by [`go vet`].

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

Exception: Field names *may* be omitted in test tables when there are 3 or
fewer fields.

```go
tests := []struct{
  op Operation
  want string
}{
  {Add, "add"},
  {Subtract, "subtract"},
}
```

#### Omit Zero Value Fields in Structs

When initializing structs with field names, omit fields that have zero values
unless they provide meaningful context. Otherwise, let Go set these to zero
values automatically.

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

This helps reduce noise for readers by omitting values that are default in
that context. Only meaningful values are specified.

Include zero values where field names provide meaningful context. For example,
test cases in [Test Tables](#test-tables) can benefit from names of fields
even when they are zero-valued.

```go
tests := []struct{
  give string
  want int
}{
  {give: "0", want: 0},
  // ...
}
```

#### Use `var` for Zero Value Structs

When all the fields of a struct are omitted in a declaration, use the `var`
form to declare the struct.

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

This differentiates zero valued structs from those with non-zero fields
similar to the distinction created for [map initialization], and matches how
we prefer to [declare empty slices][Declaring Empty Slices].

[map initialization]: #initializing-maps

#### Initializing Struct References

Use `&T{}` instead of `new(T)` when initializing struct references so that it
is consistent with the struct initialization.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
sval := T{Name: "foo"}

// inconsistent
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

### Initializing Maps

Prefer `make(..)` for empty maps, and maps populated
programmatically. This makes map initialization visually
distinct from declaration, and it makes it easy to add size
hints later if available.

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
var (
  // m1 is safe to read and write;
  // m2 will panic on writes.
  m1 = map[T1]T2{}
  m2 map[T1]T2
)
```

</td><td>

```go
var (
  // m1 is safe to read and write;
  // m2 will panic on writes.
  m1 = make(map[T1]T2)
  m2 map[T1]T2
)
```

</td></tr>
<tr><td>

Declaration and initialization are visually similar.

</td><td>

Declaration and initialization are visually distinct.

</td></tr>
</tbody></table>

Where possible, provide capacity hints when initializing
maps with `make()`. See
[Specifying Map Capacity Hints](#specifying-map-capacity-hints)
for more information.

On the other hand, if the map holds a fixed list of elements,
use map literals to initialize the map.

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


The basic rule of thumb is to use map literals when adding a fixed set of
elements at initialization time, otherwise use `make` (and specify a size hint
if available).

### Format Strings outside Printf

If you declare format strings for `Printf`-style functions outside a string
literal, make them `const` values.

This helps `go vet` perform static analysis of the format string.

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

### Naming Printf-style Functions

When you declare a `Printf`-style function, make sure that `go vet` can detect
it and check the format string.

This means that you should use predefined `Printf`-style function
names if possible. `go vet` will check these by default. See [Printf family]
for more information.

[Printf family]: https://golang.org/cmd/vet/#hdr-Printf_family

If using the predefined names is not an option, end the name you choose with
f: `Wrapf`, not `Wrap`. `go vet` can be asked to check specific `Printf`-style
names but they must end with f.

```shell
$ go vet -printfuncs=wrapf,statusf
```

See also [go vet: Printf family check].

[go vet: Printf family check]: https://kuzminva.wordpress.com/2017/11/07/go-vet-printf-family-check/

## Patterns

### Test Tables

Use table-driven tests with [subtests] to avoid duplicating code when the core
test logic is repetitive.

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

Test tables make it easier to add context to error messages, reduce duplicate
logic, and add new test cases.

We follow the convention that the slice of structs is referred to as `tests`
and each test case `tt`. Further, we encourage explicating the input and output
values for each test case with `give` and `want` prefixes.

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

Parallel tests, like some specialized loops (for example, those that spawn
goroutines or capture references as part of the loop body),
must take care to explicitly assign loop variables within the loop's scope to
ensure that they hold the expected values.

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

In the example above, we must declare a `tt` variable scoped to the loop
iteration because of the use of `t.Parallel()` below.
If we do not do that, most or all tests will receive an unexpected value for
`tt`, or a value that changes as they're running.

### Functional Options

Functional options is a pattern in which you declare an opaque `Option` type
that records information in some internal struct. You accept a variadic number
of these options and act upon the full information recorded by the options on
the internal struct.

Use this pattern for optional arguments in constructors and other public APIs
that you foresee needing to expand, especially if you already have three or
more arguments on those functions.

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

// Open creates a connection.
func Open(
  addr string,
  opts ...Option,
) (*Connection, error) {
  // ...
}
```

</td></tr>
<tr><td>

The cache and logger parameters must always be provided, even if the user
wants to use the default.

```go
db.Open(addr, db.DefaultCache, zap.NewNop())
db.Open(addr, db.DefaultCache, log)
db.Open(addr, false /* cache */, zap.NewNop())
db.Open(addr, false /* cache */, log)
```

</td><td>

Options are provided only if needed.

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

Our suggested way of implementing this pattern is with an `Option` interface
that holds an unexported method, recording options on an unexported `options`
struct.

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

// Open creates a connection.
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

Note that there's a method of implementing this pattern with closures but we
believe that the pattern above provides more flexibility for authors and is
easier to debug and test for users. In particular, it allows options to be
compared against each other in tests and mocks, versus closures where this is
impossible. Further, it lets options implement other interfaces, including
`fmt.Stringer` which allows for user-readable string representations of the
options.

See also,

- [Self-referential functions and the design of options]
- [Functional options for friendly APIs]

  [Self-referential functions and the design of options]: https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
  [Functional options for friendly APIs]: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis

<!-- TODO: replace this with parameter structs and functional options, when to
use one vs other -->

## Linting

More importantly than any "blessed" set of linters, lint consistently across a
codebase.

We recommend using the following linters at a minimum, because we feel that they
help to catch the most common issues and also establish a high bar for code
quality without being unnecessarily prescriptive:

- [errcheck] to ensure that errors are handled
- [goimports] to format code and manage imports
- [golint] to point out common style mistakes
- [govet] to analyze code for common mistakes
- [staticcheck] to do various static analysis checks

  [errcheck]: https://github.com/kisielk/errcheck
  [goimports]: https://godoc.org/golang.org/x/tools/cmd/goimports
  [golint]: https://github.com/golang/lint
  [govet]: https://golang.org/cmd/vet/
  [staticcheck]: https://staticcheck.io/


### Lint Runners

We recommend [golangci-lint] as the go-to lint runner for Go code, largely due
to its performance in larger codebases and ability to configure and use many
canonical linters at once. This repo has an example [.golangci.yml] config file
with recommended linters and settings.

golangci-lint has [various linters] available for use. The above linters are
recommended as a base set, and we encourage teams to add any additional linters
that make sense for their projects.

[golangci-lint]: https://github.com/golangci/golangci-lint
[.golangci.yml]: https://github.com/uber-go/guide/blob/master/.golangci.yml
[various linters]: https://golangci-lint.run/usage/linters/
