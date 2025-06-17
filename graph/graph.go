package graph

import (
	"edgesim/edge"
	"fmt"
	"math/rand"
	"strconv"
)

type Node struct {
	Id string
}

type Line struct {
	Id    int
	Val   float64
	NodeA *Node
	NodeB *Node
}

type Graph struct {
	Id    string
	Lines []Line
	Nodes []Node
}

func AddNode(graph *Graph, id string) {
	var node Node = Node{id}
	graph.Nodes = append(graph.Nodes, node)
}

func FindNodeById(graph *Graph, id string) *Node {

	var returnNode *Node
	for _, node := range graph.Nodes {
		if node.Id == id {
			returnNode = &node
			return returnNode
		}
	}
	return returnNode

}

func AddLine(graph *Graph, nodeA *Node, nodeB *Node, id int, val float64) {
	var line Line = Line{id, val, nodeA, nodeB}
	graph.Lines = append(graph.Lines, line)
}

func PrintGraph(graph *Graph) {
	fmt.Println("Graph {")

	for _, line := range graph.Lines {
		fmt.Printf("  %s -- %s [len = %d, label = \"%d\"] \n", line.NodeA.Id, line.NodeB.Id, line.Val, line.Val)
	}
	fmt.Println("}")
}

func PrintGraphUsingReader(graph *Graph, nodeId int) {
	fmt.Println("Graph {")
	fmt.Printf("  {node [style=filled,color=skyblue] e%d}\n", nodeId)
	for _, line := range graph.Lines {
		fmt.Printf("  e%d -- e%d [len = %d, label = \"%d\"] \n", line.NodeA.Id, line.NodeB.Id, line.Val, line.Val)
	}
	fmt.Println("}")
}

func PrintNetworkGraph(graph *Graph) {
	fmt.Println("Graph {")
	fmt.Printf("  {node [shape=octagon, style=filled,color=green] e%d}\n", len(graph.Nodes)-1)
	for _, line := range graph.Lines {
		fmt.Printf("  e%d -- e%d [len = %d, label = \"%d\"] \n", line.NodeA.Id, line.NodeB.Id, line.Val, line.Val)
	}
	fmt.Println("}")
}

func PrintAffinityOverallGraph(graph *Graph) {
	fmt.Println("Graph {")
	for _, line := range graph.Lines {
		fmt.Printf("  e%d -- e%d [len = %d, label = \"%.2f\"] \n", line.NodeA.Id, line.NodeB.Id, line.Val/5, float32(line.Val)/100)
	}
	fmt.Println("}")
}

func PrintAffinityGraph(graph *Graph, nodeId string) {
	fmt.Println("Graph {")
	fmt.Printf("  {node [style=filled,color=yellow] e%s}\n", nodeId)
	for _, line := range graph.Lines {
		fmt.Printf("  %s -- %s [len = %f, label = \" %.2f\"] \n", line.NodeA.Id, line.NodeB.Id, line.Val/5, line.Val)
	}
	fmt.Println("}")
}

func GenerateRandomGraph(graph *Graph, numOfNode int) {

	cnt := numOfNode - 1
	randomValue := 0
	rangeOfValue := 10
	numOfLine := 1

	for i := 1; i <= numOfNode; i++ {
		AddNode(graph, strconv.Itoa(i))
	}
	for i := 1; i <= numOfNode; i++ {
		for j := 1; j <= cnt; j++ {
			randomValue = rand.Intn(rangeOfValue) + 1
			AddLine(graph, &graph.Nodes[i-1], &graph.Nodes[i+j-1], numOfLine, float64(randomValue))
			numOfLine++
		}
		cnt--
	}
}

func GenerateRandomGraphWithNR(graph *Graph, numOfNode int) []int {
	numOfNode++
	cnt := numOfNode - 1
	randomValue := 0
	rangeOfValue := 10
	numOfLine := 1
	var nearestRouterOverhead []int

	for i := 1; i <= numOfNode; i++ {
		AddNode(graph, strconv.Itoa(i))
	}
	for i := 1; i <= numOfNode; i++ {
		for j := 1; j <= cnt; j++ {
			randomValue = rand.Intn(rangeOfValue) + 1
			AddLine(graph, &graph.Nodes[i-1], &graph.Nodes[i+j-1], numOfLine, float64(randomValue))
			if i+j == numOfNode {
				nearestRouterOverhead = append(nearestRouterOverhead, randomValue)
			}
			numOfLine++
		}
		cnt--
	}
	return nearestRouterOverhead
}

func findNodeById(graph *Graph, Id string) *Node {
	for _, node := range graph.Nodes {
		if node.Id == Id {
			return &node
		}
	}
	return nil
}

