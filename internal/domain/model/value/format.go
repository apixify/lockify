package value

type FileFormat string

const (
	Json   FileFormat = "json"
	DotEnv FileFormat = "dotenv"
)

func NewFileFormat(value string) FileFormat {
	return FileFormat(value)
}

func (fileFormat FileFormat) String() string {
	return string(fileFormat)
}

func (fileFormat FileFormat) IsJson() bool {
	return fileFormat == Json
}

func (fileFormat FileFormat) IsDotEnv() bool {
	return fileFormat == DotEnv
}

func (fileFormat FileFormat) IsValid() bool {
	return fileFormat.IsJson() || fileFormat.IsDotEnv()
}
