## Robot Javelin solution

This project solves the December 2025 Jane Street Puzzle [Robot Javelin](https://www.janestreet.com/puzzles/robot-javelin-index/).

I solved this problem exactly on paper - using this project to test parts of my solution stage by stage, using Monte Carlo methods.

Therefore, there is a lot of mess in this file and it may be slightly hard to read. Some functions are also not super useful. And much of this can likely be optimised.

If you compile and run main.go as is, you will get the maximum probability of victory for the player who does not have access to the bit-spying technology, to around 5 d.p. This is not sufficient on its own for the puzzle, which requires 10 d.p. However, it provided me with enough confidence in my exact answer that I stopped there.

On paper, I broadly solved the puzzle as follows:

### Start with variant 1.

> This variant is where neither player can see anything about what the other player has done until they compare their final score at the end of the game.

> The strategy for variant 1 has to be that players will re-roll if their first roll is below a threshold, `t`.

> My initial guess here was that the Nash equilibrium occurs when the `expected value = t` - in other words, that a player re-rolls if their first roll is below their expected value. This equation yields a `t* = 0.61803 = (sqrt5 - 1 )/ 2`.

> I couldn't convince myself that this logic necessarily implies a Nash equilibrium, so derived the equations for the probability of win for a given player as a function of the threshold used by both players.

> Under these equations, if one player picks a value of `t != t*`, then the other player can always find a threshold at which their probability of winning is greater than `0.5`. However, if one player chooses `t = t*`, the other player is guaranteed to have a win probability of less than `0.5`, or equal to `0.5` if they also pick `t*`. This is the definition of a Nash equilibrium!

### Now move to variant 2.

> Now that one player (the spy) can pick a threshold and tell whether the other is above or below it on their first go, and they also assume the other player is using the Nash equilibrium as their threshold, one player has an advantage. They have an insight into whether the other player is likely to re-roll or not, and so have greater certainty over what the final distribution of the other player's score is likely to be.

> Intuition suggests that the player with this spying power should choose the threshold they spy on the other player at, `t_spy`, as equal to the threshold the other player will be using (`t*`). This is because the spy will gain the maximum certainty about what the other player is going to do on the next turn. Choosing another other value of `t_spy` means there will remain some uncertainty about whether the other player will re-roll or not. For example, having `t_spy < t*` means that, if the spy finds out the other player rolled above `t_spy`, they can't be sure whether the other player will re-roll or stick, and so the strategy taken by the spy will naturally be sub-optimal.

> However, this logic alone wasn't sufficient to convince me that this was the optimum value for `t_spy`, so I again went through some tedious algebra to calculate the probability of winning for the spy as a function of `t_spy` and of a new threshold which the spy will use to use to decide whether to re-roll or not, all under the assumption that the other player uses `t*`. This ends up showing that the optimal strategy for the spy is to choose `t_spy = t*`, and when the spying shows that the other player rolls below `t*` on their first roll, the spy should re-roll if they have a first roll below `t_down = 0.5`, and if the other player has a first roll above `t*`, they should increase their re-roll threshold to `t_up = 0.69098... = (5 - sqrt5) / 4`. This gives a probability of winning for the spy of `p = (13 - 4 * sqrt5) / 8 = 0.5069660...`.

### Finally, variant 3.

> So, the spied on has become aware they are surveilled. Knowing that their opponent now has two different strategies depending on their own first roll, the spied on player should adopt two different strategies, depending on how their first roll went.

> By this point, I had become used to doing what felt like endless algebra, and also didn't really have a great intuition for what the best strategy for the spied on player should be - so I just calculated the probability functions.

> This yields two new thresholds for the player being spied on - a `t_low = 7 / 12` and `t_high = t*`. These produce a combined probability of victory for the spy player (using them just for consistency) of `p = 5 * sqrt5 / 16 - 37 / 192 = 0.5060629...`.

> This means that the spied on player now has a better chance of winning, but it is still less than a `0.5`. In fact, their knowledge of being spied on only allows them to increase their win probability by `349/192 - 13 * sqrt5 / 16 = 0.0009031...`. 

> But this makes intuitive sense. The spy's advantage is stronger as it gains higher quality / more information about the state of the game. The spy knows with certainty what the first roll of the other player was and the other player is simply doing what they can to minimise the impact of that extra information - it is reasonable to assume that the spy's advantage could only realistically be countered if the other player gains some information on the state of the spy's rolls.

> Oh, and the solution to the puzzle is the complement of that final probability - `229/192 - 5 * sqrt5 / 16 = 0.4939370...`