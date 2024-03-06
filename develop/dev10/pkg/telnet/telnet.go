package telnet

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	t "github.com/matjam/telnet"
	flag "github.com/spf13/pflag"
)

func Connect() {
	// получаем параметры
	var sec time.Duration
	flag.DurationVarP(&sec, "timeout", "t", 10, "choose timeout")

	flag.Parse()

	host := flag.Arg(0)
	if host == "" {
		log.Fatal("host is undefined")
	}

	port := ""
	if len(flag.Args()) > 1 {
		port = flag.Arg(1)
		_, err := strconv.Atoi(port)
		if err != nil {
			log.Fatal("wrong port")
		}
	}

	connect(sec, host, port)
}

func connect(sec time.Duration, host, port string) {
	addr := host
	if host != "" {
		addr += fmt.Sprintf(":%s", port)
	}

	// завершаем, если нажали CTRL+D
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT)

	for {
		select {
		case <-sigs:
			return
		default:
		}
		conn, err := t.Dial(addr)
		if err != nil {
			log.Fatal(err)
		}
		conn.SetDeadline(time.Now().Add(sec))

		var res string
		fmt.Scanln(&res)

		if len(res) > 0 {
			conn.Write([]byte(res))
		}

		var buf []byte
		conn.Read(buf)
		if len(buf) > 0 {
			fmt.Println(string(buf))
		}
		conn.Close()
	}
}
