package initializer

import (
	"context"
	"sort"

	"github.com/whoisnian/glb/logger"
	"github.com/whoisnian/myRecord/global"
	"github.com/whoisnian/myRecord/schema"
)

func ApplySchema() {
	files, err := schema.FS.ReadDir(".")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Debug("Found ", len(files), " sql files")

	sort.Slice(files, func(i, j int) bool { return files[i].Name() < files[j].Name() })

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		logger.Debug("Read and exec '", file.Name(), "'...")
		data, err := schema.FS.ReadFile(file.Name())
		if err != nil {
			logger.Fatal(err)
		}

		if _, err = global.Pool.Exec(context.Background(), string(data)); err != nil {
			logger.Fatal(err)
		}
	}
}
