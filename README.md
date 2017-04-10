# gigo(t)

> Cuisse d’agneau, de chevreuil, coupée pour être mangée.
>
> `"Un bon gigot d’agneau ."`

[![travis Status](https://travis-ci.org/mh-cbon/gigo.svg?branch=master)](https://travis-ci.org/mh-cbon/gigo)[![appveyor Status](https://ci.appveyor.com/api/projects/status/github/mh-cbon/gigo?branch=master&svg=true)](https://ci.appveyor.com/project/mh-cbon/gigo)
[![Go Report Card](https://goreportcard.com/badge/github.com/mh-cbon/gigo)](https://goreportcard.com/report/github.com/mh-cbon/gigo)

[![GoDoc](https://godoc.org/github.com/mh-cbon/gigo?status.svg)](http://godoc.org/github.com/mh-cbon/gigo)


go generate super charged on steroids


## Example

in input it takes


###### > demo.gigo.go
```go
// +build gigo

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
<:range $a := .Args> func (s <:$.Name>Slice) FindBy<:$a>(<:$a> <:$.ArgType $a>) (<:$.Name>,bool) {
  for i, item := range s.items {
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

func (s <:.Name>Slice) Index(search <:.Name>) int {
  for i, item := range s.items {
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
// +build gigo

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
 func (s TodoSlice) FindByName(Name string) (Todo,bool) {
  for i, item := range s.items {
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


func (s TodoSlice) Index(search Todo) int {
  for i, item := range s.items {
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
  if i:= s.Index(item); i > -1 {
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
 func (m MutexedTodoSlice) FindByName(Name string)  (Todo,bool) {
  // lock them all
  lock.Lock()
  defer lock.Unlock()
  // invoke embedded type
  m.embed.FindByName(Name)
}
 func (m MutexedTodoSlice) Push(item Todo)  int {
  // lock them all
  lock.Lock()
  defer lock.Unlock()
  // invoke embedded type
  m.embed.Push(item)
}
 func (m MutexedTodoSlice) Index(search Todo)  int {
  // lock them all
  lock.Lock()
  defer lock.Unlock()
  // invoke embedded type
  m.embed.Index(search)
}
 func (m MutexedTodoSlice) RemoveAt(i index)  int {
  // lock them all
  lock.Lock()
  defer lock.Unlock()
  // invoke embedded type
  m.embed.RemoveAt(i)
}
 func (m MutexedTodoSlice) Remove(item Todo)  int {
  // lock them all
  lock.Lock()
  defer lock.Unlock()
  // invoke embedded type
  m.embed.Remove(item)
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

Or you can dump the tokenizer output

###### $ go run main.go -symbol Push dump demo.gigo.go
```sh
-> FuncDecl 11 tokens                    56:0   nlToken              "\n"
                                         57:0   commentLineToken     "// create new Method Push of type ."
                                         57:35  nlToken              "\n"
                                         58:0   funcToken            "func"
                                         58:4   wsToken              " "
 -> PropsBlockDecl 3 tokens              58:5   parenOpenToken       "("
  -> PropDecl 3 tokens                  
   => IdentifierDecl 1 token             58:6   wordToken            "s"
                                         58:7   wsToken              " "
   -> ExpressionDecl 1 tokens           
    -> IdentifierDecl 2 tokens          
     -> BodyBlockDecl 4 tokens           58:8   TplOpenToken         "<:"
                                         58:10  DotToken             "."
                                         58:11  wordToken            "Name"
     <- BodyBlockDecl                    58:15  TplCloseToken        ">"
    <- IdentifierDecl                    58:16  wordToken            "Slice"
   <- ExpressionDecl 1 tokens           
  <- PropDecl 3 tokens                  
 <- PropsBlockDecl                       58:21  parenCloseToken      ")"
                                         58:22  wsToken              " "
 => IdentifierDecl 1 token               58:23  wordToken            "Push"
 -> PropsBlockDecl 3 tokens              58:27  parenOpenToken       "("
  -> PropDecl 3 tokens                  
   => IdentifierDecl 1 token             58:28  wordToken            "item"
                                         58:32  wsToken              " "
   -> ExpressionDecl 1 tokens           
    -> IdentifierDecl 1 tokens          
     -> BodyBlockDecl 4 tokens           58:33  TplOpenToken         "<:"
                                         58:35  DotToken             "."
                                         58:36  wordToken            "Name"
     <- BodyBlockDecl                    58:40  TplCloseToken        ">"
    <- IdentifierDecl 1 tokens          
   <- ExpressionDecl 1 tokens           
  <- PropDecl 3 tokens                  
 <- PropsBlockDecl                       58:41  parenCloseToken      ")"
 -> PropsBlockDecl 1 tokens             
  -> ExpressionDecl 1 tokens            
   -> IdentifierDecl 2 tokens            58:42  wsToken              " "
   <- IdentifierDecl                     58:43  IntToken             "int"
  <- ExpressionDecl 1 tokens            
 <- PropsBlockDecl 1 tokens             
 -> BodyBlockDecl 11 tokens              58:46  wsToken              " "
                                         58:47  BraceOpenToken       "{"
                                         58:48  nlToken              "\n"
                                         59:0   wsToken              " "
                                         59:1   wsToken              " "
  -> ExpressionDecl 4 tokens            
   -> IdentifierDecl 3 tokens            59:2   wordToken            "s"
                                         59:3   DotToken             "."
   <- IdentifierDecl                     59:4   wordToken            "items"
                                         59:9   wsToken              " "
                                         59:10  assignToken          "="
   -> CallExpr 2 tokens                 
    -> IdentifierDecl 2 tokens           59:11  wsToken              " "
    <- IdentifierDecl                    59:12  wordToken            "append"
    -> CallExprBlock 6 tokens            59:18  parenOpenToken       "("
     -> ExpressionDecl 1 tokens         
      -> IdentifierDecl 3 tokens         59:19  wordToken            "s"
                                         59:20  DotToken             "."
      <- IdentifierDecl                  59:21  wordToken            "items"
     <- ExpressionDecl 1 tokens         
                                         59:26  CommaToken           ","
                                         59:27  wsToken              " "
     -> ExpressionDecl 1 tokens         
      => IdentifierDecl 1 token          59:28  wordToken            "item"
     <- ExpressionDecl 1 tokens         
    <- CallExprBlock                     59:32  parenCloseToken      ")"
   <- CallExpr 2 tokens                 
  <- ExpressionDecl 4 tokens            
                                         59:33  nlToken              "\n"
                                         60:0   wsToken              " "
                                         60:1   wsToken              " "
  -> ReturnDecl 4 tokens                 60:2   returnToken          "return"
                                         60:8   wsToken              " "
   -> ExpressionDecl 1 tokens           
    -> CallExpr 2 tokens                
     => IdentifierDecl 1 token           60:9   wordToken            "len"
     -> CallExprBlock 3 tokens           60:12  parenOpenToken       "("
      -> ExpressionDecl 1 tokens        
       -> IdentifierDecl 3 tokens        60:13  wordToken            "s"
                                         60:14  DotToken             "."
       <- IdentifierDecl                 60:15  wordToken            "items"
      <- ExpressionDecl 1 tokens        
     <- CallExprBlock                    60:20  parenCloseToken      ")"
    <- CallExpr 2 tokens                
   <- ExpressionDecl 1 tokens           
  <- ReturnDecl                          60:21  nlToken              "\n"
 <- BodyBlockDecl                        61:0   BraceCloseToken      "}"
<- FuncDecl 11 tokens
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
