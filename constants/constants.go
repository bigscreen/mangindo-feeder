package constants

const (
	WorkerName         = "mangindo-feeder-worker"
	WorkerDefaultQueue = "mangindo-worker-default"

	ServerError              = "origin server error:"
	InvalidJSONResponseError = "invalid JSON response from origin server"

	GetMangaListCommand   = "GetMangaListCommand"
	GetChapterListCommand = "GetChapterListCommand"
	GetContentListCommand = "GetContentListCommand"

	GetMangasApiPath   = "/mangindo/v1/mangas"
	GetChaptersApiPath = "/mangindo/v1/mangas/{title_id}/chapters"
	GetContentsApiPath = "/mangindo/v1/mangas/{title_id}/chapters/{chapter}/contents"

	TitleIdKeyParam = "title_id"
	ChapterKeyParam = "chapter"

	MangaCacheExpirationInMn   = 60
	ChapterCacheExpirationInMn = 30
	ContentCacheExpirationInMn = 60 * 48

	SetMangaCacheJob   = "SetMangaCacheJob"
	SetChapterCacheJob = "SetChapterCacheJob"
	SetContentCacheJob = "SetContentCacheJob"

	JobArgTitleId = "JobArg_TitleId"
	JobArgChapter = "JobArg_Chapter"
)
