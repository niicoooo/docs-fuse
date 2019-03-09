package docsFuseLibCommon

import (
	"docs-fuse/lib/docsConn"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"golang.org/x/net/context"

	log "github.com/sirupsen/logrus"
)

type fSNodeFileDocFile struct {
	conn   *docsConn.DocsConn
	fileId string
	inode  uint64
}

func NewFSNodeFileDocFile(conn *docsConn.DocsConn, fileId string, inode uint64) fs.Node {
	return fSNodeFileDocFile{
		conn:   conn,
		fileId: fileId,
		inode:  inode,
	}
}

func (this fSNodeFileDocFile) Attr(ctx context.Context, a *fuse.Attr) error {
	log.Printf("fSNodeFileDocFile/Attr: %+v", this.fileId)
	data, err := this.conn.GetFileData(this.fileId) //TODO: optimize
	if err != nil {
		log.Printf("fSNodeFileDocFile/Attr: %+v", err)
		return err
	}
	a.Inode = this.inode
	a.Mode = 0444
	a.Size = uint64(len(data))
	return nil
}

func (this fSNodeFileDocFile) ReadAll(ctx context.Context) ([]byte, error) {
	log.Printf("fSNodeFileDocFile/ReadAll: %+v", this.fileId)
	data, err := this.conn.GetFileData(this.fileId)
	if err != nil {
		log.Printf("fSNodeFileDocFile/ReadAll: %+v", err)
		return nil, err
	}
	return data, nil
}
