package main

import (
	"docs-fuse/lib/docsConn"
	"docs-fuse/lib/docsFuseLibCommon"

	"strings"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"

	log "github.com/sirupsen/logrus"
)

type fSDirTags struct {
	conn    *docsConn.DocsConn
	tagId   string
	nodeDir docsFuseLibCommon.FSNodeDirGeneric
}

func NewFSDirTags(conn *docsConn.DocsConn, tagId string, inode uint64) fs.Node {
	return fSDirTags{
		conn:    conn,
		tagId:   tagId,
		nodeDir: docsFuseLibCommon.NewFSNodeDirGeneric(conn, tagId, inode),
	}
}

func (this fSDirTags) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("fSDirTags/Attr")
	return this.nodeDir.Attr(ctx, a)
}

func (this fSDirTags) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	list, err := this.conn.GetTags()
	if err != nil {
		log.Printf("fSDirTags/ReadDirAll error: %+v", err)
		return nil, err
	}
	log.Printf("fSDirTags/ReadDirAll")

	this.nodeDir.ResetItem()
	this.nodeDir.AddItem(NewFSDirTagsDocsDir, this.tagId, "docs", fuse.DT_Dir)
	for _, v := range list.Tags {
		if strings.Compare(v.Parent, this.tagId) == 0 {
			this.nodeDir.AddItem(NewFSDirTags, v.Id, v.Name, fuse.DT_Dir)
		}
	}
	return this.nodeDir.ReadDirAll(ctx)
}

func (this fSDirTags) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Printf("fSDirTags/Lookup (found): %+v", name)
	return this.nodeDir.Lookup(ctx, name)
}
