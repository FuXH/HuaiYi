package weather

import (
	"chat_service/repository/remote/net"
	"chat_service/tool_function"
	"chat_service/util"
	"context"
	"encoding/json"
	"fmt"
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
		fmt.Println(err)
		return ""
	}

	// 1、获取当前地点
	province, city, country := "广东省", "深圳市", "龙华区"

	// 2、获取当前天气
	weather := getTodayWeather(province, city, country)
	return weather
}

// WeatherInfo 天气信息
type WeatherInfo struct {
	Observe    *Observe               `json:"observe"`     // 实时观测数据
	Forecast1H map[string]*Forecast1H `json:"forecast_1h"` // 1小时预测数据
	Index      *Index                 `json:"index"`       // 天气指数
	Alarm      map[string]*AlarmInfo  `json:"alarm"`       // 警告
	Tips       *Tips                  `json:"tips"`        // 提示
	Rise       map[string]*Rise       `json:"rise"`        // 日出日落
	Air        *Air                   `json:"air"`         // 空气质量
}

type Observe struct {
	Degree            string `json:"degree"`
	Humidity          string `json:"humidity"`
	Precipitation     string `json:"precipitation"`
	Pressure          string `json:"pressure"`
	UpdateTime        string `json:"update_time"`
	Weather           string `json:"weather"`
	WeatherBgPag      string `json:"weather_bg_pag"`
	WeatherCode       string `json:"weather_code"`
	WeatherColor      string `json:"weather_color"`
	WeatherShort      string `json:"weather_short"`
	WindDirection     string `json:"wind_direction"`
	WindDirectionName string `json:"wind_direction_name"`
	WindPower         string `json:"wind_power"`
}

type Forecast1H struct {
	Degree        string `json:"degree"`
	UpdateTime    string `json:"update_time"`
	Weather       string `json:"weather"`
	WeatherCode   string `json:"weather_code"`
	WeatherShort  string `json:"weather_short"`
	WindDirection string `json:"wind_direction"`
	WindPower     string `json:"wind_power"`
}

type Index struct {
	Airconditioner *IndexInfo `json:"airconditioner"`
	Allergy        *IndexInfo `json:"allergy"`
	Chill          *IndexInfo `json:"chill"`
	Clothes        *IndexInfo `json:"clothes"`
	Comfort        *IndexInfo `json:"comfort"`
	Dry            *IndexInfo `json:"dry"`
	Drying         *IndexInfo `json:"drying"`
	Heatstroke     *IndexInfo `json:"heatstroke"`
	Makeup         *IndexInfo `json:"makeup"`
	Mood           *IndexInfo `json:"mood"`
	Morning        *IndexInfo `json:"morning"`
	Sports         *IndexInfo `json:"sports"`
	Sunscreen      *IndexInfo `json:"sunscreen"`
	Tourism        *IndexInfo `json:"tourism"`
	Traffic        *IndexInfo `json:"traffic"`
	Ultraviolet    *IndexInfo `json:"ultraviolet"`
	Umbrella       *IndexInfo `json:"umbrella"`
}

type IndexInfo struct {
	Detail string `json:"detail"`
	Info   string `json:"info"`
	Name   string `json:"name"`
}

type AlarmInfo struct {
	City       string `json:"city"`
	Detail     string `json:"detail"`
	Name       string `json:"name"`
	Province   string `json:"province"`
	TypeCode   string `json:"type_code"`
	TypeName   string `json:"type_name"`
	UpdateTime string `json:"update_time"`
}

type Tips struct {
	Observe map[string]string `json:"observe"`
}

type Rise struct {
	Sunrise string `json:"sunrise"`
	Sunset  string `json:"sunset"`
	Time    string `json:"time"`
}

type Air struct {
	Aqi        int    `json:"aqi"`
	AqiLevel   int    `json:"aqi_level"`
	AqiName    string `json:"aqi_name"`
	Co         string `json:"co"`
	No2        string `json:"no2"`
	O3         string `json:"o3"`
	Pm10       string `json:"pm10"`
	Pm25       string `json:"pm2.5"`
	So2        string `json:"so2"`
	UpdateTime string `json:"update_time"`
}

func getTodayWeather(province, city, country string) string {
	url := "https://wis.qq.com/weather/common"
	params := map[string]interface{}{
		"source":       "pc",
		"weather_type": "observe|forecast_1h|index|alarm|tips|rise|air",
		"province":     province,
		"city":         city,
		"country":      country,
	}
	rsp := &struct {
		Data    *WeatherInfo `json:"data"`
		Message string       `json:"message"`
		Status  int          `json:"status"`
	}{}
	if err := net.HttpClientGet(context.Background(), url, params, nil, rsp); err != nil {
		fmt.Println("getTodayWeather fail, err: ", err)
		return ""
	}
	body, _ := json.Marshal(rsp)
	fmt.Println("rsp: ", rsp)
	return string(body)
}
