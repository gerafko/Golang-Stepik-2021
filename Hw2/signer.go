package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type ByString []string

func main() {

}

func ExecutePipeline(input ...job) {
	fmt.Println("ExecutePipeline")
	wg := &sync.WaitGroup{}
	ch := make([]chan interface{}, len(input)+1)
	for i := range ch {
		ch[i] = make(chan interface{}, 100)
		fmt.Println(i, " channel")
	}
	for i, f := range input {
		fmt.Println(i, "func")
		wg.Add(1)
		go funcWorkLayer(ch[i], ch[i+1], wg, f)
	}
	wg.Wait()
}

func funcWorkLayer(in, out chan interface{}, wg *sync.WaitGroup, f job) {
	f(in, out)
	close(out)
	wg.Done()
}

func SingleHash(in, out chan interface{}) {
	fmt.Println("	SingleHash")
	mu := &sync.Mutex{}
	wgS := &sync.WaitGroup{}
	for input := range in {
		if n, ok := input.(int); ok {
			wgS.Add(1)
			go func(n int, out chan interface{}, wgS *sync.WaitGroup, mu *sync.Mutex) {
				crc32Chan1 := make(chan string)
				crc32Chan2 := make(chan string)
				mu.Lock()
				md5 := DataSignerMd5(strconv.Itoa(n))
				mu.Unlock()
				go AsyncDataSignerCrc32(strconv.Itoa(n), crc32Chan1)
				go AsyncDataSignerCrc32(md5, crc32Chan2)
				parms := <-crc32Chan1 + "~" + <-crc32Chan2
				out <- parms
				wgS.Done()
			}(n, out, wgS, mu)
		}
	}
	wgS.Wait()
}

func MultiHash(in, out chan interface{}) {
	fmt.Println("	MultiHash")
	wgM := &sync.WaitGroup{}
	for line := range in {
		wgM.Add(1)
		go func(line interface{}, wgM *sync.WaitGroup, out chan interface{}) {
			if val, ok := line.(string); ok {
				wgMasync := &sync.WaitGroup{}
				combined := make([]string, 6)
				mu := &sync.Mutex{}
				for i := 0; i < 6; i++ {
					wgMasync.Add(1)
					go func(combined []string, index int, val string, wgMasync *sync.WaitGroup) {
						crc32Chan := make(chan string)
						go AsyncDataSignerCrc32(strconv.Itoa(index)+val, crc32Chan)
						mu.Lock()
						combined[index] = <-crc32Chan
						mu.Unlock()
						wgMasync.Done()
					}(combined, i, val, wgMasync)
				}
				wgMasync.Wait()
				out <- strings.Join(combined, "")
			}
			wgM.Done()
		}(line, wgM, out)
	}
	wgM.Wait()
}
func CombineResults(in, out chan interface{}) {
	finalSlice := []string{}
	for lines := range in {
		if val, ok := lines.(string); ok {
			finalSlice = append(finalSlice, val)
		}
	}
	sort.Sort(ByString(finalSlice))
	fmt.Println(strings.Join(finalSlice, `_`))
	out <- strings.Join(finalSlice, `_`)
}

func (a ByString) Len() int           { return len(a) }
func (a ByString) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByString) Less(i, j int) bool { return a[i] < a[j] }

func AsyncDataSignerCrc32(n string, out chan string) {
	out <- DataSignerCrc32(n)
}
