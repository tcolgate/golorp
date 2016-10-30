# Much TODO:

## Lexer

- Fix float input to match standard
- Fix string escapes to match standard
- Implement better syntax error handling and reporting.

## Parser

- Implement better syntax error handling and reporting.
- Update ops table as op directives are seen

## Everythig else

- Implement something
- Intrinsics (built in terms, e.g. cons, unify etc.)

# Numbers

Golog tries the Ivy type arbitrary precision numbers using math/big.
This seems like a good idea to me. Will try and do similar.

# Lists

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

- It seems like an interesting task
- Serializability, we can potentially save the state
- Transpiling, we might be able to compile WAM to go, which
  could be fun.

Disadvantages:

- Probably overkill for what I actually wanted to achieve (a simple
  interpreter)
- Potentially much more complicated than a straight unification
  algorithm (I think golog does that in very little code
- Manual garbage colleciton?

# Blue Sky

I'll likely get bored before doing any of this but.

- Module system taken from Go's approachi
  - `golorp get`
  - G`OLORPPATH,` `golorp run`
  - `golorp build` compile to a binary, via go (equiv to swipl -c
- `golorp fmt`
- cut
- spawn/1 , using go routines.
