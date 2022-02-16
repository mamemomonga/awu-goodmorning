package goodmorning

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// 便利ツール: https://mholt.github.io/json-to-go/
// 便利ツール: github.com/davecgh/go-spew/spew

type WeatherJWAPlace struct {
	Name       string   `json:"name"`
	EnName     string   `json:"enName"`
	OfficeName string   `json:"officeName"`
	Children   []string `json:"children"`
}

type WeatherJWAArea struct {
	Centers  map[string]WeatherJWAPlace `json:"centers"`
	Offices  map[string]WeatherJWAPlace `json:"offices"`
	Class10s map[string]WeatherJWAPlace `json:"class10s"`
	Class15s map[string]WeatherJWAPlace `json:"class15s"`
	Class20s map[string]WeatherJWAPlace `json:"class20s"`
}

func (t *GoodMorning) GetWeatherArea() (data WeatherJWAArea, err error) {
	err = nil

	resp, err := http.Get("http://www.jma.go.jp/bosai/common/const/area.json")
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, _ := ioutil.ReadAll(resp.Body)

	if err = json.Unmarshal(b, &data); err != nil {
		log.Fatalf("JSON Unmarshal error: %V", err)
	}
	return data, nil
}

func (t *GoodMorning) WeatherAreaSelect(callback func(string)) (err error) {

	type refDT struct {
		name     string
		code     string
		children []string
	}

	area, err := t.GetWeatherArea()
	app := tview.NewApplication()
	root := tview.NewTreeNode("場所を選択").
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	// 子供を追加する処理
	// appender := func(tg *tview.TreeNode, cd string, ar WeatherJWAPlace) {
	// 	node := tview.NewTreeNode(fmt.Sprintf("  %s [%s]", ar.Name, cd)).
	// 		SetReference(refDT{code: cd, children: ar.Children, name: ar.Name}).
	// 		SetSelectable(true)
	// 	tg.AddChild(node)
	// }

	// 追加する処理
	add := func(target *tview.TreeNode, parent string) {
		// ルートメニュー(Centers)
		// 並び順: コード番号でソートする
		if parent == "" {
			centers := make([]string, len(area.Centers))
			idx := 0
			for cd := range area.Centers {
				centers[idx] = cd
				idx++
			}
			sort.Strings(centers)
			for _, cd := range centers {
				ar := area.Centers[cd]
				node := tview.NewTreeNode(fmt.Sprintf("  %s [%s]", ar.Name, cd)).
					SetReference(refDT{code: cd, children: ar.Children, name: ar.Name}).
					SetSelectable(true)
				target.AddChild(node)
			}
			return
		}
		// サブメニュー(Offices)
		// 並び順: CentersのChildren順に準ずる
		for cd, ar := range area.Offices {
			if cd == parent {
				node := tview.NewTreeNode(fmt.Sprintf("  %s [%s]", ar.Name, cd)).
					SetReference(refDT{code: cd, children: nil, name: ar.Name}).
					SetSelectable(true)
				target.AddChild(node)
			}
		}
	}

	add(root, "")

	// 選択されたときの処理
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			app.Stop()
			return
		}
		rf := ref.(refDT)
		children := node.GetChildren()
		if len(children) == 0 {
			if len(rf.children) == 0 {
				// 子に追加する処理がなければ終了
				app.Stop()
				callback(rf.code)
			} else {
				for _, c := range rf.children {
					add(node, c)
				}
			}
		} else {
			node.SetExpanded(!node.IsExpanded())
		}
	})

	if err := app.SetRoot(tree, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
	return
}
