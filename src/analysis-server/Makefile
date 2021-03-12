TARGET_DIR := ${ZBS_BIN_PATH}
ifeq (${TARGET_DIR},)
    TARGET_DIR=.
endif
API_TARGET := $(TARGET_DIR)/zbs-server
CLI_TARGET := $(TARGET_DIR)/zbs-cli
DEPENDCY := $(shell find . -name "*.go")
API_SOURCE:= api/main.go
CLI_SOURCE:= cli/main.go

#BRANCH=`git rev-parse --abbrev-ref --symbolic-full-name @{u}`
BRANCH=`git branch | sed -n -e 's/^\* \(.*\)/\1/p'`

define get_git_commit_id
    if [ -d $(1) ]; then cd $(1);git rev-parse HEAD;fi
endef

ZBS_COMMON=`$(call get_git_commit_id, "../zbs-common")`
ZBS_CLIENT=`$(call get_git_commit_id, "../zbs-client")`
ZBS_GATEWAY=`$(call get_git_commit_id, "../zbs-gateway")`
ZBS_SCHEDULER=`$(call get_git_commit_id, "../zbs-scheduler")`
ZBS_SERVER=`$(call get_git_commit_id, "../zbs-server")`
ZBS_STORAGE=`$(call get_git_commit_id, "../zbs-storage")`
ZBS_WORKER=`$(call get_git_commit_id, "../zbs-worker")`
ZBS_OPEN_API=`$(call get_git_commit_id, "../zbs-openapi")`

COMMIT_ID=`git rev-parse HEAD`

DATAFMT=`date "+%Y%m%d%H%M%S"`
GOVERSION=`go version | awk '{print $3}' `
LDFLAGS=-ldflags "-X common/tag.ZBS_CLIENT_VER=${ZBS_CLIENT} -X common/tag.ZBS_COMMON_VER=${ZBS_COMMON} \
        -X common/tag.ZBS_GATEWAY_VER=${ZBS_GATEWAY} -X common/tag.ZBS_SCHEDULER_VER=${ZBS_SCHEDULER} \
        -X common/tag.ZBS_SERVER_VER=${ZBS_SERVER} -X common/tag.ZBS_STORAGE_VER=${ZBS_STORAGE} \
        -X common/tag.ZBS_WORKER_VER=${ZBS_WORKER} -X common/tag.ZBS_BUILD_TIME=$(DATAFMT) \
        -X 'common/tag.GO_VERSION=${GOVERSION}' -X common/tag.ZBS_BUILD_VERSION=${COMMIT_ID}--${BRANCH} \
        -X common/tag.ZBS_OPENAPI_VER=${ZBS_OPEN_API} "

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
