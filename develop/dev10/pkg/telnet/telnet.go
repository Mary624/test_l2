package telnet

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

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

	err := connect(sec, host, port)
	if err != nil {
		log.Fatal(err)
	}
}

func connect(sec time.Duration, host, port string) error {
	addr := host
	if host != "" {
		addr += fmt.Sprintf(":%s", port)
	}

	// завершаем, если нажали CTRL+D
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGQUIT)

	for {
		conn, err := net.Dial("tcp", addr)
		if err != nil {
			fmt.Println("can't connect")
			fmt.Println(err)
			return nil
		}

		select {
		case <-sigs:
			return nil
		default:
			handleClient(conn, sec)
		}
	}
}

func handleClient(conn net.Conn, sec time.Duration) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(time.Second * time.Duration(sec)))

	buf := make([]byte, 0, 4096)
	tmp := make([]byte, 256)
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		buf = append(buf, tmp[:n]...)

	}

	fmt.Println(string(buf))
}
