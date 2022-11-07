package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Define a struct for you collector that contains pointers
// to prometheus descriptors for each metric you wish to expose.
// Note you can also include fields of other types if they provide utility
// but we just won't be exposing them as metrics.
type powerballCollector struct {
	powerballMetric *prometheus.Desc
}

// You must create a constructor for you collector that
// initializes every descriptor and returns a pointer to the collector
func newPowerballCollector() *powerballCollector {
	return &powerballCollector{
		powerballMetric: prometheus.NewDesc("powerball_metric",
			"Rolls you some new numbers for winningz",
			[]string{"ball1", "ball2", "ball3", "ball4", "ball5", "pb"}, nil,
		),
	}
}

// Each and every collector must implement the Describe function.
// It essentially writes all descriptors to the prometheus desc channel.
func (collector *powerballCollector) Describe(ch chan<- *prometheus.Desc) {

	//Update this section with the each metric you create for a given collector
	ch <- collector.powerballMetric
}

func duplicateRollCheck(rollIn int, roll1 int, roll2 int, roll3 int, roll4 int, roll5 int) bool {
	if rollIn == roll1 || rollIn == roll2 || rollIn == roll3 || rollIn == roll4 || rollIn == roll5 {
		return true
	}
	return false
}

func roll(roll1 int, roll2 int, roll3 int, roll4 int, roll5 int) int {

	time.Sleep(time.Millisecond * 100) // sleeping adds randomness to seed
	rand.Seed(time.Now().UnixNano())

	min := 1
	max := 69
	roll := rand.Intn(max-min) + min

	//Roll with loop
	if roll1 == 0 {
		//roll ball 1
		roll = rand.Intn(max-min) + min
	} else if roll2 == 0 {
		//roll ball 2
		for {
			roll = rand.Intn(max-min) + min
			if duplicateRollCheck(roll, roll1, roll2, roll3, roll4, roll5) == false {
				return roll
			} else {
				log.Warn("Duplicate found!, rerolling...")
			}
		}

	} else if roll3 == 0 {
		//roll ball 3
		for {
			roll = rand.Intn(max-min) + min
			if duplicateRollCheck(roll, roll1, roll2, roll3, roll4, roll5) == false {
				return roll
			} else {
				log.Warn("Duplicate found!, rerolling...")
			}
		}
	} else if roll4 == 0 {
		//roll ball 4
		for {
			roll = rand.Intn(max-min) + min
			if duplicateRollCheck(roll, roll1, roll2, roll3, roll4, roll5) == false {
				return roll
			} else {
				log.Warn("Duplicate found!, rerolling...")
			}
		}
	} else {
		//rolling ball 5
		for {
			roll = rand.Intn(max-min) + min
			if duplicateRollCheck(roll, roll1, roll2, roll3, roll4, roll5) == false {
				return roll
			} else {
				fmt.Println("Duplicate found! Rerolling...")
			}
		}
	}

	return roll
}

// Collect implements required collect function for all prometheus collectors
func (collector *powerballCollector) Collect(ch chan<- prometheus.Metric) {
	//1-69 //1-26
	//choose 5 numbers then choose one powerball

	sum := 0
	for i := 1; i < 20; i++ {
		sum += i

		ball_roll1 := roll(0, 0, 0, 0, 0)
		ball_roll2 := roll(ball_roll1, 0, 0, 0, 0)
		ball_roll3 := roll(ball_roll1, ball_roll2, 0, 0, 0)
		ball_roll4 := roll(ball_roll1, ball_roll2, ball_roll3, 0, 0)
		ball_roll5 := roll(ball_roll1, ball_roll2, ball_roll3, ball_roll4, 0)
		sortRolls := []int{ball_roll1, ball_roll2, ball_roll3, ball_roll4, ball_roll5}

		sort.Ints(sortRolls)
		//fmt.Println(sortRolls)

		ball_roll1 = sortRolls[0]
		ball_roll2 = sortRolls[1]
		ball_roll3 = sortRolls[2]
		ball_roll4 = sortRolls[3]
		ball_roll5 = sortRolls[4]

		numbers := &PowerBall{
			ball1: ball_roll1,
			ball2: ball_roll2,
			ball3: ball_roll3,
			ball4: ball_roll4,
			ball5: ball_roll5,
			pb:    (rand.Intn(26-1) + 1),
		}

		//Write latest value for each metric in the prometheus metric channel.
		//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
		m1 := prometheus.MustNewConstMetric(collector.powerballMetric, prometheus.GaugeValue, float64(i), strconv.Itoa(numbers.ball1), strconv.Itoa(numbers.ball2), strconv.Itoa(numbers.ball3), strconv.Itoa(numbers.ball4), strconv.Itoa(numbers.ball5), strconv.Itoa(numbers.pb))
		ch <- m1
		//fmt.Println(numbers)

	}

}

func main() {
	powerball := newPowerballCollector()
	prometheus.MustRegister(powerball)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":9101", nil))
}

type PowerBall struct {
	ball1 int
	ball2 int
	ball3 int
	ball4 int
	ball5 int
	pb    int
}
