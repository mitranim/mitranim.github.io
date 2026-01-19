{{mdToToc .MdTpl}}

## Intro

My coding background is in recent languages such as Go, Rust, Clojure, Python, JS, etc. They come with large standard libraries, automatic memory management, generics, type inference, fancy data structures, and more. Perhaps they are better described as "batteries included" languages.

Always wondered what it's like to program in C, which provides none of that. C does come with a sizable standard library, but it mostly deals with the OS and IO. After a month of adjustment, I got surprisingly comfortable!

## Modules

Most C projects have complicated build systems which describe the file dependency graph, instead of having it expressed inside source files. `clangd` also tends to require external configuration in such codebases.

Turns out there's a better way, at least for small projects. Stick `#pragma once` into every file, and always explicitly include `.c` and `.h` files from each other, starting from the main file:

```d
// lib_0.c
#pragma once

// lib_1.c
#pragma once
#include "./lib_0.c"

// main.c
#pragma once
#include "./lib_0.c"
#include "./lib_1.c"
```

Then point the compiler at the main file, without having to fiddle with the build system, list other input files, or configure "include" flags. `clangd` works out of the box with no configuration. Using relative paths ensures path resolution, without any ambiguities.

```sh
cc main.c -o main
```

## Errors

C seems to suggest that errors should be integers. In reality, we want errors to be strings. But we'd rather not have to `free` them, which is noisy and easy to forget. Fortunately, there's a simple solution: statically allocated buffers. This lets you return error strings with abandon:

```d
typedef char const *Err;

char thread_local ERR_BUF[4096];

#define errf(fmt, ...)                                      \
  (                                                         \
    snprintf(ERR_BUF, 4096, fmt __VA_OPT__(,) __VA_ARGS__), \
    ERR_BUF                                                 \
  )

// Caller doesn't need to free!
Err some_func() {return errf("some err: %d", 123);}
```

Error checking can be tersed with a simple "try" macro:

```d
#define try(expr) {Err err = (expr); if (err) return err;}

Err some_func() {
  try(other_func());
  try(other_func());
  try(other_func());
  return nullptr;
}
```

Compiler can remind you to check errors:

```d
#define MUST_USE __attribute((warn_unused_result))

MUST_USE typedef char const *Err;
```

The stdlib supports backtraces. They're not as precise as in "modern" languages; you get only procedure names with instruction offsets, without file / row / col info, but it's often good enough.

Putting it all together:

```details | d | error stuff — click to expand
#define MUST_USE __attribute((warn_unused_result))

MUST_USE typedef char const *Err;

char thread_local ERR_BUF[4096];
void thread_local *BT_BUF[256] = {};
int  thread_local BT_BUF_LEN   = 0;

#define errf(fmt, ...)                                      \
  (                                                         \
    backtrace_capture(),                                    \
    snprintf(ERR_BUF, 4096, fmt __VA_OPT__(,) __VA_ARGS__), \
    ERR_BUF                                                 \
  )

void backtrace_capture() {
  BT_BUF_LEN = backtrace(BT_BUF, arr_cap(BT_BUF));
}

void backtrace_print() {
  if (BT_BUF_LEN) backtrace_symbols_fd(BT_BUF, BT_BUF_LEN, STDERR_FILENO);
}

Err some_func() {try(func()); try(func()); try(func()); return nullptr;}

int main() {
  Err err = some_func();
  if (!err) return 0;

  fprintf(stderr, "error: %s\n", err);
  backtrace_print();
  return 1;
}
```

## Numbers

To keep the language portable to weird systems, the standard defines weird numeric types. Their names are kind of insane. "Short"? "long"? "double"? "long _long_"? The words don't mean anything anymore. Is `size_t` the same width as `uintptr_t`? Is `off_t` same or larger? Is `sizeof(void*)` same as `sizeof(size_t)`? The standard has a lot to say, which mostly amounts to "implementation dependent".

Fortunately, on non-weird architectures, the answer tends to be "it's word-sized", and `char` is 8 bits. One can define a small set of sane number types and stick to those. Also, use `-funsigned-char`. The aliases below are mostly superfluous, but I find them easier to type than `*_t`.

