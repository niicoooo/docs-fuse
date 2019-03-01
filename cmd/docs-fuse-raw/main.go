package main

import (
	"../../lib/docsConn"

	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
	"golang.org/x/net/context"

	log "github.com/sirupsen/logrus"
)

/* *** FS *** */

type FS struct {
	conn *docsConn.DocsConn
}

func (this FS) Root() (fs.Node, error) {
	log.Printf("FS/Root")
	return FSRootDir{
		conn: this.conn,
	}, nil
}

/* *** FSRootDir *** */

type FSRootDir struct {
	conn *docsConn.DocsConn
}

func (FSRootDir) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("FSRootDir/Attr")
	a.Inode = 1
	a.Mode = os.ModeDir | 0555
	return nil
}

func (this FSRootDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	list, err := this.conn.GetDocumentList()
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
	name = name[:strings.Index(name, " ")]
	log.Printf("FSRootDir/Lookup: %+v", name)
	return FSDocumentDir{
		conn:  this.conn,
		name:  name,
		inode: fs.GenerateDynamicInode(1, name),
	}, nil
}

/* *** FSDocumentDir *** */

type FSDocumentDir struct {
	conn  *docsConn.DocsConn
	name  string
	inode uint64
}

func (this FSDocumentDir) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("FSDocumentDir/Attr: %+v", this.name)
	a.Inode = this.inode
	a.Mode = os.ModeDir | 0555
	return nil
}

func (this FSDocumentDir) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	list, _, err := this.conn.GetFileList(this.name)
	if err != nil {
		log.Printf("FSDocumentDir/ReadDirAll: %+v", err)
		return nil, err
	}
	log.Printf("FSDocumentDir/ReadDirAll: %s", this.name)
	dirs := make([]fuse.Dirent, 0, len(list.Files))
	for _, v := range list.Files {
		var name string
		name = v.Id + " " + v.Name
		dirs = append(dirs, fuse.Dirent{
			Inode: fs.GenerateDynamicInode(this.inode, v.Id),
			Name:  name,
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
		return FSFileJsonFile{
			conn:  this.conn,
			name:  this.name,
			inode: fs.GenerateDynamicInode(this.inode, name),
		}, nil
	}
	if strings.Compare(name, "data.json") == 0 {
		return FSDocJsonFile{
			conn:  this.conn,
			name:  this.name,
			inode: fs.GenerateDynamicInode(this.inode, name),
		}, nil
	}
	name = name[:strings.Index(name, " ")]
	return FSFileFile{
		conn:  this.conn,
		name:  name,
		inode: fs.GenerateDynamicInode(this.inode, name),
	}, nil
}

/* *** FSDocJsonFile *** */

type FSDocJsonFile struct {
	conn  *docsConn.DocsConn
	name  string
	inode uint64
}

func (this FSDocJsonFile) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("FSDocJsonFile/Attr: %+v", this.name)
	data, err := this.conn.GetDocument(this.name) //TOSO: optimier
	if err != nil {
		log.Printf("FSDocJsonFile/Attr: %+v", err)
		return err
	}
	a.Inode = this.inode
	a.Mode = 0444
	a.Size = uint64(len(data))
	return nil
}

func (this FSDocJsonFile) ReadAll(ctx context.Context) ([]byte, error) {
	log.Printf("FSDocJsonFile/ReadAll: %+v", this.name)
	data, err := this.conn.GetDocument(this.name)
	if err != nil {
		log.Printf("FSDocJsonFile/ReadAll: %+v", err)
		return nil, err
	}
	return data, nil
}

/* *** FSFileJsonFile *** */

type FSFileJsonFile struct {
	conn  *docsConn.DocsConn
	name  string
	inode uint64
}

func (this FSFileJsonFile) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("FSFileJsonFile/Attr: %+v", this.name)
	_, data, err := this.conn.GetFileList(this.name) //TOSO: optimier
	if err != nil {
		log.Printf("FSFileJsonFile/Attr: %+v", err)
		return err
	}
	a.Inode = this.inode
	a.Mode = 0444
	a.Size = uint64(len(data))
	return nil
}

func (this FSFileJsonFile) ReadAll(ctx context.Context) ([]byte, error) {
	log.Printf("FSFileJsonFile/ReadAll: %+v", this.name)
	_, data, err := this.conn.GetFileList(this.name)
	if err != nil {
		log.Printf("FSFileJsonFile/ReadAll: %+v", err)
		return nil, err
	}
	return data, nil
}

/* *** FSFileFile *** */

type FSFileFile struct {
	conn  *docsConn.DocsConn
	name  string
	inode uint64
}

func (this FSFileFile) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("FSFileFile/Attr: %+v", this.name)
	data, err := this.conn.GetFileData(this.name)
	if err != nil {
		log.Printf("FSFileFile/Attr: %+v", err)
		return err
	}
	a.Inode = this.inode
	a.Mode = 0444
	a.Size = uint64(len(data))
	return nil
}

func (this FSFileFile) ReadAll(ctx context.Context) ([]byte, error) {
	log.Printf("FSFileFile/ReadAll: %+v", this.name)
	data, err := this.conn.GetFileData(this.name)
	if err != nil {
		log.Printf("FSFileFile/ReadAll: %+v", err)
		return nil, err
	}
	return data, nil
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
		fmt.Printf("Error: Unable to connect!\n")
		os.Exit(1)
	}

	filesys := &FS{
		conn: conn,
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

	err = fs.Serve(c, filesys)
	if err != nil {
		log.Fatal(err)
	}

	<-c.Ready
	if err := c.MountError; err != nil {
		log.Fatal(err)
	}
}
