package timex

import (
	"time"
	_ "time/tzdata"
)

func Init() {
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic("时区加载失败: " + err.Error())
	}
	time.Local = loc
}
