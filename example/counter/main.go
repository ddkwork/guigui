// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2024 The Guigui Authors

package main

import (
	"fmt"
	"image"
	"io/fs"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/hajimehoshi/guigui"
	"github.com/hajimehoshi/guigui/basicwidget"
	"github.com/hajimehoshi/guigui/layout"
)

type Root struct {
	guigui.DefaultWidget

	background  basicwidget.Background
	resetButton basicwidget.Button
	decButton   basicwidget.Button
	incButton   basicwidget.Button
	counterText basicwidget.Text

	counter int
}

func (r *Root) Build(context *guigui.Context, appender *guigui.ChildWidgetAppender) error {
	if files := ebiten.DroppedFiles(); files != nil {
		go func() {
			if err := fs.WalkDir(files, ".", func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return err
				}
				if d.IsDir() {
					return nil
				}
				println(path)

				return nil
			}); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}()
	}

	appender.AppendChildWidgetWithBounds(&r.background, context.Bounds(r))

	r.counterText.SetSelectable(true)
	r.counterText.SetBold(true)
	r.counterText.SetHorizontalAlign(basicwidget.HorizontalAlignCenter)
	r.counterText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
	r.counterText.SetScale(4)
	r.counterText.SetValue(fmt.Sprintf("%d", r.counter))

	r.resetButton.SetText("Reset")
	r.resetButton.SetOnUp(func() {
		r.counter = 0
	})
	context.SetEnabled(&r.resetButton, r.counter != 0)

	r.decButton.SetText("Decrement")
	r.decButton.SetOnUp(func() {
		r.counter--
	})

	r.incButton.SetText("Increment")
	r.incButton.SetOnUp(func() {
		r.counter++
	})

	u := basicwidget.UnitSize(context)
	gl := layout.GridLayout{
		Bounds: context.Bounds(r).Inset(u),
		Heights: []layout.Size{
			layout.FlexibleSize(1),
			layout.FixedSize(u),
		},
		RowGap: u,
	}
	appender.AppendChildWidgetWithBounds(&r.counterText, gl.CellBounds(0, 0))
	{
		gl := layout.GridLayout{
			Bounds: gl.CellBounds(0, 1),
			Widths: []layout.Size{
				layout.FixedSize(6 * u),
				layout.FlexibleSize(1),
				layout.FixedSize(6 * u),
				layout.FixedSize(6 * u),
			},
			ColumnGap: u / 2,
		}
		appender.AppendChildWidgetWithBounds(&r.resetButton, gl.CellBounds(0, 0))
		appender.AppendChildWidgetWithBounds(&r.decButton, gl.CellBounds(2, 0))
		appender.AppendChildWidgetWithBounds(&r.incButton, gl.CellBounds(3, 0))
	}

	return nil
}

func main() {
	op := &guigui.RunOptions{
		Title:         "Counter",
		WindowMinSize: image.Pt(600, 300),
		RunGameOptions: &ebiten.RunGameOptions{
			ApplePressAndHoldEnabled: true,
		},
	}
	if err := guigui.Run(&Root{}, op); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
