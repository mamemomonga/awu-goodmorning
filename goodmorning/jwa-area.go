package goodmorning

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

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

func (t *GoodMorning) weatherAreaFlatten(s WeatherJWAArea) (data map[string]WeatherJWAPlace) {
	data = map[string]WeatherJWAPlace{}
	for c, p := range s.Centers {
		data[c] = p
	}
	for c, p := range s.Offices {
		data[c] = p
	}
	//	for c, p := range s.Class10s {
	//		data[c] = p
	//	}
	//	for c, p := range s.Class15s {
	//		data[c] = p
	//	}
	//	for c, p := range s.Class20s {
	//		data[c] = p
	//	}
	return
}

func (t *GoodMorning) WeatherAreaSelect(callback func(string)) (err error) {

	type refDT struct {
		name     string
		code     string
		children []string
	}

	arearaw, err := t.GetWeatherArea()
	area := t.weatherAreaFlatten(arearaw)

	app := tview.NewApplication()
	root := tview.NewTreeNode("場所を選択").
		SetColor(tcell.ColorRed)
	tree := tview.NewTreeView().
		SetRoot(root).
		SetCurrentNode(root)

	appender := func(tg *tview.TreeNode, cd string, ar WeatherJWAPlace) {
		node := tview.NewTreeNode(fmt.Sprintf("  %s [%s]", ar.Name, cd)).
			SetReference(refDT{code: cd, children: ar.Children, name: ar.Name}).
			SetSelectable(true)
		tg.AddChild(node)
	}

	add := func(target *tview.TreeNode, parent string) bool {
		hasChild := false
		if parent == "" {
			for cd, ar := range arearaw.Centers {
				appender(target, cd, ar)
				hasChild = true
			}
		} else {
			for cd, ar := range area {
				if cd == parent {
					appender(target, cd, ar)
					hasChild = true
				}
			}
		}
		return hasChild
	}
	add(root, "")
	tree.SetSelectedFunc(func(node *tview.TreeNode) {
		ref := node.GetReference()
		if ref == nil {
			app.Stop()
			return // Selecting the root node does nothing.
		}
		rf := ref.(refDT)
		children := node.GetChildren()
		last := false
		if len(children) == 0 {
			if len(rf.children) == 0 {
				last = true
			}
			for _, c := range rf.children {
				if !add(node, c) {
					last = true
				}
			}
			if last {
				app.Stop()
				callback(rf.code)
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
