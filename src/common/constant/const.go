package constant

import (
	"errors"
	"time"
)

// MAX_LEN 10
const (
	ID_PREFIX_HOST              = "host-"
	ID_PREFIX_POOL              = "pool-"
	ID_PREFIX_RACK              = "rack-"
	ID_PREFIX_DISK              = "disk-"
	ID_PREFIX_RG                = "rg-"
	ID_PREFIX_PROXY             = "proxy-"
	ID_PREFIX_HEARTBEAT         = "hb-"
	ID_PREFIX_REPLICA           = "replica-"
	ID_PREFIX_VOLUME            = "vol-"
	ID_PREFIX_SNAPSHOT          = "snapshot-"
	ID_PREFIX_VOLUME_SYSTEM     = "vol-sys-"
	ID_PREFIX_VOLUME_ATTACHMENT = "vol-tach-"
	ID_PREFIX_TASK              = "task-"
	ID_PREFIX_VOLUME_TASK       = "vol-task-"
	ID_PREFIX_SNAPSHOT_TASK     = "ss-task-"
	ID_GATEWAY_REQ              = "gw-"
	ID_PREFIX_ADMINTASK         = "admin-task-"
	ID_PREFIX_RT                = "replication-task-"
	ID_PREFIX_QUOTA             = "quota-"
	ID_PREFIX_CACHE_LIST        = "zbs-proxy-"
)

const (
	HOST_STATUS_DOWN = 0
	HOST_STATUS_UP   = 1
)

const (
	HOST_ADMINSTATUS_DOWN = 0
	HOST_ADMINSTATUS_UP   = 1
	HOST_ADMINSTATUS_GRAY = 2
)

const (
	NAME_MAXL     = 32
	NAME_MINL     = 1
	DESC_MAXL     = 256
	HOST_PORT     = 30000
	PORT_INTERVAL = 2
)

const (
	Order_Desc = 0
	Order_Asc  = 1
)

const (
	KB = 1 << 10
	MB = KB << 10
	GB = MB << 10
	TB = GB << 10
	PB = TB << 10
)

const (
	RAFT_ROLE_INVALID   = 0
	RAFT_ROLE_CANDIDATE = 1
	RAFT_ROLE_LEADER    = 2
	RAFT_ROLE_FOLLOWER  = 3
)
const (
	TIMEOUT = 5
)

//Message used between client&data_node, leader&follower
const (
	ZBS_MAGIC                      uint8 = 0xEB
	DISK_SECTOR_SIZE               int64 = 4096
	DIO_BUF_ALIGNMENT              int64 = 512
	HEADER_LENGTH                  int64 = 160 //sizeof(Request)
	RAFT_TERM_OFFSET                     = HEADER_LENGTH
	RAFT_TERM_LEN                        = 8
	RAFT_LOGINDEX_OFFSET                 = RAFT_TERM_OFFSET + RAFT_TERM_LEN
	RAFT_LOGINDEX_LEN                    = 8
	RAFT_LOGLENGTH_OFFSET                = RAFT_LOGINDEX_OFFSET + RAFT_LOGINDEX_LEN
	RAFT_LOGLENGTH_LEN                   = 8
	PADDING_END_OFFSET                   = RAFT_LOGLENGTH_OFFSET + RAFT_LOGLENGTH_LEN
	PADDING_DATA_LENGTH                  = RAFT_TERM_LEN + RAFT_LOGINDEX_LEN + RAFT_LOGLENGTH_LEN
	REQUEST_MSG_END_OFFSET               = PADDING_END_OFFSET + PADDING_DATA_LENGTH
	HANDLE_BYTES_COUNT                   = 8
	INT32_BYTES_COUNT                    = 4
	INT64_BYTES_COUNT                    = 8
	RECONNECT_FOLLOWER_RETRY_TIMES       = 3
	GROUPID_STRING_LENGTH                = 20
	DISKID_STRING_LENGTH                 = 20
	CHUNKID_STRING_LENGTH                = 28
	MAX_PACKET_DATA_SIZE                 = 1024 * 128
	MAX_PACKET_SIZE                      = DISK_SECTOR_SIZE + MAX_PACKET_DATA_SIZE
	DEFAULT_OBJ_SIZE                     = 1 * GB
	SNAPSHOT_BLOCK_SIZE                  = 4 * MB
)

