# Numbers

Golog tries the Ivy type arbitrary precision numbers using math/big.
This seems like a good idea to me. Will try and do similar.

# lists
prolog lists are lisp style cons pairs BUT, the WAM has optimisations
for flattening lists for more efficient storage.

At the moment the parser produces a series of cons/2 terms (using
cons/0 for [] ), Originall it just produced, flat lists, but representing
tail seemed tricky there.

We can always flatten them again later for the WAM if needed

# To WAM or not to WAM
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

