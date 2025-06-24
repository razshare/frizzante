package main

import "embed"

//go:embed app/dist
//go:embed makefile
var emb embed.FS
