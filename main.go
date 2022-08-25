package Ferrum

import (
	"fmt"
	"github.com/wissance/stringFormatter"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	osSignal := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(osSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	app := Create()
	res, err := app.Start()
	if !res {
		msg := stringFormatter.Format("An error occurred during starting application, error is: {0}", err.Error())
		fmt.Println(msg)
	}

	go func() {
		sig := <-osSignal
		//logging.InfoLog(stringFormatter.Format("Got signal from OS: {0}", sig))
		fmt.Println(stringFormatter.Format("Got signal from OS: {0}", sig))
		done <- true
	}()
	<-done

	res, err = app.Stop()
	if !res {
		msg := stringFormatter.Format("An error occurred during stopping application, error is: {0}", err.Error())
		fmt.Println(msg)
	}

}
