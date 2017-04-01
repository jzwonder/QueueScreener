package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Person struct {
	Age      int
	Inserted bool
}

func NewPerson(r *rand.Rand) Person {
	age := r.Intn(61) + 20 //age 20-80
	person := Person{
		Age:      age,
		Inserted: age < 60,
	}
	return person
}

type Queue []Person

func (q *Queue) Print(str string) {
	fmt.Printf("%s\t\t", str)
	for _, person := range *q {
		fmt.Printf("[%d]", person.Age)
	}
	fmt.Print("\n")
}

func (q *Queue) Enqueue(person Person) string {
	if !person.Inserted {
		age, i := person.Age, len(*q)-1
		for ; i >= 0; i-- {
			if (*q)[i].Age > age {
				break
			}
		}
		*q = append(*q, Person{})
		copy((*q)[i+2:], (*q)[i+1:])
		(*q)[i+1] = person
		return "insert"
	} else {
		*q = append(*q, person)
		return "append"
	}
}

func (q *Queue) Init(length int) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for i := 0; i < length; i++ {
		person := NewPerson(r)
		q.Enqueue(person)
	}
}

func Screen(scr chan<- bool, done <-chan struct{}) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for {
		time.Sleep(3 * time.Second)
		scr <- r.Intn(10) > 0
		select {
		case <-done:
			return
		default:
		}
	}
}

func Enqueue(enq chan<- Person, done <-chan struct{}) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for {
		randSec := 1 + r.Intn(6)
		time.Sleep(time.Duration(randSec) * time.Second)
		enq <- NewPerson(r)
		select {
		case <-done:
			return
		default:
		}
	}
}

func main() {
	scr, enq := make(chan bool), make(chan Person)
	done := make(chan struct{})
	var queue Queue
	queue.Init(10)
	queue.Print("init")
	go Enqueue(enq, done)
	go Screen(scr, done)
	var counter int
	for {
		str := ""
		select {
		case person := <-enq:
			str = queue.Enqueue(person)
		case passed := <-scr:
			if passed {
				queue = queue[1:]
				str = "pass"
			} else {
				person := queue[0]
				copy(queue[:len(queue)-1], queue[1:])
				queue[len(queue)-1] = person
				str = "requeue"
			}
			counter++
			if counter == 60 { //for running 3 min
				close(done)
				return
			}
		}
		queue.Print(str)
	}
}