```details | d | num.h — click to expand
#pragma once
#include <stddef.h>
#include <stdint.h>
#include <wchar.h>

typedef size_t  Uint;
typedef ssize_t Sint;

typedef uint8_t U8;
typedef int8_t  S8;

typedef uint16_t U16;
typedef int16_t  S16;

typedef uint32_t U32;
typedef int32_t  S32;

typedef uint64_t U64;
typedef int64_t  S64;

typedef float  F32;
typedef double F64;

#define FMT_UINT "%zu"
#define FMT_SINT "%zd"
```

## Type inference

C23 adds the `auto` variable type. Very handy when coming from Go and Rust where local type inference is a norm. You may need to instruct the compiler to enable C23 features. In Sublime Text, I defined a snippet `let` which expands to `const auto`, making this easy to type:

```d
const auto var0 = some_func();
const auto var1 = another_func();
```

## Generics

Even without the `_Generic` macro, you can go a _long_ way towards generic data structures and operations via carefully written macros, backed by support procedures as necessary. Excerpt from my C replica of Go slices:

```d
#define list_of(Elem) \
  struct {            \
    Elem *dat;        \
    Uint len;         \
    Uint cap;         \
  }

typedef list_of(Uint) Uint_list;
typedef list_of(Sint) Sint_list;

#define list_append(tar, ...)                                \
  ({                                                         \
    const auto ptr = (tar);                                  \
    list_reserve_more((List_head*)ptr, sizeof(ptr->dat[0])); \
    ptr->dat[ptr->len++] = (__VA_ARGS__);                    \
  })

int main() {
  defer(list_deinit) Uint_list uints = {};
  defer(list_deinit) Sint_list sints = {};
  list_append(&uints, 123);
  list_append(&sints, 234);
}
```

`*ptr++` all over the place may be considered idiomatic, but I'd rather read _words_. These thin veneers can also check bounds, avoiding buffer overflows.

Not having dicts/maps in the standard library was a bit of a worry, but they turned out easy to implement in a semi-generic fashion (mine is under 200 LoC). Nice open source generic data structure libraries exist, but they're overkill for many apps.

## Resource management

