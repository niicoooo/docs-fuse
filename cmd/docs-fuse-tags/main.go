package main

import (
	"docs-fuse/lib/docsConn"

	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"

	log "github.com/sirupsen/logrus"
)

/* *** FS *** */

type FS struct {
	conn *docsConn.DocsConn
}

func (this FS) Root() (fs.Node, error) {
	log.Printf("FS/Root")
	return NewFSDirTags(this.conn, "", 1), nil
}

/* *** Main *** */

func main() {
	flags := flag.NewFlagSet(os.Args[0], 0)
	addr := flags.String("a", "http://127.0.0.1:8100", "Server URL")
	login := flags.String("u", "admin", "User")
	password := flags.String("p", "admin", "Password")
	dir := flags.String("d", "mnt", "Mount path")
	help := flags.Bool("h", false, "Print help")
	verbose := flags.Bool("v", false, "Verbose")

	flags.Parse(os.Args[1:])
	if *help {
		flags.Usage()
		return
	}
	fmt.Printf("Usage: %s -h to display help\n", os.Args[0])
	fmt.Printf("Connecting: %s@%s\n", *login, *addr)
	fmt.Printf("Ctrl-C or umount to stop\n")

	if *verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.WarnLevel)
	}

	conn, err := docsConn.NewDocsConn(*addr, *login, *password)
	if err != nil {
		fmt.Printf("Error - unable to connect: %v\n", err)
		os.Exit(1)
	}

	c, err := fuse.Mount(
		*dir,
		fuse.FSName("Docs-fuse"),
		fuse.ReadOnly(),
		fuse.VolumeName(*login+"@"+*addr),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		_ = <-sigs
		fmt.Printf("\nUnmounting\n")
		for ok := true; ok; ok = (fuse.Unmount(*dir) != nil) {
		}
	}()

	err = fs.Serve(c, &FS{
		conn: conn,
	})
	if err != nil {
		log.Fatal(err)
	}

	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}
