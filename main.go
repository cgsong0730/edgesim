package main

import sim "edgesim/simulation"

func main() {
	// 60 * 60 * 24 -> 86400, 3600

	// Example
	//sim.ImagePullingVisualization(5, 5, 2, 10000, 1000)

	// Experiment 1 - 다양한 엣지 서버 수에 따른 성능 - 평균 응답 시간

	// NumOfEdgeServer - 10
	go sim.ImagePullingSimulation(10, 5, 3, 43200, 1800, 1)
	go sim.ImagePullingSimulation(10, 5, 3, 43200, 1800, 2)
	go sim.ImagePullingSimulation(10, 5, 3, 86400, 3600, 3)

	// NumOfEdgeServer - 20
	go sim.ImagePullingSimulation(20, 5, 3, 43200, 1800, 1)
	go sim.ImagePullingSimulation(20, 5, 3, 43200, 1800, 2)
	go sim.ImagePullingSimulation(20, 5, 3, 43200, 1800, 3)

	// NumOfEdgeServer - 30
	go sim.ImagePullingSimulation(30, 5, 3, 43200, 1800, 1)
	go sim.ImagePullingSimulation(30, 5, 3, 43200, 1800, 2)
	go sim.ImagePullingSimulation(30, 5, 3, 43200, 1800, 3)

	// Experiment 2 - 다양한 서브 그룹의 수 k에 따른 성능

	// k - 5
	go sim.ImagePullingSimulation(30, 5, 5, 43200, 1800, 1)
	go sim.ImagePullingSimulation(30, 5, 5, 43200, 1800, 2)
	go sim.ImagePullingSimulation(30, 5, 5, 43200, 1800, 3)

	// k - 10
	go sim.ImagePullingSimulation(30, 5, 10, 43200, 1800, 1)
	go sim.ImagePullingSimulation(30, 5, 10, 43200, 1800, 2)
	go sim.ImagePullingSimulation(30, 5, 10, 43200, 1800, 3)

	// k - 15
	go sim.ImagePullingSimulation(30, 5, 15, 43200, 1800, 1)
	go sim.ImagePullingSimulation(30, 5, 15, 43200, 1800, 2)
	go sim.ImagePullingSimulation(30, 5, 15, 43200, 1800, 3)

}
