package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/petermattis/goid"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Info struct {
	TestA string `json:"-"`
	TestB string `json:"test_b"`
	TestC int    `json:"test_c"`
}

type Info2 struct {
	Name string `json:"test_b"`
	Age  int    `json:"test_c"`
}

func twoSum(nums [5]int, target int) []int {
	result := []int{}
	m := make(map[int]int)
	for i, k := range nums {
		if value, exist := m[target-k]; exist {
			result = append(result, value)
			result = append(result, i)
		}
		m[k] = i
	}
	return result
}

func getWorker(waitCh chan int, symbol int, wg *sync.WaitGroup) (next chan int) {
	notify := make(chan int)
	wg.Add(1)
	go func(waitCh chan int) {
		defer func() {
			wg.Done()
		}()
		for d := range waitCh {
			if d >= 30 {
				break
			}
			fmt.Println("goroutine:", symbol, "print", d+1)
			notify <- d + 1
		}
		close(notify)
		fmt.Println("goroutine: finish", symbol)
	}(waitCh)
	return notify
}
func main3() {
	var counter Counter
	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 100000; j++ {
				counter.Incr()
			}
		}()
	}
	wg.Wait()
	fmt.Println(counter.count)
}

func main() {
	timer := time.NewTimer(5 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	resultChan := make(chan bool)
	slice := []int{1, 2, 3, 10, 999, 8, 345, 7, 98, 33, 66, 77, 88, 68, 96}
	size := 3
	sliceLen := len(slice)
	target := 345
	//开启goroutine
	for i := 0; i < sliceLen; i += size {
		end := i + size
		if end >= sliceLen  {
			fmt.Println(sliceLen - 1)
			end = sliceLen -1
		}
		key := i
		fmt.Println(i, end)
		fmt.Println("cap->", cap(slice), "len->" ,len(slice))
		//go SearchTarget(ctx, slice[key:end], target, resultChan)
		go func(ctx context.Context,target int, slice []int, resultchan  chan bool) {
			fmt.Println("key->", key, "len->", end)
			SearchTarget(ctx, slice[key:end], target, resultChan)
		}(ctx, target, slice[key:end], resultChan)
	}

	//处理结果
	select {
	case <-timer.C:
		fmt.Println("time out")
		cancel()
	case <-resultChan:
		fmt.Println("found it !")
		cancel()
	}
	time.Sleep(time.Second * 2)
}

func SearchTarget(ctx context.Context, data []int, target int, resultChan chan bool)  {
	for _, v := range  data {
		select {
		case  <- ctx.Done():
			fmt.Println( "Task cancelded! \n")
			return
		default:
			fmt.Println( "v->", v)
			time.Sleep(time.Millisecond * 1500)
			if target == v {
				resultChan <- true
				return
			}
		}
	}
}

func main5() {
	str := "hello word"
	fmt.Println(str[1])
	var slice = make([]int, 0, 10)
	slice = append(slice, 1, 2, 3)
	work(slice)
	slice = append(slice, 200)
	fmt.Println(slice)
	time.Sleep(3 * time.Second)
}

func work(a []int) {
	a = append(a, 100)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(ss []int) {
		time.Sleep(2 * time.Second)
		ss = append(ss, 300)
		fmt.Println("s->", ss)
		wg.Done()
	}(a)
	wg.Wait()
}

func main4() {
	s := []string{"a", "b", "c"}
	fmt.Printf("main 内存地址：%p\n", s)
	hello(s)
	fmt.Println(s)

}

func hello(s []string) {
	fmt.Printf("hello 内存地址：%p\n", s)
	s[0] = "j"
}

func reverString(s string) string {
	str := []rune(s)
	l := len(str)
	for i := 0; i < l/2; i++ {
		str[i], str[l-1-i] = str[l-1-i], str[i]
	}
	return string(str)
}

func f() (r int) {
	r = 1
	func(r int) {
		r = r + 5
	}(r)
	return
}

func f3() (r int) {
	defer func(r *int) {
		*r = *r + 5
	}(&r)
	return 1
}

func foo(l sync.Mutex) {
	fmt.Println("in foo")
	l.Lock()
	fmt.Println(GoId())
	bar(l)
	l.Unlock()
}
func bar(l sync.Mutex) {
	l.Lock()
	fmt.Println(GoId())
	fmt.Println("in bar")
	l.Unlock()
}

type RecusiveMutex struct {
	sync.Mutex
	owner     int64
	recursion int32
}

func (m *RecusiveMutex) Lock() {
	gid := goid.Get()

	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	atomic.StoreInt64(&m.owner, gid)
	m.recursion = 1
}

func (m *RecusiveMutex) Unlock() {
	gid := goid.Get()
	if (atomic.LoadInt64(&m.owner)) != gid {
		panic(fmt.Sprintf("wrong the owner (%d) : %d!", m.owner, gid))
	}
	m.recursion--
	if m.recursion != 0 {
		return
	}
	atomic.StoreInt64(&m.owner, -1)
	m.Mutex.Unlock()
}

func GoId() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idFiled := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine"))[0]
	id, err := strconv.Atoi(idFiled)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id: %v", err))
	}
	return id
}

type C struct {
	sync.Mutex
	Count int
}

type Counter struct {
	CounterType int

	Mu sync.Mutex

	count int
}

func (c *Counter) Incr() {
	c.Mu.Lock()
	c.count++
	c.Mu.Unlock()
}

func (c *Counter) Count() int {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	return c.count
}

func main2() {
	wg := sync.WaitGroup{}
	start := make(chan int)
	lastCh := start
	for i := 0; i < 3; i++ {
		lastCh = getWorker(lastCh, i+1, &wg)
	}
	start <- 0
	for v := range lastCh {
		start <- v
	}
	close(start)
	wg.Wait()
}

func main1() {
	s := "Hello World!"
	b := []byte(s)

	sEnc := base64.StdEncoding.EncodeToString(b)
	fmt.Printf("enc=[%s]\n", sEnc)

	sDec, err := base64.StdEncoding.DecodeString(sEnc)
	if err != nil {
		fmt.Printf("base64 decode failure, error=[%v]\n", err)
	} else {
		b := string(sDec)
		fmt.Println(b)
	}
	panic(123)

	wg := sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	wg.Add(1)
	go handle(ctx, 1500*time.Millisecond, &wg)
	select {
	case <-ctx.Done():
		fmt.Println("main", ctx.Err())
	}
	wg.Wait()
	fmt.Println("main done")
}

func handle(ctx context.Context, duration time.Duration, wg *sync.WaitGroup) {
	select {
	case <-ctx.Done():
		fmt.Println("handle:", ctx.Err())
		wg.Done()
	case <-time.After(duration):
		fmt.Println("process request with", duration)
		wg.Done()
	}
}

func slic(s *[]int) *[]int {
	news := *s
	news[0] = 100
	news = append(news, 3, 4, 5, 6)
	return &news
}
