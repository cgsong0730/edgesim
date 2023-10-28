package main

import (
	sim "edgesim/simulation"
	"runtime"
	"sync"
)

func main() {
	// 60 * 60 * 24 -> 86400, 3600

	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := new(sync.WaitGroup)
	wg.Add(18)

	// Example
	//sim.ImagePullingVisualization(5, 5, 2, 10000, 1000)

	// Experiment 1 - 다양한 엣지 서버 수에 따른 성능 - 평균 응답 시간

	// NumOfEdgeServer - 10
	go func() {
		sim.ImagePullingSimulation(10, 5, 3, 43200, 1800, 1)
		wg.Done()
	}()
	go func() {
		sim.ImagePullingSimulation(10, 5, 3, 43200, 1800, 2)
		wg.Done()
	}()
	go func() {
		sim.ImagePullingSimulation(10, 5, 3, 86400, 3600, 3)
		wg.Done()
	}()

	// NumOfEdgeServer - 20
	go func() {
		sim.ImagePullingSimulation(20, 5, 3, 43200, 1800, 1)
		wg.Done()
	}()
	go func() {
		sim.ImagePullingSimulation(20, 5, 3, 43200, 1800, 2)
		wg.Done()
	}()
	go func() {
		sim.ImagePullingSimulation(20, 5, 3, 43200, 1800, 3)
		wg.Done()
	}()

	// NumOfEdgeServer - 30
	go func() {
		sim.ImagePullingSimulation(30, 5, 3, 43200, 1800, 1)
		wg.Done()
	}()
	go func() {
		sim.ImagePullingSimulation(30, 5, 3, 43200, 1800, 2)
		wg.Done()
	}()
	go func() {
		sim.ImagePullingSimulation(30, 5, 3, 43200, 1800, 3)
		wg.Done()
	}()

	// Experiment 2 - 다양한 서브 그룹의 수 k에 따른 성능

	// k - 5
	//sim.ImagePullingSimulation(30, 5, 5, 43200, 1800, 1)
	//sim.ImagePullingSimulation(30, 5, 5, 43200, 1800, 2)
	//sim.ImagePullingSimulation(30, 5, 5, 43200, 1800, 3)

	// k - 10
	//sim.ImagePullingSimulation(30, 5, 10, 43200, 1800, 1)
	//sim.ImagePullingSimulation(30, 5, 10, 43200, 1800, 2)
	//sim.ImagePullingSimulation(30, 5, 10, 43200, 1800, 3)

	// k - 15
	//sim.ImagePullingSimulation(30, 5, 15, 43200, 1800, 1)
	//sim.ImagePullingSimulation(30, 5, 15, 43200, 1800, 2)
	//sim.ImagePullingSimulation(30, 5, 15, 43200, 1800, 3)

	wg.Wait()
}
