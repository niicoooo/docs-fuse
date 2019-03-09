package main

import (
	"docs-fuse/lib/docsConn"
	"docs-fuse/lib/docsFuseLibCommon"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"

	log "github.com/sirupsen/logrus"
)

type fSNodeDirDoc struct {
	conn    *docsConn.DocsConn
	docId   string
	nodeDir docsFuseLibCommon.FSNodeDirGeneric
}

func NewFSDocumentDir(conn *docsConn.DocsConn, docId string, inode uint64) fs.Node {
	return fSNodeDirDoc{
		conn:    conn,
		docId:   docId,
		nodeDir: docsFuseLibCommon.NewFSNodeDirGeneric(conn, docId, inode),
	}
}

func (this fSNodeDirDoc) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("fSNodeDirDoc/Attr: %+v", this.docId)
	return this.nodeDir.Attr(ctx, a)
}

func (this fSNodeDirDoc) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	list, _, err := this.conn.GetFileList(this.docId)
	if err != nil {
		log.Printf("fSNodeDirDoc/ReadDirAll: %+v", err)
		return nil, err
	}
	log.Printf("fSNodeDirDoc/ReadDirAll: %s", this.docId)

	this.nodeDir.ResetItem()
	this.nodeDir.AddItem(docsFuseLibCommon.NewFSNodeFileDataJson, this.docId+"d", "data.json", fuse.DT_File)
	this.nodeDir.AddItem(docsFuseLibCommon.NewFSNodeFileFilesJson, this.docId+"f", "files.json", fuse.DT_File)
	for _, v := range list.Files {
		this.nodeDir.AddItem(docsFuseLibCommon.NewFSNodeFileDocFile, v.Id, v.Name, fuse.DT_File)
	}

	return this.nodeDir.ReadDirAll(ctx)
}

func (this fSNodeDirDoc) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Printf("fSNodeDirDoc/Lookup: %+v", name)
	return this.nodeDir.Lookup(ctx, name)
}
