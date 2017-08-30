package stor

import (
	"github.com/zero-os/0-stor/client/meta"
)

// SingleShardWriter is stor.client which implements block.Writer
// and only work against single client
type SingleShardWriter struct {
	cli   Client
	Shard string
}

func NewSingleShardWriter(conf Config, org, namespace, iyoToken string) (*SingleShardWriter, error) {
	cli, err := NewClientWithToken(&conf, org, namespace, iyoToken)
	if err != nil {
		return nil, err
	}
	return &SingleShardWriter{
		cli:   cli,
		Shard: conf.Shard,
	}, nil
}

func (ssw *SingleShardWriter) WriteBlock(key, val []byte, md *meta.Meta) (*meta.Meta, error) {
	_, err := ssw.cli.ObjectCreate(key, val, nil)
	if err != nil {
		return md, err
	}
	md.SetKey(key)
	if err := md.SetShardSlice([]string{ssw.Shard}); err != nil {
		return md, err
	}
	return md, nil
}