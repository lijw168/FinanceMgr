package options

type CreateSnapshotOptions struct {
	VolumeId     string `json:"volume_id"`
	TenantId     string `json:"tenant_id"`
	ClientToken  string `json:"clientToken"`
	SnapshotName string `json:"name"`
	Description  string `json:"description"`
}

type RestoreSnapshotOptions struct {
	SnapshotId string `json:"snapshot_id"`
	TenantId   string `json:"tenant_id"`
}

type DescribeSnapshotOptions struct {
	TenantId   string `json:"tenant_id"`
	SnapshotId string `json:"snapshot_id"`
}

type DeleteSnapshotOptions struct {
	TenantId   string `json:"tenant_id"`
	SnapshotId string `json:"snapshot_id"`
}

type ChangeSnapshotOptions struct {
	TenantId     string `json:"tenant_id"`
	SnapshotId   string `json:"id"`
	SnapshotName string `json:"name"`
	Description  string `json:"description"`
	Share        int8   `json:"share"`
}

type ConvertImageToSnapshotOptions struct {
	TenantId      string  `json:"tenant_id"`
	SnapshotName  string  `json:"name"`
	PublicImage   bool    `json:"public_image"`
	ImageFormat   string  `json:"image_format"`
	ImageLocation string  `json:"image_location"`
	ImageType     string  `json:"image_type"`
	ImageId       string  `json:"image_id"`
	ImageHash     string  `json:"image_hash"`
	AzName        *string `json:"az_name"`
}

type CopySnapshotListOptions struct {
	Region      string   `json:"region"`
	TenantId    string   `json:"tenant_id"`
	SnapshotIds []string `json:"snapshot_ids"`
	DstTenantId string   `json:"dst_tenant_id"`
}
