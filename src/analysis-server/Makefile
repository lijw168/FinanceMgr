TARGET_DIR := ${FINANCE_BIN_PATH}
ifeq (${TARGET_DIR},)
    TARGET_DIR=.
endif
API_TARGET := $(TARGET_DIR)/analysis_server
CLI_TARGET := $(TARGET_DIR)/analysis_cli
DEPENDCY := $(shell find . -name "*.go")
API_SOURCE:= analysis-server/api/main
CLI_SOURCE:= cli/main.go

#BRANCH=`git rev-parse --abbrev-ref --symbolic-full-name @{u}`
BRANCH=`git branch | sed -n -e 's/^\* \(.*\)/\1/p'`

define get_git_commit_id
    if [ -d $(1) ]; then cd $(1);git rev-parse HEAD;fi
endef

COMMON=`$(call get_git_commit_id, "../common")`
ANALYSIS_SERVER=`$(call get_git_commit_id, "../analysis-server")`

COMMIT_ID=`git rev-parse HEAD`

DATAFMT=`date "+%Y%m%d%H%M%S"`
GOVERSION=`go version | awk '{print $3}' `

LDFLAGS=-ldflags "-X common/tag.WEB_SERVER_VERSION=${ANALYSIS_SERVER} -X common/tag.COMMON_VER=${COMMON} \
                  -X common/tag.FINANCE_BUILD_TIME=$(DATAFMT) -X 'common/tag.GO_VERSION=${GOVERSION}' \
                  -X common/tag.FINANCE_BUILD_VERSION=${COMMIT_ID}--${BRANCH}"

DBG_FLAGS=-gcflags "-N -l"

ALL: $(API_TARGET) $(CLI_TARGET)

$(API_TARGET):$(DEPENDCY) 
	@echo ${API_TARGET}
	@go build ${LDFLAGS} ${DBG_FLAGS} -o $(API_TARGET) $(API_SOURCE)
$(CLI_TARGET):$(DEPENDCY) 
	@echo ${CLI_TARGET}
	@go build ${LDFLAGS} ${DBG_FLAGS} -o $(CLI_TARGET) $(CLI_SOURCE)
clean:
	@rm -f $(API_TARGET) $(CLI_TARGET)
