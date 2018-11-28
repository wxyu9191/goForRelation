package main

/**
	1, capture the message, the extract the source ip and destination ip
	2, pass the source ip and destination ip to loadRelation.go to search there node_path
 */
import (
	_ "net"
	"fmt"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"strings"
	loadR "relationrate/loadRelation"
	"relationrate/show"
	"log"
	"os"
	"strconv"
	"net/http"
)

//var goal map[string][]string
var goal = make(map[string][]string, 10000)
//var goalIp map[string]int

//func loadWorker() (newGoal map[string]int) {
//	i := 0
//	c := cron.New()
//	spec := "0 0 * * * ?"
//	c.AddFunc(spec, func() {
//		log.Println("running days:", i)
//		goal = make(map[string]int)
//	})
//	c.Start()
//	select {}
//
//	return goal
//}

//var goal = make(map[string]int) //频率

func httpServ()  {
	http.HandleFunc("/", show.ResponseIp)
	//fmt.Println("----------走到这里了333")
	//time.Sleep(1000)
	errInfo := http.ListenAndServe(":19968", nil)
	if errInfo != nil {
		log.Fatal("ListenAndServe: ", errInfo)
	}
}

func main() {

	go httpServ()

	fmt.Println("0000-----------------0000")

	file, err := os.Create("test.log")
	if err != nil {
		log.Fatalln("fail to create test.log file!")
	}
	logger := log.New(file, "", log.LstdFlags|log.Llongfile)

	log.Println("packet start...")
	logger.Println("packet start...")

	//申请空间

	deviceName := "eth0"
	snapLen := int32(65535)

	log.Printf("device: %v, snapLen: %v", deviceName, snapLen)
	logger.Printf("device: %v, snapLen: %v", deviceName, snapLen)

	//打开网络接口，开始抓取在线数据
	handle, err := pcap.OpenLive(deviceName, snapLen, true, pcap.BlockForever)
	if err != nil {
		fmt.Printf("pcap open live failed: %v", err)
		return
	}

	//defer handle.Close()

	//抓包
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packetSource.NoCopy = true

	//goal = loadWorker() //定时清空 csv

	for packet := range packetSource.Packets() { //循环读取每条数据流
		if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
			log.Println("unexpected packet")
			logger.Println("unexpected packet")
			continue
		}

		//fmt.Printf("packet type is %T\n", packet)
		fmt.Println("**********源ip和目的ip***********")
		//fmt.Printf("packet: %v\n", packet)
		members := strings.Split(packet.String(), " ")

		srcIp := strings.Split(members[56], "=")[1]
		desIp := strings.Split(members[57], "=")[1]

		fmt.Printf("src ip: %v\n", srcIp)
		fmt.Printf("des ip: %v\n", desIp)

		//判断是否是合法ip
		if legalIp(srcIp) == false || legalIp(desIp) == false {
			fmt.Println("ip 格式错误")
			continue
		}

		//ip 获取源Ip和目的ip的服务列表
		resSrc, err := loadR.FetchIp(srcIp)
		resDes, err := loadR.FetchIp(desIp)
		if err != nil {
			log.Println(err)
		}
		//查看是否有服务
		if resSrc == nil || resDes == nil {
			continue
		}
		fmt.Println("源IP对应的服务:", resSrc)
		fmt.Println("目的IP对应的服务:", resDes)

		//做个判断 看是不是跨机房的调用
		var tmpIp string
		across := check(srcIp, desIp)
		fmt.Println("+++++++++++", across)
		if across == true {
			tmpIp = srcIp + "---" + desIp
			fmt.Println("ip组合为: ", tmpIp)
		}

		if value, ok := goal[tmpIp]; ok {
			//value的第一个值 表示为彼此对应服务的次数
			fmt.Println("$$$$$$$$$不是第一次出现$$$$$$$$")
			newCount, err := strconv.Atoi(value[0])
			if err != nil {
				fmt.Println(err)
				fmt.Println(value[0] + "字符串转换成整数失败")
			} else {
				newCount += 1
				value[0] = strconv.Itoa(newCount)
			}
		} else {
			//只有在第一次将 IP对 放入map的时候，添加相应的服务
			fmt.Println("@@@@@@@@第一次出现@@@@@@@@@")
			goal[tmpIp] = make([]string, 1)
			goal[tmpIp][0] = "1"
			goal[tmpIp] = appendService(goal[tmpIp], resSrc, resDes)
		}

		//定时清空缓存

		fmt.Println(goal)
		//进行排序,并返回结果
		show.Showers(goal)
	}
	defer handle.Close()
}

func legalIp(s string) (ok bool) {
	stars := strings.Split(s, ".")
	for i := range stars {
		tmp, _ := strconv.Atoi(stars[i])
		if (tmp <= 255 && 0 <= tmp) && (len(stars) == 4) {
			return true
		} else {
			return false
		}
	}
	return ok
}

func appendService(res []string, resSrc []string, resDes []string) (newRes []string) {
	if len(resSrc) != 0 {
		for i := range resSrc {
			newRes = append(res, resSrc[i])
		}
	}

	if len(resDes) != 0 {
		for j := range resDes {
			newRes = append(newRes, resDes[j])
		}
	}

	return newRes
}

//判断是不是跨机房服务， 其中一个是136，一个不是，就是跨机房
func check(ip1 string, ip2 string) (checkCross bool) {
	//if ip1 and ip2 is cross, check = true
	ip1s := strings.Split(ip1, ".")
	ip2s := strings.Split(ip2, ".")
	if (ip1s[1] == "136" && ip2s[1] != "136") || (ip1s[1] != "136" && ip2s[1] == "136") {
		checkCross = true
	} else {
		checkCross = false
	}
	return checkCross
}

//判断哪些服务是跨机房服务，也就是源p对应的服务信息于目的ip对应的服务信息取并集
//func checker(s []string, s2 []string) (goalService []string) {
//	for _, value := range s {
//		for _, value2 := range s2 {
//			if strings.Compare(value, value2) == 0 {
//				goalService = append(goalService, value)
//			}
//		}
//	}
//	return goalService
//}
