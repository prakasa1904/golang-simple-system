package internal

import (
	"embed"
	"os"

	"github.com/devetek/go-core/mdfs"
)

var (
	// Update this embed when create html in different structure

	//go:embed views/pages/admin/**/*.html views/pages/**/*.html views/components/**/admin/*.html views/components/**/*.html
	tmpls embed.FS

	Template = mdfs.New(tmpls, "", os.Getenv("ENV"))
)
