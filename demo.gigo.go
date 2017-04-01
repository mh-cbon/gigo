package main

type Todo struct {
  Name string
  Done bool
}

type Todos implements<:Mutexed (Slice, .Todo "Name")> {
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
