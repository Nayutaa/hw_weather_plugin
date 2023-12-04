package weather

import (
	_ "embed"
	"errors"
	"fmt"
	"hw_weather_plugin/Draw"
	"hw_weather_plugin/api"
	stringsPkg "hw_weather_plugin/utils/strings"
	"strconv"
	"strings"
	"time"
)

var (
	//go:embed 湿度.png
	humidityPNG []byte
)

func DerawImage(cityID, host, weatherKey string) ([]byte, error) {
	//获取一言
	oneSentence, err := api.GetOneSentence()
	if err != nil {
		return nil, err
	}
	// 处理一言
	sentences := strings.Split(oneSentence.Hitokoto, "，")
	if len(sentences) != 2 {
		return nil, errors.New("一言接口获取失败,数据不符合要求：\n" + oneSentence.Hitokoto)
	}
	// 获取天气
	var weatherInfo api.WeatherResp
	if weatherKey == "" {
		weatherInfo, err = api.GetWeather(cityID)
	} else {
		r1, _, err := api.GetCurrentWeather(cityID, host, weatherKey)
		if err != nil {
			return nil, err
		}
		r2, _, err := api.GetWeatherIndex(cityID, host, weatherKey)
		if err != nil {
			return nil, err
		}
		weatherInfo.Parse(&r1, &r2)
	}
	if err != nil {
		return nil, err
	}
	draw, err := Draw.NewCanvas(128, 296, Draw.GetRGBA(255, 255, 255, 255))
	if err != nil {
		return nil, err
	}
	// 一言
	draw.DrawTextVertical(sentences[0], 16, Draw.GetRGBA(0, 0, 0, 255), 108, 0)
	sentences[1] = strings.ReplaceAll(sentences[1], "。", "")
	sentences[1] = strings.ReplaceAll(sentences[1], "？", "")
	draw.DrawTextVertical(sentences[1], 16, Draw.GetRGBA(0, 0, 0, 255), 90, 29)
	draw.DrawBox(88, 23, 1, 130, Draw.GetRGBA(0, 0, 0, 255))
	// 天气

	// 天气情况
	// 矩形背景宽度无字12,每一个字符加16,x原始38
	tmpInt := stringsPkg.GetStrLen(weatherInfo.WeatherStatus.Text)
	draw.DrawRoundedBox(38-(float64(tmpInt)*16/2), 5, 12+(float64(tmpInt)*16), 19, 3, Draw.GetRGBA(0, 0, 0, 255))
	draw.DrawText(weatherInfo.WeatherStatus.Text, 16, Draw.GetRGBA(255, 255, 255, 255), 44-(tmpInt*16/2), 5)

	// 天气图标
	draw.DrawWeatherIcon(weatherInfo.WeatherStatus.Icon, 48, Draw.GetRGBA(0, 0, 0, 255), 19, 32)
	// 温度
	draw.DrawText(weatherInfo.WeatherStatus.Temp+"°C", 25, Draw.GetRGBA(0, 0, 0, 255), 19, 82)

	// 空气质量
	draw.DrawText("空气质量", 12.5, Draw.GetRGBA(0, 0, 0, 255), 5, 116)
	// 矩形背景宽度无字6,每一个字符加12,x原始68
	tmpInt = stringsPkg.GetStrLen(weatherInfo.WeatherIndexs.Air.Category)
	draw.DrawRoundedBox(68-(float64(tmpInt)*12/2), 116, 6+(float64(tmpInt)*12), 15, 3, Draw.GetRGBA(0, 0, 0, 255))
	draw.DrawText(weatherInfo.WeatherIndexs.Air.Category, 12, Draw.GetRGBA(255, 255, 255, 255), 70-(tmpInt*12/2), 116)

	// 湿度
	draw.DrawImageData(humidityPNG, 7, 132)
	draw.DrawText(weatherInfo.WeatherStatus.Humidity+"%", 12, Draw.GetRGBA(0, 0, 0, 255), 25, 132)
	// 湿度进度条
	draw.DrawRoundedBox(5, 147, 80, 8, 3, Draw.GetRGBA(0, 0, 0, 255))
	// 内填充
	draw.DrawRoundedBox(6, 148, 78, 6, 3, Draw.GetRGBA(255, 255, 255, 255))
	// 进度
	tmpInt, _ = strconv.Atoi(weatherInfo.WeatherStatus.Humidity)
	draw.DrawRoundedBox(7, 149, 76*(float64(tmpInt)/100), 4, 3, Draw.GetRGBA(0, 0, 0, 255))
	draw.DrawBox(3, 158, 121, 1, Draw.GetRGBA(0, 0, 0, 255))

	// 日期
	timeNow := time.Now()
	dayStr := fmt.Sprintf("%d月%d日", timeNow.Month(), timeNow.Day())
	tmpInt = stringsPkg.GetStrLen(dayStr)
	// 图像宽度128,每一个字符加16.5,x原始64
	draw.DrawText(dayStr, 16, Draw.GetRGBA(0, 0, 0, 255), 64-(4*16/2), 160)
	// 星期
	// 总外框
	draw.DrawRoundedBox(5, 180, 116, 19, 3, Draw.GetRGBA(0, 0, 0, 255))
	// 内填充
	draw.DrawRoundedBox(6, 181, 114, 17, 3, Draw.GetRGBA(255, 255, 255, 255))
	weekday := int(timeNow.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekday--
	for i, s := range []string{"一", "二", "三", "四", "五", "六", "日"} {
		color := Draw.GetRGBA(0, 0, 0, 255)
		if weekday == i {
			color = Draw.GetRGBA(255, 255, 255, 255)
			w := 16.0
			// 画背景
			if weekday == 0 {
				// 周一
				draw.DrawRoundedBox(5, 180, 6, 19, 3, Draw.GetRGBA(0, 0, 0, 255))
			}
			if weekday == 6 {
				// 周日
				w = 14
				draw.DrawRoundedBox(115, 180, 6, 19, 3, Draw.GetRGBA(0, 0, 0, 255))
			}
			draw.DrawBox(float64(8+i*16), 180, w, 19, Draw.GetRGBA(0, 0, 0, 255))
		}
		draw.DrawText(s, 12, color, 10+(i*8+(8*i)), 182)
	}
	// 今日待办
	draw.DrawText("今日待办", 12.5, Draw.GetRGBA(0, 0, 0, 255), 38, 202)
	draw.DrawBox(3, 217, 121, 1, Draw.GetRGBA(0, 0, 0, 255))
	// 预留代办内容
	draw.DrawText("无", 12.5, Draw.GetRGBA(0, 0, 0, 255), 57, 247)
	return draw.SaveBytes()
}
