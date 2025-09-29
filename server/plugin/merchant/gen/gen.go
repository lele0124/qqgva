package main

import (
	"path/filepath"

	"github.com/flipped-aurora/gin-vue-admin/server/plugin/merchant/model"
	"gorm.io/gen"
)

func main() {
	g := gen.NewGenerator(gen.Config{OutPath: filepath.Join("..", "..", "..", "merchant", "blender", "model", "dao"), Mode: gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface})
	g.ApplyBasic(new(model.Merchant))
	g.Execute()
}
