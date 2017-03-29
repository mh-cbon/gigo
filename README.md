# gigo

go generate on steroids, in input it takes

```go

type Todo struct {
  Name string
  Done bool
}

type Todos implements<Mutexed (Slice .Todo "Name")> {
  // it reads as a mutexed list of todo,
  // where Name is an additionnal arg to define a FindByName method.
}

func (t *Todos) Hello(){fmt.Println("Hello")}


template Mutexed<.Name> struct {
  lock *sync.Mutex
  embed <.Name>
}

<range $m := .Methods> func (m Mutexed<$.Name>) <$m.Name>(<$m.Params>) <$m.Out> {
  lock.Lock()
  defer lock.Unlock()
  m.embed.<$m.GetName>(<$m.Args>)
}

template <.Name>Slice struct {
  items []<.Name>
}

<range $a := .Args> func (m <$.Name>Slice) FindBy<$a>(<$a> <$.ArgType $a>) (<$.Name>,bool) {
  for i, items := range s.items {
    if item.<$a> == <$a> {
      return item, true
    }
  }
  return {}<$.Name>, false
}

func (s <.Name>Slice) Push(item <.Name>) int {
  s.items = append(s.items, item)
  return len(s.items)
}

func (s <.Name>Slice) Index(item <.Name>) int {
  for i, items := range s.items {
    if item == search {
      return i
    }
  }
  return -1
}

func (s <.Name>Slice) RemoveAt(i index) int {
	s.items = append(s.items[:i], s.items[i+1:]...)
}

func (s <.Name>Slice) Remove(item <.Name>) int {
  if i:= s.Index(item); i > -1 {
    s.RemoveAt(i)
    return i
  }
  return -1
}
```

It produces

```go

type Todo struct {
  Name string
  Done bool
}


type TodoSlice struct {
  items []Todo
}


func (s TodoSlice) Push(item Todo) int {
  s.items = append(s.items, item)
  return len(s.items)
}
 func (m TodoSlice) FindByName(Name string) (Todo,bool) {
  for i, items := range s.items {
    if item.<$a> == <$a> {
      return item, true
    }
  }
  return {}<$.Name>, false
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
  if i:= s.Index(item); i > -1 {
    s.RemoveAt(i)
    return i
  }
  return -1
}
// while this is compatible with its local contracts,
// it will work and still takes advantages of concrete types exported by consumed package.

type MutexedTodoSlice struct {
  lock *sync.Mutex
  embed TodoSlice
}


 func (m MutexedTodoSlice)  FindByName((Name string))  (Todo,bool) {
  lock.Lock()
  defer lock.Unlock()
  m.embed.<$m.GetName>(<$m.Args>)
}
 func (m MutexedTodoSlice)  Push((item Todo))  int {
  lock.Lock()
  defer lock.Unlock()
  m.embed.<$m.GetName>(<$m.Args>)
}
 func (m MutexedTodoSlice)  Index((item Todo))  int {
  lock.Lock()
  defer lock.Unlock()
  m.embed.<$m.GetName>(<$m.Args>)
}
 func (m MutexedTodoSlice)  RemoveAt((i index))  int {
  lock.Lock()
  defer lock.Unlock()
  m.embed.<$m.GetName>(<$m.Args>)
}
 func (m MutexedTodoSlice)  Remove((item Todo))  int {
  lock.Lock()
  defer lock.Unlock()
  m.embed.<$m.GetName>(<$m.Args>)
}


type Todos struct {
	MutexedTodoSlice
  // it reads as a mutexed list of todo.
}

func (t *Todos) Hello(){fmt.Println("Hello")}


```

Still some work to be done, but you got the idea!
