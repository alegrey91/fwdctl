BINARY_NAME=fwdctl
BINARY_DIR=./bin

build: create-bin-dir
	go mod download
	go build \
		-v \
		-o ${BINARY_DIR}/${BINARY_NAME} \
		.

build-cover: create-bin-dir
	go mod download
	go build \
		-v \
		-cover \
		-o ${BINARY_DIR}/${BINARY_NAME} \
		.

build-gh: create-bin-dir
ifndef GITHUB_REF_NAME
	$(error GITHUB_REF_NAME is undefined)
endif
	go mod download
	go build \
		-v \
      		-ldflags="-s -w -X 'github.com/alegrey91/fwdctl/internal/constants.Version=${GITHUB_REF_NAME }'" \
		-o ${BINARY_DIR}/${BINARY_NAME} \
		.


install:
	cp ${BINARY_DIR}/${BINARY_NAME} /usr/local/bin/

create-bin-dir:
	mkdir -p ${BINARY_DIR}

clean:
	rm -rf ${BINARY_DIR}
