package web

import "embed"

// FS holds the static assets served to the browser (index.html and app.js).
// Embedding the whole directory rather than only index.html keeps the script
// tag <script src="/app.js"> resolvable, so the page always loads its logic.
//
//go:embed index.html app.js
var FS embed.FS
