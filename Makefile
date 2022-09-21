TEMPLATEFILES=$(wildcard slides/template/*)
.PHONY: dev-server
dev-server:
	@cd slides && bs serve

slides/dist/presentation.html: $(TEMPLATEFILES) slides/presentation.md
	@cd slides && bs export

.PHONY: slides
slides: slides/dist/presentation.html
