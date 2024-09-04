package weather

import (
	"chat_service/repository/remote/net"
	"chat_service/tool_function"
	"chat_service/util"
	"context"
	"encoding/json"
	"strconv"
	"strings"
)

const (
	FuncName = "get_current_weather"
	FuncDesc = "查询龙华的实时天气预报，包括日期、城市、温度、湿度、降水概率、能见度、穿衣建议。"
	//FuncDesc      = "查询所在地区的穿衣建议。"
	ArgumentsDesc = "{}"

	CityNameFile = "./city_name.txt"
	CurPosition  = "龙华"
)

// Weather 天气handler
type Weather struct {
	FuncName string
	FuncDesc string
	ArgsDesc string
	cityMap  map[string]int
}

func init() {
	cityMap := make(map[string]int)
	cityList := util.ReadFileByLine(CityNameFile)
	for _, val := range cityList {
		tmp := strings.Split(val, ",")
		if len(tmp) != 2 {
			continue
		}
		num, _ := strconv.Atoi(tmp[1])
		cityMap[tmp[0]] = num
	}
	handler := &Weather{
		FuncName: FuncName,
		FuncDesc: FuncDesc,
		ArgsDesc: ArgumentsDesc,
		cityMap:  cityMap,
	}

	tool_function.RegisterFunction(FuncName, handler)
}

func (p *Weather) GetInfo() (string, string, string) {
	return p.FuncName, p.FuncDesc, p.ArgsDesc
}

func (p *Weather) Call(input string) string {
	args := &struct{}{}
	if err := json.Unmarshal([]byte(input), args); err != nil {
		return ""
	}

	// 1、获取当前地点
	id := p.cityMap[CurPosition]
	id = 4833

	// 2、获取当前天气
	weather := getTodayWeather(id)
	return weather
}

// WeatherInfo 天气信息
type WeatherInfo struct {
	CityName   string          `json:"cityname"`
	Today      string          `json:"today"`
	MaxTemp    string          `json:"max_temp"`
	MinTemp    string          `json:"min_temp"`
	NowTemp    string          `json:"now_temp"`
	Weather    string          `json:"weather"`
	Aqi        string          `json:"aqi"`
	Wind       string          `json:"wind"`
	Humidity   string          `json:"humidity"`
	Air        string          `json:"air"`
	Visibility string          `json:"visibility"`
	Shzhishu   []*WeatherIndex `json:"shzhishu"`
}

// WeatherIndex 天气指数
type WeatherIndex struct {
	Title string `json:"title"`
	Des   string `json:"des"`
}

func getTodayWeather(cityID int) string {
	url := "http://send.wxbus163.cn/weather/getToday"
	req := &struct {
		CityID int `json:"cityid"`
	}{
		CityID: cityID,
	}
	params := map[string]interface{}{
		"cityid": cityID,
	}
	rsp := &struct {
		Status int          `json:"status"`
		Msg    string       `json:"msg"`
		List   *WeatherInfo `json:"list"`
	}{}
	if err := net.HttpClientPost(context.Background(), url, params, nil, req, rsp); err != nil {
		return ""
	}
	body, _ := json.Marshal(rsp.List)
	return string(body)
}
