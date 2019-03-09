package docsFuseLibCommon

import (
	"docs-fuse/lib/docsConn"

	"fmt"
	"os"
	"strings"
	"sync"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"

	log "github.com/sirupsen/logrus"
)



type FSNodeDirGeneric struct {
	conn  *docsConn.DocsConn
	docId string
	inode uint64

	ctx *fSNodeDirGenericCtx
}

type child struct {
	Generator func(conn *docsConn.DocsConn, docId string, parentInode uint64) fs.Node
	docId     string
	Name      string
	Type      fuse.DirentType
}

type fSNodeDirGenericCtx struct {
	name2id  map[string]string
	id2child map[string]child
	mux      sync.Mutex
}

func NewFSNodeDirGeneric(conn *docsConn.DocsConn, docId string, inode uint64) FSNodeDirGeneric {
	return FSNodeDirGeneric{
		conn:  conn,
		docId: docId,
		inode: inode,
		ctx: &fSNodeDirGenericCtx{
			name2id:  make(map[string]string),
			id2child: make(map[string]child),
		},
	}
}

func (this FSNodeDirGeneric) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("FSNodeDirGeneric/Attr: %+v", this.docId)
	a.Inode = this.inode
	a.Mode = os.ModeDir | 0555
	return nil
}

func (this FSNodeDirGeneric) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	log.Printf("FSNodeDirGeneric/ReadDirAll: %s", this.docId)

	this.ctx.mux.Lock()
	defer this.ctx.mux.Unlock()

	dirs := make([]fuse.Dirent, 0, len(this.ctx.id2child))
	for _, v := range this.ctx.id2child {
		dirs = append(dirs, fuse.Dirent{
			Inode: fs.GenerateDynamicInode(this.inode, v.docId),
			Name:  v.Name,
			Type:  v.Type,
		},
		)
	}
	return dirs, nil
}

func (this FSNodeDirGeneric) AddItem(Generator func(conn *docsConn.DocsConn, docId string, parentInode uint64) fs.Node, docId string, OriginalName string, Type fuse.DirentType) {
	this.ctx.mux.Lock()
	defer this.ctx.mux.Unlock()

	name := OriginalName

	if _, ok := this.ctx.name2id[name]; ok {
		for i := 2; ok; i++ {
			if index := strings.LastIndex(OriginalName, "."); index == -1 {
				name = fmt.Sprintf("%s(%v)", OriginalName, i)
			} else {
				name = fmt.Sprintf("%s(%v)%s", OriginalName[:index], i, OriginalName[index:])
			}
			_, ok = this.ctx.name2id[name]
		}
	}

	this.ctx.name2id[name] = docId
	this.ctx.id2child[docId] = child{
		Generator: Generator,
		docId:     docId,
		Name:      name,
		Type:      Type,
	}

}

func (this FSNodeDirGeneric) ResetItem() {
	this.ctx.name2id = make(map[string]string)
	this.ctx.id2child = make(map[string]child)
}

func (this FSNodeDirGeneric) Lookup(ctx context.Context, name string) (fs.Node, error) {
	this.ctx.mux.Lock()
	defer this.ctx.mux.Unlock()
	if val, ok := this.ctx.name2id[name]; ok {
		log.Printf("FSNodeDirGeneric/Lookup (found): %+v", name)
		return this.ctx.id2child[val].Generator(this.conn, val, this.inode), nil
	} else {
		log.Printf("FSNodeDirGeneric/Lookup (not found) filename: %+v", name)
		return nil, fuse.ENOENT
	}
}
