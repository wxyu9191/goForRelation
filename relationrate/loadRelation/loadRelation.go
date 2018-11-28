package loadRelation

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"strings"
)

var resultant = make(map[string][]string)

func FetchIp(ip string) (res []string, err error) {
	if res, ok := resultant[ip]; ok {
		return res, err
	}
	var srcUrl = "http://axe.yidian-inc.com/api/v1/service/resource/?ip=" + ip

	var resSrc []byte

	responseSrc, err := http.Get(srcUrl)
	if err != nil {
		return res, err
	}

	defer responseSrc.Body.Close()

	resSrc, _ = ioutil.ReadAll(responseSrc.Body)

	var results map[string]interface{}//从接口得到的返回关键字为results的列表结构
	json.Unmarshal(resSrc, &results)

	if v, ok := results["results"]; ok {//遍历results
		ws := v.([]interface{})
		for _, wsItem := range ws { //i result
			//tmp_res := make([]string, len(ws))
			node := wsItem.(map[string]interface{})
			if tmp, ok := node["node_path"]; ok {//找到每个result结构中的node_path字段
				nodePath := tmp.(string)
				ans := strings.Split(nodePath, "/")
				length := len(ans)
				res = append(res, ans[length-2])//切割后倒数第二个字段就是服务的名称
			}

		}
	}
	resultant[ip] = res
	fmt.Println("匹配服务树读取到的服务:", resultant[ip])
	//返回根据解析到的数据ip而得到的的服务列表
	return res, err
}
