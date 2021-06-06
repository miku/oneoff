// Solution of a Problem in Concurrent Program Control
//
// Dijkstra, 1965: http://rust-class.org/static/classes/class19/dijkstra.pdf
// https://i.imgur.com/weHFuip.png
//
// Notes:
//
// * if b[i] is true, then the program is not "looping", e.g. not waiting to enter a critical section
// * if b[i] is false (Li0), the program entered "looping" stage, before the critical section
//
// k looks like a "selector", the "if" statement in Li1 like a two-stage "check"
//
//     First we want k and i be the same; which only happens, if current
//     program k is not waiting (i.e. b[k] is true).
//
//     If the program at k is not waiting, we can jump to the next branch (Li4).
//
//     Here, we want all other programs not in the critical section; also
//     setting c[i] to false to indicate, that we want to run.
//
package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"time"
)

const (
	Enter = "\u001b[32mE\u001b[0m"
	Exit  = "\u001b[31mX\u001b[0m"
	False = "\u001b[34mF\u001b[0m"
	True  = "\u001b[37mT\u001b[0m"
)

var (
	N    = flag.Int("n", 3, "number of processes")
	T    = flag.Duration("t", 1000*time.Millisecond, "exit simulation after t (e.g. 1s, 100ms, ...)")
	Cdur = flag.Duration("C", 10*time.Millisecond, "critical section duration")
	Rdur = flag.Duration("R", 100*time.Millisecond, "remainder duration")
)

func main() {
	flag.Parse()
	// The integer k will satisfy 1 <= k <= N, b[i] and c[i] will only be set
	// by the ith computer; they will be inspected by others. [...] all
	// boolean arrays mentioned set to true. The starting value of k is
	// immaterial.
	var (
		k int = 0
		b     = make([]bool, *N)
		c     = make([]bool, *N)
	)
	setTrue(b)
	setTrue(c)

	// The program for the ith computer [...]
	computer := func(i int) {
		fmt.Printf("[%d] R .. k=%v, b=%v, c=%v\n", i, k, Btoa(b), Btoa(c))
	Li0:
		fmt.Printf("[%d] Li0  k=%v, b=%v, c=%v\n", i, k, Btoa(b), Btoa(c))
		b[i] = false
	Li1:
		if k != i {
			c[i] = true
			// Li3
			if b[k] {
				k = i
			}
			goto Li1
		} else {
			// Li4
			c[i] = false
			for j := 0; j < *N; j++ {
				if j != i && !c[j] {
					goto Li1
				}
			}
		}
		// Critical section, only at most one "computer" will be in this section at any time.
		fmt.Printf("[%d] %s >> k=%v, b=%v, c=%v\n", i, Enter, k, Btoa(b), Btoa(c))
		time.Sleep(*Cdur)
		fmt.Printf("[%d] %s << k=%v, b=%v, c=%v\n", i, Exit, k, Btoa(b), Btoa(c))
		// Critical section ends.

		c[i] = true
		b[i] = true

		fmt.Printf("[%d] Rem  k=%v, b=%v, c=%v\n", i, k, Btoa(b), Btoa(c))
		time.Sleep(time.Duration(rand.Int63n(Rdur.Milliseconds())) * time.Millisecond)
		goto Li0
	}
	for i := 0; i < *N; i++ {
		go computer(i)
	}
	time.Sleep(*T)
	fmt.Println("[X] timeout")
}

func setTrue(b []bool) {
	for i := 0; i < len(b); i++ {
		b[i] = true
	}
}

func Btoa(b []bool) string {
	var buf bytes.Buffer
	for _, v := range b {
		if v {
			io.WriteString(&buf, True)
		} else {
			io.WriteString(&buf, False)
		}
	}
	return buf.String()
}

