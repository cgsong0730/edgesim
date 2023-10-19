package edge

import (
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

var RegistryServerList []RegistryServer

type ContainerImage struct {
	Id   int
	Size int
}

type EdgeServer struct {
	Id               int
	NumOfImage       int
	MaxCacheSize     int
	CurrentCacheSize int
	LocalImages      []ContainerImage
	HitCount         int
	MissCount        int
	RegistryServers  []RegistryServer
	FirstRegistry    EdgeRegistryServer
	SecondRegistry   EdgeRegistryServer
	History          []int
	AffinityOverhead int
	NetworkOverhead  int
}

type RegistryServer struct {
	Id       int
	Images   []ContainerImage
	Overhead int
}

type EdgeRegistryServer struct {
	NodeId            int
	MaxNumOfImage     int
	CurrentNumOfImage int
	Images            []ContainerImage
}

func CleanCache(s *EdgeServer) {

	s.NumOfImage = 0
	s.CurrentCacheSize = 0
	s.LocalImages = nil

}

func DownloadImage(s *EdgeServer, img ContainerImage) error {
	var err error

	if s.MaxCacheSize >= s.CurrentCacheSize+img.Size {
		s.NumOfImage++
		s.CurrentCacheSize += img.Size
		s.LocalImages = append(s.LocalImages, img)
		return nil
	} else {
		return err
	}
}

func ImagePulling(s *EdgeServer, i int, nr int) {

	var str string
	var b []byte

	f, err := os.OpenFile("./result.txt", os.O_APPEND|os.O_RDWR, 0755)
	defer f.Close()
	check(err)

	//var err error
	s.History = append(s.History, i)

	// non-pulling
	for _, img := range s.LocalImages {
		if i == img.Id {
			//fmt.Println("non-pulling")
			s.HitCount++

			str = fmt.Sprintf("%d\n", 0)
			b = []byte(str)
			_, err = f.Write(b)
			check(err)

			return
			//return nil
		}
	}

	// pulling - edge registry server 1
	for _, img := range s.FirstRegistry.Images {
		if i == img.Id {
			//fmt.Println("pulling - first edge server :", s.NetworkOverhead)

			str = fmt.Sprintf("%d\n", s.NetworkOverhead)
			b = []byte(str)
			_, err = f.Write(b)
			check(err)

			s.MissCount++
			if DownloadImage(s, img) != nil {
				return
				//return nil
			} else {
				return
				//return err
			}
		}
	}

	// pulling - edge registry server 2
	for _, img := range s.SecondRegistry.Images {
		if i == img.Id {
			//fmt.Println("pulling - second edge server :", s.AffinityOverhead)
			//fmt.Println("pulling - second edge server :", s.NetworkOverhead)

			str = fmt.Sprintf("%d\n", s.NetworkOverhead)
			b = []byte(str)
			_, err = f.Write(b)
			check(err)

			s.MissCount++
			if DownloadImage(s, img) != nil {
				//return nil
				return
			} else {
				//return err
				return
			}
		}
	}

	// pulling - retrieves container image from multiple registry servers.
	for _, r := range s.RegistryServers {
		for _, img := range r.Images {
			if i == img.Id {
				//fmt.Println("pulling - remote registry server:", r.Overhead+nr)

				str = fmt.Sprintf("%d\n", r.Overhead+nr)
				b = []byte(str)
				_, err = f.Write(b)
				check(err)

				s.MissCount++
				if DownloadImage(s, img) != nil {
					//return nil
					return
				} else {
					//return err
					return
				}
			}
		}
	}

	//return err
}

func Init() {
	var img ContainerImage
	var r RegistryServer
	var index int = 0

	for i := 1; i <= 10; i++ {
		r.Id = i
		for j := 1 + index*100; j <= (index+1)*100; j++ {
			img.Id = j
			img.Size = 300
			r.Images = append(r.Images, img)
		}
		RegistryServerList = append(RegistryServerList, r)
		index++
		r.Images = nil
	}
}

func InitRegistryServer(rs *RegistryServer, numOfContainer int, sizeOfContainer int, firstId int, lastId int) {
	for i := firstId; i <= lastId; i++ {
		img := ContainerImage{
			i,
			sizeOfContainer,
		}
		rs.Images = append(rs.Images, img)
	}
}

// func CreateEdgeRegistryServer(edgeServerList []EdgeServer, leaderId int) EdgeRegistryServer {
func CreateEdgeRegistryServer(edgeServer EdgeServer, leaderId int, numOfImage int) EdgeRegistryServer {

	edgeRegistryServer := EdgeRegistryServer{
		NodeId:            leaderId,
		MaxNumOfImage:     500,
		CurrentNumOfImage: 0,
		Images:            nil,
	}

	//for i := 0; i < numOfImage; i++ {
	//	img := ContainerImage{
	//		edgeServer.History[i],
	//		300,
	//	}
	//	if edgeRegistryServer.CurrentNumOfImage < edgeRegistryServer.MaxNumOfImage {
	//		edgeRegistryServer.Images = append(edgeRegistryServer.Images, img)
	//		edgeRegistryServer.CurrentNumOfImage++
	//	}
	//}

	for _, imageId := range edgeServer.History {
		img := ContainerImage{
			imageId,
			300,
		}
		if edgeRegistryServer.CurrentNumOfImage < edgeRegistryServer.MaxNumOfImage {
			edgeRegistryServer.Images = append(edgeRegistryServer.Images, img)
			edgeRegistryServer.CurrentNumOfImage++
		}
	}

	//for _, edge := range edgeServerList {
	//	if edge.Id == leaderId {
	//		for _, imageId := range edge.History {
	//			img := ContainerImage{
	//				imageId,
	//				300,
	//			}
	//			if edgeRegistryServer.CurrentNumOfImage < edgeRegistryServer.MaxNumOfImage {
	//				edgeRegistryServer.Images = append(edgeRegistryServer.Images, img)
	//				edgeRegistryServer.CurrentNumOfImage++
	//			}
	//		}
	//	}
	//}

	return edgeRegistryServer
}
