package edge

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
	Id                      int
	Name                    string
	NumOfImage              int
	MaxCacheSize            int
	CurrentCacheSize        int
	LocalImages             []ContainerImage
	HitCount                int
	MissCount               int
	RegistryServers         []RegistryServer
	FirstRegistry           EdgeRegistryServer
	SecondRegistry          EdgeRegistryServer
	FirstRegistryBandwidth  float64
	SecondRegistryBandwidth float64
	History                 []int
	AffinityOverhead        int
	NetworkOverhead         int
}

type RegistryServer struct {
	Id       int
	Images   []ContainerImage
	Overhead int
}

type EdgeRegistryServer struct {
	NodeId            string
	MaxNumOfImage     int
	CurrentNumOfImage int
	Images            []ContainerImage
}

func CleanCache(s *EdgeServer) {

	s.NumOfImage = 0
	s.CurrentCacheSize = 0
	s.LocalImages = nil

}

func DownloadImage(s *EdgeServer, img ContainerImage) {

	if s.MaxCacheSize >= s.CurrentCacheSize+img.Size {
		s.NumOfImage++
		s.CurrentCacheSize += img.Size
		s.LocalImages = append(s.LocalImages, img)
	} else {
		CleanCache(s)
		s.NumOfImage++
		s.CurrentCacheSize += img.Size
		s.LocalImages = append(s.LocalImages, img)
	}

}

//func ImagePullingB1(s *EdgeServer, i int, nr int) int {
//
//	//var err error
//	s.History = append(s.History, i)
//
//	// non-pulling
//	for _, img := range s.LocalImages {
//		if i == img.Id {
//			//fmt.Println("non-pulling")
//			s.HitCount++
//
//			//str = fmt.Sprintf("%d\n", 0)
//			//b = []byte(str)
//			//_, err = f.Write(b)
//			//check(err)
//
//			return 0
//		}
//	}
//
//	// pulling - retrieves container image from multiple registry servers.
//	for _, r := range s.RegistryServers {
//		for _, img := range r.Images {
//			if i == img.Id {
//				//fmt.Println("pulling - remote registry server:", r.Overhead+nr)
//
//				//str = fmt.Sprintf("%d\n", r.Overhead+nr)
//				//b = []byte(str)
//				//_, err = f.Write(b)
//				//check(err)
//
//				s.MissCount++
//				if DownloadImage(s, img) != nil {
//					//return nil
//					return r.Overhead + nr
//				} else {
//					//return err
//					return r.Overhead + nr
//				}
//			}
//		}
//	}
//
//	//return err
//	return 0
//}

//func ImagePullingB2(s *EdgeServer, i int, nr int) int {
//
//	//var err error
//	s.History = append(s.History, i)
//
//	// non-pulling
//	for _, img := range s.LocalImages {
//		if i == img.Id {
//			//fmt.Println("non-pulling")
//			s.HitCount++
//
//			//str = fmt.Sprintf("%d\n", 0)
//			//b = []byte(str)
//			//_, err = f.Write(b)
//			//check(err)
//
//			return 0
//		}
//	}
//
//	// pulling - edge registry server 1
//	for _, img := range s.FirstRegistry.Images {
//		if i == img.Id {
//			//fmt.Println("pulling - first edge server :", s.NetworkOverhead)
//
//			//str = fmt.Sprintf("%d\n", s.NetworkOverhead)
//			//b = []byte(str)
//			//_, err = f.Write(b)
//			//check(err)
//
//			s.MissCount++
//			if DownloadImage(s, img) != nil {
//				return s.NetworkOverhead
//				//return nil
//			} else {
//				return s.NetworkOverhead
//				//return err
//			}
//		}
//	}
//
//	// pulling - retrieves container image from multiple registry servers.
//	for _, r := range s.RegistryServers {
//		for _, img := range r.Images {
//			if i == img.Id {
//				//fmt.Println("pulling - remote registry server:", r.Overhead+nr)
//
//				//str = fmt.Sprintf("%d\n", r.Overhead+nr)
//				//b = []byte(str)
//				//_, err = f.Write(b)
//				//check(err)
//
//				s.MissCount++
//				if DownloadImage(s, img) != nil {
//					//return nil
//					return r.Overhead + nr
//				} else {
//					//return err
//					return r.Overhead + nr
//				}
//			}
//		}
//	}
//
//	//return err
//	return 0
//}

