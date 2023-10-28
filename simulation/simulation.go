package simulation

import (
	"edgesim/edge"
	graph "edgesim/graph"
	"edgesim/weightedrand"
	"fmt"
	"math/rand"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ImagePullingSimulation(numOfEdgeServer int, numOfRegistryServer int, numOfSubgroup int, numOfPulling int, managementInterval int, baseLine int) {
	var RegistryServerList []edge.RegistryServer
	var firstRegistryServer edge.EdgeRegistryServer
	var secondRegistryServer edge.EdgeRegistryServer
	edgeServerList := []edge.EdgeServer{}
	var nearestRouterOverhead []int
	var no graph.Graph // network overhead
	fineName := "./result.txt"
	var file *os.File
	sum := 0
	cnt := 0

	_, err := os.Stat(fineName)
	if os.IsNotExist(err) {
		file, _ = os.Create(fineName)
	} else {
		file, _ = os.OpenFile(fineName, os.O_APPEND|os.O_RDWR, 0755)
	}

	for i := 0; i < numOfRegistryServer; i++ {
		rs := edge.RegistryServer{i + 1, nil, rand.Intn(10) + 90}
		edge.InitRegistryServer(&rs, 100, 500, 1+i*100, 100+i*100)
		RegistryServerList = append(RegistryServerList, rs)
	}

	for i := 1; i <= numOfEdgeServer; i++ {
		edgeServer := edge.EdgeServer{
			Id:               i,
			NumOfImage:       0,
			MaxCacheSize:     5000,
			CurrentCacheSize: 0,
			LocalImages:      nil,
			RegistryServers:  RegistryServerList,
		}
		edgeServerList = append(edgeServerList, edgeServer)
		nearestRouterOverhead = append(nearestRouterOverhead, 0)
	}

	nearestRouterOverhead = graph.GenerateRandomGraphWithNR(&no, numOfEdgeServer)
	no.Nodes = no.Nodes[:len(no.Nodes)-1]
	nsubgraphs := graph.ClusterGraph(&no, numOfSubgroup)

	for _, subgraph := range nsubgraphs {
		leader := graph.ElectReaderUsingOverhead(&subgraph)
		firstRegistryServer = edge.CreateEdgeRegistryServer(edgeServerList[leader-1], leader, numOfPulling/managementInterval)
		for index, _ := range edgeServerList {
			//edge.FirstRegistry = firstRegistryServer
			for _, line := range no.Lines {
				if line.NodeA.Id == index && line.NodeB.Id == leader || line.NodeA.Id == leader && line.NodeB.Id == index {
					edgeServerList[index].NetworkOverhead = line.Val
				}
			}
			edgeServerList[index].FirstRegistry = firstRegistryServer
		}
	}

	for j := 1; j <= numOfPulling; j++ {
		var af graph.Graph // affinity
		var weightList []weightedrand.Choice[int, int]

		for i := 0; i < numOfEdgeServer; i++ {
			if i%numOfRegistryServer == 0 {
				weightList = append(weightList, weightedrand.NewChoice(i+1, 3))
			} else {
				weightList = append(weightList, weightedrand.NewChoice(i+1, 1))
			}
			chooser, _ := weightedrand.NewChooser(weightList...)
			result := chooser.Pick()

			var pullingTime = 0
			if baseLine == 1 {
				pullingTime = edge.ImagePullingB1(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])
			} else if baseLine == 2 {
				pullingTime = edge.ImagePullingB2(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])
			} else if baseLine == 3 {
				pullingTime = edge.ImagePullingB3(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])
			} else {
				pullingTime = edge.ImagePulling(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])
			}

			sum += pullingTime
			cnt++

			if j%500 == 0 {
				edge.CleanCache(&edgeServerList[i])
			}
		}

		if j%managementInterval == 0 {

			graph.GenerateAffinityGraph(&af, edgeServerList)
			asubgraphs := graph.ClusterGraph(&af, numOfSubgroup)
			for _, subgraph := range asubgraphs {
				//fmt.Println("affinity subgraph")
				leader := graph.ElectReaderUsingAffinity(&subgraph)
				secondRegistryServer = edge.CreateEdgeRegistryServer(edgeServerList[leader-1], leader, numOfPulling/managementInterval)
				for index, _ := range edgeServerList {
					for _, line := range af.Lines {
						if line.NodeA.Id == index && line.NodeB.Id == leader || line.NodeA.Id == leader && line.NodeB.Id == index {
							edgeServerList[index].AffinityOverhead = line.Val
						}
					}
					edgeServerList[index].SecondRegistry = secondRegistryServer
				}
			}

			if j%managementInterval*3 == 0 {
				for _, edge := range edgeServerList {
					edge.History = nil
				}
			}
		}
	}
	str := fmt.Sprintf("%d %d %d %d\n", numOfEdgeServer, numOfSubgroup, baseLine, sum/cnt)
	b := []byte(str)
	_, err = file.Write(b)
	check(err)
}

