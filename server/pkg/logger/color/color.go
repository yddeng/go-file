package color

import (
	"fmt"
	"strings"
)

type Color int

const (
	BlackFont Color = 30 + iota
	RedFont
	GreenFont
	YellowFont
	BlueFont
	PurpleFont
	AzureFont
	WhiteFont
	Padding
)
const (
	BlackBg Color = 40 + iota
	RedBg
	GreenBg
	YellowBg
	BlueBg
	PurpleBg
	AzureBg
	WhiteBg
)

type format struct {
	format string
}

func (c Color) Format(str interface{}, with ...Color) *format {
	return &format{format: c.Echo(str, with...)}
}
func (f format) Exec(vals ...interface{}) string {
	return fmt.Sprintf(f.format, vals...)
}
func (c Color) Echo(str interface{}, with ...Color) string {
	l := len(with)
	if l == 0 {
		return fmt.Sprintf("\033[%vm%v\033[0m", c, str)
	}
	with = append(with, c)
	mode := make([]string, l)
	p := ""
	for i := 0; i < l; i++ {
		if with[i] == Padding {
			p = " "
		} else {
			mode[i] = fmt.Sprintf("%v", with[i])
		}
	}
	return fmt.Sprintf("\033[%vm%v%v%v\033[0m", strings.Join(mode, ";"), p, str, p)
}
