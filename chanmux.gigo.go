// +build gigo

package main

// see https://dave.cheney.net/2016/11/13/do-not-fear-first-class-functions
// P: Let’s talk about actors

type Todos implements<:Slice .Todo>{}

/*
type Todo struct {
  Name string
  Done bool
}

type Todos struct {
  TodosSlice
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
*/

type MuxTodos implements<:ChanMuxer .Todos>{}

/*


type MuxTodos struct {
  TodosChanMuxer
}

type TodosChanMuxer struct {
  ops chan func(Todos)
  stop chan bool
  started chan bool
}

// for every method of ., create a new method on ChanMux
func (m TodosChanMuxer) RemoveAt(i index) int {
  res := make(chan int)
  m.ops <- func(embed Todos) {
    res <- embed.RemoveAt(i)
  }
  return <-res
}

func (m *TodosChanMuxer) loop() {
  embed := &Todos{}
  for {
    select {
    case op:=<-m.ops:
      op(embed)
    case s:=<-m.stop:
      return
    default:
      m.started<-true
    }
  }
}

func (m *TodosChanMuxer) Start() bool {
  m.loop()
  return <-m.started
}

func (m *TodosChanMuxer) Stop() bool {
  s := <-m.stop
  return s
}
*/


template <:.Name>ChanMuxer struct {
  ops chan func(<:.Name>)
  stop chan bool
}

// for every method of ., create a new method on ChanMux
<:range $m := .Methods> func (m *<:.Name>ChanMuxer) <:$m.Name>(<:$m.GetArgsBlock | joinexpr ",">) <:$m.Out> {
  res := make(chan []interface{})
  m.ops <- func(embed <:.Name>) {
    <:$m.Out | joinexpr ","> := embed.<:$m.GetName>(<:$m.GetArgsNames | joinexpr ",">)
    res <- []interface{<:$m.Out | joinexpr ",">}
  }
  ret := <-res
  // how to return ret ?
  return <:$m.Out | joinexpr ",">
}

func (m *<:.Name>ChanMuxer) loop() {
  embed := &<:.Name>{}
  for {
    select {
    case op:=<-m.ops:
      op(embed)
    case s:=<-m.stop:
      return
    }
  }
}

func (m *<:.Name>ChanMuxer) Start()  {
  m.loop()
}

func (m *<:.Name>ChanMuxer) Stop()  {
  m.stop<-true
}