//func ImagePullingB3(s *EdgeServer, i int, nr int) int {
//
//	//var err error
//	s.History = append(s.History, i)
//
//	// non-pulling
//	for _, img := range s.LocalImages {
//		if i == img.Id {
//			//fmt.Println("non-pulling")
//			s.HitCount++
//
//			//str = fmt.Sprintf("%d\n", 0)
//			//b = []byte(str)
//			//_, err = f.Write(b)
//			//check(err)
//
//			return 0
//		}
//	}
//
//	// pulling - edge registry server 1
//	for _, img := range s.FirstRegistry.Images {
//		if i == img.Id {
//			//fmt.Println("pulling - first edge server :", s.NetworkOverhead)
//
//			//str = fmt.Sprintf("%d\n", s.NetworkOverhead)
//			//b = []byte(str)
//			//_, err = f.Write(b)
//			//check(err)
//
//			s.MissCount++
//			if DownloadImage(s, img) != nil {
//				return s.NetworkOverhead
//				//return nil
//			} else {
//				return s.NetworkOverhead
//				//return err
//			}
//		}
//	}
//
//	// pulling - edge registry server 2
//	for _, img := range s.SecondRegistry.Images {
//		if i == img.Id {
//			//fmt.Println("pulling - second edge server :", s.AffinityOverhead)
//			//fmt.Println("pulling - second edge server :", s.NetworkOverhead)
//
//			//str = fmt.Sprintf("%d\n", s.NetworkOverhead)
//			//b = []byte(str)
//			//_, err = f.Write(b)
//			//check(err)
//
//			s.MissCount++
//			if DownloadImage(s, img) != nil {
//				//return nil
//				return s.NetworkOverhead
//			} else {
//				//return err
//				return s.NetworkOverhead
//			}
//		}
//	}
//
//	// pulling - retrieves container image from multiple registry servers.
//	for _, r := range s.RegistryServers {
//		for _, img := range r.Images {
//			if i == img.Id {
//				//fmt.Println("pulling - remote registry server:", r.Overhead+nr)
//
//				//str = fmt.Sprintf("%d\n", r.Overhead+nr)
//				//b = []byte(str)
//				//_, err = f.Write(b)
//				//check(err)
//
//				s.MissCount++
//				if DownloadImage(s, img) != nil {
//					//return nil
//					return r.Overhead + nr
//				} else {
//					//return err
//					return r.Overhead + nr
//				}
//			}
//		}
//	}
//
//	//return err
//	return 0
//}

//func ImagePulling(s *EdgeServer, i int, nr int) int {
//
//	//var err error
//	s.History = append(s.History, i)
//
//	// non-pulling
//	for _, img := range s.LocalImages {
//		if i == img.Id {
//			//fmt.Println("non-pulling")
//			s.HitCount++
//
//			//str = fmt.Sprintf("%d\n", 0)
//			//b = []byte(str)
//			//_, err = f.Write(b)
//			//check(err)
//
//			return 0
//		}
//	}
//
//	// pulling - edge registry server 1
//	for _, img := range s.FirstRegistry.Images {
//		if i == img.Id {
//			//fmt.Println("pulling - first edge server :", s.NetworkOverhead)
//
//			//str = fmt.Sprintf("%d\n", s.NetworkOverhead)
//			//b = []byte(str)
//			//_, err = f.Write(b)
//			//check(err)
//
//			s.MissCount++
//			if DownloadImage(s, img) != nil {
//				return s.NetworkOverhead
//				//return nil
//			} else {
//				return s.NetworkOverhead
//				//return err
//			}
//		}
//	}
//
//	//pulling - edge registry server 2
//	for _, img := range s.SecondRegistry.Images {
//		if i == img.Id {
//			//fmt.Println("pulling - second edge server :", s.AffinityOverhead)
//			//fmt.Println("pulling - second edge server :", s.NetworkOverhead)
//
//			//str = fmt.Sprintf("%d\n", s.NetworkOverhead)
//			//b = []byte(str)
//			//_, err = f.Write(b)
//			//check(err)
//
//			s.MissCount++
//			if DownloadImage(s, img) != nil {
//				//return nil
//				return s.NetworkOverhead
//			} else {
//				//return err
//				return s.NetworkOverhead
//			}
//		}
//	}
//
//	// pulling - retrieves container image from multiple registry servers.
//	for _, r := range s.RegistryServers {
//		for _, img := range r.Images {
//			if i == img.Id {
//				//fmt.Println("pulling - remote registry server:", r.Overhead+nr)
//
//				//str = fmt.Sprintf("%d\n", r.Overhead+nr)
//				//b = []byte(str)
//				//_, err = f.Write(b)
//				//check(err)
//
//				s.MissCount++
//				if DownloadImage(s, img) != nil {
//					//return nil
//					return r.Overhead + nr
//				} else {
//					//return err
//					return r.Overhead + nr
//				}
//			}
//		}
//	}
//
//	//return err
//	return 0
//}

