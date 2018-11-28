package show

import (
	"sort"
	"fmt"
	"strconv"
	"net/http"
		"encoding/json"
	)

type ValSorter struct {
	Keys   []string	   `json:"ipCouple"`
	Values []int	   `json:"count"`
}

var result = &ValSorter{} //最终展示的结果

func Showers(goal map[string][]string) {

	//goal = map[string]int{"alpha": 34, "bravo": 56, "charlie": 23,
	//	"delta": 87, "echo": 56, "foxtrot": 12,
	//	"golf": 34, "hotel": 16, "indio": 87,
	//	"juliet": 65, "kili": 43, "lima": 98}

	fmt.Println("Show the sort of service between sh and wj: ")

	res := NewValSorter(goal)
	res.Sort()
	//fmt.Printf("%v\n", *res)

	result = res
	//fmt.Printf("%v\n", *result)

	//writers(result)
}

func ResponseIp(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("----------走到这里了")
	//time.Sleep(1000)
	body, err := json.Marshal(result)
	if err != nil{
		panic(err.Error())
	}
	fmt.Fprintf(w, string(body))
}

//打印csv
//func writers(res *ValSorter) {
//	length := len(res.Keys)
//
//	ticker := time.NewTicker(time.Second * 10)
//	go func() {
//		for range ticker.C {
//			strTime := time.Now().Format("2006-01-02_15:04:05")
//			csvFile, err := os.Create("./" + strTime + ".csv")
//			if err != nil {
//				panic(err)
//			}
//			csvWriter := csv.NewWriter(csvFile)
//
//			for i := 0; i < length; i++{
//				sev := res.Keys[i]
//				count := strconv.Itoa(res.Values[i])
//
//				e := csvWriter.Write([]string{sev, count})
//				if e != nil {
//					fmt.Print("write failure")
//				}
//				//fmt.Println(sev, count)
//			}
//
//			csvWriter.Flush()
//			err2 := csvWriter.Error()
//			if err2 != nil {
//				fmt.Print("write failure2")
//			}
//			csvFile.Close()
//		}
//	}()
//
//}

//对map类型数据按value进行排序
func NewValSorter(m map[string][]string) *ValSorter {
	vs := &ValSorter{
		Keys:   make([]string, 0, len(m)),
		Values: make([]int, 0, len(m)),
	}
	for k, v := range m {
		vs.Keys = append(vs.Keys, k)
		count, _ := strconv.Atoi(v[0])
		vs.Values = append(vs.Values, count)
	}
	return vs
}

func (vs *ValSorter) Sort() {
	sort.Sort(vs)
}

func (vs *ValSorter) Len() int           { return len(vs.Values) }
func (vs *ValSorter) Less(i, j int) bool { return vs.Values[i] > vs.Values[j] }
func (vs *ValSorter) Swap(i, j int) {
	vs.Values[i], vs.Values[j] = vs.Values[j], vs.Values[i]
	vs.Keys[i], vs.Keys[j] = vs.Keys[j], vs.Keys[i]
}
