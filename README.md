gigo

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


###### $ go run main.go
```sh
package main

type Todo struct {
  Name string
  Done bool
}// a template to generate a type Slice of .
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

Still some work to be done, but you got the idea!

## Changes

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
