package main

import (
  "bufio"
	"flag"
	"fmt"
	"github.com/jlaffaye/goftp"
	"log"
	"os"
)

var (
	host     = flag.String("host", "", "FTP server")
	username = flag.String("username", "admin", "Username FTP server")
	password = flag.String("password", "somepassword", "Password for FTP server")
	help     = flag.Bool("h", true, "show this help")
)

func usage() {
	PrintErr("usage: ftpclient [flags]", "")
	flag.PrintDefaults()
	os.Exit(2)
}

func PrintErr(str string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, str, a)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *help {
		usage()
	}

	c, err := ftp.Connect(*host)
	defer c.Quit()

	if err != nil {
		log.Fatalf("error connecting %v", err)
	}
	err = c.Login(*username, *password)
	if err != nil {
		log.Fatalf("error %v", err)
	}
	cd, _ := c.CurrentDir()

	fmt.Printf("current directory %v\n", cd)
	en, _ := c.List(cd)

	fmt.Println("List files and folders in current directory")

	for _, v := range en {
		fmt.Println(v.Name)
	}

	var f *os.File
	if f, err = os.Open("ftpz.txt"); err != nil {
		log.Fatalf("error %v", err)
	}

	defer f.Close()
	reader := bufio.NewReader(f)

	err = c.ChangeDir(cd + "/received/")
	if err != nil {
		log.Fatalf("error %v", err)
	}
	cwd, _ := c.CurrentDir()
	fmt.Printf("current directory %v\n", cwd)
	flist, _ := c.List(cwd)
	for _, v := range flist {
		fmt.Println(v.Name)
	}
	err = c.Stor(cwd+"/mysample.txt", reader)
	if err != nil {
		log.Fatalf("error %v", err)
	}

}
