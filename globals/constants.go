package globals

const KB = 1024
const MB = 1024 * KB
const GB = 1024 * MB
const TB = 1024 * GB
const PB = 1024 * TB
const EB = 1024 * PB
const CodegenModHint = "//gen:mod"
const CodegenModsHint = "//gen:mods"
const RenderScriptFormat = `
if(!module){
	var module={exports:{}};
}
(function(){
	%s
	return render
})()
`