var (
	ERR_INVALIDARGUMENT       = errors.New("Invalid Argument")
	ERR_INVALIDVERSION        = errors.New("Invalid Version")
	ERR_SYSCOMMAND            = errors.New("Execute Command Error")
	ERR_BADMAGIC              = errors.New("Bad Magic")
	ERR_BADCHECKSUM           = errors.New("Bad CheckSum")
	ERR_CHANFULL              = errors.New("Channel is full")
	ERR_READLOGEND            = errors.New("Read To LogEnd")
	ERR_FOLLOWER_UNCONNECTED  = errors.New("Follower Unconnected")
	ERR_FOLLOWER_UNAVAILABLE  = errors.New("Follower Unconnected")
	ERR_APPEND_REJECT         = errors.New("Append Reject")
	ERR_NOTREADY              = errors.New("Not Ready")
	ERR_NOTALLOW              = errors.New("Not Allow")
	ERR_EXIST                 = errors.New("Already exist")
	ERR_NOTEXIST              = errors.New("Not exist")
	ERR_INVALIDOBJFORMAT      = errors.New("Invalid object format")
	ERR_TCP_CONNECTION_NIL    = errors.New("tcp connection nil")
	ERR_UNSUPPORT             = errors.New("Not Support")
	ERR_NOFREEBLOCK           = errors.New("NoFreeBlockToAlloc")
	ERR_CHECKSUM_FAIL         = errors.New("CheckSumFail")
	ERR_BLOCKNOTALLOC         = errors.New("BlockNotAlloc")
	ERR_METADATADEV_TOO_SMALL = errors.New("MetadataDevTooSmall")
	ERR_WRITE_DATA_TOO_LARGE  = errors.New("WriteDataTooLarge")
	ERR_OUTOFBOUND            = errors.New("Out of Bound")
	ERR_EIO                   = errors.New("IO ERROR")
	ERR_NOSPACELEFT           = errors.New("No Space Left")
	ERR_FORMAT                = errors.New("Format Error")
	ERR_FULL_ZERO             = errors.New("Full Zero")
	ERR_SESSION               = errors.New("Session Error")
	ERR_NOMEM                 = errors.New("No Memory")
	ERR_BUSY                  = errors.New("Error Busy")
	ERR_TIMEOUT               = errors.New("Operation timeout")
)

const (
	_ = iota
	ERRNO_OK
	ERRNO_ERROR                       // unknown
	ERRNO_EIO                         // IO error
	ERRNO_APPEND_REJECT               //
	ERRNO_INVALID_VERSION             // term missmatch
	ERRNO_INVALID_ARG                 // init state
	ERRNO_NOTEXIST                    // object/group not exists
	ERRNO_TIMEOUT                     // client
	ERRNO_TMP_ERR                     //
	ERRNO_FULL_ZERO                   // block do not map
	ERRNO_CONGESTION                  // io congestion
	ERRNO_SESSION                     // session mismatch
	ERRNO_BUSY                        // device busy
	ERRNO_SNAP_NO_CHANGE              // snapshot block not changed
	ERRNO_SNAP_NO_CHANGE_OR_FULL_ZERO // snapshot block not changed
	ERRNO_NOT_ALLOW                   // Illegal operation
)

func ERR_TO_ERRNO(err error) uint8 {
	if err == nil {
		return ERRNO_OK
	}

	switch err {
	case ERR_NOTEXIST:
		return ERRNO_NOTEXIST
	case ERR_FULL_ZERO:
		return ERRNO_FULL_ZERO
	default:
		return ERRNO_EIO
	}
}

