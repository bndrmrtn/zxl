# Zx (/ziː.ɛks/) Language

Zex is a simple programming language
that is designed to be easy to use and understand.

## About

Zex is a weakly-typed, interpreted language that is designed to be easy to use and understand.
Zex is built in Go as a learning project and is not intended to be used in production.

## Features

- Weakly-typed
- Interpreted
- Zex Blocks

## Installation

To install Zex, you need to have Go installed on your system.

```bash
go install github.com/bndrmrtn/zxl@latest
```

## Usage

To run a Zex program, you can use the `zexlang` command.

```bash
zexlang run <file>
```

Flags can be used to cache or debug the program.

```bash
zexlang run <file> --cache --debug
```

## Examples

#### Basic Hello World program:

```zex
println("Hello, World!");
```

#### Define a variable and print it:

```zex
let x = 10;
println(x);
```

#### Define a function and call it:

```zex
fn add(a, b) {
  return a + b;
}

let result = add(10, 20);
println(result);
```

#### Define a block and use it:
(this is currently not working as intended, but it is the intended syntax)

```zex
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
```zex
namespace main;

import("other.zx");
other.printHello();
```

File `other.zx`:
```zex
namespace other;

fn printHello() {
  println("Hello from other.zx!");
}
```

## Support

Support my work by giving this project a star.

- [PayPal](https://www.paypal.me/instasiteshu)
- [Ko-Fi](https://ko-fi.com/bndrmrtn)
