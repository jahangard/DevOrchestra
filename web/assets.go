package web

import "embed"

// FS contains the standalone offline frontend.
//go:embed index.html app.js style.css
var FS embed.FS
