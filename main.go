package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Ishan27g/sentri/hook"
	"github.com/Ishan27g/sentri/internal"
	"github.com/sirupsen/logrus"
)

var name = "testApp"

func init() {
	tf := new(logrus.TextFormatter)
	tf.DisableTimestamp = true
	logrus.SetFormatter(tf)
}

func main() {

	logrus.SetOutput(hook.Hook(name, os.Stdout))

	for i := 0; i < 5; i++ {
		logrus.Println(fmt.Sprintf("my app log -%d", i))
	}
	<-time.After(1 * time.Second)
	fmt.Println(internal.GetPool().Count())
}
