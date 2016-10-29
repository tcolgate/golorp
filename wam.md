Some thoughts on implementation.

Implementing the WAMM directly seems potentially sub-optimal.
We have effective garbage collection in Go, so potentially
should be able to leverage that.

It may be that there is no actual advantage to implementing
the WAM, but some rough guesses to why we should:

- Serializability, we can potentially save the state
- Transpiling, we might be able to compile WAM to go, which
  could be fun.
- It seems like an interesting task

Disadvantages:
- Probably overkill for what I actually want to achieve (a simple
  interpreter)
- Potentially much more complicated than a straight for unification
  algorithm (I think golog does that very little code

