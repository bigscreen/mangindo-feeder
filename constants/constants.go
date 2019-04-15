package constants

const (
	WorkerName = "mangindo-feeder-worker"

	ServerError              = "origin server error:"
	InvalidJSONResponseError = "invalid JSON response from origin server"

	GetMangaListCommand = "GetMangaListCommand"
	GetChapterListCommand = "GetChapterListCommand"
	GetContentListCommand = "GetContentListCommand"

	GetMangasApiPath = "/mangindo/v1/mangas"
	GetChaptersApiPath = "/mangindo/v1/chapters"

	TitleIdKeyParam = "title_id"
)
