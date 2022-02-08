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
	timeToDead   int
	timeToSleep  int
	timeToEat    int
	lastEating   int64
	startSim     int64
	leftFork     *sync.Mutex
	rightFork    *sync.Mutex
	number       int
	messageMutex *sync.Mutex
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func printMesage(messageType string, number int, messageMutex *sync.Mutex, startSim int64) {
	messageMutex.Lock()
	fmt.Println(time.Now().UnixMilli()-startSim, "Philosopher", number, "is", messageType)
	if messageType == "dead" {
		os.Exit(1)
	}
	messageMutex.Unlock()

}

func monitor(philo *[]Philo) {
	for {
		for _, value := range *philo {
			if (time.Now().UnixMilli() - value.lastEating) > int64(time.Duration(value.timeToDead)) {
				printMesage("dead", value.number, value.messageMutex, value.startSim)
			}
		}

	}
}

func philoLiveCycle(philo *Philo) {
	for {
		var firstFork *sync.Mutex
		var secondFork *sync.Mutex
		if philo.number/2 == 0 {
			firstFork = philo.leftFork
			secondFork = philo.rightFork
		} else {
			firstFork = philo.rightFork
			secondFork = philo.leftFork
		}

		firstFork.Lock()
		printMesage("take a fork", philo.number, philo.messageMutex, philo.startSim)
		secondFork.Lock()
		printMesage("take a fork", philo.number, philo.messageMutex, philo.startSim)
		printMesage("eating", philo.number, philo.messageMutex, philo.startSim)
		philo.lastEating = time.Now().UnixMilli()
		time.Sleep(time.Millisecond * time.Duration(philo.timeToEat))
		firstFork.Unlock()
		secondFork.Unlock()
		printMesage("is sleeping", philo.number, philo.messageMutex, philo.startSim)
		time.Sleep(time.Millisecond * time.Duration(philo.timeToSleep))
		printMesage("thinking", philo.number, philo.messageMutex, philo.startSim)
	}
}

func main() {
	t := time.Now().UnixMilli()
	args := os.Args[1:]
	if len(args) != 4 {
		log.Fatal("wrong numbers of arg : <philoCount>, <timeToDead> ,<timeToSleep> ,<time_to_eat>")
	}

	var wg sync.WaitGroup
	philoCount, err := strconv.Atoi(args[0])
	checkError(err)
	timeToDead, err := strconv.Atoi(args[1])
	checkError(err)
	timeToSleep, err := strconv.Atoi(args[2])
	checkError(err)
	timeToEat, err := strconv.Atoi(args[3])
	var messageMutex sync.Mutex
	var mutex = make([]sync.Mutex, philoCount)
	var philosophers = make([]Philo, philoCount)
	for i := 0; i < philoCount; i++ {
		philosophers[i] = Philo{
			number:       i + 1,
			timeToDead:   timeToDead,
			timeToSleep:  timeToSleep,
			timeToEat:    timeToEat,
			leftFork:     &mutex[i],
			messageMutex: &messageMutex,
			startSim:     t,
			lastEating:   t,
		}
		if i == (philoCount - 1) {
			philosophers[i].rightFork = &mutex[0]
		} else {
			philosophers[i].rightFork = &mutex[i+1]
		}
	}
	for i := 0; i < philoCount; i++ {
		wg.Add(1)
		go philoLiveCycle(&philosophers[i])
		defer wg.Done()
	}
	go monitor(&philosophers)
	wg.Wait()
	fmt.Println("end main()")
}
