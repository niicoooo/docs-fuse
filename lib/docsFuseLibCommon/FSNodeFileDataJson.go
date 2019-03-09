package docsFuseLibCommon

import (
	"docs-fuse/lib/docsConn"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"

	log "github.com/sirupsen/logrus"
)

type fSNodeFileDataJson struct {
	conn  *docsConn.DocsConn
	docId string
	inode uint64
}

func NewFSNodeFileDataJson(conn *docsConn.DocsConn, docId string, inode uint64) fs.Node {
	return fSNodeFileDataJson{
		conn:  conn,
		docId: docId[:36],
		inode: inode,
	}
}

func (this fSNodeFileDataJson) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("fSNodeFileDataJson/Attr: %+v", this.docId)
	data, err := this.conn.GetDocument(this.docId) //TOSO: optimier
	if err != nil {
		log.Printf("fSNodeFileDataJson/Attr: %+v", err)
		return err
	}
	a.Inode = this.inode
	a.Mode = 0444
	a.Size = uint64(len(data))
	return nil
}

func (this fSNodeFileDataJson) ReadAll(ctx context.Context) ([]byte, error) {
	log.Printf("fSNodeFileDataJson/ReadAll: %+v", this.docId)
	data, err := this.conn.GetDocument(this.docId)
	if err != nil {
		log.Printf("fSNodeFileDataJson/ReadAll: %+v", err)
		return nil, err
	}
	return data, nil
}
