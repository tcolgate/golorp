likes('sam
things ',orange).
likes(sam,apple).
likes('sam',ham).
likes('sam','eggs \'n ham').
likes('sam',cheese).  likes('sam',crackers).
llikes('sam',boiledham).
b/things('tristan').
===>('tristan').
then().
then2.
p([H|T],H,T).


move(1,X,Y,_) :-
  write('Move top disk from '),
  write(X),
  write(' to '),
  write(Y),
  nl.

move(N,X,Y,Z) :-
  N>1,
  M is N-1,
  move(M,X,Z,Y),
  move(1,X,Y,_),
  move(M,Z,Y,X).

my(π,λ).


