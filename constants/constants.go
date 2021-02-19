package constants

const (
	WorkerName         = "mangindo-feeder-worker"
	WorkerDefaultQueue = "mangindo-worker-default"

	ServerError              = "origin server error:"
	InvalidJSONResponseError = "invalid JSON response from origin server"

	GetMangaListCommand   = "GetMangaListCommand"
	GetChapterListCommand = "GetChapterListCommand"
	GetContentListCommand = "GetContentListCommand"

	GetMangasAPIPath   = "/mangindo/v1/mangas"
	GetChaptersAPIPath = "/mangindo/v1/mangas/{title_id}/chapters"
	GetContentsAPIPath = "/mangindo/v1/mangas/{title_id}/chapters/{chapter}/contents"

	TitleIDKeyParam = "title_id"
	ChapterKeyParam = "chapter"

	MangaCacheExpirationInMn   = 60
	ChapterCacheExpirationInMn = 30
	ContentCacheExpirationInMn = 60 * 48

	SetMangaCacheJob   = "SetMangaCacheJob"
	SetChapterCacheJob = "SetChapterCacheJob"
	SetContentCacheJob = "SetContentCacheJob"

	JobArgTitleID = "JobArg_TitleId"
	JobArgChapter = "JobArg_Chapter"

	NullText = "null"
)
