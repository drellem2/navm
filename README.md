# navm

To run demo:
```
./build.sh
demo interpret "1 2 3 * +"
```
or
`demo compile "1 2 3 * +"`

To assemble asm output on apple silicon:

```
as out.s -o out.o
ld out.o -o a.out -l System -syslibroot $(xcrun -sdk macosx --show-sdk-path)  -e _start -arch arm64
```

## Cross-compilation
Not quite ready yet since there's only one backend. Backends will use the same target triples as `zig cc`.
Example: `zig cc -target x86_64-linux-musl`

## Gotchas
- The stack pointer register is not yet supported by the interpreter.
