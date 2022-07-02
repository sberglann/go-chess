package main

import (
	"sync"
)

var wg sync.WaitGroup

func main() {
	/*
		var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
		flag.Parse()
		if *cpuprofile != "" {
			f, err := os.Create(*cpuprofile)
			if err != nil {
				log.Fatal(err)
			}
			pprof.StartCPUProfile(f)
			defer pprof.StopCPUProfile()
		}
	*/

	StartServer()
}
