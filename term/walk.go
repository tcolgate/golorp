package term

// WalkDepthFirst performs a depth first pre-order, left-to-right walk,
// calling function tf with each node.
func WalkDepthFirst(pre, post func(Term), t Term) {
	if pre != nil {
		pre(t)
	}
	switch st := t.(type) {
	case *Callable:
		for _, at := range st.Args() {
			WalkDepthFirst(pre, post, at)
		}
	case Variable:
	default:
	}
	if post != nil {
		post(t)
	}
}
