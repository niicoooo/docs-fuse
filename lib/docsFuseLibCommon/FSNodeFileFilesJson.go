package docsFuseLibCommon

import (
	"docs-fuse/lib/docsConn"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"

	log "github.com/sirupsen/logrus"
)

type fSNodeFileFilesJson struct {
	conn  *docsConn.DocsConn
	docId string
	inode uint64
}

func NewFSNodeFileFilesJson(conn *docsConn.DocsConn, docId string, inode uint64) fs.Node {
	return fSNodeFileFilesJson{
		conn:  conn,
		docId: docId[:36],
		inode: inode,
	}
}

func (this fSNodeFileFilesJson) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("fSNodeFileFilesJson/Attr: %+v", this.docId)
	_, data, err := this.conn.GetFileList(this.docId) //TODO: optimize
	if err != nil {
		log.Printf("fSNodeFileFilesJson/Attr: %+v", err)
		return err
	}
	a.Inode = this.inode
	a.Mode = 0444
	a.Size = uint64(len(data))
	return nil
}

func (this fSNodeFileFilesJson) ReadAll(ctx context.Context) ([]byte, error) {
	log.Printf("fSNodeFileFilesJson/ReadAll: %+v", this.docId)
	_, data, err := this.conn.GetFileList(this.docId)
	if err != nil {
		log.Printf("fSNodeFileFilesJson/ReadAll: %+v", err)
		return nil, err
	}
	return data, nil
}