func ImagePullingVisualization(numOfEdgeServer int, numOfRegistryServer int, numOfSubgroup int, numOfPulling int, managementInterval int) {
	var RegistryServerList []edge.RegistryServer
	var firstRegistryServer edge.EdgeRegistryServer
	var secondRegistryServer edge.EdgeRegistryServer
	edgeServerList := []edge.EdgeServer{}
	var nearestRouterOverhead []int
	var no graph.Graph // network overhead
	fileList := []*os.File{}

	for i := 1; i <= numOfEdgeServer; i++ {
		fileName := fmt.Sprintf("./result%d.txt", i)
		var file *os.File

		_, err := os.Stat(fileName)
		if os.IsNotExist(err) {
			file, _ = os.Create(fileName)
		} else {
			file, _ = os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0755)
		}

		fileList = append(fileList, file)
		defer file.Close()
	}

	for i := 0; i < numOfRegistryServer; i++ {
		rs := edge.RegistryServer{i + 1, nil, rand.Intn(10) + 90}
		edge.InitRegistryServer(&rs, 100, 500, 1+i*100, 100+i*100)
		RegistryServerList = append(RegistryServerList, rs)
	}

	for i := 1; i <= numOfEdgeServer; i++ {
		edgeServer := edge.EdgeServer{
			Id:               i,
			NumOfImage:       0,
			MaxCacheSize:     5000,
			CurrentCacheSize: 0,
			LocalImages:      nil,
			RegistryServers:  RegistryServerList,
		}
		edgeServerList = append(edgeServerList, edgeServer)
		nearestRouterOverhead = append(nearestRouterOverhead, 0)
	}

	nearestRouterOverhead = graph.GenerateRandomGraphWithNR(&no, numOfEdgeServer)
	fmt.Println("network overhead graph")
	graph.PrintNetworkGraph(&no)

	no.Nodes = no.Nodes[:len(no.Nodes)-1]
	nsubgraphs := graph.ClusterGraph(&no, numOfSubgroup)

	for _, subgraph := range nsubgraphs {
		fmt.Println("network overhead subgraph")
		leader := graph.ElectReaderUsingOverhead(&subgraph)
		graph.PrintGraphUsingReader(&subgraph, leader)
		firstRegistryServer = edge.CreateEdgeRegistryServer(edgeServerList[leader-1], leader, numOfPulling/managementInterval)
		for index, _ := range edgeServerList {
			//edge.FirstRegistry = firstRegistryServer
			for _, line := range no.Lines {
				if line.NodeA.Id == index && line.NodeB.Id == leader || line.NodeA.Id == leader && line.NodeB.Id == index {
					edgeServerList[index].NetworkOverhead = line.Val
				}
			}
			edgeServerList[index].FirstRegistry = firstRegistryServer
		}
	}

	for j := 1; j <= numOfPulling; j += 1 {
		var af graph.Graph // affinity
		var weightList []weightedrand.Choice[int, int]

		for i := 0; i < numOfEdgeServer; i++ {
			if i%numOfRegistryServer == 0 {
				weightList = append(weightList, weightedrand.NewChoice(i+1, 3))
			} else {
				weightList = append(weightList, weightedrand.NewChoice(i+1, 1))
			}
			chooser, _ := weightedrand.NewChooser(weightList...)
			result := chooser.Pick()

			pullingTime := edge.ImagePulling(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])

			if j%10 == 0 {
				str := fmt.Sprintf("%d %d\n", j, pullingTime)
				b := []byte(str)
				_, err := fileList[i].Write(b)
				check(err)
			}

			if j%500 == 0 {
				edge.CleanCache(&edgeServerList[i])
			}
		}
		if j%managementInterval == 0 {
			graph.GenerateAffinityGraph(&af, edgeServerList)
			asubgraphs := graph.ClusterGraph(&af, numOfSubgroup)
			fmt.Println("affinity graph")
			graph.PrintAffinityOverallGraph(&af)
			for _, subgraph := range asubgraphs {
				fmt.Println("affinity subgraph")
				leader := graph.ElectReaderUsingAffinity(&subgraph)
				graph.PrintAffinityGraph(&subgraph, leader)
				secondRegistryServer = edge.CreateEdgeRegistryServer(edgeServerList[leader-1], leader, numOfPulling/managementInterval)
				for index, _ := range edgeServerList {
					for _, line := range af.Lines {
						if line.NodeA.Id == index && line.NodeB.Id == leader || line.NodeA.Id == leader && line.NodeB.Id == index {
							edgeServerList[index].AffinityOverhead = line.Val
						}
					}
					edgeServerList[index].SecondRegistry = secondRegistryServer
				}
			}
		}

		// 주기적으로 edge server의 history를 정리
		if j%managementInterval*3 == 0 {
			for _, edge := range edgeServerList {
				edge.History = nil
			}
		}
	}
}
