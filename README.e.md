# {{.Name}}(t)

> Cuisse d’agneau, de chevreuil, coupée pour être mangée.
>
> `"Un bon gigot d’agneau ."`

{{template "badge/travis" .}}{{template "badge/appveyor" .}}{{template "badge/goreport" .}}{{template "badge/godoc" .}}

{{pkgdoc}}

## Example

in input it takes

{{file "demo.gigo.go"}}

It produces

{{cli "go" "run" "main.go" "gen" "demo.gigo.go"}}

You can also get a specific symbol
{{cli "go" "run" "main.go" "-symbol" "Push" "gen" "demo.gigo.go"}}

Or you can dump the tokenizer output
{{cli "go" "run" "main.go" "-symbol" "Push" "dump" "demo.gigo.go"}}

Or get it to string after tokenization
{{cli "go" "run" "main.go" "-symbol" "Push" "str" "demo.gigo.go"}}

## Changes

#### Cli

Added cli features to gen, dump and output results.

#### Fixed body parsing and printing

Now when a template func is encountered

```go
<:range $m := .Methods> func (m Mutexed<:$.Name>) <:$m.Name>(<:$m.GetArgsBlock | joinexpr ",">) <:$m.Out> {
  lock.Lock()
  defer lock.Unlock()
  m.embed.<:$m.GetName>(<:$m.GetArgsNames | joinexpr ",">)
}
```

Its body is evaluated, and some helpers have been added to properly display it.

__before__
```go
 func (m MutexedTodoSlice)  Push((item Todo))  int {
  lock.Lock()
  defer lock.Unlock()
  m.embed. Push((item Todo))
}
```
__after__
```go
 func (m MutexedTodoSlice)  Push(item Todo)  int {
  lock.Lock()
  defer lock.Unlock()
  m.embed. Push(item)
}
```


#### Added nice error support

```go
package tomate

type tomate struct qsdqd{} // bad
```

```sh
unexpected token
In file=<noname> At=3:19
Found=wordToken wanted=[bracketOpenToken]

...
5  package tomate
6  type tomate struct qsdqd{}
   ---------------------↑
...
```
