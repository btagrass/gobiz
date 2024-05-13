package uat

// import (
// 	"strings"

// 	"github.com/go-vgo/robotgo"
// 	"github.com/samber/lo"
// )

// type Desktop struct {
// }

// func NewDesktop() (*Desktop, error) {
// 	return &Desktop{}, nil
// }

// func (d *Desktop) Click(x, y int, button string, double bool) error {
// 	if button == "" {
// 		button = "left"
// 	}
// 	robotgo.MoveClick(x, y, button, double)
// 	return nil
// }

// func (d *Desktop) DragDrop(x, y int) error {
// 	robotgo.DragSmooth(x, y)
// 	return nil
// }

// func (d *Desktop) Press(key string) error {
// 	keys := strings.Split(key, "::")
// 	args := lo.Map(keys[:len(keys)-1], func(item string, _ int) any {
// 		return item
// 	})
// 	robotgo.KeyPress(keys[len(keys)-1], args...)
// 	return nil
// }

// func (d *Desktop) Scroll(x, y int, direction string) error {
// 	if direction == "" {
// 		direction = "down"
// 	}
// 	robotgo.Move(x, y)
// 	robotgo.ScrollDir(1, direction)
// 	return nil
// }
