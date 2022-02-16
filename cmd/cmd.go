package main

import "github.com/mamemomonga/awu-goodmorning/goodmorning"

func main() {
	gm := goodmorning.NewGoodMorning()
	// 140000は神奈川県の気象情報
	gm.Print("140000")
}
