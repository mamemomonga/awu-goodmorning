package goodmorning

//パッケージインポート
import (
	"strconv"
	"time"
)

type DateInfo struct {
	Date       string
	Time       string
	RemainDays string
}

// DateInfo is 日付の情報
func (t *GoodMorning) GetDateInfo() DateInfo {

	//日付を取得
	tm := time.Now()
	//曜日を日本語に変換
	weekday := [...]string{"日", "月", "火", "水", "木", "金", "土"}
	//接尾辞を付けて変数に格納
	weekday_format := weekday[tm.Weekday()] + "曜日"
	var (
		nowyear int = tm.Year()
		// nowmonth = t.Month()
		// nowday int = t.Day()
		nowyearday int = tm.YearDay()
	)

	//残り日数を計算
	var daynum int
	if t.isLeapYear(nowyear) {
		//うるう年の日数
		daynum = 366
	} else {
		//うるう年でない年の日数
		daynum = 365
	}
	//残りの日数
	lastyearday := int(daynum - nowyearday)

	return DateInfo{
		Date:       tm.Format(fmtDate) + weekday_format,
		Time:       tm.Format(fmtTime),
		RemainDays: strconv.Itoa(lastyearday),
	}
}

// うるう年かどうか判定する
func (t *GoodMorning) isLeapYear(year int) bool {
	if year%400 == 0 { // 400で割り切れたらうるう年
		return true
	} else if year%100 == 0 { // 100で割り切れたらうるう年じゃない
		return false
	} else if year%4 == 0 { // 4で割り切れたらうるう年
		return true
	} else {
		return false
	}
}
