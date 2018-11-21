OUT_DIR="./terraform/.temp"
FUNC_DIR="./functions"
TMPL_DIR="./templates"

help:
	@echo "Usage: make [FUNC]"
	@echo "Where FUNC is the cloudfunction name (build) "

build:
	@printf "packaging \"build\" ... "
	@mkdir -p $(OUT_DIR)/build
	@cp -r $(FUNC_DIR)/* $() $(TMPL_DIR) $(OUT_DIR)/build
	@echo "done"
