package main

import (
	"flag"
	"log"
	"os"

	"github.com/mamemomonga/awu-goodmorning/don"
	"github.com/mamemomonga/awu-goodmorning/goodmorning"
)

const (
	_ = iota
	RUNMODE_KANAGAWA
	RUNMODE_SELECT
	RUNMODE_TOOT
)

func main() {
	mode := 0
	{
		var (
			flgKngw bool
			flgSel  bool
			flgDon  bool
		)
		flag.BoolVar(&flgKngw, "k", false, "神奈川の気象情報を表示")
		flag.BoolVar(&flgSel, "s", false, "地域を選んで気象情報を表示")
		flag.BoolVar(&flgDon, "m", false, "地域を選んでマストドンに投稿")
		flag.Parse()
		if flgKngw {
			mode = RUNMODE_KANAGAWA
		} else if flgSel {
			mode = RUNMODE_SELECT
		} else if flgDon {
			mode = RUNMODE_TOOT
		}
	}

	gm := goodmorning.NewGoodMorning()
	switch mode {
	case RUNMODE_KANAGAWA:
		// 140000は神奈川県の気象情報
		gm.Print("140000")
	case RUNMODE_SELECT:
		gm.WeatherAreaSelect(func(code string) {
			gm.Print(code)
		})
	case RUNMODE_TOOT:
		gm.WeatherAreaSelect(func(code string) {
			err := don.TootIt(gm.String(code))
			if err != nil {
				log.Fatal(err)
			}
		})
	default:
		flag.PrintDefaults()
		os.Exit(1)
	}
}
