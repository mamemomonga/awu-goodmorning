package goodmorning

//パッケージインポート
import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"time"
)

//お天気構造体
type WeatherJWA struct {
	PublishingOffice string    `json:"publishingOffice"`
	ReportDatetime   time.Time `json:"reportDatetime"`
	TargetArea       string    `json:"targetArea"`
	// HeadlineText     string    `json:"headlineText"`
	Text string `json:"text"`
}

type Weather struct {
	AreaName string
	Text     string
	Publish  string
}

// FetchWeather is 天気情報の取得
func (t *GoodMorning) GetWeather(pathCode string) (res Weather, err error) {
	err = nil
	//気象庁のjsonから取得

	url := fmt.Sprintf("https://www.jma.go.jp/bosai/forecast/data/overview_forecast/%s.json", pathCode)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	byteArray, _ := ioutil.ReadAll(resp.Body)

	jsonBytes := ([]byte)(byteArray)
	data := new(WeatherJWA)

	if err = json.Unmarshal(jsonBytes, data); err != nil {
		fmt.Println("JSON Unmarshal error:", err)
		return
	}
	//二重改行を消去
	rep := regexp.MustCompile(`\n\n`)
	data.Text = rep.ReplaceAllString(data.Text, "\n")

	rep2 := regexp.MustCompile(`　`)
	data.Text = rep2.ReplaceAllString(data.Text, "")

	res.AreaName = data.TargetArea
	res.Text = data.Text
	res.Publish = data.ReportDatetime.Format(fmtDate) + data.ReportDatetime.Format(fmtTime) + data.PublishingOffice

	return
}
