package main

import (
	hyentity "chat_service/entity/hunyuan_msg"
	"fmt"
	"reflect"

	_ "chat_service/logic/role/alfred"
	_ "chat_service/logic/role/translator"
	_ "chat_service/tool_function"
)

type CityInfo struct {
	ID       int    `json:"id"`
	Cityname string `json:"cityname"`
}

func main() {
	fmt.Println("怀义怀义，知我心意！")

	//msg := &hyentity.HyMessage{
	//	Role:    hyentity.HyRoleUser,
	//	Content: "龙华今天穿什么衣服？",
	//}
	//_, err := translator.Role.Chat("", msg, hyentity.NewChatConfig())
	//fmt.Println(err)
	//chatCfg := hyentity.NewChatConfig()
	//chatCfg.ToolChoice = "custom"
	//_, err := role.GetRole("Alfred").Chat("", msg, chatCfg)

	//err := role.GetRole("Alfred").Do("今天穿什么衣服？")
	//fmt.Println(err)
	rsp := &hyentity.HyChatRsp{}
	fmt.Println(reflect.TypeOf(rsp), reflect.ValueOf(rsp))

	//w := &weather.Weather{}
	//str := w.Call("{}")
	//fmt.Println("str: ", str)

	//err := role.GetRole("Translator").Do("今天是9月5号了。")
	//fmt.Println(err)

	//url := "http://send.wxbus163.cn/weather/getCityList"
	//rsp := &struct {
	//	Code int         `json:"code"`
	//	List []*CityInfo `json:"list"`
	//}{}
	//_ = net.HttpClientGet(context.Background(), url, nil, nil, rsp)
	//res := make([]string, 0)
	//for _, val := range rsp.List {
	//	res = append(res, fmt.Sprintf("%d,%s\n", val.ID, val.Cityname))
	//}
	//WriteFile("city_name.txt", res)

	/*
		//
		//tc, _ := tcvectordb.NewTCVectorDB(url, user, key)
		//ctx := context.Background()
		//tc.ListDatabase(ctx)
	*/

	/*


		cli := hunyuan.NewHyClient(secretID, secretKey)
		msg := &hyentity.HyMessage{
			Role:    hyentity.HyRoleUser,
			Content: "今天天气怎样?",
		}
		cfg := &hyentity.HyChatConfig{
			IsStream: false,
			RspChan:  make(chan string),
		}
		tool := hyentity.NewHyTool(tool_function.FunctionList["GetCurWeather"])

		//非流式
		rsp, _ := cli.Chat([]*hyentity.HyMessage{msg}, cfg, []*hyentity.HyTool{tool})
		fmt.Println("rsp: ", rsp)
		fmt.Println("rsp2: ", rsp.Choices[0])
	*/

	/*
		// 流式
		//rspChan := make(chan string)
		//go func() {
		//	_ = cli.Chat(msg, &entity.HyChatConfig{
		//		IsStream: true,
		//		RspChan:  rspChan,
		//		GapTime:  time.Second / 5,
		//	})
		//}()
		//
		//index := 0
		//for data := range rspChan {
		//	fmt.Println("rspChan: ", index, data)
		//	index++
		//}
	*/
}