const (
	_ = iota //'0' is reserved
	// the highest bit is 1 indicate RSP_, 0 indicate REQ_

	REQ_READ = 0x1 // client Read data node
	RSP_READ = 0x81

	REQ_WRITE = 0x2 // client Write data node
	RSP_WRITE = 0x82

	REQ_BACKUP_WRITE = 0x3 // raft leader append write Follower
	RSP_BACKUP_WRITE = 0x83

	// raft follower in migration reads data from leader
	REQ_MG_READ = 0x4
	RSP_MG_READ = 0x84

	REQ_APPEND_ENTRY_RPC = 0x5 // raft leader append RPC header
	RSP_APPEND_ENTRY     = 0x85

	REQ_LEADER_CONFIRM = 0x6 // raft leader confirm 'leader' role
	RSP_LEADER_CONFIRM = 0x86

	// raft leader sends migration cmd to follower
	REQ_MIGRATION = 0x7
	RSP_MIGRATION = 0x87

	REQ_DELETE_OBJECT = 0x8 // client deletes object
	RSP_DELETE_OBJECT = 0x88

	// raft leader sends delete object to follower
	REQ_BACKUP_DELETE_OBJECT = 0x9
	RSP_BACKUP_DELETE_OBJECT = 0x89

	// client login storage node
	REQ_LOGIN = 0xA
	RSP_LOGIN = 0xA1

	// proxy client reads bitmap of object from zbs-proxy
	REQ_BMP_READ = 0xB
	RSP_BMP_READ = 0x8B

	// Get object's snap version list
	REQ_READ_SNAP_VERSION = 0xC
	RSP_READ_SNAP_VERSION = 0x8C

	// delete snapshot object
	REQ_DELETE_SNAP_OBJECT = 0xD
	RSP_DELETE_SNAP_OBJECT = 0x8D

	// ping-pong request to detect network healthy
	REQ_PING = 0xE
	RSP_PING = 0x8E

	// check data whether exists in disk
	REQ_CHECK_DATA_EXIST = 0xF
	RSP_CHECK_DATA_EXIST = 0x8F

	//Notice : opcode is uint8 request.Request
)

//volume actions: action="" 表示该字段无效
const (
	AttachingAction        = "attaching"
	DetachingAction        = "detaching"
	CreatingSnapshotAction = "creating-snapshot"
	GarbageDeletingAction  = "garbage-deleting"
	RestoringVolumeAction  = "restoring-volume"
)

//volume status
const (
	InvalidVolumeStatus  = -1
	CreatingVolume       = 0  //The volume is being created.
	AvailableVolume      = 1  //The volume is ready to attach to an instance.
	InUseVolume          = 2  //The volume is attached to an instance.
	DeleteingVolume      = 3  //The volume is being deleted.
	DeletedVolume        = 4  //The volume has been deleted.
	ExtendingVolume      = 5  //The volume is being extended
	RestoringVolume      = 6  //A backup is being restored to the volume.
	ErrorCreatingVolume  = 7  //A volume creation error occurred.
	ErrorExtendingVolume = 8  //An error occurred while attempting to extend a volume.
	ErrorDeletingVolume  = 9  //A volume deletion error occurred.
	ErrorRestoringVolume = 10 //A backup restoration error occurred.
	RecycledVolume       = 11 // volume recycled
	ErrorRecycledVolume  = 12 // volume recycled error
)

//snapshot status
const (
	CreatingSnapshot      = 20 //The snapshot is being created.
	AvailableSnapshot     = 21 //The snapshot is ready to restore a volume
	InUseSnapshot         = 22 //The snapshot is working(restore/createVolume)
	DeleteingSnapshot     = 23 //The snapshot is being deleted.
	DeletedSnapshot       = 24 //The snapshot is being deleted.
	ErrorCreatingSnapshot = 25 //A snapshot created error occurred.
	ErrorDeletingSnapshot = 26 //A snapshot deletion error occurred.
	ErrorRecycledSnapshot = 27 // a snapshot recycled error
	RecyclingSnapshot     = 28 // a snapshot is recycling
	RecycledSnapshot      = 29 // a snapshot recycled
	ErrorCopyingSnapshot  = 40 // Copy snapshot error
	ErrorUploadSnapshot   = 41
	PreCreatingSnapshot   = 42 // PreCreatingSnapshot -> CreatingSnapshot
)

