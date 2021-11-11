package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"
)

type Philo struct {
	time_to_dead  int
	time_to_sleep int
	time_to_eat   int
	leftFork      *sync.Mutex
	rightFork     *sync.Mutex
	number        int
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

/*func monitor(philo []Philo) {

}*/

func philoLiveCycle(ch chan Philo, philo *Philo) {
	if philo.number/2 == 0 {
		time.Sleep(1 * time.Second)
	}
	for {

		philo.rightFork.Lock()
		philo.leftFork.Lock()
		fmt.Println("Philosopher ", philo.number, "is sleeping")
		time.Sleep(time.Millisecond * time.Duration(philo.time_to_eat))
		philo.rightFork.Unlock()
		philo.leftFork.Unlock()
		fmt.Println("Philosopher ", philo.number, "is sleeping")
		time.Sleep(time.Millisecond * time.Duration(philo.time_to_sleep))
		fmt.Println("Philosopher ", philo.number, "is thinking")
	}
	ch <- *philo
}

func main() {
	args := os.Args[1:]
	if len(args) != 4 {
		log.Fatal("wrong numbers of arg : <philo_count>, <time_to_dead> ,<time_to_sleep> ,<time_to_eat>")
	}

	philo_count, err := strconv.Atoi(args[0])
	checkError(err)
	timeToDead, err := strconv.Atoi(args[1])
	checkError(err)
	timeToSleep, err := strconv.Atoi(args[2])
	checkError(err)
	timeToEat, err := strconv.Atoi(args[3])
	var channels = make([]chan Philo, philo_count)
	var mutex = make([]sync.Mutex, philo_count)
	var philosophers = make([]Philo, philo_count)
	for i := 0; i < philo_count; i++ {
		channels[i] = make(chan Philo)
	}
	for i := 0; i < philo_count; i++ {
		philosophers[i] = Philo{
			number:        i + 1,
			time_to_dead:  timeToDead,
			time_to_sleep: timeToSleep,
			time_to_eat:   timeToEat,
			leftFork:      &mutex[i],
		}
		if i == philo_count-1 {
			philosophers[i].rightFork = &mutex[0]
		} else {
			philosophers[i].rightFork = &mutex[i+1]
		}
	}
	for i := 0; i < philo_count; i++ {
		go philoLiveCycle(channels[i], &philosophers[i])
		fmt.Println(<-channels[i])
	}
	fmt.Println("end main()")
}
