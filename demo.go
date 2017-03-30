package main

// Let s say we provide an api to manage TODOS
// an http rest api.
// Let s manage it in memory for the sake of example.
// In such situation, the backend server must maintain in memory
// a list of data, currently, the list of todos.
// It seems a super simple task,
// but because of concurrency it is not that much,
// and even for an exeprienced developer, the mistake might happen quickly,
// in both case the result is desastrous,
// the app become buggy, unresponsive and behave weirdly.
//
// Lets code an example of what such shared list would be


// As an API provider I declare concrete types.
type Todo struct {
  Name string
  Done bool
}

// I want to expose a list type which let consumer acces todo items in a tread safe way.
type TodoList struct {
  items []todo
  lock mutex.Sync // <- the lock
}

// To avoid deadlock problem, the methods of the todoList
// will exist in two version

// - a private version, non TS, it can call other non TS methods,
func (t todoList) push(t todo) {
  // lock.Lock() // you must not lock here, if you d do so, you might easily
  // get a deadlock becasue you tried to lock an already locked mutex.
  // to avoid that, carefully make use of other nonTS methods
  /// ... code
  if true {

  }
}

// - a public version, totally TS, it must only take care to call only nonTS versions
func (t todoList) Push(t todo) {
  lock.Lock()
  defer lock.Unlock() // the thread safety mechanism must be added every where, and carefully unlocked
  return t.push(t) // call non ts methods to not deadlock.
}


/*
// doing so you can avoid deadlock issues for the consumer,
// and its not too dificult to write for the provider,
// but it still require cares and attention.

// Lets jump into implements type templater!
*/

// implements type templater are concrete type made of virtual types,
// It only dynamically declares (really just evaluate the templates at pre-static analisys phase) the signature
// of the named type Todos.
// It is exactly the same, and absolutely no more than manually declare the Todos type.
// But its cool becasue it improves coding experience,
// it mixins virutal types which are type generator to populate the static type declaration.

// In our example, we might want to do something similar to this,
// - defines a concrete type to name a []Todo
// - Apply it Slice capability to provide methods such Push/Index/Remove ect
// - Ensure it exposes a public API which is totally thread safe

type TodoSlice struct {
  items []Todo
}


 func (m TodoSlice) FindByName(Name string) (Todo,bool) {
  for i, items := range s.items {
    if item.<$a> == <$a> {
      return item, true
    }
  }
  return {}<$.Name>, false
}


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
  if i:= s.Index(item); i > -1 {
    s.RemoveAt(i)
    return i
  }
  return -1
}
// the programmer fixed a tricky problem
// in a glance!

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

// or symply
// type Todos implements<Mutexed (Slice .Todo)>
// type Todos implements<.Todo | Slice | Mutexed>
//-
// It should probably be something like this in real apps
// type Todos implements<Slice .Todo>
// type TodosManager struct {Items Todos}
// type MutexedTodosManager implements<Mutexed .TodosManager>

// The equivalent go code result should be
/*
type SliceTodos struct {
  items []Todo
}
func (s SliceTodos)Push(t Todo) int{...}

type MutexedTodos struct {
  embed SliceTodos
}
func (s MutexedTodos)Push(t Todo) int{
  lock.Lock()
  defer lock.Unlock()
  return s.embed.Push(t)
}

type Todos struct {
  MutexedTodos
}
*/

// For a consumer it means it can access
// a concrete Todos type
// with a well defined set of methods.
// thus, let the consumer emit a contract againt this type.
type todosProvider interface {
  Push(Todo)
  Remove(Todo)
  // consumer does not need more methods in this demo.
}
type Consumer struct {
  todos todosProvider
}
func (c Consumer) AddTodo(t Todo) error {
  if t.Name=="" {
    return fmt.Errof("Todo name must not be nil")
  }
  c.todos.Push(todo) // Really cool thing here is
  // that the mutexed capability of the slice is controlled at type declaration.
  // if the consumer would like to use a non mutexed,
  // it simply delcares a local type such type myTodos implements<Slice .Todo>,
  // then inject a new instance of it into the consumer isntance.
  // Boom, improved api.
  return nil
}
func (c Consumer) RmTodo(t Todo) error {
  if i := c.todos.Remove(); i == -1 {
    return fmt.Errof("Todo was not found %v", t)
  }
  return nil
}

// Wait, go can go generate, why is this any better ?
/*
This is very similar to go generate, but, go generate has few problems
- go generate can do anything, it also means it is not generic enough to provide a dsl-like
- go generate does not intervein in the static analysis phase, it works before that => no completion, no analysis, just generate
- its not fun to use you need to take care of various little things (out location, go gen command, and as of today, write the generator)

implements type templater is no more than a templater,
but it is exactly what s needed to improve programming experience of go,
it does not hurt the go type system
it helps to factorize and produce better code with less bug
*/


// other ideas

// trait ?

// type Formatter interface{
//   Format(f State, c rune)
// }

/*
template Dumper trait {
}
func (s Dumper) Format(f State, c rune)  {
  <range $p := $.Props>
    <.Name>
  <end>
  io.WriteString(os.Stdout, "build a pretty printed string of s")
}

type PrettyTodos implements<Dumper .Todo>

type MyType implements<Mutexed (Slice .PrettyTodos "Name")> {
}
*/

/*
// define a template func
<define> func nameInTemplate(a astThing, w out, args ...string)error {

}
*/

// define new keywords
// - to open a resource
f, err := open os.Open(...)
// translate to
f, err := os.Open(...)
defer f.Close()
// - must keyword
f := must open os.Open(...)
// translate to
f, err := os.Open(...)
if err != nil {
  ... // this depends of the context,
  // if the method does not return arguments => panic
  // if the method returns arguments => return ...
}
defer f.Close()

// - the ... keyword in return ...err
// where it declares default value for
// out values with ..., insert provided err appropriately
func p() (int, bool, *A, B){
  return ...
}
// translate to
func p() (int, bool, *A, B){
  return 0, false, nil, B{}
}

// - the or keyword in must call() or smting
f := must os.Open(...) or panic
// translate to
f, err := os.Open(...)
if err != nil {
  panic(err)
}
//-
f, err := must os.Open(...) or return ...err
// translate to
f, err := os.Open(...)
if err != nil {
  return err
}
//- going further
f, err := must open os.Open(...) or
  return ...log.log("file", "open", err)

// translate to
f, err := os.Open(...)
if err != nil {
  return log.log("file", "open", err)
}
defer f.Close()

// - interface combiner
func p(value SomePusher+Committer) {}

// which would translate into a new interface declaration such as
type genXXXer interface {
  SomePusher
  Committer
}
// and translate the func signature to
func p(value genXXXer) {}

// - const name
const (
	numberToken lexer.TokenType = iota
	wsToken
)
var xName = constname(numberToken)

// translate to
const (
	numberToken lexer.TokenType = iota
	wsToken
)
var xName = "numberToken"
