package epubserver

const ReaderJs = `
document.onkeydown = function(e) {
	if (e.which == 37) {
		if (componentPrev) {
			window.location = '/' + encodeURI(componentPrev);
		}
	}
	else if (e.which == 39) {
		if (componentNext) {
			window.location = '/' + encodeURI(componentNext);
		}
	}
};
`