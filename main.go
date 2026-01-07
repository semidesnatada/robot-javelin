package main

import (
	"fmt"
	"math/rand"
	"slices"
	"sync"
	"time"
)

func main() {
    t_nash := 0.61803398875
    t_spy_optimal := t_nash
    t1_low_optimal := 0.5
    t1_high_optimal := 0.690983005625
    t2_low_optimal := 7.0 / 12.0
    t2_high_optimal := t_nash

    t_0 := time.Now()
    
    numWorkers := 8 // equivalent ot number of cores
    totalIterations := int64(10_000_000_000)
    iterationsPerWorker := totalIterations / int64(numWorkers)
    
    var v3_p1_wins, v3_p2_wins, v2_p1_wins, v2_p2_wins int64
    var wg sync.WaitGroup
    
	// we need a lot of iterations to get a high precision result from Monte Carlo
	// this function uses parallelisation (using workers on different CPU cores)
	// to speed up the process

	// we need to use local counters for each worker to avoid collisions when
	// adding results together. we could use a mutex for extra protection, but 
	// this doesn't seem necessary for this number of workers.
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()

            var local_v3_p1, local_v3_p2, local_v2_p1, local_v2_p2 int64
            
            start := int64(workerID) * iterationsPerWorker
            end := start + iterationsPerWorker
            if workerID == numWorkers-1 {
                end = totalIterations
            }
            
            for i := start; i < end; i++ {
                m := match{
                    player_1_round_1: rand.Float64(),
                    player_1_round_2: rand.Float64(),
                    player_2_round_1: rand.Float64(),
                    player_2_round_2: rand.Float64(),
                }
                
                if test_variant_3(m, t_spy_optimal, t1_low_optimal, t1_high_optimal, t2_low_optimal, t2_high_optimal) {
                    local_v3_p1++
                } else {
                    local_v3_p2++
                }
                
                if test_variant_2(m, t_spy_optimal, t1_low_optimal, t1_high_optimal, t_nash) {
                    local_v2_p1++
                } else {
                    local_v2_p2++
                }
                
            }
            
            // add local results to global counter
            v3_p1_wins += local_v3_p1
            v3_p2_wins += local_v3_p2
            v2_p1_wins += local_v2_p1
            v2_p2_wins += local_v2_p2
        }(w)
    }
    
    wg.Wait()
    
    v2_wp := float64(v2_p1_wins) / (float64(v2_p1_wins) + float64(v2_p2_wins))
    v3_wp := float64(v3_p1_wins) / (float64(v3_p1_wins) + float64(v3_p2_wins))
    
    fmt.Printf("under variant 2 rules, player 1 has a win probability of %.10f \n", v2_wp)
    fmt.Printf("under variant 3 rules, player 1 has a win probability of %.10f \n", v3_wp)
    fmt.Printf("  this is a difference of %f \n", v2_wp-v3_wp)

	fmt.Printf("\n Finished after %v\n", time.Since(t_0))
}

// this commented out function is the non-parallised monte carlo solution.

// func main() {

// 	// ms := generate_matches(1_000_000_000)

// 	// identify_variant_3_best_sub_thresholds(ms)

// 	t_nash := 0.61803398875
// 	t_spy_optimal := t_nash
// 	t1_low_optimal := 0.5
// 	t1_high_optimal := 0.690983005625
// 	t2_low_optimal := 7.0/12.0
// 	t2_high_optimal := t_nash

// 	// v2_wp := get_variant_2_win_probability(ms, t_spy_optimal, t1_low_optimal, t1_high_optimal, t_nash)
// 	// v3_wp := get_variant_3_win_probability(ms, t_spy_optimal, t1_low_optimal, t1_high_optimal, t2_low_optimal, t2_high_optimal)

// 	t_0 := time.Now()

// 	var v3_p1_wins, v3_p2_wins, v2_p1_wins, v2_p2_wins int
// 	for i := 0; i < 100_000_000_000; i++ {
// 		m := match{
// 			player_1_round_1: rand.Float64(),
// 			player_1_round_2: rand.Float64(),
// 			player_2_round_1: rand.Float64(),
// 			player_2_round_2: rand.Float64(),
// 		}

