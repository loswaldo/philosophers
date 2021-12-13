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
	if philo.number/2 == 0 {
		time.Sleep(70 * time.Nanosecond)
	}
	for {

		philo.rightFork.Lock()
		printMesage("take a fork", philo.number, philo.messageMutex, philo.startSim)
		philo.leftFork.Lock()
		printMesage("take a fork", philo.number, philo.messageMutex, philo.startSim)
		printMesage("eating", philo.number, philo.messageMutex, philo.startSim)
		time.Sleep(time.Millisecond * time.Duration(philo.timeToEat))
		philo.lastEating = time.Now().UnixMilli()
		philo.rightFork.Unlock()
		printMesage("put a fork", philo.number, philo.messageMutex, philo.startSim)
		philo.leftFork.Unlock()
		printMesage("put a fork", philo.number, philo.messageMutex, philo.startSim)
		printMesage("is sleeping", philo.number, philo.messageMutex, philo.startSim)
		time.Sleep(time.Millisecond * time.Duration(philo.timeToSleep))
		printMesage("thinking", philo.number, philo.messageMutex, philo.startSim)
	}
}

func main() {
	t := time.Now().UnixMilli()
	args := os.Args[1:]
	if len(args) != 4 {
		log.Fatal("wrong numbers of arg : <philo_count>, <timeToDead> ,<timeToSleep> ,<time_to_eat>")
	}

	var wg sync.WaitGroup
	philo_count, err := strconv.Atoi(args[0])
	checkError(err)
	timeToDead, err := strconv.Atoi(args[1])
	checkError(err)
	timeToSleep, err := strconv.Atoi(args[2])
	checkError(err)
	timeToEat, err := strconv.Atoi(args[3])
	var messageMutex sync.Mutex
	var mutex = make([]sync.Mutex, philo_count)
	var philosophers = make([]Philo, philo_count)
	for i := 0; i < philo_count; i++ {
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
		if i == (philo_count - 1) {
			philosophers[i].rightFork = &mutex[0]
		} else {
			philosophers[i].rightFork = &mutex[i+1]
		}
	}
	for i := 0; i < philo_count; i++ {
		wg.Add(1)
		go philoLiveCycle(&philosophers[i])
		defer wg.Done()
	}
	go monitor(&philosophers)
	wg.Wait()
	fmt.Println("end main()")
}
