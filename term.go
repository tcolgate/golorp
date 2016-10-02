package golorp

type Atom string

type Number int64

type Functor struct {
	Name  string
	Arity int
	Args  []Cell
}

type Cell struct {
	fun Functor
	ref Ref
}

type Ref *Cell
