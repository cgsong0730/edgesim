package main

import sim "edgesim/simulation"

func main() {
	// 60 * 60 * 24 -> 86400, 3600

	// Example
	sim.ImagePullingSimulation(10, 5, 3, 10000, 1000)

	// Experiment 1
	//sim.ImagePullingSimulation(10, 5, 3, 86400, 3600)
	//sim.ImagePullingSimulation(20, 5, 3, 86400, 3600)
	//sim.ImagePullingSimulation(30, 5, 3, 86400, 3600)

	// Experiment 2
	//sim.ImagePullingSimulation(30, 5, 5, 86400, 3600)
	//sim.ImagePullingSimulation(30, 5, 10, 86400, 3600)
	//sim.ImagePullingSimulation(30, 5, 15, 86400, 3600)

}