func ImagePullingWithData(s *EdgeServer, i int, numOfPulling int) (float64, int) {

	s.History = append(s.History, i)

	// non-pulling
	for _, img := range s.LocalImages {
		if i == img.Id {
			s.HitCount++
			return 0, 0
		}
	}

	// pulling - edge registry server 1
	for _, img := range s.FirstRegistry.Images {
		if i == img.Id {
			s.MissCount++
			DownloadImage(s, img)
			return float64(img.Size) / s.FirstRegistryBandwidth, 1
		}
	}

	// pulling - edge registry server 2
	for _, img := range s.SecondRegistry.Images {
		if i == img.Id {
			s.MissCount++
			DownloadImage(s, img)
			return float64(img.Size) / s.SecondRegistryBandwidth, 2
		}
	}

	// pulling - retrieves container image from multiple registry servers.
	for _, r := range s.RegistryServers {
		for _, img := range r.Images {
			if i == img.Id {
				s.MissCount++
				DownloadImage(s, img)
				return float64(r.Overhead), 3
			}
		}
	}

	return 0, 4
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

func CreateEdgeRegistryServer(edgeServerSubList []*EdgeServer, edgeServer *EdgeServer, leaderId string, numOfImage int) EdgeRegistryServer {

	edgeRegistryServer := EdgeRegistryServer{
		NodeId:            leaderId,
		MaxNumOfImage:     numOfImage,
		CurrentNumOfImage: 0,
		Images:            nil,
	}

	var ImageIdList []int

	edgeNum := len(edgeServerSubList)
	for _, edgeMember := range edgeServerSubList {

		if len(edgeMember.History) > 0 && len(edgeMember.History) > numOfImage/edgeNum {
			for i := len(edgeMember.History) - 1; i >= len(edgeMember.History)-numOfImage/edgeNum; i-- {
				ImageIdList = append(ImageIdList, edgeMember.History[i])
			}
			if len(edgeMember.History) > numOfImage/edgeNum {
				edgeMember.History = edgeMember.History[:len(edgeMember.History)-numOfImage/edgeNum] // 뒤에서 n개 제거
			} else {
				edgeMember.History = []int{} // n 이상이면 전체 제거
			}
		}
	}

	for _, imgId := range ImageIdList {
		img := ContainerImage{
			imgId,
			500,
		}
		if edgeRegistryServer.CurrentNumOfImage < edgeRegistryServer.MaxNumOfImage {
			edgeRegistryServer.Images = append(edgeRegistryServer.Images, img)
			edgeRegistryServer.CurrentNumOfImage++
		}
	}

	//for _, imgId := range edgeServer.History {
	//
	//	img := ContainerImage{
	//		imgId,
	//		300,
	//	}
	//	if edgeRegistryServer.CurrentNumOfImage < edgeRegistryServer.MaxNumOfImage {
	//		edgeRegistryServer.Images = append(edgeRegistryServer.Images, img)
	//		edgeRegistryServer.CurrentNumOfImage++
	//	}
	//}

	return edgeRegistryServer
}
