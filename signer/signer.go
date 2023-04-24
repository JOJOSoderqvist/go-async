package main

import (
	"sort"
	"strconv"
	"sync"
)

func ExecutePipeline(jobs ...job) {
	var a = make(chan interface{})
	a <- 10
	a <- "as"

	for _, j := range jobs {

	}
}

func SingleHash(in, out chan interface{}) {
	var md5Wg sync.WaitGroup
	var crc32Wg sync.WaitGroup
	var md5OutputChan = make(chan interface{}, len(in))
	var crc32OutputChan = make(chan interface{}, len(in))
	var finalChan = make(chan interface{}, len(in))
	data := <-in

	//init crc23Calc
	crc32Wg.Add(1)
	go crc32Calc(&crc32Wg, data.(string), crc32OutputChan)
	crc32Wg.Wait()

	// Calculating md5
	md5Wg.Add(1)
	go md5Calc(&md5Wg, data.(string), md5OutputChan)
	md5Wg.Wait()

	crc32Wg.Add(1)
	for value := range md5OutputChan {
		go crc32Calc(&crc32Wg, value.(string), finalChan)
	}
	crc32Wg.Wait()

	// final output
	for i := 0; i < len(crc32OutputChan); i++ {
		crc32Value := <-crc32OutputChan
		convValue := <-finalChan
		out <- crc32Value.(string) + "~" + convValue.(string)
	}
}

func md5Calc(wg *sync.WaitGroup, data string, out chan<- interface{}) {
	out <- DataSignerMd5(data)
	wg.Done()
}

func crc32Calc(wg *sync.WaitGroup, data string, out chan<- interface{}) {
	out <- DataSignerCrc32(data)
	wg.Done()
}

func MultiHash(in, out chan interface{}) {
	var toCalc, result string
	var tempChan = make(chan interface{}, len(in))
	var wg sync.WaitGroup
	wg.Add(len(in))
	for i := 0; i < 6; i++ {
		inputData := <-in
		toCalc += strconv.Itoa(i) + inputData.(string)
		crc32Calc(&wg, toCalc, tempChan)
	}
	wg.Wait()
	for value := range tempChan {
		result += value.(string)
	}
	out <- result
}

func CombineResults(in, out chan interface{}) {
	var result []int
	for value := range in {
		result = append(result, value.(int))
	}
	sort.Ints(result)
	for i := 0; i < len(result)-1; i++ {
		out <- strconv.Itoa(result[i]) + "-"
	}
	out <- strconv.Itoa(result[len(result)])
}
