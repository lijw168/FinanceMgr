package options

type CreateVolumeOptions struct {
	TenantId       string `json:"tenant_id"`
	VolumeName     string `json:"volume_name"`
	Size           uint64 `json:"size"`
	VolumeTypeName string `json:"name"`
	AzName         string `json:"az_name"`
	Description    string `json:"description"`
	SnapshotId     string `json:"snapshot_id"`
	Bootable       uint8  `json:"tag"`
	ClientToken    string `json:"clientToken"`
}
type ChangeVolumeOptions struct {
	TenantId    string  `json:"tenant_id"`
	VolumeId    string  `json:"id"`
	VolumeName  *string `json:"name"`
	Description *string `json:"description"`
}
type ResizeVolumeOptions struct {
	TenantId string `json:"tenant_id"`
	VolumeId string `json:"volume_id"`
	Size     uint64 `json:"size"`
}
type DeleteVolumeOptions struct {
	TenantId string `json:"tenant_id"`
	VolumeId string `json:"volume_id"`
}

type DescribeVolumeOptions struct {
	TenantId string `json:"tenant_id"`
	VolumeId string `json:"volume_id"`
}

type AttachVolumeOptions struct {
	TenantId        string `json:"tenant_id"`
	HostIp          string `json:"host_ip"`
	VolumeId        string `json:"volume_id"`
	InstanceUuid    string `json:"instance_uuid"`
	InstanceType    string `json:"instance_type"`
	MultiAttachment bool   `json:"multi_attachment"`
	AttachMode      string `json:"attach_mode"`
}

type DetachVolumeOptions struct {
	TenantId     string `json:"tenant_id"`
	VolumeId     string `json:"volume_id"`
	AttachmentId string `json:"attachment_id"`
}
type DescribeAttachmentOptions struct {
	TenantId     string `json:"tenant_id"`
	AttachmentId string `json:"attachment_id"`
}

type AdminVolumeOptions struct {
	TenantId string `json:"tenant_id"`
	VolumeId string `json:"volume_id"`
}