//volume-attachment status
const (
	VolumeAttaching   = 30 //The volume is attaching to an instance.
	VolumeDetaching   = 31 //The volume is detaching from an instance.
	VolumeAttached    = 32 //The volume has been attached an instance
	VolumeDetached    = 33 //The volume has been detached an instance
	ErrorVolumeAttach = 34 //attached error occurred
	ErrorVolumeDetach = 35 //detached error occurred
)

var StatusMessage = map[int]string{
	InvalidVolumeStatus:   "invalid",
	CreatingVolume:        "creating",
	AvailableVolume:       "available",
	InUseVolume:           "in-use",
	DeleteingVolume:       "deleting",
	ExtendingVolume:       "extending",
	RestoringVolume:       "restoring",
	ErrorCreatingVolume:   "error_create",
	ErrorDeletingVolume:   "error_delete",
	ErrorRestoringVolume:  "error_restore",
	ErrorExtendingVolume:  "error_extend",
	DeletedVolume:         "deleted",
	PreCreatingSnapshot:   "creating",
	CreatingSnapshot:      "creating",
	AvailableSnapshot:     "available",
	InUseSnapshot:         "in-use",
	DeleteingSnapshot:     "deleting",
	DeletedSnapshot:       "deleted",
	ErrorCreatingSnapshot: "error_create",
	ErrorDeletingSnapshot: "error_delete",
	ErrorCopyingSnapshot:  "error_create",
	ErrorUploadSnapshot:   "error_create",
	VolumeAttached:        "attached",
	VolumeDetached:        "detached",
	ErrorVolumeAttach:     "error_attach",
	ErrorVolumeDetach:     "error_detach",
	VolumeAttaching:       "attaching",
	VolumeDetaching:       "detaching",
	RecycledVolume:        "deleted", // inner-status same as deleted
	ErrorRecycledVolume:   "deleted", // inner-status same as deleted
	RecycledSnapshot:      "deleted", // inner-status same as deleted
	ErrorRecycledSnapshot: "deleted", // inner-status same as deleted
	RecyclingSnapshot:     "deleted", // inner-status same as deleted
}

//task type
const (
	VolumeTask      = "volume_task"
	ReplicationTask = "replication_task"
	SnapshotTask    = "snapshot_task"
	SnapRecycleTask = "snap_recycle_task"
)

//task status
const (
	DB_TASK_STATUS_INVALID    = -1
	DB_TASK_STATUS_PREPARE    = 0
	DB_TASK_STATUS_BEGIN      = 1
	DB_TASK_STATUS_INPROGRESS = 2
	DB_TASK_STATUS_SUCCESS    = 3
	DB_TASK_STATUS_FAILED     = 4
	DB_TASK_STATUS_OMIT       = 5
)

//ebs task status
const (
	DB_EBS_TASK_STATUS_BEGIN      = 5
	DB_EBS_TASK_STATUS_INPROGRESS = 6
	DB_EBS_TASK_STATUS_SUCCESS    = 7
	DB_EBS_TASK_STATUS_FAILED     = 8
)

//volume_task action
const (
	EmptyVolumeTask       = 0
	CreateVolumeTask      = 1
	DeleteVolumeTask      = 2
	AttachVolumeTask      = 3
	DetachVolumeTask      = 4
	ResizeVolumeTask      = 5 // TODO
	CreateVolumeErrorTask = 6
	DeleteVolumeErrorTask = 7
)

//snapshot type
const (
	GENERIC_SNAPSHOT       = 0
	PUBLIC_IMAGE_SNAPSHOT  = 1
	PRIVATE_IMAGE_SNAPSHOT = 2
	COPIED_SNAPSHOT        = 3
)
const (
	IMAGE_POOL_ID = "pool-image"
)