// 		if test_variant_3(m, t_spy_optimal, t1_low_optimal, t1_high_optimal, t2_low_optimal, t2_high_optimal) {
// 			v3_p1_wins ++
// 		} else {
// 			v3_p2_wins ++
// 		}
// 		if test_variant_2(m, t_spy_optimal, t1_low_optimal, t1_high_optimal, t_nash) {
// 			v2_p1_wins ++
// 		} else {
// 			v2_p2_wins ++
// 		}
// 		if i % 1_000_000_000 == 0 {
// 			fmt.Printf("reached %d after %v\n", i, time.Since(t_0))
// 		}
// 	}

// 	v2_wp := float64(v2_p1_wins) / (float64(v2_p1_wins) + float64(v2_p2_wins))
// 	v3_wp := float64(v3_p1_wins) / (float64(v3_p1_wins) + float64(v3_p2_wins))

// 	fmt.Printf("under variant 2 rules, player 1 has a win probability of %.10f \n", v2_wp)
// 	fmt.Printf("under variant 3 rules, player 1 has a win probability of %.10f \n", v3_wp)
// 	fmt.Printf("  this is a difference of %f \n", v2_wp-v3_wp)

// }

func identify_variant_3_best_sub_thresholds(ms []match) {

	t_nash := 0.61803398875
	t_spy_optimal := t_nash
	t1_low_optimal := 0.5
	t1_high_optimal := 0.690983005625
	t2_low_optimal := 7.0 / 12.0
	// t2_high_optimal := t_nash

	// fmt.Println()

	// for t2_low := 0.583; t2_low < 0.584; t2_low += 0.0001 {

	// 	v2_wp := get_variant_2_win_probability(ms, t_spy_optimal, t1_low_optimal, t1_high_optimal, t_nash)
	// 	v3_wp := get_variant_3_win_probability(ms, t_spy_optimal, t1_low_optimal, t1_high_optimal, t2_low, t2_high_optimal)

	// 	fmt.Printf("Player 2 benefits by %f for a t2_low of %f\n", 10_000 * (v2_wp-v3_wp), t2_low)
	// }
	// fmt.Println()

	fmt.Println()


	for t2_high := t_nash; t2_high <= 1.0; t2_high += 0.01 {

		v2_wp := get_variant_2_win_probability(ms, t_spy_optimal, t1_low_optimal, t1_high_optimal, t_nash)
		v3_wp := get_variant_3_win_probability(ms, t_spy_optimal, t1_low_optimal, t1_high_optimal, t2_low_optimal, t2_high)

		fmt.Printf("Player 2 benefits by %f for a t2_high of %f\n", 10_000 * (v2_wp-v3_wp), t2_high)
	}
	fmt.Println()
	

}

func identify_variant_2_best_sub_thresholds(ms []match) {

	t_nash := 0.61803398875
	
	fmt.Println()
	p1_above := 0.690983005625
	for p1_below := 0.498; p1_below < 0.502; p1_below += 0.0001 {
		wp := get_variant_2_win_probability(ms, t_nash, p1_below, p1_above, t_nash)
		fmt.Printf("player 1 wins with probability %f when choosing a lower t of %f with opponent at N* and spy_t = N*\n", wp, p1_below)
	}
	fmt.Println()
	
	fmt.Println()
	p1_below := 0.5
	for p1_above := 0.68500; p1_above < 0.7000; p1_above += 0.0001 {
		wp := get_variant_2_win_probability(ms, t_nash, p1_below, p1_above, t_nash)
		fmt.Printf("player 1 wins with probability %f when choosing an upper t of %f with opponent at N* and spy_t = N*\n", wp, p1_above)
	}
	fmt.Println()

}

