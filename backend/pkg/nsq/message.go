package nsq

type UploadImageMessage struct {
	ImageUrl string
}

type ImageTagsMessage struct {
	ImageUrl string
	Tags     []string
}