//snapshot task action
const (
	CreateSnapshotTask             = 1 << 0
	DeleteSnapshotTask             = 1 << 1
	RestoreSnapshotTask            = 1 << 2
	CreateVolByRestoreSnapshotTask = 1 << 3
	ConvertImageToTask             = 1 << 4
	CopySnapshotTask               = 1 << 5
)
const (
	URL_PREFIX_RNODE_PUT_STAT  = "recoverynodeputstat"
	URL_PREFIX_CNODE_PUT_STAT  = "chunknodeputstat"
	URL_PREFIX_VMMGR_PUT_STAT  = "volumemgrputstat"
	URL_PREFIX_VOLUME_PUT_STAT = "volumeputstat"
)

//recycle task action
const (
	RecycleStartTask = 1 << 0
	RecycleEndTask   = 1 << 1
)

// proxy
const (
	PROXY_TYPE_INVALID           = 0
	PROXY_TYPE_REPLICATION_GROUP = 1
	PROXY_TYPE_DISK              = 2
	PROXY_TYPE_VOLUME            = 3
	PROXY_TYPE_ZBS_PROXY_LIST    = 4
)

const (
	// 不需要更新Version
	ZBS_PROXY_LIST_INVALID_VERSION = -1
	// 无效的Version
	ZBS_PROXY_LIST_BASE_VERSION = 0
)

const (
	CONFIG_KEY_ZBS_PROXY_VERSION = "zbs-proxy-version"
)

const (
	LB_ACTION_GETDETAIL     = "GetDiskDetail"
	LB_ACTION_HEARTBEAT     = "HeartBeat"
	LB_ACTION_REPLICASTATUS = "ReplicaStatus"
)

type DiskStatusType int

const (
	DISK_STATUS_UP   = DiskStatusType(1)
	DISK_STATUS_DOWN = DiskStatusType(2)
)

type DiskAdminStatusType int

const (
	DISK_ADMIN_STATUS_IN  = DiskAdminStatusType(1)
	DISK_ADMIN_STATUS_OUT = DiskAdminStatusType(2)
)

type ReplicationTaskActionType int

const (
	REPLICATION_TASK_ACTION_RESCHEDULE      = ReplicationTaskActionType(1)
	REPLICATION_TASK_ACTION_LEADER_TRANSFER = ReplicationTaskActionType(2)
	REPLICATION_TASK_ACTION_MOVE_REPLICA    = ReplicationTaskActionType(3)
)

//admin task opt_type
const (
	_ = iota
	ADMIN_ADD
	ADMIN_DEL
	ADMIN_UP
	ADMIN_DOWN
)

//admin task type
const (
	_ = iota
	POOL_TASK
	RACK_TASK
	HOST_TASK
	DISK_TASK
)

const (
	_ = iota
	POOL_ADD
	POOL_DEL
	RACK_ADD
	RACK_DEL
	HOST_ADD
	HOST_DEL
	DISK_ADD
	DISK_DEL
	DISK_DOWN
	DISK_UP
	RESCHEDULE
	LEADER_TRANSFER
	MOVE_REPLICA
	DELETE_REPLICA
	ADD_REPLICA
	POOL_ENABLE
	POOL_DISABLE
)

const (
	_ = iota
	REPLICA_STATUS_OK
	REPLICA_STATUS_INVALID
	REPLICA_STATUS_MIG_SRC // migration source
	REPLICA_STATUS_MIG_DST // migration destination
	REPLICA_STATUS_OLD_MIG_SRC
)

const (
	HB_REPLICA_STATUS_RECOVERING  = 1 << iota
	HB_REPLICA_STATUS_IN_ELECTION = 1 << iota
)

const (
	DEFAULT_REPLICA_COUNT = 3
)

const (
	BASE_PORT    = 30000
	MANAGE_PORT  = BASE_PORT
	STORAGE_PORT = BASE_PORT + 1000
	CLIENT_PORT  = BASE_PORT + 2000
	TRACE_PORT   = BASE_PORT + 3000
)

const (
	POOL_TYPE_INVALID = iota
	POOL_TYPE_EBS
	POOL_TYPE_PRERELEASE
	POOL_TYPE_GRAYSCALE
	POOL_TYPE_PRODUCTION
	POOL_TYPE_SPECIAL
)

