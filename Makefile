OUT_DIR="./terraform/.temp"
FUNC_DIR="./functions"
TMPL_DIR="./templates"
SECRET_DIR="./terraform/.secret"

help:
	@echo "Usage: make [FUNC]"
	@echo "Where FUNC is the cloud function name (build) "

build:
	@printf "packaging \"build\" ... "
	@mkdir -p $(OUT_DIR)/build
	@cp -r $(FUNC_DIR)/* $(SECRET_DIR)/function-service-account.json $(TMPL_DIR) $(OUT_DIR)/build
	@echo "done"
