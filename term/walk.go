package term

// WalkDepthFirst performs a depth first pre-order, left-to-right walk,
// calling function tf with each node.
func WalkDepthFirst(tf func(Term), t Term) {
	tf(t)
	switch st := t.(type) {
	case *Callable:
		for _, at := range st.Args() {
			WalkDepthFirst(tf, at)
		}
	case Variable:
	default:
	}
}
