### mutualdep

Handling mutual dependencies in Golang

相信不少 `Gopher` 在写 `Golang` 程序都遇到过 `import cycle not allowed` 问题，本人最近研读 [go-ethereum](https://github.com/ethereum/go-ethereum) 源码时，发现定义 `interface` 也能解决此问题， 还能解决连分包都不能解决的情况， 并且比分包更加简单快捷。下面逐个讲解 `分包` 和 `定义接口` 这两种方法。

# 1. 应用场景

假设有如下使用场景：

1. `A` 是应用程序的框架级结构体，在 `A` 包含子模块 `B` 和 `C` 的指针；
2. `B` 为了方便的使用应用的其他子模块（比如 `C` ）功能，所以在其结构体包含了 `A` 的指针；
3. `C` 要调用 `A` 包中的某个方法；

# 2. 代码实现

其程序大致如下：

* `package a` 代码如下：

```go
package a

import (
	"fmt"

	"github.com/ggq89/mutualdep/b"
	"github.com/ggq89/mutualdep/c"
)

type A struct {
	Pb *b.B
	Pc *c.C
}

func New(ic int) *A {
	a := &A{
		Pc: c.New(ic),
	}

	a.Pb = b.New(a)

	return a
}

func Printf(v int) {
	fmt.Printf("%v", v)
}
```

* `package b` 代码如下：

```go
package b

import (
	"github.com/ggq89/mutualdep/a"
)

type B struct {
	Pa *a.A
}

func New(a *a.A) *B {
	return &B{
		Pa: a,
	}
}

func (b *B) DisplayC() {
	b.Pa.Pc.Show()
}
```

* `package c` 代码如下：

```go
package c

import "github.com/ggq89/mutualdep/a"

type C struct {
	Vc int
}

func New(i int) *C {
	return &C{
		Vc: i,
	}
}

func (c *C) Show() {
	a.Printf(c.Vc)
}
```

`package a` 依赖 `package b` 和 `package c`，同时 `package b` 依赖 `package a` 、 `package c` 也依赖 `package a` 。

`main` 函数代码如下：

```go
package main

import "github.com/ggq89/mutualdep/a"

func main() {
	a := a.New(3)
	a.Pb.DisplayC()
}
```

编译时就会报错如下：

```bash
import cycle not allowed
package main
	imports github.com/ggq89/mutualdep/a
	imports github.com/ggq89/mutualdep/b
	imports github.com/ggq89/mutualdep/a
```

# 3. 定义接口

现在的问题是：

```bash
A depends on B
B depends on A
```

对于 `A struct` 和 `B struct` 有彼此的指针这种相互依赖问题，可以使用定义接口的方法解决，具体步骤如下：

* 在 `package b` 中 定义 `a interface ` ； 将 `b` 所有使用到结构体 `a` 的变量和方法的地方全部转化成 使用接口 `a` 的方法；在 `a interface ` 中补充缺少的方法；

经过上面的步骤处理后， `package b` 代码如下：

```go
package b

import (
	"github.com/ggq89/mutualdep/c"
)

type B struct {
	Pa a
}

type a interface {
	GetC() *c.C
}

func New(a a) *B {
	return &B{
		Pa:a,
	}
}

func (b *B) DisplayC() {
	b.Pa.GetC().Show()
}
```

* 在 `package a` 中补充可能缺少的方法；

处理后， `package a` 中的代码如下：

```go
package a

import (
	"fmt"

	"github.com/ggq89/mutualdep/b"
	"github.com/ggq89/mutualdep/c"
)

type A struct {
	Pb *b.B
	Pc *c.C
}

func New(ic int) *A {
	a := &A{
		Pc:c.New(ic),
	}

	a.Pb = b.New(a)

	return a
}

func (a *A)GetC() *c.C {
	return a.Pc
}

func Printf(v int)  {
	fmt.Printf("%v", v)
}
```

# 4. 拆分包

再次编译，提示如下：

```
import cycle not allowed
package main
	imports github.com/ggq89/mutualdep/a
	imports github.com/ggq89/mutualdep/b
	imports github.com/ggq89/mutualdep/c
	imports github.com/ggq89/mutualdep/a
```

现在是另一个相互依赖问题：

```bash
A depends on C
C depends on A
```

与前面的相互依赖不同，前面的依赖是由于 `A struct` 和 `B struct` 有彼此的指针导致的，属于硬相互依赖；

而这里是由于 `package c` 中的方法调用 `package a` 中的方法引起的，属于软相互依赖；

* 这种相互依赖可以通过将方法拆分到另一个包的方式来解决；在拆分包的过程中，可能会将结构体的方法转化为普通的函数；

引入 `package f` ， 将方法迁移到 `f` 中 :

```go
package f

import "fmt"

func Printf(v int) {
	fmt.Printf("%v", v)
}
```

方法移动到 `package f` 后， `package a` 的代码如下：

```go
package a

import (
	"github.com/ggq89/mutualdep/b"
	"github.com/ggq89/mutualdep/c"
)

type A struct {
	Pb *b.B
	Pc *c.C
}

func New(ic int) *A {
	a := &A{
		Pc: c.New(ic),
	}

	a.Pb = b.New(a)

	return a
}

func (a *A) GetC() *c.C {
	return a.Pc
}
```

`package c`随之改成调用`package f`，其代码如下：

```go
package c

import (
	"github.com/ggq89/mutualdep/a/f"
)

type C struct {
	Vc int
}

func New(i int) *C {
	return &C{
		Vc: i,
	}
}

func (c *C) Show() {
	f.Printf(c.Vc)
}
```

现在依赖关系如下：

```bash
A depends on B and C
B depends on C
C depends on F
```

至此，两种包相互依赖关系都得以解决。

# 5. 总结

1. 对于软相互依赖，利用分包的方法就能解决，有些函数导致的相互依赖只能通过分包解决；分包能细化包的功能；

2. 对于硬相互依赖只能通过定义接口的方法解决；定义接口能提高包的独立性，同时也提高了追踪代码调用关系的难度；

---

参考文章：

* golang不允许循环import问题("import cycle not allowed") : http://ju.outofmemory.cn/entry/230115
* golang解决import cycle not allowed的一种思路 : https://studygolang.com/articles/10582?fr=sidebar
* golang中("import cycle not allowed")错误 : https://blog.csdn.net/skh2015java/article/details/53943784
