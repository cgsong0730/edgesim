package experiment

import (
	"edgesim/simulation"
	"runtime"
	"sync"
)

func FirstExperiment() {

	// 60 * 60 * 24 -> 86400, 3600
	numOfSim := 6

	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := new(sync.WaitGroup)
	wg.Add(numOfSim) // num of sim

	go func() {
		simulation.ImagePullingSimulationWithFile("mkrp", 5, 50, true, true)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("mkrp", 10, 50, true, true)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("mkrp", 15, 50, true, true)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("comm", 5, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("comm", 10, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("comm", 15, 50, false, false)
		wg.Done()
	}()

	wg.Wait()
}

func SecondExperiment() {

	// 60 * 60 * 24 -> 86400, 3600
	numOfSim := 10

	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := new(sync.WaitGroup)
	wg.Add(numOfSim) // num of sim

	go func() {
		simulation.ImagePullingSimulationWithRandomGraph("mkrp", 3, 10, true, true)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithRandomGraph("mkrp", 3, 20, true, true)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithRandomGraph("mkrp", 3, 30, true, true)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithRandomGraph("mkrp", 3, 40, true, true)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithRandomGraph("mkrp", 3, 50, true, true)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithRandomGraph("comm", 3, 10, true, true)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithRandomGraph("comm", 3, 20, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithRandomGraph("comm", 3, 30, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithRandomGraph("comm", 3, 40, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithRandomGraph("comm", 3, 50, false, false)
		wg.Done()
	}()

	wg.Wait()
}

func ThirdExperiment() {

	// 60 * 60 * 24 -> 86400, 3600
	numOfSim := 23

	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := new(sync.WaitGroup)
	wg.Add(numOfSim) // num of sim

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 3, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 4, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 5, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 6, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 7, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 8, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 9, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 10, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 11, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 12, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 13, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 14, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 15, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 16, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 17, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 18, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 19, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 20, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 21, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 22, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 23, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 24, 50, false, false)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithFile("snrp", 25, 50, false, false)
		wg.Done()
	}()

	wg.Wait()
}
