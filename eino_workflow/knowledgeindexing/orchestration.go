package knowledgeindexing

import (
	"context"

	"github.com/cloudwego/eino/components/document"
	"github.com/cloudwego/eino/compose"
)

func BuildKnowledgeIndexing(ctx context.Context) (r compose.Runnable[document.Source, []string], err error) {
	const (
		Loader1              = "Loader1"
		DocumentTransformer2 = "DocumentTransformer2"
		Indexer1             = "Indexer1"
	)
	g := compose.NewGraph[document.Source, []string]()
	loader1KeyOfLoader, err := newLoader(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddLoaderNode(Loader1, loader1KeyOfLoader)
	documentTransformer2KeyOfDocumentTransformer, err := newDocumentTransformer(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddDocumentTransformerNode(DocumentTransformer2, documentTransformer2KeyOfDocumentTransformer)
	indexer1KeyOfIndexer, err := newIndexer(ctx)
	if err != nil {
		return nil, err
	}
	_ = g.AddIndexerNode(Indexer1, indexer1KeyOfIndexer)
	_ = g.AddEdge(compose.START, Loader1)
	_ = g.AddEdge(Indexer1, compose.END)
	_ = g.AddEdge(Loader1, DocumentTransformer2)
	_ = g.AddEdge(DocumentTransformer2, Indexer1)
	r, err = g.Compile(ctx, compose.WithGraphName("KnowledgeIndexing"), compose.WithNodeTriggerMode(compose.AllPredecessor))
	if err != nil {
		return nil, err
	}
	return r, err
}
