package goodmorning

//パッケージインポート
import (
	"fmt"
	"log"
)

// GoodMorning is 構造体
type GoodMorning struct {
}

// NewGoodMorning is コンストラクタ
func NewGoodMorning() (t *GoodMorning) {
	t = new(GoodMorning)
	return t
}

func (t *GoodMorning) Say() {
	di := t.GetDateInfo()

	// 140000は神奈川県の気象情報
	wa, err := t.GetWeather("140000")
	if err != nil {
		log.Fatal(err)
	}

	//区切り線
	fmt.Println(fmtHR)
	//1行目
	fmt.Println("今日は" + di.Date + "です。" + "今年は残り" + di.RemainDays + "日です")
	//2行目
	fmt.Println("現在の時刻は" + di.RemainDays + "です。")
	//区切り線
	fmt.Println(fmtHR)

	fmt.Println(wa.AreaName + "の天気")
	fmt.Println(wa.Text)
	fmt.Println("(" + wa.Publish + "発表）")

	//区切り線
	fmt.Println(fmtHR)
}
