package knowledgeindexing

import (
	"context"

	"github.com/cloudwego/eino-ext/components/indexer/milvus"
	"github.com/cloudwego/eino/components/indexer"
)

func newMilvusIndexer(ctx context.Context) (idr indexer.Indexer, err error) {
	indexer, err := milvus.NewIndexer(ctx, &milvus.IndexerConfig{})
	if err != nil {
		return nil, err
	}
	return indexer, nil
}
