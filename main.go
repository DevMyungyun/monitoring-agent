package main
import (
	"fmt"
	"github.com/robfig/cron/v3"
)
func main() {
	fmt.Println("hello world")
	c := cron.New()
	c.AddFunc("1 * * * *", func() { fmt.Println("Every hour on the half hour") })
	c.Start()
}