BUILD_DIR ?= build

build-%:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$* ./cmd/$*

run-%: build-% 
	$(BUILD_DIR)/$*
