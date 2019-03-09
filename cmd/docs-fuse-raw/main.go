package main

import (
	"docs-fuse/lib/docsConn"
	"docs-fuse/lib/docsFuseLibCommon"

	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"

	log "github.com/sirupsen/logrus"
)

/* *** FS *** */

type FS struct {
	ctx *FSCtx
}

func (this FS) Root() (fs.Node, error) {
	log.Printf("FS/Root")
	return FSRootDir{
		ctx: this.ctx,
	}, nil
}

type FSCtx struct {
	conn *docsConn.DocsConn
}

/* *** FSRootDir *** */

type FSRootDir struct {
	ctx *FSCtx
}

func (FSRootDir) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("FSRootDir/Attr")
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	return nil
}

func (this FSRootDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	list, err := this.ctx.conn.GetDocumentList()
	if err != nil {
		log.Printf("FSRootDir/ReadDirAll: %+v", err)
		return nil, err
	}
	log.Printf("FSRootDir/ReadDirAll")
	dirs := make([]fuse.Dirent, 0, len(list.Documents))
	for _, v := range list.Documents {
		dirs = append(dirs, fuse.Dirent{
			Inode: fs.GenerateDynamicInode(1, v.Id),
			Name:  v.Id + " " + v.Title,
			Type:  fuse.DT_Dir,
		},
		)
	}
	return dirs, nil
}

func (this FSRootDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Printf("FSRootDir/Lookup: %+v", name)
	return FSDocumentDir{
		ctx:   this.ctx,
		docId: name[:strings.Index(name, " ")],
		inode: fs.GenerateDynamicInode(1, name[:strings.Index(name, " ")]),
	}, nil
}

/* *** FSDocumentDir *** */

type FSDocumentDir struct {
	ctx   *FSCtx
	docId string
	inode uint64
}

func (this FSDocumentDir) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("FSDocumentDir/Attr: %+v", this.docId)
	a.Inode = this.inode
	a.Mode = os.ModeDir | 0555
	return nil
}

func (this FSDocumentDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	list, _, err := this.ctx.conn.GetFileList(this.docId)

	if err != nil {
		log.Printf("FSDocumentDir/ReadDirAll: %+v", err)
		return nil, err
	}
	log.Printf("FSDocumentDir/ReadDirAll: %s", this.docId)
	dirs := make([]fuse.Dirent, 0, len(list.Files))

	for _, v := range list.Files {
		dirs = append(dirs, fuse.Dirent{
			Inode: fs.GenerateDynamicInode(this.inode, v.Id),
			Name:  v.Id + " " + v.Name,
			Type:  fuse.DT_File,
		},
		)
	}
	files := [2]string{"files.json", "data.json"}
	for _, v := range files {
		dirs = append(dirs, fuse.Dirent{
			Inode: fs.GenerateDynamicInode(this.inode, v),
			Name:  v,
			Type:  fuse.DT_File,
		})
	}
	return dirs, nil
}

func (this FSDocumentDir) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Printf("FSDocumentDir/Lookup: %+v", name)
	if strings.Compare(name, "files.json") == 0 {
		return docsFuseLib.NewFSNodeFileFilesJson(this.ctx.conn, this.docId, fs.GenerateDynamicInode(this.inode, "files.json")), nil
	}
	if strings.Compare(name, "data.json") == 0 {
		return docsFuseLib.NewFSNodeFileDataJson(this.ctx.conn, this.docId, fs.GenerateDynamicInode(this.inode, "data.json")), nil
	}
	return docsFuseLib.NewFSNodeFileDocFile(this.ctx.conn, name[:strings.Index(name, " ")], fs.GenerateDynamicInode(this.inode, name[:strings.Index(name, " ")])), nil
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
		fuse.FSName("Docs"),
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
		ctx: &FSCtx{
			conn: conn,
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}
