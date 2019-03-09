package main

import (
	"docs-fuse/lib/docsConn"
	"docs-fuse/lib/docsFuseLibCommon"

	"fmt"
	"strings"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"

	log "github.com/sirupsen/logrus"
)

type fSDirTagsDocs struct {
	conn    *docsConn.DocsConn
	tagId   string
	nodeDir docsFuseLibCommon.FSNodeDirGeneric
}

func NewFSDirTagsDocsDir(conn *docsConn.DocsConn, tagId string, inode uint64) fs.Node {
	return fSDirTagsDocs{
		conn:    conn,
		tagId:   tagId,
		nodeDir: docsFuseLibCommon.NewFSNodeDirGeneric(conn, tagId, inode),
	}
}

func (this fSDirTagsDocs) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("fSDirTagsDocs/Attr")
	return this.nodeDir.Attr(ctx, a)
}

func (this fSDirTagsDocs) ReadDirAll(ctx context.Context) ([]fuse.Dirent, error) {
	var list2 *docsConn.DocumentList
	var err error
	if strings.Compare(this.tagId, "") == 0 {
		list2, err = this.conn.GetDocumentList()
	} else {
		list, err := this.conn.GetTags()
		if err != nil {
			log.Printf("fSDirTagsDocs/ReadDirAll tag error: %+v", err)
			return nil, err
		}

		tagName := ""
		for _, v := range list.Tags {
			if strings.Compare(v.Id, this.tagId) == 0 {
				tagName = v.Name
			}
		}
		if strings.Compare(tagName, "") == 0 {
			log.Printf("fSDirTagsDocs/ReadDirAll tag not found: %+v", this.tagId)
			return nil, fmt.Errorf("Tag no found!")
		}

		list2, err = this.conn.GetDocumentListByTag(tagName)
	}
	if err != nil {
		log.Printf("fSDirTagsDocs/ReadDirAll error: %+v", err)
		return nil, err
	}
	log.Printf("fSDirTagsDocs/ReadDirAll")

	this.nodeDir.ResetItem()
	for _, v := range list2.Documents {
		this.nodeDir.AddItem(NewFSDocumentDir, v.Id, v.Title, fuse.DT_Dir)
	}
	return this.nodeDir.ReadDirAll(ctx)

}

func (this fSDirTagsDocs) Lookup(ctx context.Context, name string) (fs.Node, error) {
	log.Printf("fSDirTagsDocs/Lookup (found): %+v", name)
	return this.nodeDir.Lookup(ctx, name)
}