func verify_modified_probabilities(ms []match) {
	// can use this function to set parameters for the spy threshold, the relative thresholds for p1
	// for checking against the derived probability equations

	// 0.5(sqrt(5)-1) -- nash equilibrium
	t_nash := 0.61803398875
	
	// modified game variable threshold for player 1
	t_low := 0.58821
	t_high := 0.826432
	// boundary across which player one can vary strategies
	t_spy := 0.808

	// t_spy_optimal := t_nash
	// t_low_optimal := 0.5
	// t_high_optimal := 0.690983005625

	// assume player 2 continues to use nash equilibrium
	win_prob := get_variant_2_win_probability(ms, t_spy, t_low, t_high, t_nash)

	// in the modified game, we expect the maximum of win_prob to be 0.50696601125 = (13-4*sqrt(5))/8

	fmt.Println()
	fmt.Printf("player 1 has a win chance of %f, with the following params:\n spy threshold: %f \n lower_t: %f \n upper t: %f\nand opponent playing at Nash equilibrium\n", win_prob, t_spy, t_low, t_high)
	fmt.Printf("player 2's win chance is therefore %f\n", 1-win_prob)
	fmt.Println()
}

func verify_unmodified_best_responses() {
	ms := generate_matches(10_000_000)

	for p2t := 0.4; p2t < 0.9; p2t += 0.01 {
		win_prob, best_response := identify_best_response(ms, p2t)
		fmt.Printf("Player 1 has a win percentage of %f at a best response threshold of %f against a player 2 with threshold %f\n", win_prob, best_response, p2t)
	}
}

func identify_best_response(ms []match, p2t float64) (float64, float64) {
	// returns a win probability for player 1 and their
	//  best response threshold, for a given player 2 choice
	// this has verified the derived probability equations

	var max_p, best_t float64

	// fmt.Println()
	// narrowed down this range by inspection
	for p1t := 0.0; p1t <= 1.0; p1t += 0.001 {
		win_percent := get_variant_1_win_probability(ms, p1t, p2t)
		if win_percent > max_p {
			max_p = win_percent
			best_t = p1t
		}
		// fmt.Printf("Player 1 has a win percentage of %f when playing with a threshold of %f against a Nash equilibrium playing p2\n", win_percent, p1t)
	}
	// fmt.Println()

	return max_p, best_t
}

func verify_nash_equilibrium() {
	ms := generate_matches(100_000_000)

	// 0.5(sqrt(5)-1) -- nash equilibrium
	p2t := 0.61803398875

	fmt.Println()
	// narrowed down this range by inspection
	for p1t := 0.595; p1t <= 0.635; p1t += 0.0001 {
		win_percent := get_variant_1_win_probability(ms, p1t, p2t)
		fmt.Printf("Player 1 has a win percentage of %f when playing with a threshold of %f against a Nash equilibrium playing p2\n", win_percent, p1t)
	}
	fmt.Println()
}

func get_variant_3_win_probability(ms []match, p1_spy_t, p1t_below, p1t_above, p2t_below, p2t_above float64) float64 {
	// returns the win probability for player 1 given a set of
	// simulated matches and a threshold for re-rolls for each player

	var p1_wins, p2_wins int

	for _, m := range ms {
		if test_variant_3(m, p1_spy_t, p1t_below, p1t_above, p2t_below, p2t_above) {
			p1_wins ++
		} else {
			p2_wins ++
		}
	}

	return float64(p1_wins) / (float64(p1_wins) + float64(p2_wins))

}

func get_variant_2_win_probability(ms []match, p1_spy_t, p1t_below, p1t_above, p2t float64) float64 {
	// returns the win probability for player 1 given a set of
	// simulated matches and a threshold for re-rolls for each player

	var p1_wins, p2_wins int

	for _, m := range ms {
		if test_variant_2(m, p1_spy_t, p1t_below, p1t_above, p2t ) {
			p1_wins ++
		} else {
			p2_wins ++
		}
	}

	return float64(p1_wins) / (float64(p1_wins) + float64(p2_wins))

}

