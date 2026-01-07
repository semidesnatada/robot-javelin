## Robot Javelin solution

This project solves the December 2025 Jane Street Puzzle [Robot Javelin](https://www.janestreet.com/puzzles/robot-javelin-index/).

I solved this problem exactly on paper - using this project to test parts of my solution stage by stage, using Monte Carlo methods.

Therefore, there is a lot of mess in this file and it may be slightly hard to read. Some functions are also not super useful. And much of this can likely be optimised.

If you compile and run main.go as is, you will get the maximum probability of victory for the player who does not have access to the bit-spying technology, to around 5 d.p. This is not sufficient on its own for the puzzle, which requires 10 d.p. However, it provided me with enough confidence in my exact answer that I stopped there.

On paper, I broadly solved the puzzle as follows:

### Start with "variant 1"

> This variant is where neither player can see anything about what the other player has done until they compare their final score at the end of the game.

> The strategy for variant 1 has to be that players will re-roll if their first roll is below a threshold, `t`.

> My initial guess here was that the Nash equilibrium occurs when the `expected value = t` - in other words, that a player re-rolls if their first roll is below their expected value. This equation yields a `t* = 0.61803 = (sqrt5 - 1 )/ 2`.

> I couldn't convince myself that this value for t necessarily implies a Nash equilibrium, so derived the equations for the probability of win for a given player as a function of the threshold used by both players. If one player picks a value of `t != t*`, then the other player can always find a threshold at which their probability of winning is greater than 0.5. However, if one player chooses `t = t*`, the other player is guaranteed to have a win probability of less than 0.5, or equal to 0.5 if they also pick `t*`. This is the definition of a Nash equilibrium!

### Now move to variant 2

> Now that the other player can 