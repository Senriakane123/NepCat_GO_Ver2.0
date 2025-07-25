package GO_Tyne_UI

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func Tyne_UI_Init() {
	myApp := app.New()
	myWindow := myApp.NewWindow("登录界面")

	username := widget.NewEntry()
	username.SetPlaceHolder("请输入用户名")

	password := widget.NewPasswordEntry()
	password.SetPlaceHolder("请输入密码")

	loginButton := widget.NewButton("登录", func() {
		if username.Text == "admin" && password.Text == "123456" {
			dialog.ShowInformation("登录成功", "欢迎你："+username.Text, myWindow)
		} else {
			dialog.ShowError(fmt.Errorf("用户名或密码错误"), myWindow)
		}
	})

	settingsButton := widget.NewButton("设置", func() {
		dialog.ShowInformation("设置", "打开设置界面", myWindow)
		// 可弹出新窗口或配置弹窗
	})

	myWindow.SetContent(container.NewVBox(
		username,
		password,
		loginButton,
		settingsButton,
	))

	myWindow.Resize(fyne.NewSize(300, 200))
	myWindow.ShowAndRun()
}