// pool status
const (
	POOL_DISABLED = iota
	POOL_ENABLED
)

const (
	UNLIMITED = "unlimited"
)

const (
	MEDIA_TYPE_HDD = "hdd"
	MEDIA_TYPE_SSD = "ssd"
)

const (
	QuotaTypeVolume   = "volume"
	QuotaTypeSnapshot = "snapshot"
)

const (
	DISK_LOCAL_PORT_START = 40000
	DISK_LOCAL_PORT_END   = 65000
	DISK_USAGE_THROTTLE   = 90
)

const (
	VOL_LOCAL_PORT_START = 25000
	VOL_LOCAL_PORT_END   = 65000
)

const (
	MAX_OPEN_FD_LIMIT = 100000
)

//Flag bit of Request.flag
const (
	// req sent by zbs-client when volume is in snapshoting.
	REQ_FLAG_SNAPPING_IO = uint64(1 << 0)

	// req sent by zbs-worker side read specified version data or version map.
	REQ_FLAG_SNAP_ADMIN_IO = uint64(1 << 1)

	// req sent by zbs-client when volume is lazy volume.
	REQ_FLAG_LAZYING_IO = uint64(1 << 2)

	// req used inner zbs-storage to indicate no version update for this IO.
	REQ_FLAG_INNER_NOUPDATE_SNAP_VERSION = uint64(1 << 3)

	// flag for verify special io req.
	REQ_FLAG_VERIFY = uint64(1 << 4)
)

const (
	REQUEST_VERSION_1_0 = 10
	REQUEST_VERSION_2_0 = 20
)

const (
	SnapBlkShiftSize  = 22
	SnapObjectSize    = 1024 * 1024 * 1024
	SnapBlkSize       = 1 << SnapBlkShiftSize // 4MiB
	SnapBlkNum        = SnapObjectSize / SnapBlkSize
	SnapUploadSize    = 64 * 1024 * 1024
	SnapMetaMagic     = 0x1F2E3D4C
	SnapCompressAlg   = ""
	SnapEncryption    = ""
	SnapBPRspDataSize = 256
)

//zbs-cli cache action
const (
	ZBS_PROXY_ADD    = 1
	ZBS_PROXY_DELETE = 2
	ZBS_PROXY_LIST   = 3
)

const (
	ExclusiveSnapshot = 0
	SharedSnapshot    = 1
)

const (
	BUF_CHECK_SUM_CRC = iota
	BUF_CHECK_SUM_MD5
)

const (
	SV_ACTION_APPLYCOPYSNL  = "ApplyCopySnl"
	SV_ACTION_FINISHCOPYSNL = "FinishCopySnl"
)

const (
	REMOTE_REQUEST_TIMEOUT = time.Second * 3
)

const (
	SnapActionNone            = 0
	SnapActionCreating        = 1
	SnapActionGarbageDeleting = 2
	SnapActionRestoring       = 3
)

const (
	// AWS s3 sdk needs a region which is useless in JSS, so define a default one
	DEFAULT_JSS_REGION = "default_jss_region"
)
const (
	IMG_TYPE_QCOW2         = "qcow2"
	IMG_TYPE_RAW           = "raw"
	ANALYSIS_URL_S3_TPYE   = "s3"
	ANALYSIS_URL_FILE_TPYE = "file"
)

const (
	SNAP_META_VERSION_OLD = 0 // previous version: ebs, zbs 1.0, zbs 2.0
	SNAP_META_VERSION_3_0 = 1 // zbs 3.0
)

const (
	CLUSTER_NAME_EBS_PREFIX  = "cluster_ebs_"
	CLUSTER_NAME_PARENT_TASK = "cluster_snapshot_batch"
)

// cluster version
const (
	CLUSTER_VERSION_INVALID = 0
	CLUSTER_VERSION_ZBS_2   = 2
	CLUSTER_VERSION_ZBS_3   = 3
)

const (
	CLUSTER_STATUS_DISABLE = 0
	CLUSTER_STATUS_DEFAULT = 1
)
