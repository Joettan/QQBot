package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type WeatherService struct {
}

var (
	W *WeatherService
)

const (
	weatherUrl = "https://api.qweather.com/v7/weather/3d?key=2d86ba4c63a04e0885032397e2a4795a&location=%s"
	cityUrl    = "https://geoapi.qweather.com/v2/city/lookup?key=2d86ba4c63a04e0885032397e2a4795a&location=%s"
)

type LocationResponse struct {
	Code     string     `json:"code"`
	Location []Location `json:"location"`
	Refer    Refer      `json:"refer"`
}

type Location struct {
	Name      string `json:"name"`
	ID        string `json:"id"`
	Lat       string `json:"lat"`
	Lon       string `json:"lon"`
	Adm2      string `json:"adm2"`
	Adm1      string `json:"adm1"`
	Country   string `json:"country"`
	Tz        string `json:"tz"`
	UtcOffset string `json:"utcOffset"`
	IsDst     string `json:"isDst"`
	Type      string `json:"type"`
	Rank      string `json:"rank"`
	FxLink    string `json:"fxLink"`
}

// 定义结构体来反序列化上述数据
type WeatherResp struct {
	Code       string  `json:"code"`
	Daily      []Daily `json:"daily"`
	FxLink     string  `json:"fxLink"`
	Refer      Refer   `json:"refer"`
	UpdateTime string  `json:"updateTime"`
}

type Daily struct {
	Cloud          string `json:"cloud"`
	FxDate         string `json:"fxDate"`
	Humidity       string `json:"humidity"`
	IconDay        string `json:"iconDay"`
	IconNight      string `json:"iconNight"`
	MoonPhase      string `json:"moonPhase"`
	MoonPhaseIcon  string `json:"moonPhaseIcon"`
	Moonrise       string `json:"moonrise"`
	Moonset        string `json:"moonset"`
	Precip         string `json:"precip"`
	Pressure       string `json:"pressure"`
	Sunrise        string `json:"sunrise"`
	Sunset         string `json:"sunset"`
	TempMax        string `json:"tempMax"`
	TempMin        string `json:"tempMin"`
	TextDay        string `json:"textDay"`
	TextNight      string `json:"textNight"`
	UvIndex        string `json:"uvIndex"`
	Vis            string `json:"vis"`
	Wind360Day     string `json:"wind360Day"`
	Wind360Night   string `json:"wind360Night"`
	WindDirDay     string `json:"windDirDay"`
	WindDirNight   string `json:"windDirNight"`
	WindScaleDay   string `json:"windScaleDay"`
	WindScaleNight string `json:"windScaleNight"`
	WindSpeedDay   string `json:"windSpeedDay"`
	WindSpeedNight string `json:"windSpeedNight"`
}

type Refer struct {
	License []string `json:"license"`
	Sources []string `json:"sources"`
}

func InitWeatherService() {
	W = &WeatherService{}
}

func (W WeatherService) FetchWeather(location string) (*WeatherResp, error) {

	resp, err := http.Get(fmt.Sprintf(weatherUrl, location))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var weatherResp WeatherResp
	err = json.Unmarshal(body, &weatherResp)
	if err != nil {
		return nil, err
	}

	return &weatherResp, nil
}

func (W *WeatherService) FetchLocation(location string) (string, error) {

	resp, err := http.Get(fmt.Sprintf(cityUrl, location))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var locationResp LocationResponse
	err = json.Unmarshal(body, &locationResp)
	if err != nil {
		return "", err
	}
	if len(locationResp.Location) == 0 {
		return "", nil
	}

	return locationResp.Location[0].ID, nil
}