func get_variant_1_win_probability(ms []match, p1t, p2t float64) float64 {
	// returns the win probability for player 1 given a set of
	// simulated matches and a threshold for re-rolls for each player

	var p1_wins, p2_wins int

	for _, m := range ms {
		if test_variant_1(m, p1t, p2t) {
			p1_wins ++
		} else {
			p2_wins ++
		}
	}

	return float64(p1_wins) / (float64(p1_wins) + float64(p2_wins))

}

func test_variant_3(m match, p1_spy_t, p1t_below, p1t_above, p2t_below, p2t_above float64) bool {
	// modified javelin game, player 1 can choose a spy-threshold 
	// and finds out whether player 2's first go was above or below this spy-threshold
	// player 1 can therefore choose two separate thresholds for deciding whether to 
	// re-roll, depending on whether p2's first go is above or below p1_spy

	// in this third variant, p2 is aware that p1 has different strategies for when
	// p2's first shot is above or below the t_Nash, so p2 now also has two different strategies.
	// p1 is not aware that p2 is doing this though, so uses its original thresholds

	// returns true if player 1 ones, false if player 2 wins

	var p1_final, p2_final float64

	// p1's decision now depends on the spying they do
	if m.player_2_round_1 < p1_spy_t {
		if m.player_1_round_1 < p1t_below {
			p1_final = m.player_1_round_2
		} else {
			p1_final = m.player_1_round_1
		}
	} else {
		if m.player_1_round_1 < p1t_above {
			p1_final = m.player_1_round_2
		} else {
			p1_final = m.player_1_round_1
		}
	}

	// p2's decision making is unchanged from variant 1
	if m.player_2_round_1 < p1_spy_t {
		if m.player_2_round_1 < p2t_below {
			p2_final = m.player_2_round_2
		} else {
			p2_final = m.player_2_round_1
		}
	} else {
		if m.player_2_round_1 < p2t_above {
			p2_final = m.player_2_round_2
		} else {
			p2_final = m.player_2_round_1
		}
	}

	return p1_final > p2_final
}

func test_variant_2(m match, p1_spy_t, p1t_below, p1t_above, p2t float64) bool {
	// modified javelin game, player 1 can choose a spy-threshold 
	// and finds out whether player 2's first go was above or below this spy-threshold
	// player 1 can therefore choose two separate thresholds for deciding whether to 
	// re-roll, depending on whether p2's first go is above or below p1_spy

	// p2t is the thresholds for re-roll for player 2
	// returns true if player 1 ones, false if player 2 wins

	var p1_final, p2_final float64

	// p1's decision now depends on the spying they do
	if m.player_2_round_1 < p1_spy_t {
		if m.player_1_round_1 < p1t_below {
			p1_final = m.player_1_round_2
		} else {
			p1_final = m.player_1_round_1
		}
	} else {
		if m.player_1_round_1 < p1t_above {
			p1_final = m.player_1_round_2
		} else {
			p1_final = m.player_1_round_1
		}
	}

	// p2's decision making is unchanged from variant 1
	if m.player_2_round_1 < p2t {
		p2_final = m.player_2_round_2
	} else {
		p2_final = m.player_2_round_1
	}

	return p1_final > p2_final
}

func test_variant_1(m match, p1t, p2t float64) bool {
	// unmodified javelin game, players do not know 
	// their opponent's strategy or result until the 
	// game is over.
	// p1t and p2t are the thresholds for re-roll for each player
	// returns true if player 1 ones, false if player 2 wins

	var p1_final, p2_final float64

	if m.player_1_round_1 < p1t {
		p1_final = m.player_1_round_2
	} else {
		p1_final = m.player_1_round_1
	}

	if m.player_2_round_1 < p2t {
		p2_final = m.player_2_round_2
	} else {
		p2_final = m.player_2_round_1
	}

	return p1_final > p2_final
}

func generate_matches(n int) []match {

	out := make([]match, n)

	for i := 0; i < n; i++ {
		out[i] = match{
			player_1_round_1: rand.Float64(),
			player_1_round_2: rand.Float64(),
			player_2_round_1: rand.Float64(),
			player_2_round_2: rand.Float64(),
		}
	}

	return out

}

