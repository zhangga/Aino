package transformer_test

import (
	"context"
	"testing"

	"github.com/cloudwego/eino-ext/components/document/transformer/splitter/markdown"
	"github.com/cloudwego/eino/schema"
)

func TestMarkdownTransformer(t *testing.T) {
	ctx := context.Background()
	transformer, err := markdown.NewHeaderSplitter(ctx, &markdown.HeaderConfig{
		Headers: map[string]string{
			"##": "",
		},
	})
	if err != nil {
		t.Fatalf("failed to create markdown transformer: %v", err)
	}

	markdownDoc := &schema.Document{
		Content: `
# Title
## Section 1
This is the first section.

## Section 2
This is the second section.

## Section 3
This is the third section.
`,
	}

	// 转换文档
	transformedDocs, err := transformer.Transform(ctx, []*schema.Document{markdownDoc})
	if err != nil {
		t.Fatalf("failed to transform document: %v", err)
	}

	// 输出转换后的文档内容
	for i, doc := range transformedDocs {
		t.Logf("Transformed Document %d Content:\n%s", i+1, doc.Content)
	}
}
