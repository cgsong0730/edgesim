package simulation

import (
	"bufio"
	"edgesim/edge"
	graph "edgesim/graph"
	"edgesim/weightedrand"
	_ "edgesim/weightedrand"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	_ "regexp"
	"strconv"
	_ "strconv"
	"strings"
	_ "strings"
)

type LeaderOfCluster struct {
	Id     string
	Leader string
}

type WeightOfTwoNode struct {
	NodeA  string
	NodeB  string
	Weight float64
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func JSONLinesToMaps(path string) ([]map[string]any, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var out []map[string]any
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		var m map[string]any
		if err := json.Unmarshal(sc.Bytes(), &m); err != nil {
			return nil, err
		}
		out = append(out, m)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func getString(rec map[string]any, key string) (string, bool) {
	v, ok := rec[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	return s, ok
}

func getStringSlice(rec map[string]any, key string) ([]string, bool) {
	raw, ok := rec[key].([]any)
	if !ok {
		return nil, false
	}
	out := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			out = append(out, s)
		} else {
			return nil, false // 타입 섞여 있으면 실패 처리
		}
	}
	return out, true
}

func ImagePullingSimulationWithFile(baseLine string, numOfCluster int, numOfNode int, useAffinity bool, useNetworkOverhead bool) {

	var RegistryServerList []edge.RegistryServer
	var edgeServerList []*edge.EdgeServer
	var edgeServerSubList []*edge.EdgeServer

	//var firstRegistryServer edge.EdgeRegistryServer
	//var secondRegistryServer edge.EdgeRegistryServer
	//var nearestRouterOverhead []int

	var cl []graph.Graph     // clustered network
	var ll []LeaderOfCluster // leader list
	var wl []WeightOfTwoNode // network bandwidth

	crecordList := make(map[int][]graph.Graph)
	lrecordList := make(map[int][]LeaderOfCluster)
	wrecordList := make(map[int][]WeightOfTwoNode)

	numOfRegistryServer := 5
	//numOfEdgeServer := 50
	//numOfPulling := 5
	numOfPulling := 100 // 200

	weightList := make(map[int][]weightedrand.Choice[int, int], numOfRegistryServer)

	cities := []string{"Aachen", "Augsburg", "Bayreuth", "Berlin", "Bielefeld", "Braunschweig", "Bremen", "Bremerhaven", "Chemnitz", "Darmstadt", "Dortmund", "Dresden", "Duesseldorf", "Erfurt", "Essen", "Flensburg", "Frankfurt", "Freiburg", "Fulda", "Giessen", "Greifswald", "Hamburg", "Hannover", "Kaiserslautern", "Karlsruhe", "Kassel", "Kempten", "Kiel", "Koblenz", "Koeln", "Konstanz", "Leipzig", "Magdeburg", "Mannheim", "Muenchen", "Muenster", "Norden", "Nuernberg", "Oldenburg", "Osnabrueck", "Passau", "Regensburg", "Saarbruecken", "Schwerin", "Siegen", "Stuttgart", "Trier", "Ulm", "Wesel", "Wuerzburg"}

	var clustersFilePath string
	var leadersFilePath string

	var file *os.File
	fineName := "./result.txt"

	if baseLine == "mkrp" {

		if numOfCluster == 5 {
			clustersFilePath = "mkrp-c-5.txt"
			leadersFilePath = "mkrp-l-5.txt"
		} else if numOfCluster == 10 {
			clustersFilePath = "mkrp-c-10.txt"
			leadersFilePath = "mkrp-l-10.txt"
		} else if numOfCluster == 15 {
			clustersFilePath = "mkrp-c-15.txt"
			leadersFilePath = "mkrp-l-15.txt"
		}

	} else if baseLine == "comm" {

		if numOfCluster == 5 {
			clustersFilePath = "comm-c-5.txt"
			leadersFilePath = "comm-l-5.txt"
		} else if numOfCluster == 10 {
			clustersFilePath = "comm-c-10.txt"
			leadersFilePath = "comm-l-10.txt"
		} else if numOfCluster == 15 {
			clustersFilePath = "mkrp-c-15.txt"
			leadersFilePath = "mkrp-l-15.txt"
		}

	}

	//clustersFilePath := "comm-c-5.txt"
	//leadersFilePath := "comm-l-5.txt"

	var nweightsFilePath string

	// Initializing random value weights
	for i := 1; i <= numOfRegistryServer; i++ {
		var weights []weightedrand.Choice[int, int]

		for j := 1; j <= numOfRegistryServer; j++ {

			if i == j {
				weights = append(weights, weightedrand.NewChoice(j, 5))
			} else {
				weights = append(weights, weightedrand.NewChoice(j, 1))
			}
		}

		weightList[i] = weights
	}

	// Initializing remote registry servers
	for i := 0; i < numOfRegistryServer; i++ {
		rs := edge.RegistryServer{i + 1, nil, rand.Intn(3) + 9} // network overhead
		edge.InitRegistryServer(&rs, 100, 500, 1+i*100, 100+i*100)
		RegistryServerList = append(RegistryServerList, rs)
	}

	// Initializing edge servers
	for i, city := range cities {
		edgeServer := edge.EdgeServer{
			Id:               i,
			Name:             city,
			NumOfImage:       0,
			MaxCacheSize:     50000,
			CurrentCacheSize: 0,
			LocalImages:      nil,
			RegistryServers:  RegistryServerList,
		}
		edgeServerList = append(edgeServerList, &edgeServer)
	}

	// parsing weight file
	for i := 0; i < 288; i++ {
		nweightsFilePath = "weights/weight-" + strconv.Itoa(i+1) + ".txt"
		var wn WeightOfTwoNode
		f, err := os.Open(nweightsFilePath) // 예: 위 데이터를 저장한 파일
		if err != nil {
			panic(err)
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)

		for scanner.Scan() {
			line := scanner.Text()
			var from, to string
			var w float64

			lineTemp := strings.ReplaceAll(line, "(", "")
			lineTemp = strings.ReplaceAll(lineTemp, ")", "")
			lineTemp = strings.ReplaceAll(lineTemp, " ", "")
			parts := strings.Split(lineTemp, ",")

			for i, part := range parts {
				if i == 0 {
					from = part
				} else if i == 1 {
					to = part
				} else if i == 2 {
					w, _ = strconv.ParseFloat(part, 64)
				}
			}
			wn.NodeA = from
			wn.NodeB = to
			wn.Weight = w
			//fmt.Printf("%s, %s, %.4f\n", from, to, w)
			wl = append(wl, wn)

		}
		wrecordList[i] = wl
		wl = nil
	}

	// parsing clusters file
	cRecords, err := JSONLinesToMaps(clustersFilePath)
	if err != nil {
		log.Fatalf("Fail to parsing file: %v", err)
	}

	// parsing leaders file
	lRecords, err := JSONLinesToMaps(leadersFilePath)
	if err != nil {
		log.Fatalf("Fail to parsing file: %v", err)
	}

	for i, rec := range lRecords {
		var lc LeaderOfCluster
		for k := range rec {
			if city, ok := getString(rec, k); ok {
				lc.Id = k
				lc.Leader = city
			}
			ll = append(ll, lc)
		}
		lrecordList[i] = ll
		ll = nil
	}

	var nodeA *graph.Node
	var nodeB *graph.Node
	for i, rec := range cRecords {

		fmt.Println("Record ", i)

		for k := range rec {
			var subgraph graph.Graph
			subgraph.Id = k
			if tags, ok := getStringSlice(rec, k); ok {
				for _, tag := range tags {
					graph.AddNode(&subgraph, tag)
				}
			}
			for j, weight := range wrecordList[i] {

				for l, node1 := range subgraph.Nodes {
					if weight.NodeA == node1.Id {
						//isNodeA = true
						nodeA = &subgraph.Nodes[l]
						for m, node2 := range subgraph.Nodes {
							if weight.NodeB != node1.Id && weight.NodeB == node2.Id {
								//isNodeB = true
								nodeB = &subgraph.Nodes[m]

								line := graph.Line{j, weight.Weight, nodeA, nodeB}
								subgraph.Lines = append(subgraph.Lines, line)
								continue
							}
						}
					}
				}
			}

			cl = append(cl, subgraph)
			//graph.PrintGraph(&subgraph)
		}

		//nearestRouterOverhead = append(nearestRouterOverhead, 0)
		crecordList[i] = cl
		cl = nil
	}

	sum := float64(0)
	cnt := float64(0)
	avg := float64(0)
	//numOfSubgroup := 10

	type0Cnt := float64(0)
	type1Cnt := float64(0)
	type2Cnt := float64(0)
	type3Cnt := float64(0)

	error1Cnt := float64(0)
	error2Cnt := float64(0)

	errorCnt := float64(0)

	//for i := 0; i < 288; i++ {
	for i := 0; i < 96; i++ {

		fmt.Printf("recored: %d \n", i)

		cl = crecordList[i]
		ll = lrecordList[i]
		wl = wrecordList[i]

		var firstRegistryServer edge.EdgeRegistryServer
		// first k registry deployment
		for _, subgraph := range cl {
			var leaderId string
			var cLeaderEdgeServer *edge.EdgeServer
			for _, leader := range ll {
				if subgraph.Id == leader.Id {
					//println("subgraph:", subgraph.Id)
					//println("leader Id:", leader.Id)
					//println("leader:", leader.Leader)
					leaderId = leader.Leader
				}
			}

			// edge server list in cluster + networkbandwith
			for _, edgeServer := range edgeServerList {

				for _, node := range subgraph.Nodes {
					if node.Id == edgeServer.Name {
						edgeServerSubList = append(edgeServerSubList, edgeServer)
					}
				}
				if leaderId == edgeServer.Name {
					cLeaderEdgeServer = edgeServer
				}
			}

			firstRegistryServer = edge.CreateEdgeRegistryServer(edgeServerSubList, cLeaderEdgeServer, leaderId, 300)

			for _, edgeServer := range edgeServerSubList {

				edgeServer.FirstRegistryBandwidth = float64(118.7)
				for _, line := range subgraph.Lines {
					if edgeServer.Name == line.NodeA.Id && cLeaderEdgeServer.Name == line.NodeB.Id ||
						cLeaderEdgeServer.Name == line.NodeA.Id && edgeServer.Name == line.NodeB.Id {
						edgeServer.FirstRegistryBandwidth = line.Val
					}
				}
				edgeServer.FirstRegistry = firstRegistryServer
			}
			//edgeServerSubList = []*edge.EdgeServer{}
			edgeServerSubList = edgeServerSubList[:0]
		}

		// second k registry deployment
		var af graph.Graph // affinity
		var secondRegistryServer edge.EdgeRegistryServer

		graph.GenerateAffinityGraph(&af, edgeServerList)

		//asubgraphs := graph.ClusterGraph(&af, numOfSubgroup)
		var aleaderId string
		for _, subgraph := range cl {
			var aLeaderEdgeServer *edge.EdgeServer
			aleaderId = graph.ElectReaderUsingAffinity(&subgraph)

			for _, edgeServer := range edgeServerList {

				for _, node := range subgraph.Nodes {
					if node.Id == edgeServer.Name {
						edgeServerSubList = append(edgeServerSubList, edgeServer)
					}
				}
				if aleaderId == edgeServer.Name {
					aLeaderEdgeServer = edgeServer
				}
			}

			secondRegistryServer = edge.CreateEdgeRegistryServer(edgeServerSubList, aLeaderEdgeServer, aleaderId, 300)
			//for index, _ := range edgeServerList {
			//	edgeServerList[index].SecondRegistry = secondRegistryServer
			//}

			for _, edgeServer := range edgeServerSubList {

				edgeServer.SecondRegistryBandwidth = float64(118.7)
				for _, line := range subgraph.Lines {
					if edgeServer.Name == line.NodeA.Id && aLeaderEdgeServer.Name == line.NodeB.Id ||
						aLeaderEdgeServer.Name == line.NodeA.Id && edgeServer.Name == line.NodeB.Id {
						edgeServer.SecondRegistryBandwidth = line.Val
					}
				}
				edgeServer.SecondRegistry = secondRegistryServer
			}

		}

		for j := 1; j <= numOfPulling; j += 1 {

			fmt.Printf("pulling: %d \n", j)

			for k, edgeServer := range edgeServerList {

				chooser, err := weightedrand.NewChooser(weightList[j%numOfRegistryServer+1]...)
				if err != nil {
					log.Fatalf("pulling %d: %v", k, err)
				}
				result := chooser.Pick()

				// 랜덤으로 아이디 생성
				//pullingTime := edge.ImagePullingWithData(edgeServer, j, 0)

				requestId := rand.Intn(100) + 1 + (result-1)*100
				pullingTime, pullingType := edge.ImagePullingWithData(edgeServer, requestId, useAffinity, useNetworkOverhead)
				//fmt.Printf("pullTime: %.2f %d \n", pullingTime, pullingType)

				if !math.IsInf(pullingTime, 0) {

					if pullingType == 0 {
						type0Cnt++
					} else if pullingType == 1 {
						type1Cnt++
					} else if pullingType == 2 {
						type2Cnt++
					} else if pullingType == 3 {
						type3Cnt++
					}
					sum += pullingTime
					cnt++
				} else {
					fmt.Printf("error: %d \n", pullingType)

					if pullingType == 1 {
						error1Cnt++
					} else if pullingType == 2 {
						//graph.PrintAffinityGraph(&af, aleaderId)
						error2Cnt++
					}
					errorCnt++
				}
			}

			if i%100 == 0 {
				for _, edgeServer := range edgeServerList {
					edge.CleanCache(edgeServer)
				}
			}
			//fmt.Printf("sum: %f \n", sum)
		}
	} // of records

	avg = sum / cnt
	fmt.Printf("avg: %f \n", avg)

	fmt.Printf("type0Cnt: %f \n", (type0Cnt/cnt)*100)
	fmt.Printf("type1Cnt: %f \n", (type1Cnt/cnt)*100)
	fmt.Printf("type2Cnt: %f \n", (type2Cnt/cnt)*100)
	fmt.Printf("type3Cnt: %f \n", (type3Cnt/cnt)*100)

	fmt.Printf("error1Cnt: %f \n", (error1Cnt/errorCnt)*100)
	fmt.Printf("error2Cnt: %f \n", (error2Cnt/errorCnt)*100)

	fmt.Printf("errorCnt: %f \n", (errorCnt/(cnt+errorCnt))*100)

	file, _ = os.OpenFile(fineName, os.O_APPEND|os.O_RDWR, 0755)
	_, ferr := os.Stat(fineName)

	if os.IsNotExist(ferr) {
		file, _ = os.Create(fineName)
	} else {
		file, _ = os.OpenFile(fineName, os.O_APPEND|os.O_RDWR, 0755)
	}

	// baseLine string, numOfCluster int, numOfNode int, useAffinity bool, useNetworkOverhead bool
	str := fmt.Sprintf("%s \t %d \t %d \t %t \t %t \t %.2f \n", baseLine, numOfCluster, numOfNode, useAffinity, useNetworkOverhead, avg)
	b := []byte(str)
	_, err = file.Write(b)
	check(err)
}

//func ImagePullingSimulationWithNetworkx(numOfEdgeServer int, numOfRegistryServer int, numOfSubgroup int, numOfPulling int, managementInterval int, baseLine int) {
//	var RegistryServerList []edge.RegistryServer
//	var firstRegistryServer edge.EdgeRegistryServer
//	var secondRegistryServer edge.EdgeRegistryServer
//	edgeServerList := []edge.EdgeServer{}
//	var nearestRouterOverhead []int
//	var no graph.Graph // network overhead
//	var file *os.File
//	sum := 0
//	cnt := 0
//
//	fineName := "./result.txt"
//	_, err := os.Stat(fineName)
//	if os.IsNotExist(err) {
//		file, _ = os.Create(fineName)
//	} else {
//		file, _ = os.OpenFile(fineName, os.O_APPEND|os.O_RDWR, 0755)
//	}
//
//	for i := 0; i < numOfRegistryServer; i++ {
//		rs := edge.RegistryServer{i + 1, nil, rand.Intn(10) + 90}
//		edge.InitRegistryServer(&rs, 100, 500, 1+i*100, 100+i*100)
//		RegistryServerList = append(RegistryServerList, rs)
//	}
//
//	for i := 1; i <= numOfEdgeServer; i++ {
//		edgeServer := edge.EdgeServer{
//			Id:               i,
//			NumOfImage:       0,
//			MaxCacheSize:     5000,
//			CurrentCacheSize: 0,
//			LocalImages:      nil,
//			RegistryServers:  RegistryServerList,
//		}
//		edgeServerList = append(edgeServerList, edgeServer)
//		nearestRouterOverhead = append(nearestRouterOverhead, 0)
//	}
//
//	nearestRouterOverhead = graph.GenerateRandomGraphWithNR(&no, numOfEdgeServer)
//	no.Nodes = no.Nodes[:len(no.Nodes)-1]
//	nsubgraphs := graph.ClusterGraph(&no, numOfSubgroup)
//
//	for _, subgraph := range nsubgraphs {
//		leader := graph.ElectReaderUsingOverhead(&subgraph)
//		firstRegistryServer = edge.CreateEdgeRegistryServer(edgeServerList[leader-1], leader, numOfPulling/managementInterval)
//		for index, _ := range edgeServerList {
//			//edge.FirstRegistry = firstRegistryServer
//			for _, line := range no.Lines {
//				if line.NodeA.Id == index && line.NodeB.Id == leader || line.NodeA.Id == leader && line.NodeB.Id == index {
//					edgeServerList[index].NetworkOverhead = line.Val
//				}
//			}
//			edgeServerList[index].FirstRegistry = firstRegistryServer
//		}
//	}
//
//	for j := 1; j <= numOfPulling; j++ {
//		var af graph.Graph // affinity
//		var weightList []weightedrand.Choice[int, int]
//
//		for i := 0; i < numOfEdgeServer; i++ {
//			if i%numOfRegistryServer == 0 {
//				weightList = append(weightList, weightedrand.NewChoice(i+1, 3))
//			} else {
//				weightList = append(weightList, weightedrand.NewChoice(i+1, 1))
//			}
//			chooser, _ := weightedrand.NewChooser(weightList...)
//			result := chooser.Pick()
//
//			var pullingTime = 0
//			if baseLine == 1 {
//				pullingTime = edge.ImagePullingB1(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])
//			} else if baseLine == 2 {
//				pullingTime = edge.ImagePullingB2(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])
//			} else if baseLine == 3 {
//				pullingTime = edge.ImagePullingB3(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])
//			} else {
//				pullingTime = edge.ImagePulling(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])
//			}
//
//			sum += pullingTime
//			cnt++
//
//			if j%500 == 0 {
//				edge.CleanCache(&edgeServerList[i])
//			}
//		}
//
//		if j%managementInterval == 0 {
//
//			graph.GenerateAffinityGraph(&af, edgeServerList)
//			asubgraphs := graph.ClusterGraph(&af, numOfSubgroup)
//			for _, subgraph := range asubgraphs {
//				//fmt.Println("affinity subgraph")
//				leader := graph.ElectReaderUsingAffinity(&subgraph)
//				secondRegistryServer = edge.CreateEdgeRegistryServer(edgeServerList[leader-1], leader, numOfPulling/managementInterval)
//				for index, _ := range edgeServerList {
//					for _, line := range af.Lines {
//						if line.NodeA.Id == strconv.Itoa(index && line.NodeB.Id == leader || line.NodeA.Id == leader && line.NodeB.Id == index {
//							edgeServerList[index].AffinityOverhead = line.Val
//						}
//					}
//					edgeServerList[index].SecondRegistry = secondRegistryServer
//				}
//			}
//
//			if j%managementInterval*3 == 0 {
//				for _, edge := range edgeServerList {
//					edge.History = nil
//				}
//			}
//		}
//	}
//	str := fmt.Sprintf("%d %d %d %d\n", numOfEdgeServer, numOfSubgroup, baseLine, sum/cnt)
//	b := []byte(str)
//	_, err = file.Write(b)
//	check(err)
//}

// func ImagePullingSimulationWithGraph(numOfPulling int, managementInterval int, baseLine int) {
//
//		var RegistryServerList []edge.RegistryServer
//		var firstRegistryServer edge.EdgeRegistryServer
//		//var secondRegistryServer edge.EdgeRegistryServer
//		edgeServerList := []edge.EdgeServer{}
//		var nearestRouterOverhead []int
//		var no []graph.Graph // network overhead
//		var nodeList []*graph.Node
//		numOfRegistryServer := 15
//		numOfEdgeServer := 30
//		fineName := "./result.txt"
//		var file *os.File
//		sum := 0
//		cnt := 0
//
//		_, err := os.Stat(fineName)
//		if os.IsNotExist(err) {
//			file, _ = os.Create(fineName)
//		} else {
//			file, _ = os.OpenFile(fineName, os.O_APPEND|os.O_RDWR, 0755)
//		}
//
//		for i := 0; i < numOfRegistryServer; i++ {
//			rs := edge.RegistryServer{i + 1, nil, rand.Intn(10) + 90}
//			edge.InitRegistryServer(&rs, 100, 500, 1+i*100, 100+i*100)
//			RegistryServerList = append(RegistryServerList, rs)
//		}
//
//		// 파일 열기
//		file_n, err := os.Open("n2_15")
//		if err != nil {
//			fmt.Println("can't open file-", err)
//			return
//		}
//		defer file_n.Close()
//
//		file_w, err := os.Open("w2_15")
//		if err != nil {
//			fmt.Println("can't open file-", err)
//			return
//		}
//		defer file_n.Close()
//
//		file_e, err := os.Open("e2_15")
//		if err != nil {
//			fmt.Println("can't open file-", err)
//			return
//		}
//		defer file_e.Close()
//
//		var fileContentN string
//		scannerN := bufio.NewScanner(file_n)
//		for scannerN.Scan() {
//			fileContentN += scannerN.Text()
//		}
//		if err := scannerN.Err(); err != nil {
//			fmt.Println("can't read file-", err)
//			return
//		}
//
//		var fileContentW string
//		scannerW := bufio.NewScanner(file_w)
//		for scannerW.Scan() {
//			fileContentW += scannerW.Text()
//		}
//		if err := scannerW.Err(); err != nil {
//			fmt.Println("can't read file-", err)
//			return
//		}
//
//		var fileContentE string
//		scannerE := bufio.NewScanner(file_e)
//		for scannerE.Scan() {
//			fileContentE += scannerE.Text()
//		}
//		if err := scannerE.Err(); err != nil {
//			fmt.Println("can't read file-", err)
//			return
//		}
//
//		jsonCompatibleDataN := strings.ReplaceAll(fileContentN, "{", "[")
//		jsonCompatibleDataN = strings.ReplaceAll(jsonCompatibleDataN, "}", "]")
//
//		jsonCompatibleDataW := strings.ReplaceAll(fileContentW, "(", "[")
//		jsonCompatibleDataW = strings.ReplaceAll(jsonCompatibleDataW, ")", "]")
//
//		var parsedDataN [][]int
//		if err := json.Unmarshal([]byte(jsonCompatibleDataN), &parsedDataN); err != nil {
//			fmt.Println("JSON 파싱 중 오류:", err)
//			return
//		}
//
//		var parsedDataW [][]int
//		if err := json.Unmarshal([]byte(jsonCompatibleDataW), &parsedDataW); err != nil {
//			fmt.Println("JSON 파싱 중 오류:", err)
//			return
//		}
//
//		re := regexp.MustCompile(`\d+`)
//		matches := re.FindAllString(fileContentE, -1)
//
//		// 정수 슬라이스 생성
//		var parsedDataE []int
//		for _, match := range matches {
//			num, err := strconv.Atoi(match)
//			if err != nil {
//				fmt.Println("문자열을 정수로 변환하는 데 실패했습니다:", err)
//				return
//			}
//			parsedDataE = append(parsedDataE, num)
//		}
//
//		for _, group := range parsedDataN {
//			var subgraph graph.Graph
//			for _, value := range group {
//				graph.AddNode(&subgraph, value)
//				edgeServer := edge.EdgeServer{
//					Id:               value,
//					NumOfImage:       0,
//					MaxCacheSize:     5000,
//					CurrentCacheSize: 0,
//					LocalImages:      nil,
//					RegistryServers:  RegistryServerList,
//				}
//				edgeServerList = append(edgeServerList, edgeServer)
//				nearestRouterOverhead = append(nearestRouterOverhead, 0)
//			}
//			no = append(no, subgraph)
//		}
//
//		for _, g := range no {
//			for _, node := range g.Nodes {
//				nodeList = append(nodeList, &node)
//			}
//		}
//
//		for _, g := range no {
//			var nodeA *graph.Node
//			var nodeB *graph.Node
//			var weight int
//			numOfLine := 1
//
//			for _, group := range parsedDataW {
//				for j, value := range group {
//					if j == 0 {
//						nodeA = graph.FindNodeById(&g, value)
//					}
//					if j == 1 {
//						nodeB = graph.FindNodeById(&g, value)
//					}
//					if j == 2 {
//						weight = value
//					}
//				}
//				if nodeA != nil && nodeB != nil {
//					graph.AddLine(&g, nodeA, nodeB, numOfLine, weight)
//					numOfLine++
//				}
//			}
//			graph.PrintGraph(&g)
//		}
//
//		for i, _ := range no {
//			leader := parsedDataE[i]
//			for j, edgeNode := range edgeServerList {
//				if edgeNode.Id == j {
//					firstRegistryServer = edge.CreateEdgeRegistryServer(edgeNode, leader, numOfPulling/managementInterval)
//				}
//			}
//			for j, _ := range edgeServerList {
//				edgeServerList[j].FirstRegistry = firstRegistryServer
//			}
//		}
//
//		for j := 1; j <= numOfPulling; j++ {
//			var weightList []weightedrand.Choice[int, int]
//
//			for i := 0; i < numOfEdgeServer; i++ {
//				if i%numOfRegistryServer == 0 {
//					weightList = append(weightList, weightedrand.NewChoice(i+1, 3))
//				} else {
//					weightList = append(weightList, weightedrand.NewChoice(i+1, 1))
//				}
//				chooser, _ := weightedrand.NewChooser(weightList...)
//				result := chooser.Pick()
//
//				var pullingTime = 0
//				pullingTime = edge.ImagePullingB2(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])
//				sum += pullingTime
//				cnt++
//
//				if j%500 == 0 {
//					edge.CleanCache(&edgeServerList[i])
//				}
//			}
//		}
//		str := fmt.Sprintf("%d %d %d %d\n", numOfEdgeServer, numOfRegistryServer, baseLine, sum/cnt)
//		b := []byte(str)
//		_, err = file.Write(b)
//		check(err)
//	}

//func ImagePullingSimulation(numOfEdgeServer int, numOfRegistryServer int, numOfSubgroup int, numOfPulling int, managementInterval int, baseLine int) {
//	var RegistryServerList []edge.RegistryServer
//	var firstRegistryServer edge.EdgeRegistryServer
//	var secondRegistryServer edge.EdgeRegistryServer
//	edgeServerList := []edge.EdgeServer{}
//	var nearestRouterOverhead []int
//	var no graph.Graph // network overhead
//	fineName := "./result.txt"
//	var file *os.File
//	sum := 0
//	cnt := 0
//
//	_, err := os.Stat(fineName)
//	if os.IsNotExist(err) {
//		file, _ = os.Create(fineName)
//	} else {
//		file, _ = os.OpenFile(fineName, os.O_APPEND|os.O_RDWR, 0755)
//	}
//
//	for i := 0; i < numOfRegistryServer; i++ {
//		rs := edge.RegistryServer{i + 1, nil, rand.Intn(10) + 90}
//		edge.InitRegistryServer(&rs, 100, 500, 1+i*100, 100+i*100)
//		RegistryServerList = append(RegistryServerList, rs)
//	}
//
//	for i := 1; i <= numOfEdgeServer; i++ {
//		edgeServer := edge.EdgeServer{
//			Id:               i,
//			NumOfImage:       0,
//			MaxCacheSize:     5000,
//			CurrentCacheSize: 0,
//			LocalImages:      nil,
//			RegistryServers:  RegistryServerList,
//		}
//		edgeServerList = append(edgeServerList, edgeServer)
//		nearestRouterOverhead = append(nearestRouterOverhead, 0)
//	}
//
//	nearestRouterOverhead = graph.GenerateRandomGraphWithNR(&no, numOfEdgeServer)
//	no.Nodes = no.Nodes[:len(no.Nodes)-1]
//	nsubgraphs := graph.ClusterGraph(&no, numOfSubgroup)
//
//	for _, subgraph := range nsubgraphs {
//		leader := graph.ElectReaderUsingOverhead(&subgraph)
//		firstRegistryServer = edge.CreateEdgeRegistryServer(edgeServerList[leader-1], leader, numOfPulling/managementInterval)
//		for index, _ := range edgeServerList {
//			//edge.FirstRegistry = firstRegistryServer
//			for _, line := range no.Lines {
//				if line.NodeA.Id == index && line.NodeB.Id == leader || line.NodeA.Id == leader && line.NodeB.Id == index {
//					edgeServerList[index].NetworkOverhead = line.Val
//				}
//			}
//			edgeServerList[index].FirstRegistry = firstRegistryServer
//		}
//	}
//
//	for j := 1; j <= numOfPulling; j++ {
//		var af graph.Graph // affinity
//		var weightList []weightedrand.Choice[int, int]
//
//		for i := 0; i < numOfEdgeServer; i++ {
//			if i%numOfRegistryServer == 0 {
//				weightList = append(weightList, weightedrand.NewChoice(i+1, 3))
//			} else {
//				weightList = append(weightList, weightedrand.NewChoice(i+1, 1))
//			}
//			chooser, _ := weightedrand.NewChooser(weightList...)
//			result := chooser.Pick()
//
//			var pullingTime = 0
//			if baseLine == 1 {
//				pullingTime = edge.ImagePullingB1(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])
//			} else if baseLine == 2 {
//				pullingTime = edge.ImagePullingB2(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])
//			} else if baseLine == 3 {
//				pullingTime = edge.ImagePullingB3(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])
//			} else {
//				pullingTime = edge.ImagePulling(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])
//			}
//
//			sum += pullingTime
//			cnt++
//
//			if j%500 == 0 {
//				edge.CleanCache(&edgeServerList[i])
//			}
//		}
//
//		if j%managementInterval == 0 {
//
//			graph.GenerateAffinityGraph(&af, edgeServerList)
//			asubgraphs := graph.ClusterGraph(&af, numOfSubgroup)
//			for _, subgraph := range asubgraphs {
//				//fmt.Println("affinity subgraph")
//				leader := graph.ElectReaderUsingAffinity(&subgraph)
//				secondRegistryServer = edge.CreateEdgeRegistryServer(edgeServerList[leader-1], leader, numOfPulling/managementInterval)
//				for index, _ := range edgeServerList {
//					for _, line := range af.Lines {
//						if line.NodeA.Id == index && line.NodeB.Id == leader || line.NodeA.Id == leader && line.NodeB.Id == index {
//							edgeServerList[index].AffinityOverhead = line.Val
//						}
//					}
//					edgeServerList[index].SecondRegistry = secondRegistryServer
//				}
//			}
//
//			if j%managementInterval*3 == 0 {
//				for _, edge := range edgeServerList {
//					edge.History = nil
//				}
//			}
//		}
//	}
//	str := fmt.Sprintf("%d %d %d %d\n", numOfEdgeServer, numOfSubgroup, baseLine, sum/cnt)
//	b := []byte(str)
//	_, err = file.Write(b)
//	check(err)
//}

//func ImagePullingVisualization(numOfEdgeServer int, numOfRegistryServer int, numOfSubgroup int, numOfPulling int, managementInterval int) {
//	var RegistryServerList []edge.RegistryServer
//	var firstRegistryServer edge.EdgeRegistryServer
//	var secondRegistryServer edge.EdgeRegistryServer
//	edgeServerList := []edge.EdgeServer{}
//	var nearestRouterOverhead []int
//	var no graph.Graph // network overhead
//	fileList := []*os.File{}
//
//	for i := 1; i <= numOfEdgeServer; i++ {
//		fileName := fmt.Sprintf("./edge%d.txt", i)
//		var file *os.File
//
//		_, err := os.Stat(fileName)
//		if os.IsNotExist(err) {
//			file, _ = os.Create(fileName)
//		} else {
//			file, _ = os.OpenFile(fileName, os.O_APPEND|os.O_RDWR, 0755)
//		}
//
//		fileList = append(fileList, file)
//		defer file.Close()
//	}
//
//	for i := 0; i < numOfRegistryServer; i++ {
//		rs := edge.RegistryServer{i + 1, nil, rand.Intn(10) + 90}
//		edge.InitRegistryServer(&rs, 100, 500, 1+i*100, 100+i*100)
//		RegistryServerList = append(RegistryServerList, rs)
//	}
//
//	for i := 1; i <= numOfEdgeServer; i++ {
//		edgeServer := edge.EdgeServer{
//			Id:               i,
//			NumOfImage:       0,
//			MaxCacheSize:     5000,
//			CurrentCacheSize: 0,
//			LocalImages:      nil,
//			RegistryServers:  RegistryServerList,
//		}
//		edgeServerList = append(edgeServerList, edgeServer)
//		nearestRouterOverhead = append(nearestRouterOverhead, 0)
//	}
//
//	nearestRouterOverhead = graph.GenerateRandomGraphWithNR(&no, numOfEdgeServer)
//	fmt.Println("network overhead graph")
//	graph.PrintNetworkGraph(&no)
//
//	no.Nodes = no.Nodes[:len(no.Nodes)-1]
//	nsubgraphs := graph.ClusterGraph(&no, numOfSubgroup)
//
//	for _, subgraph := range nsubgraphs {
//		fmt.Println("network overhead subgraph")
//		leader := graph.ElectReaderUsingOverhead(&subgraph)
//		graph.PrintGraphUsingReader(&subgraph, leader)
//		firstRegistryServer = edge.CreateEdgeRegistryServer(edgeServerList[leader-1], leader, numOfPulling/managementInterval)
//		for index, _ := range edgeServerList {
//			//edge.FirstRegistry = firstRegistryServer
//			for _, line := range no.Lines {
//				if line.NodeA.Id == index && line.NodeB.Id == leader || line.NodeA.Id == leader && line.NodeB.Id == index {
//					edgeServerList[index].NetworkOverhead = line.Val
//				}
//			}
//			edgeServerList[index].FirstRegistry = firstRegistryServer
//		}
//	}
//
//	for j := 1; j <= numOfPulling; j += 1 {
//		var af graph.Graph // affinity
//		var weightList []weightedrand.Choice[int, int]
//
//		for i := 0; i < numOfEdgeServer; i++ {
//			if i%numOfRegistryServer == 0 {
//				weightList = append(weightList, weightedrand.NewChoice(i+1, 3))
//			} else {
//				weightList = append(weightList, weightedrand.NewChoice(i+1, 1))
//			}
//			chooser, _ := weightedrand.NewChooser(weightList...)
//			result := chooser.Pick()
//
//			pullingTime := edge.ImagePulling(&edgeServerList[i], (rand.Intn(100)+1)+result*100, nearestRouterOverhead[i])
//
//			if j%10 == 0 {
//				str := fmt.Sprintf("%d %d\n", j, pullingTime)
//				b := []byte(str)
//				_, err := fileList[i].Write(b)
//				check(err)
//			}
//
//			if j%500 == 0 {
//				edge.CleanCache(&edgeServerList[i])
//			}
//		}
//		if j%managementInterval == 0 {
//			graph.GenerateAffinityGraph(&af, edgeServerList)
//			asubgraphs := graph.ClusterGraph(&af, numOfSubgroup)
//			fmt.Println("affinity graph")
//			graph.PrintAffinityOverallGraph(&af)
//			for _, subgraph := range asubgraphs {
//				fmt.Println("affinity subgraph")
//				leader := graph.ElectReaderUsingAffinity(&subgraph)
//				graph.PrintAffinityGraph(&subgraph, leader)
//				secondRegistryServer = edge.CreateEdgeRegistryServer(edgeServerList[leader-1], leader, numOfPulling/managementInterval)
//				for index, _ := range edgeServerList {
//					for _, line := range af.Lines {
//						if line.NodeA.Id == index && line.NodeB.Id == leader || line.NodeA.Id == leader && line.NodeB.Id == index {
//							edgeServerList[index].AffinityOverhead = line.Val
//						}
//					}
//					edgeServerList[index].SecondRegistry = secondRegistryServer
//				}
//			}
//		}
//
//		// 주기적으로 edge server의 history를 정리
//		if j%managementInterval*3 == 0 {
//			for _, edge := range edgeServerList {
//				edge.History = nil
//			}
//		}
//		// 주기적으로 캐시를 정리
//	}
//}