func GenerateAffinityGraph(graph *Graph, edgeServers []*edge.EdgeServer) {
	hit := 0
	miss := 0
	var affinity float64
	for _, edgeServer := range edgeServers {
		AddNode(graph, edgeServer.Name)
	}

	var nodeA, nodeB *Node
	for i, edgeServerA := range edgeServers {
		nodeA = findNodeById(graph, edgeServerA.Name)
		for _, edgeServerB := range edgeServers {
			nodeB = findNodeById(graph, edgeServerB.Name)
			if edgeServerA.Id != edgeServerB.Id {
				for _, imgIdA := range edgeServerA.History {
					for _, imgIdB := range edgeServerB.History {
						if imgIdA == imgIdB {
							hit++
						} else {
							miss++
						}
					}
				}
			}
			if hit != 0 && miss != 0 {
				affinity = (float64(hit) / (float64(hit) + float64(miss))) * 100
			} else {
				affinity = 0
			}

			AddLine(graph, nodeA, nodeB, i, affinity)
			hit = 0
			miss = 0
			affinity = 0
		}
	}
	//for i := 1; i <= numOfEdge; i++ { // 1부터 numOfEdge 까지의 합
	//	for j := 1; j <= cnt; j++ {
	//		for _, startId := range edgeServers[i-1].History {
	//			hit = 0
	//			miss = 0
	//			numOfPulling = 0
	//			for _, endId := range edgeServers[i+j-1].History {
	//				numOfPulling++
	//				if startId == endId {
	//					hit++
	//				} else {
	//					miss++
	//				}
	//
	//			}
	//		}
	//		//AddLine(graph, &graph.Nodes[i-1], &graph.Nodes[i+j-1], numOfLine, (int)(hit/(hit+miss))*1000)
	//		affinity := float32(hit) / float32(hit+miss) * 100
	//		AddLine(graph, &graph.Nodes[i-1], &graph.Nodes[i+j-1], numOfLine, float64(affinity*100))
	//		//fmt.Println("affinity:", (affinity/numOfPulling)*100, "%")
	//		numOfLine++
	//	}
	//	cnt--
	//}
}

func SumOfVal(graph *Graph) int {
	sum := 0
	for _, line := range graph.Lines {
		sum += int(line.Val)
	}
	return sum
}

func ClusterGraph(graph *Graph, k int) []Graph {
	var graphs []Graph
	//var err error

	n := len(graph.Nodes) // the number of nodes in graph
	q := n / k            // the quotient of n divided by k
	r := n % k            // the remainder of n divided by k

	// k 값에 대한 조건 검사
	if n <= k {
		fmt.Println("Set k less than ", n)
	}

	// 클러스터의 수만큼 반복
	for i := 0; i < k; i++ {
		subgraph := Graph{}
		//fmt.Println("group", i)

		// 각 클러스터의 노드 수집
		for j := 1; j <= q; j++ {
			//fmt.Println("value:", i*q+j)
			node := graph.Nodes[i*q+j-1]
			subgraph.Nodes = append(subgraph.Nodes, node)
		}
		// 잔여 노드 확인
		if i == k-1 && r != 0 {
			for i := 0; i < r; i++ {
				subgraph.Nodes = append(subgraph.Nodes, graph.Nodes[q*k+i])
			}
		}

		// 글로벌 그래프의 Lines 인텍스 확인
		//for index, line := range graph.Lines {
		//	fmt.Println("index:", index, line.NodeA.Id, "->", line.NodeB.Id)
		//}

		// 서브 그래프의 노드에 대한 라인 수집
		for _, node := range subgraph.Nodes {
			sublines := []Line{}
			for _, line := range graph.Lines {
				if node.Id == line.NodeA.Id {
					sublines = append(sublines, line)
				}
			} // get node id's lines
			nodeId, err := strconv.ParseInt(node.Id, 10, 64)
			if err != nil {
				panic(err) // 또는 log.Fatal(err), panic(err) …
			}

			for l := nodeId - int64(i*q); l < int64(len(subgraph.Nodes)); l++ {
				//fmt.Println("node id:", node.Id)
				subgraph.Lines = append(subgraph.Lines, sublines[int64(len(subgraph.Nodes))-l-1])
			}
		}

		graphs = append(graphs, subgraph)
	}
	return graphs
}

func ElectReaderUsingOverhead(graph *Graph) string {

	var leaderId string
	minSum := 1000000

	for _, node := range graph.Nodes {
		sublines := []Line{}
		sum := 0
		for _, line := range graph.Lines {
			if node.Id == line.NodeA.Id || node.Id == line.NodeB.Id {
				sublines = append(sublines, line)
			}
		}
		for _, line := range sublines {
			sum += int(line.Val)
		}
		//fmt.Println("node", node.Id, "'s graph value:", sum)
		if minSum >= sum {
			leaderId = node.Id
			minSum = sum
		}
	}
	return leaderId
}

func ElectReaderUsingAffinity(graph *Graph) string {

	var leaderId string
	var maxSum float64
	maxSum = 0

	for _, node := range graph.Nodes {
		sublines := []Line{}
		var sum float64
		sum = 0

		leaderId = node.Id
		for _, line := range graph.Lines {
			if node.Id == line.NodeA.Id || node.Id == line.NodeB.Id {
				sublines = append(sublines, line)
			}
		}
		for _, line := range sublines {
			sum += line.Val
		}
		//fmt.Println("node", node.Id, "'s graph value:", sum)
		if maxSum <= sum {
			leaderId = node.Id
			maxSum = sum
		}
	}

	return leaderId
}
