# gigo

[![travis Status](https://travis-ci.org/mh-cbon/gigo.svg?branch=master)](https://travis-ci.org/mh-cbon/gigo)[![appveyor Status](https://ci.appveyor.com/api/projects/status/github/mh-cbon/gigo?branch=master&svg=true)](https://ci.appveyor.com/project/mh-cbon/gigo)
[![Go Report Card](https://goreportcard.com/badge/github.com/mh-cbon/gigo)](https://goreportcard.com/report/github.com/mh-cbon/gigo)

[![GoDoc](https://godoc.org/github.com/mh-cbon/gigo?status.svg)](http://godoc.org/github.com/mh-cbon/gigo)


go generate super charged on steroids


## Example

in input it takes


###### > demo.gigo.go
```go
package main

type Todo struct {
  Name string
  Done bool
}

type Todos implements<:Mutexed (Slice .Todo "Name")> {
  // it reads as a mutexed list of todo.
}

func (t *Todos) Hello(){fmt.Println("Hello")}

// type Todos implements<Mutexed (Slice .Todo)>
// type Todos implements<.Todo | Slice | Mutexed>
//-
// It should probably be something like this in real apps
// type Todos implements<Slice .Todo>
// type TodosManager struct {Items Todos}
// type MutexedTodosManager implements<Mutexed .TodosManager>


// a template to mutex .
template Mutexed<:.Name> struct {
  lock *sync.Mutex
  // embed the type
  embed <:.Name>
}

// for every method of ., create a new method of Mutexed
<:range $m := .Methods> func (m Mutexed<:$.Name>) <:$m.Name>(<:$m.GetArgsBlock | joinexpr ",">) <:$m.Out> {
  // lock them all
  lock.Lock()
  defer lock.Unlock()
  // invoke embedded type
  m.embed.<:$m.GetName>(<:$m.GetArgsNames | joinexpr ",">)
}

// a template to generate a type Slice of .
template <:.Name>Slice struct {
  items []<:.Name>
}

// range over args to produce new FindBy methods
<:range $a := .Args> func (m <:$.Name>Slice) FindBy<:$a>(<:$a> <:$.ArgType $a>) (<:$.Name>,bool) {
  for i, items := range s.items {
    if item.<:$a> == <:$a> {
      return item, true
    }
  }
  return <:$.Name>{}, false
}

// create new Method Push of type .
func (s <:.Name>Slice) Push(item <:.Name>) int {
  s.items = append(s.items, item)
  return len(s.items)
}

func (s <:.Name>Slice) Index(item <:.Name>) int {
  for i, items := range s.items {
    if item == search {
      return i
    }
  }
  return -1
}

func (s <:.Name>Slice) RemoveAt(i index) int {
	s.items = append(s.items[:i], s.items[i+1:]...)
}

func (s <:.Name>Slice) Remove(item <:.Name>) int {
  if i:= s.Index(item); i > -1 {
    s.RemoveAt(i)
    return i
  }
  return -1
}
```

It produces


###### $ go run main.go gen demo.gigo.go
```sh
package main

type Todo struct {
  Name string
  Done bool
}
// a template to generate a type Slice of .
type TodoSlice struct {
  items []Todo
}

// range over args to produce new FindBy methods
 func (m TodoSlice) FindByName(Name string) (Todo,bool) {
  for i, items := range s.items {
    if item.Name == Name {
      return item, true
    }
  }
  return Todo{}, false
}

// create new Method Push of type .
func (s TodoSlice) Push(item Todo) int {
  s.items = append(s.items, item)
  return len(s.items)
}


func (s TodoSlice) Index(item Todo) int {
  for i, items := range s.items {
    if item == search {
      return i
    }
  }
  return -1
}


func (s TodoSlice) RemoveAt(i index) int {
	s.items = append(s.items[:i], s.items[i+1:]...)
}


func (s TodoSlice) Remove(item Todo) int {
  if i:= s.Index(item); i :> -1 {
    s.RemoveAt(i)
    return i
  }
  return -1
}

// a template to mutex .
type MutexedTodoSlice struct {
  lock *sync.Mutex
  // embed the type
  embed TodoSlice
}

// for every method of ., create a new method of Mutexed
 func (m MutexedTodoSlice)  FindByName(Name string)  (Todo,bool) {
  // lock them all
  lock.Lock()
  defer lock.Unlock()
  // invoke embedded type
  m.embed. FindByName(Name)
}
 func (m MutexedTodoSlice)  Push(item Todo)  int {
  // lock them all
  lock.Lock()
  defer lock.Unlock()
  // invoke embedded type
  m.embed. Push(item)
}
 func (m MutexedTodoSlice)  Index(item Todo)  int {
  // lock them all
  lock.Lock()
  defer lock.Unlock()
  // invoke embedded type
  m.embed. Index(item)
}
 func (m MutexedTodoSlice)  RemoveAt(i index)  int {
  // lock them all
  lock.Lock()
  defer lock.Unlock()
  // invoke embedded type
  m.embed. RemoveAt(i)
}
 func (m MutexedTodoSlice)  Remove(item Todo)  int {
  // lock them all
  lock.Lock()
  defer lock.Unlock()
  // invoke embedded type
  m.embed. Remove(item)
}


type Todos struct {
	MutexedTodoSlice
  // it reads as a mutexed list of todo.
}

func (t *Todos) Hello(){fmt.Println("Hello")}
```

You can also get a specific symbol

###### $ go run main.go -symbol Push gen demo.gigo.go
```sh
// create new Method Push of type .
func (s TodoSlice) Push(item Todo) int {
  s.items = append(s.items, item)
  return len(s.items)
}
```

Or you can dump the otkenizer output

###### $ go run main.go -symbol Push dump demo.gigo.go
```sh
begin  *glang.FuncDecl      Tokens(10)
                                          53:  0              nlToken "\n"
                                          54:  0     commentLineToken "// create new Method Push of type ."
                                          54: 35              nlToken "\n"
                                          55:  0            funcToken "func"
                                          55:  4              wsToken " "
 begin  *glang.PropsBlockDecl Tokens(3)
                                          55:  5       parenOpenToken "("
  begin  *glang.PropDecl      Tokens(2)
   begin  *glang.IdentifierDecl Tokens(1)
                                          55:  6            wordToken "s"
   end    *glang.IdentifierDecl tokens(1)
   begin  *glang.IdentifierDecl Tokens(2)
    begin  *glang.BodyBlockDecl Tokens(4)
                                          55:  7              wsToken " "
                                          55:  8         TplOpenToken "<:"
                                          55: 10            wordToken ".Name"
                                          55: 15        TplCloseToken ">"
    end    *glang.BodyBlockDecl tokens(4)
                                          55: 16            wordToken "Slice"
   end    *glang.IdentifierDecl tokens(2)
  end    *glang.PropDecl      tokens(2)
                                          55: 21      parenCloseToken ")"
 end    *glang.PropsBlockDecl tokens(3)
 begin  *glang.IdentifierDecl Tokens(2)
                                          55: 22              wsToken " "
                                          55: 23            wordToken "Push"
 end    *glang.IdentifierDecl tokens(2)
 begin  *glang.PropsBlockDecl Tokens(3)
                                          55: 27       parenOpenToken "("
  begin  *glang.PropDecl      Tokens(2)
   begin  *glang.IdentifierDecl Tokens(1)
                                          55: 28            wordToken "item"
   end    *glang.IdentifierDecl tokens(1)
   begin  *glang.IdentifierDecl Tokens(1)
    begin  *glang.BodyBlockDecl Tokens(4)
                                          55: 32              wsToken " "
                                          55: 33         TplOpenToken "<:"
                                          55: 35            wordToken ".Name"
                                          55: 40        TplCloseToken ">"
    end    *glang.BodyBlockDecl tokens(4)
   end    *glang.IdentifierDecl tokens(1)
  end    *glang.PropDecl      tokens(2)
                                          55: 41      parenCloseToken ")"
 end    *glang.PropsBlockDecl tokens(3)
 begin  *glang.PropsBlockDecl Tokens(1)
  begin  *glang.PropDecl      Tokens(1)
   begin  *glang.IdentifierDecl Tokens(2)
                                          55: 42              wsToken " "
                                          55: 43            wordToken "int"
   end    *glang.IdentifierDecl tokens(2)
  end    *glang.PropDecl      tokens(1)
 end    *glang.PropsBlockDecl tokens(1)
 begin  *glang.BodyBlockDecl Tokens(27)
                                          55: 46              wsToken " "
                                          55: 47     bracketOpenToken "{"
                                          55: 48              nlToken "\n"
                                          56:  0              wsToken " "
                                          56:  1              wsToken " "
                                          56:  2            wordToken "s.items"
                                          56:  9              wsToken " "
                                          56: 10          assignToken "="
                                          56: 11              wsToken " "
                                          56: 12            wordToken "append"
                                          56: 18       parenOpenToken "("
                                          56: 19            wordToken "s.items"
                                          56: 26       semiColonToken ","
                                          56: 27              wsToken " "
                                          56: 28            wordToken "item"
                                          56: 32      parenCloseToken ")"
                                          56: 33              nlToken "\n"
                                          57:  0              wsToken " "
                                          57:  1              wsToken " "
                                          57:  2          returnToken "return"
                                          57:  8              wsToken " "
                                          57:  9            wordToken "len"
                                          57: 12       parenOpenToken "("
                                          57: 13            wordToken "s.items"
                                          57: 20      parenCloseToken ")"
                                          57: 21              nlToken "\n"
                                          58:  0    bracketCloseToken "}"
 end    *glang.BodyBlockDecl tokens(27)
end    *glang.FuncDecl      tokens(10)
```

Or get it to string after tokenization

###### $ go run main.go -symbol Push str demo.gigo.go
```sh
// create new Method Push of type .
func (s <:.Name>Slice) Push(item <:.Name>) int {
  s.items = append(s.items, item)
  return len(s.items)
}
```

## Changes

#### Cli

Added cli features to `gen`, `dump` and `output` results.

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
