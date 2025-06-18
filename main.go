package main

import (
	"edgesim/simulation"
	"runtime"
	"sync"
)

func main() {

	// 60 * 60 * 24 -> 86400, 3600
	//numOfSim := 1

	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := new(sync.WaitGroup)
	wg.Add(10) // num of sim

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
		simulation.ImagePullingSimulationWithRandomGraph("comm", 3, 10, false, true)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithRandomGraph("comm", 3, 20, false, true)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithRandomGraph("comm", 3, 30, false, true)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithRandomGraph("comm", 3, 40, false, true)
		wg.Done()
	}()

	go func() {
		simulation.ImagePullingSimulationWithRandomGraph("comm", 3, 50, false, true)
		wg.Done()
	}()

	//go func() {
	//	// baseLine string, useAffinity bool, numOfCluster int, numOfNode int
	//	simulation.ImagePullingSimulationWithFile("comm", 10, 50, false, true)
	//	//simulation.ImagePullingSimulationWithGraph(43200, 1800, 1) // 그래프 모델을 파일(n.txt)에서 열어서 실행
	//	wg.Done()
	//}()

	//go func() {
	//	// baseLine string, useAffinity bool, numOfCluster int, numOfNode int
	//	simulation.ImagePullingSimulationWithFile("comm", 15, 50, false, true)
	//	//simulation.ImagePullingSimulationWithGraph(43200, 1800, 1) // 그래프 모델을 파일(n.txt)에서 열어서 실행
	//	wg.Done()
	//}()

	//go func() {
	//	// baseLine string, useAffinity bool, numOfCluster int, numOfNode int
	//	simulation.ImagePullingSimulationWithFile("mkrp", 5, 50, false, true)
	//	//simulation.ImagePullingSimulationWithGraph(43200, 1800, 1) // 그래프 모델을 파일(n.txt)에서 열어서 실행
	//	wg.Done()
	//}()

	//go func() {
	//	// baseLine string, useAffinity bool, numOfCluster int, numOfNode int
	//	simulation.ImagePullingSimulationWithFile("mkrp", 10, 50, false, true)
	//	//simulation.ImagePullingSimulationWithGraph(43200, 1800, 1) // 그래프 모델을 파일(n.txt)에서 열어서 실행
	//	wg.Done()
	//}()

	//go func() {
	//	// baseLine string, useAffinity bool, numOfCluster int, numOfNode int
	//	simulation.ImagePullingSimulationWithFile("mkrp", 15, 50, false, true)
	//	//simulation.ImagePullingSimulationWithGraph(43200, 1800, 1) // 그래프 모델을 파일(n.txt)에서 열어서 실행
	//	wg.Done()
	//}()

	//sim.ImagePullingSimulation(10, 5, 3, 43200, 1800, 1) // 일반적인 시뮬레이션

	//go func() {
	//	simulation.ImagePullingSimulationWithGraph(43200, 1800, 1) // 그래프 모델을 파일(n.txt)에서 열어서 실행
	//	wg.Done()
	//}()

	// Experiment 1 이미지 풀링 가시화 실험
	//sim.ImagePullingVisualization(5, 5, 2, 10000, 1000)

	// Experiment 2 - 다양한 엣지 서버 수에 따른 성능 - 평균 응답 시간

	// NumOfEdgeServer - 10
	//go func() {
	//	simulation.ImagePullingSimulation(10, 5, 3, 43200, 1800, 1)
	//	wg.Done()
	//}()
	//go func() {
	//	sim.ImagePullingSimulation(10, 5, 3, 43200, 1800, 2)
	//	wg.Done()
	//}()
	//go func() {
	//	sim.ImagePullingSimulation(10, 5, 3, 43200, 1800, 3)
	//	wg.Done()
	//}()

	// NumOfEdgeServer - 20
	//go func() {
	//	sim.ImagePullingSimulation(20, 5, 3, 43200, 1800, 1)
	//	wg.Done()
	//}()
	//go func() {
	//	sim.ImagePullingSimulation(20, 5, 3, 43200, 1800, 2)
	//	wg.Done()
	//}()
	//go func() {
	//	sim.ImagePullingSimulation(20, 5, 3, 43200, 1800, 3)
	//	wg.Done()
	//}()

	// NumOfEdgeServer - 30
	//go func() {
	//	sim.ImagePullingSimulation(30, 5, 3, 43200, 1800, 1)
	//	wg.Done()
	//}()
	//go func() {
	//	sim.ImagePullingSimulation(30, 5, 3, 43200, 1800, 2)
	//	wg.Done()
	//}()
	//go func() {
	//	sim.ImagePullingSimulation(30, 5, 3, 43200, 1800, 3)
	//	wg.Done()
	//}()

	// Experiment 3 - 다양한 서브 그룹의 수 k에 따른 성능

	// k - 5
	//go func() {
	//	sim.ImagePullingSimulation(30, 5, 5, 43200, 1800, 1)
	//	wg.Done()
	//}()
	//go func() {
	//	sim.ImagePullingSimulation(30, 5, 5, 43200, 1800, 2)
	//	wg.Done()
	//}()
	//go func() {
	//	sim.ImagePullingSimulation(30, 5, 5, 43200, 1800, 3)
	//	wg.Done()
	//}()

	// k - 10
	//go func() {
	//	sim.ImagePullingSimulation(30, 5, 10, 43200, 1800, 1)
	//	wg.Done()
	//}()
	//go func() {
	//	sim.ImagePullingSimulation(30, 5, 10, 43200, 1800, 2)
	//	wg.Done()
	//}()
	//go func() {
	//	sim.ImagePullingSimulation(30, 5, 10, 43200, 1800, 3)
	//	wg.Done()
	//}()

	// k - 15
	//go func() {
	//	sim.ImagePullingSimulation(30, 5, 15, 43200, 1800, 1)
	//	wg.Done()
	//}()
	//go func() {
	//	sim.ImagePullingSimulation(30, 5, 15, 43200, 1800, 2)
	//	wg.Done()
	//}()
	//go func() {
	//	sim.ImagePullingSimulation(30, 5, 15, 43200, 1800, 3)
	//	wg.Done()
	//}()

	wg.Wait()
}
