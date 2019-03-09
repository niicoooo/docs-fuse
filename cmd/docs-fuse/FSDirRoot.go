package main

import (
	"docs-fuse/lib/docsConn"
	"docs-fuse/lib/docsFuseLibCommon"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"

	log "github.com/sirupsen/logrus"
)

type fSDirRoot struct {
	conn    *docsConn.DocsConn
	nodeDir docsFuseLibCommon.FSNodeDirGeneric
}

func NewFSDirRoot(conn *docsConn.DocsConn) fs.Node {
	return fSDirRoot{
		conn:    conn,
		nodeDir: docsFuseLibCommon.NewFSNodeDirGeneric(conn, "1", 1),
	}
}

func (this fSDirRoot) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("fSDirRoot/Attr")
	return this.nodeDir.Attr(ctx, a)
}

func (this fSDirRoot) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	list, err := this.conn.GetDocumentList()
	if err != nil {
		log.Printf("fSDirRoot/ReadDirAll error: %+v", err)
		return nil, err
	}
	log.Printf("fSDirRoot/ReadDirAll")

	this.nodeDir.ResetItem()
	for _, v := range list.Documents {
		this.nodeDir.AddItem(NewFSDocumentDir, v.Id, v.Title, fuse.DT_Dir)
	}
	return this.nodeDir.ReadDirAll(ctx)
}

func (this fSDirRoot) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Printf("fSDirRoot/Lookup (found): %+v", name)
	return this.nodeDir.Lookup(ctx, name)
}
