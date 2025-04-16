# Flare ✨

Flare is a simple programming language
that is designed to be easy to use and understand. [Examples](https://github.com/orgs/flarelang/repositories).

## About

Flare is a weakly-typed, interpreted language that is designed to be easy to use and understand.
Flare is built in Go as a learning project and is not intended to be used in production.

## Features

- Weakly-typed
- Interpreted
- Flare Blocks
- Threads and Concurrency

# Syntax Highlihting

Flare now has a `VSCode` plugin for highlighting code (no LSP currently). [Download now](https://marketplace.visualstudio.com/items/?itemName=flarelang.flarelang)

## Installation

To install Flare, you need to have Go installed on your system.

```bash
go install github.com/flarelang/flare@latest
```

## Usage

To run a Flare program, you can use the `flare` command.

```bash
flare run <file>
```

Flags can be used to cache or debug the program.

```bash
flare run <file> --cache --debug
```

## Examples

#### Basic Hello World program:

```flare
println("Hello, World!");
```

#### Define a variable and print it:

```flare
let x = 10;
println(x);
```

#### Define a function and call it:

```flare
fn add(a, b) {
  return a + b;
}

let result = add(10, 20);
println(result);
```

### Error handling:

```flare
error err: fail("This is a helper to throw an error");
if err != nil {
  println("error occurred:", err);
}
```

```flare
error otherErr {
  const x = 5;
  x = 6; // an error will happen
}

if otherErr != nil {
  println("error occurred:", otherErr);
}

println(x); // x is in the same scope as global so it will print 5, because
            // thats the value of x before the error occurred
```

#### Define a block and use it:

```flare
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

File `main.fl`:
```flare
namespace main;

import("other.fl");
other.printHello();
```

File `other.fl`:
```flare
namespace other;

fn printHello() {
  println("Hello from other.fl!");
}
```

### Loops

```flare
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

```flare
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

```flare
use thread;

fn doLater() {
  // some task
}

thread.spawn(doLater);
thread.sleep(1000); // Wait one second
```

### Usage of Portals

```flare
use thread;

define User {
  let name;
  let age;
  let portal;

  fn construct(name, age) {
    this.name = name;
    this.age = age;
  }

  fn intro() {
    println("Hello, I am ", this.name, "and I am ", this.age, " years old.");
    this.portal.send(true);
  }
}

let users = [];
users.append(User("John", 25));
users.append(User("Jane", 23));
users.append(User("Emily", 21));
// ...

// create a portal for communication between threads
const portal = thread.portal(users.length);
// create a custom spawner with users.length async threads
const spawner = thread.spawner(users.length);

// spawn all async methods
for user in users {
  user.portal = portal;
  spawner.spawn(user.intro);
}

// wait for all async methods to finish
for i in (users.length) {
  // portal.receive() waits until any async method sends a message to it
  portal.receive();
}

// close the portal and spawner
portal.close();
spawner.close();
```

## Support

Support my work by giving this project a star.

- [PayPal](https://www.paypal.me/instasiteshu)
- [Ko-Fi](https://ko-fi.com/bndrmrtn)
