package simulation

import (
	"fmt"
	"gopro1/edge"
	graph "gopro1/graph_cgs"
	"gopro1/weightedrand"
	"math/rand"
)

func ImagePullingSimulation(numOfEdgeServer int, numOfRegistryServer int, numOfPulling int, managementInterval int) {
	var RegistryServerList []edge.RegistryServer
	var firstRegistryServer edge.EdgeRegistryServer
	var secondRegistryServer edge.EdgeRegistryServer
	edgeServerList := []edge.EdgeServer{}

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
	}

	for j := 1; j <= numOfPulling; j++ {
		var no graph.Graph // network overhead
		var af graph.Graph // affinity

		for i := 0; i < numOfEdgeServer; i++ {
			weightList := []int{1, 1, 1, 1, 1}
			weightList[i/numOfRegistryServer] = 3
			chooser, _ := weightedrand.NewChooser(
				weightedrand.NewChoice(1, weightList[0]),
				weightedrand.NewChoice(2, weightList[1]),
				weightedrand.NewChoice(3, weightList[2]),
				weightedrand.NewChoice(4, weightList[3]),
				weightedrand.NewChoice(5, weightList[4]),
			)

			result := chooser.Pick()
			_ = edge.ImagePulling(&edgeServerList[i], (rand.Intn(100)+1)+result*100) // pulling 비용 계산하는 기능 개발
			if j%500 == 0 {
				edge.CleanCache(&edgeServerList[i])
			}

		}

		if j%managementInterval == 0 {

			// 자원 관리 - no는 일정 수준으로 변화 - 변화 수준이나 원격 레지스트리 서버에 대한 오버헤드는 모티베이션 실험 데이터를 활용
			// af는 pulling history에 따라 결정됨

			// 2단계 글러스터링 수행 > 1단계 클러스터링 후 리더 2개를 선정하고 엣지마다 레지스트리 서버 2개를 선정

			graph.GenerateRandomGraph(&no, numOfEdgeServer)
			graph.GenerateAffinityGraph(&af, edgeServerList)

			nsubgraphs := graph.ClusterGraph(&no, 3)
			fmt.Println("network overhead graph")
			graph.PrintGraph(&no)
			for _, subgraph := range nsubgraphs {
				fmt.Println("network overhead subgraph")
				leader := graph.ElectReaderUsingAffinity(&subgraph)
				graph.PrintGraphUsingReader(&subgraph, leader)
				firstRegistryServer = edge.CreateEdgeRegistryServer(edgeServerList[leader-1], leader)
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

			asubgraphs := graph.ClusterGraph(&af, 3)
			fmt.Println("affinity graph")
			graph.PrintGraph(&af)
			for _, subgraph := range asubgraphs {
				fmt.Println("affinity subgraph")
				leader := graph.ElectReaderUsingAffinity(&subgraph)
				graph.PrintAffinityGraph(&subgraph, leader)
				secondRegistryServer = edge.CreateEdgeRegistryServer(edgeServerList[leader-1], leader)
				for index, _ := range edgeServerList {
					//edge.SecondRegistry = secondRegistryServer
					for _, line := range af.Lines {
						if line.NodeA.Id == index && line.NodeB.Id == leader || line.NodeA.Id == leader && line.NodeB.Id == index {
							edgeServerList[index].AffinityOverhead = line.Val
						}
					}
					edgeServerList[index].SecondRegistry = secondRegistryServer
				}
			}

			// 주기적으로 edge server의 history를 정리
			//for _, edge := range edgeServerList {
			//	edge.History = nil
			//}
		}
	}

	//fmt.Println("e1's hit rate:", float32(edgeServer.HitCount)/float32(edgeServer.HitCount+edgeServer.MissCount), "%")
}