// Appendix A: The Proof
//
// We start b y observing that the solution is safe in the
// sense that no two computers can be in their critical section
// simultaneously. For the only way to enter its critical
// section is the performance of the compound statement
// Li4 without jmnping back to Lil, i.e., finding all other
// c's t r u e after having set its own e to false.
// The second part of the proof must show that no infinite
// "After you"-"After you"-blocking can occur; i.e., when
// none of the computers is in its critical section, of the
// computers looping (i.e., jumping back to Lil) at least
// one--and therefore exactly one--will be allowed to enter
// its critical section in due time.
// If the kth computer is not among the looping ones,
// bik] will be t r u e and the looping ones will all find k # i.
// As a result one or mnore of them will find in Li3 the Boolean
// b[k] t r u e and therefore one or more will decide to assign
// "k : = i". After the first assignment "k : = i", b[k] becomes false and no new computers can decide again to
// assign a new value to k. When all decided assignments to
// k have been performed,/c will point to one of the looping
// computers and will not change its value for the time being,
// i.e., until b[k] becomes t r u e , viz., until the kth computer
// has completed its critical section. As soon as the value of
// ]c does not change any more, the kth computer will wait
// (via the compound statement Li4) until all other c's are
// t r u e , but this situation will certainly arise, if not already
// present, because all other looping ones are forced to set
// their e t r u e , as they will find k # i. And this, the author
// believes, completes the proof.
//
// Appendix B: Race
// ==================
// WARNING: DATA RACE
// Write at 0x00c0000bc051 by goroutine 8:
//   main.main.func1()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:65 +0xf8e
//
// Previous read at 0x00c0000bc051 by goroutine 7:
//   main.Btoa()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:111 +0xb5
//   main.main.func1()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:84 +0x536
//
// Goroutine 8 (running) created at:
//   main.main()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:97 +0x504
//
// Goroutine 7 (running) created at:
//   main.main()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:97 +0x504
// ==================
// ==================
// WARNING: DATA RACE
// Write at 0x00c0000bc054 by goroutine 8:
//   main.main.func1()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:68 +0x364
//
// Previous read at 0x00c0000bc054 by goroutine 9:
//   main.Btoa()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:111 +0xb5
//   main.main.func1()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:64 +0xdb3
//
// Goroutine 8 (running) created at:
//   main.main()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:97 +0x504
//
// Goroutine 9 (running) created at:
//   main.main()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:97 +0x504
// ==================
// ==================
// WARNING: DATA RACE
// Read at 0x00c0000bc055 by goroutine 7:
//   main.Btoa()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:111 +0xb5
//   main.main.func1()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:86 +0x7fc
//
// Previous write at 0x00c0000bc055 by goroutine 9:
//   main.main.func1()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:68 +0x364
//
// Goroutine 7 (running) created at:
//   main.main()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:97 +0x504
//
// Goroutine 9 (running) created at:
//   main.main()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:97 +0x504
// ==================
// ==================
// WARNING: DATA RACE
// Write at 0x00c0000bc050 by goroutine 7:
//   main.main.func1()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:90 +0xa5d
//
// Previous read at 0x00c0000bc050 by goroutine 9:
//   main.main.func1()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:70 +0x3ea
//
// Goroutine 7 (running) created at:
//   main.main()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:97 +0x504
//
// Goroutine 9 (running) created at:
//   main.main()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:97 +0x504
// ==================
// ==================
// WARNING: DATA RACE
// Write at 0x00c0000bc048 by goroutine 8:
//   main.main.func1()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:71 +0x418
//
// Previous read at 0x00c0000bc048 by goroutine 7:
//   main.main.func1()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:86 +0x844
//
// Goroutine 8 (running) created at:
//   main.main()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:97 +0x504
//
// Goroutine 7 (running) created at:
//   main.main()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:97 +0x504
// ==================
// ==================
// WARNING: DATA RACE
// Write at 0x00c0000bc053 by goroutine 7:
//   main.main.func1()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:68 +0x364
//
// Previous read at 0x00c0000bc053 by goroutine 8:
//   main.main.func1()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:78 +0x4dd
//
// Goroutine 7 (running) created at:
//   main.main()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:97 +0x504
//
// Goroutine 8 (running) created at:
//   main.main()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:97 +0x504
// ==================
// ==================
// WARNING: DATA RACE
// Read at 0x00c0000bc052 by goroutine 7:
//   main.Btoa()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:111 +0xb5
//   main.main.func1()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:64 +0xd64
//
// Previous write at 0x00c0000bc052 by goroutine 9:
//   main.main.func1()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:90 +0xa5d
//
// Goroutine 7 (running) created at:
//   main.main()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:97 +0x504
//
// Goroutine 9 (running) created at:
//   main.main()
//       /home/tir/code/miku/18d13ab0c7de09120a26f8ebe153ad27/dijkstra65.go:97 +0x504
// ==================
// Found 7 data race(s)
// exit status 66