type match struct {
	player_1_round_1, player_1_round_2, player_2_round_1, player_2_round_2 float64
}


func check_maximum_score() {
	scores := [][]float64{}

	var max, max_ind float64

	for i := 0.45; i <= 0.72; i += 0.001 {

		new_score := check_av_score(i, 100_000)
		if new_score > max {
			max = new_score
			max_ind = i
		}

		scores = append(scores, []float64{i,new_score})

	}

	slices.SortFunc(scores, func(a, b []float64) int {if a[1] < b[1] {return 1} else {return -1}})

	fmt.Println(scores)

	fmt.Printf("maximum score occurs with threshold %f and value %f\n", max_ind, max)
}

func check_av_score(t float64, n int) float64 {

	var tot float64

	for i := 0; i < n; i++ {

		val := rand.Float64()
		if val < t {
			val = rand.Float64()
		}
		tot += val

	}

	return tot / float64(n)

}

func check_for_nash(t float64, sample_size int) float64 {
	// for a given player2 threshold, t, returns the player 1 threshold 
	// which yields equal expected value for stick and twist for player 2.

	type round struct {
		p1_val, p1_reroll, p2_val, p2_reroll float64
	}

	rounds := []round{}

	for i := 0; i < sample_size; i++ {
		p1_val := rand.Float64()
		p1_reroll := rand.Float64()
		p2_val := rand.Float64()
		p2_reroll := rand.Float64()

		rounds = append(rounds, round{
			p1_val: p1_val,
			p1_reroll: p1_reroll,
			p2_val: p2_val,
			p2_reroll: p2_reroll,
		})
	}

	nash_map := make(map[float64]float64)

	for p1_t := 0.001; p1_t <= 1.000; p1_t += 0.0001 {

		var result float64

		for _, round := range rounds {

			var p1_choice float64
			if round.p1_val < p1_t {
				p1_choice = round.p1_reroll
			} else {
				p1_choice = round.p1_val
			}

			if round.p2_val < t {
				if round.p2_reroll > p1_choice {
					result += 1.0
				}
			} else {
				if round.p2_val > p1_choice {
					result += 1.0
				}
			}

		}

		nash_map[p1_t] = result / float64(sample_size)
	}

	// fmt.Println(nash_map)

	thresholds := []float64{}

	for key := range nash_map {
		thresholds = append(thresholds, key)

	}

	slices.SortFunc(thresholds, func(a, b float64) int {
		if nash_map[a] < nash_map[b] {
			return -1
		}
		if nash_map[a] > nash_map[b] {
			return 1
		}
		return 0
	})

	// fmt.Println(thresholds[0:5])

	return thresholds[0]

}

func test_arbitrary_threshold_comparison() {

	number_of_rounds := 1_000_000_000
	player_1_threshold := 0.616
	player_2_threshold := 0.618


	player_1_wins, player_2_wins := play_n_rounds(
		number_of_rounds,
		player_1_threshold,
		player_2_threshold,
	)

	fmt.Printf(
		"With %d rounds\n p1 threshold = %f\n p2 threshold = %f\n  p1 wins = %d\n  p2 wins = %d\n",
		number_of_rounds, player_1_threshold, player_2_threshold, player_1_wins, player_2_wins,
	)

}

func play_n_rounds(n int, thresh1, thresh2 float64) (int, int) {
	// returns number of rounds won by

	var p1_win, p2_win int

	for i := 0; i < n; i++ {
		if play_javelin_round(thresh1, thresh2) {
			p1_win += 1
		} else {
			p2_win += 1
		}
	}

	return p1_win, p2_win

}

func play_javelin_round(thresh1, thresh2 float64) bool {
	// return the expected value of stick and twist for one player

	p1_val := rand.Float64()
	p2_val := rand.Float64()

	if p1_val < thresh1 {
		p1_val = rand.Float64()
	}
	if p2_val < thresh2 {
		p2_val = rand.Float64()
	}

	return p1_val > p2_val

}