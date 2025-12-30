package knowledgeindexing

import (
	"context"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/document"
	"github.com/spf13/cobra"
	"github.com/zhangga/aino/internal/eino_workflow/knowledgeindexing"
	logger "github.com/zhangga/aino/pkg/zlog"
)

var CmdRun = &cobra.Command{
	Use:   "doc",
	Short: "run the knowledgeindexing service",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	logger.InitDefaultLogger(logger.Config{FilePath: "logs/knowledgeindexing.log", Level: "debug"})
	defer logger.Sync()

	logger.Info("starting knowledgeindexing service...")
	ctx := context.Background()

	err := indexMarkdownFiles(ctx, "./configs/eino_docs")
	if err != nil {
		logger.Fatalf("index markdown files failed: %v", err)
	}
	logger.Info("knowledgeindexing service finished.")
}

func indexMarkdownFiles(ctx context.Context, dir string) error {
	runner, err := knowledgeindexing.BuildKnowledgeIndexing(ctx)
	if err != nil {
		return fmt.Errorf("build knowledgeindexing runner failed: %w", err)
	}

	// 遍历 dir 目录下的所有 markdown 文件，并进行索引
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", path, err)
		}
		if d.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			logger.Debugf("[skip] not a markdown file: %s", path)
			return nil
		}

		logger.Infof("[indexing] markdown file: %s", path)
		ids, err := runner.Invoke(ctx, document.Source{URI: path})
		if err != nil {
			return fmt.Errorf("indexing file %s failed: %w", path, err)
		}
		logger.Infof("[indexed] file: %s, document length: %d", path, len(ids))
		return nil
	})
	return err
}
