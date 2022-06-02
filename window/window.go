package window

import (
	"gobang/component"
	"gobang/gobang"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Window struct {
	fyne.Window
	isGamePending bool
	gobangImpl    component.Gobang
	btnStart      *widget.Button
	btnStop       *widget.Button
	labWinner     *canvas.Text
	pannelCnter   fyne.CanvasObject
	btnCnter      []fyne.CanvasObject
	gridCnter     *fyne.Container
	tabLevel      *container.AppTabs
	tabDisply     *container.AppTabs
	imageBoard    *fyne.StaticResource
	imageWhite    *fyne.StaticResource
	imageBlack    *fyne.StaticResource
}

func DialogWindow(title string, board *fyne.StaticResource,
	white *fyne.StaticResource, black *fyne.StaticResource,
	app fyne.App) *Window {
	wnd := &Window{
		Window:     app.NewWindow(title),
		imageBoard: board,
		imageWhite: white,
		imageBlack: black,
	}
	wnd.SetIcon(theme.FyneLogo())

	wnd.CenterOnScreen()
	wnd.Resize(fyne.NewSize(580, 640))
	wnd.SetFixedSize(true)

	wnd.initDialog()
	return wnd
}

func (wnd *Window) initDialog() {
	wnd.pannelCnter = container.NewCenter(
		wnd.initPannel(),
		wnd.initWinner(),
	)
	c := container.NewBorder(
		wnd.initTab(),
		wnd.initBottom(),
		nil,
		nil,
		wnd.pannelCnter,
	)
	wnd.SetContent(c)
}

func (wnd *Window) initTab() fyne.CanvasObject {
	wnd.tabLevel = container.NewAppTabs(
		container.NewTabItem("入门", container.NewMax()),
		container.NewTabItem("简单", container.NewMax()),
		container.NewTabItem("普通", container.NewMax()),
		container.NewTabItem("稍难", container.NewMax()),
		container.NewTabItem("困难", container.NewMax()),
		container.NewTabItem("极难", container.NewMax()),
	)
	wnd.tabLevel.SelectIndex(2)

	wnd.tabDisply = container.NewAppTabs(
		container.NewTabItemWithIcon(
			"对弈中",
			theme.ViewRefreshIcon(),
			container.NewMax(),
		),
	)
	wnd.tabDisply.Hide()

	label := canvas.NewText(
		"by:黎宏 email:theepic7@qq.com",
		color.RGBA{0, 0, 205, 205},
	)
	label.TextSize = 10
	return container.NewBorder(
		nil,
		nil,
		nil,
		label,
		container.NewVBox(wnd.tabLevel, wnd.tabDisply),
	)
}

func (wnd *Window) initPannel() fyne.CanvasObject {
	bg := canvas.NewImageFromResource(wnd.imageBoard)

	objs := make([]fyne.CanvasObject, 0, 225)
	for i := 0; i < 225; i++ {
		btn := newCustomButton("")
		btn.setTag(i)

		btn.OnTapped = func() {
			if wnd.isGamePending {
				return
			}

			if wnd.gobangImpl.GetGameWinner() != component.GobagWinnerNil {
				return
			}

			wnd.isGamePending = true

			tag := btn.getTag()
			cnter := wnd.btnCnter[tag].(*fyne.Container)
			cnter.Objects[0] = canvas.NewImageFromResource(wnd.imageWhite)

			x, y := wnd.gobangImpl.Move(int8(tag%15), int8(tag/15))
			if x != -1 && y != -1 {
				wnd.btnCnter[int(y)*15+int(x)].(*fyne.Container).Objects[0] =
					canvas.NewImageFromResource(wnd.imageBlack)
			}

			winner := wnd.gobangImpl.GetGameWinner()
			switch winner {
			case component.GobagWinnerComputer:
				wnd.labWinner.Text = "  ^_^电脑战胜了你！"
				wnd.labWinner.Show()
			case component.GobagWinnerPlayer:
				wnd.labWinner.Text = "  ^_^你战胜了电脑！"
				wnd.labWinner.Show()
			default:
			}

			wnd.gridCnter.Refresh()
			wnd.isGamePending = false
		}

		cbtn := container.NewGridWrap(fyne.NewSize(35.5, 35.5), btn)
		objs = append(objs, cbtn)
	}

	wnd.btnCnter = objs
	wnd.gridCnter = container.NewGridWithColumns(15, objs...)
	wnd.gridCnter.Hide()

	return container.NewCenter(
		container.NewGridWrap(fyne.NewSize(536, 536), bg),
		container.NewGridWrap(fyne.NewSize(530, 530), wnd.gridCnter),
	)
}

func (wnd *Window) initBottom() fyne.CanvasObject {
	wnd.btnStart = widget.NewButton("开始", func() {
		wnd.btnStart.Disable()
		wnd.btnStop.Enable()
		wnd.tabLevel.Hide()
		wnd.tabDisply.Show()

		cpannel := wnd.pannelCnter.(*fyne.Container)
		cpannel.Objects[0] = wnd.initPannel()
		cpannel.Objects[1] = wnd.initWinner()
		wnd.gridCnter.Show()

		wnd.gobangImpl = gobang.NewBoard(
			wnd.tabLevel.SelectedIndex()+2,
			component.GobangWhite,
		)
	})

	wnd.btnStop = widget.NewButton("结束", func() {
		wnd.btnStart.Enable()
		wnd.btnStop.Disable()
		wnd.tabLevel.Show()
		wnd.tabDisply.Hide()
		wnd.gridCnter.Hide()
		wnd.labWinner.Hide()

		wnd.gobangImpl = nil
	})

	wnd.btnStop.Disable()
	return container.NewCenter(container.NewHBox(wnd.btnStart, wnd.btnStop))
}

func (wnd *Window) initWinner() fyne.CanvasObject {
	wnd.labWinner = canvas.NewText("", color.RGBA{255, 0, 0, 255})
	wnd.labWinner.TextSize = 40
	wnd.labWinner.Alignment = fyne.TextAlignCenter
	wnd.labWinner.Hide()
	return wnd.labWinner
}
