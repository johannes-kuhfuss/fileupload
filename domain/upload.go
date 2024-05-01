package domain

type Upload struct {
	FileName  string
	BcDate    string
	StartTime string
	EndTime   string
	Status    string
	Uploader  string
	Size      string
}

type Uploads []Upload