Having to manually allocate and `free` is often cited as the biggest vice of C after [buffer overflows](#buffer-overflows-and-guards). In practice, I'm finding this a non-issue thanks to handy patterns:

* Deferred cleanup.
* Preallocated arenas.
* Inline buffers in structs.
* Statically allocated buffers.

### Defer

GCC and Clang support deferred variable cleanup, which lets you attach cleanup functions to variables. A bit of macro magic makes it look nicer. Write deinit functions for a few types, and you basically have RAII. Example:

```d
#define defer(fun) __attribute__((cleanup(fun)))

void file_deinit(FILE **file) {if (file && *file) fclose(*file);}

Err file_read(char *path, char **out_body, Uint *out_len) {
  defer(file_deinit) FILE *file = fopen(path, "r");
  // ... malloc a return buffer; read file ...
  *out_body = buf;
  *out_len = len;
}

void deinit_mem(void *ptr) {if (ptr) free(*(void **)ptr);}

int main() {
  defer(deinit_mem) char *body;
  size_t len;
  file_read("./readme.md", &body, &len);
}
```

Type-specific destructors can be shortened even further, while keeping our deinit macro generic:

```d
#define defer_deinit(typ) typ __attribute__((cleanup(typ##_deinit)))

void FILE_deinit(FILE **file) {if (file && *file) fclose(*file);}

int main() {
  defer_deinit(FILE) *file = fopen(path, "r");
}
```

(But mind to check `fclose` errors after _writing_!)

### Preallocated buffers

How much space do you need? Can you get away with a small fixed-size buffer? Then just put it on the stack or inside your structures; freeing is automatic:

```d
typedef struct {
  size_t len;
  char   buf[128];
} Word;

typedef struct {
  size_t len;
  char   buf[4096];
} Line;

typedef struct {
  Word word;
  Line line;
  // ...
} Reader;

int main() {
  Reader read = {};

  read_word(&read, stdin);
  puts(read.word.buf);

  read_line(&read, stdin);
  puts(read.line.buf);
}
```

Buffers can also be static. Handy when content lifetimes don't overlap. No need to allocate and free:

```d
char static thread_local ERR_MSG[4096];
```

### Preallocated arenas

When building an append-only collection of objects, you can sometimes estimate in advance how much space is enough, and allocate it just once. This is especially true if the program imposes artificial limits on object count.  Combined with deferred deinit, this spares you from worrying about freeing individual objects.

```d
void some_func(void) {
  defer(stack_deinit) Object_stack stack = {};
  stack_init(&stack);
  // Use the memory.
}
```

Unlike Go-style resizable buffers which relocate data in memory when they grow, such arenas give you _stable object pointers_, which can be important when objects cross-reference each other a lot.

## Buffer overflows and guards

When allocating large buffers, in the range of memory page size or more (16 KiB on MacOS), you can avoid off-by-small-amount overlows and underflows by, well... crashing your program, which is better than data corruption.

The trick is to `mmap` / `mprotect` memory in the shape `guard|buffer|guard`, where guards are `PROT_NONE` while the buffer is `PROT_READ|PROT_WRITE`. Stepping into the guards delivers us a segfault.

## Compiler flags

C compilers come with many built-in diagnostics which are disabled by default. Enable them for added safety.

The following should be placed in `compile_flags.txt` so that `clangd` will pick it up. This makes it work out of the box. No fidding with `compile_commangs.json`; all we need is to slurp this file into compiler flags in our makefile.

Here are the flags I currently use. Different projects may prefer different diagnostics. Mind that `-fsanitize` flags (especially `address`) have runtime overheads and should only be used in development and debugging.

```details | sh | compile_flags.txt — click to expand
-std=c23
-funsigned-char
-fsanitize=undefined,address,integer,nullability

-Weverything
-Wno-pre-c23-compat
-Wno-c++98-compat
-Wno-padded
-Wno-missing-prototypes
-Wno-poison-system-directories
-Wno-pragma-once-outside-header
-Wno-declaration-after-statement
-Wno-covered-switch-default
-Wno-unused-function
-Wno-unused-macros
-Wno-extra-semi-stmt
-Wno-gnu-statement-expression-from-macro-expansion
-Wno-unsafe-buffer-usage
-Wno-pre-c11-compat
-Wno-shadow
-Wno-unreachable-code-return
-Wno-gnu-label-as-value
-Wno-empty-translation-unit
-Wno-c++-compat
-Wno-format-pedantic
-Wno-documentation-html
-Wno-gnu-empty-struct
-Werror=return-type
```

## Debugging

Although we've successfully acquired backtraces for "handled" errors, it's still easy to crash without a trace. (Common experience when writing a slightly buggy JIT compiler 😅.) What to do?

One quick and dirty solution is to fire up `lldb` (or your preferred debugger) and just run the program until it crashes. The debugger preserves the last state, letting you inspect registers, memory, and the backtrace. Having thus identified the faulty code, we can just `printf`-debug the hell out of it, which is often faster.

Shopped around for other debuggers and disassemblers, and didn't find anything better than `lldb` for MacOS, at least among the free offerings. Sometimes it's handier to hop into Xcode for its GUI frontend to `lldb`, which can show more stuff at once without having to type commands all the time. `gdb` doesn't seem to work on my system, Ghidra is horrible, and Radare2 uses Capstone whose Arm64 assembler / disassembler is a buggy liar.

## Tooling

C is known to need a lot of external tools. Fortunately, it was easy to avoid what I feared the worst: convoluted build systems with cmake / autoconf / automake etc.; [`include`](#modules) does the job just fine.

One thing I still really miss is semantically meaningful syntax highlighting, like the one I [wrote](https://github.com/sublimehq/Packages/pull/1662) for Go, where all declarations and modifiers are usefully scoped without special-casing built-ins. Probably end up rewriting C for Sublime at some point.

## Conclusion

I picked C for writing a Forth compiler (more on that in another post) because it provides decent control over the ABI (via assembly), and ready access to some OS APIs needed for JIT engines (`mmap` and more). Plus, I just wanted to learn the language and its patterns.

Now if someone pointed a gun (or hired me) and told me to write a web server in C, I wouldn't bat an eye. But even without that, I would consider it depending on the use case.
