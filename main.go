package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
	"unsafe"

	"github.com/creack/pty"
	"github.com/gliderlabs/ssh"
)

func setWinsize(f *os.File, w, h int) {
	syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(syscall.TIOCSWINSZ),
		uintptr(unsafe.Pointer(&struct{ h, w, x, y uint16 }{uint16(h), uint16(w), 0, 0})))
}

func main() {
	ssh.Handle(func(s ssh.Session) {
		f, err := os.OpenFile("./ssh.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Println(err)
		}
		defer f.Close()
		// Append the remote address, and the command line, to the log file.
		io.WriteString(f, fmt.Sprintf("%s %s %s\n", s.User(), s.RemoteAddr(), s.RawCommand()))

		cmd := exec.Command("bash", "-c", "HOME=/root LC_ALL=C asciinema rec recs/$(cat /dev/urandom| head -n1 | md5sum -z | cut -c -10) -c \"docker run --rm -itu nobody --cpus 0.05 --memory 25Mb --network none minimal\"")
		ptyReq, winCh, isPty := s.Pty()
		if isPty {
			cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", ptyReq.Term))
			f, err := pty.Start(cmd)
			if err != nil {
				panic(err)
			}
			go func() {
				for win := range winCh {
					setWinsize(f, win.Width, win.Height)
				}
			}()
			go func() {
				io.Copy(f, s) // stdin
			}()
			io.Copy(s, f) // stdout
			cmd.Wait()
		} else {
            out := ""
            switch s.RawCommand() {
            case "uname -s -m":
                out = "Linux x86_64\n"
            case `echo -e \x6F\x6B`:
                out = "ok\n"
            case "cat /etc/os-release":
                out = `PRETTY_NAME="Debian GNU/Linux 11 (bullseye)"
NAME="Debian GNU/Linux"
VERSION_ID="11"
VERSION="11 (bullseye)"
VERSION_CODENAME=bullseye
ID=debian
HOME_URL="https://www.debian.org/"
SUPPORT_URL="https://www.debian.org/support"
BUG_REPORT_URL="https://bugs.debian.org/"
`
            }
			io.WriteString(s, out)
			s.Exit(1)
		}
	})

	log.Println("starting ssh server on port 2222...")
	log.Fatal(ssh.ListenAndServe(":22", nil))
}
