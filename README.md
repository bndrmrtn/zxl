# Zx (/ziː.ɛks/) Language

Zx is a simple programming language
that is designed to be easy to use and understand.

## About

Zx is a weakly-typed, interpreted language that is designed to be easy to use and understand.
Zx is built in Go as a learning project and is not intended to be used in production.

## Features

- Weakly-typed
- Interpreted
- Zx Blocks
- Threads and Concurrency

## Installation

To install Zx, you need to have Go installed on your system.

```bash
go install github.com/bndrmrtn/zxl@latest
```

## Usage

To run a Zx program, you can use the `zxl` command.

```bash
zxl run <file>
```

Flags can be used to cache or debug the program.

```bash
zxl run <file> --cache --debug
```

## Examples

#### Basic Hello World program:

```zxl
println("Hello, World!");
```

#### Define a variable and print it:

```zxl
let x = 10;
println(x);
```

#### Define a function and call it:

```zxl
fn add(a, b) {
  return a + b;
}

let result = add(10, 20);
println(result);
```

#### Define a block and use it:

```zxl
define MyBlock {
  let x = 10;

  fn construct(value) {
    this.x = value;
  }
}

let block = MyBlock(20);
println(block.x);
```

### Using namespaces:

File `main.zx`:
```zxl
namespace main;

import("other.zx");
other.printHello();
```

File `other.zx`:
```zxl
namespace other;

fn printHello() {
  println("Hello from other.zx!");
}
```

### Loops

```zxl
for i in range(10) {
  println(i);
}

for i in range([10, 20]) {
  println(i);
}

for i in range([10, 20, 2]) {
  println(i);
}

for letter in "Hello" {
  println(letter);
}

for i in 7 {
  println(i);
}

for i in [5, 6, 7] {
  println(i);
}

let i = 0;
while i < 10 {
  println(i);
  i = i + 1;
}
```

### Arrays

```zxl
use iter;

let myArr = array {
  name: "John",
  "age": 30,
  city: "New York",
};

for row in iter.Array(myArr) {
    println(row.key + ":", row.value);
}
```

### Concurrency

```zxl
use thread;

fn doLater() {
  // some task
}

thread.spawn(doLater);
thread.sleep(1000); // Wait one second
```

```zxl
use thread;

define MyWorker {
  let portal;
  fn construct(portal) {
    this.portal = portal;
  }

  fn run() {
    while true {
      const data = this.portal.receive();
      if data != false {
        println("Received data:", message);
      }
    }
  }
}

const portal = thread.portal();
const worker = MyWorker(portal);
thread.spawn(worker.run);

for i in range(10) {
  portal.send(i);
}

thread.sleep(1000);
```

## Support

Support my work by giving this project a star.

- [PayPal](https://www.paypal.me/instasiteshu)
- [Ko-Fi](https://ko-fi.com/bndrmrtn)
