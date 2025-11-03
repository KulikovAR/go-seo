package usecases

const (
	ErrorSiteExists   = "SITE_EXISTS"
	ErrorSiteNotFound = "SITE_NOT_FOUND"
	ErrorSiteCreation = "SITE_CREATION_FAILED"
	ErrorSiteDeletion = "SITE_DELETION_FAILED"
	ErrorSiteFetch    = "SITE_FETCH_FAILED"

	ErrorKeywordExists   = "KEYWORD_EXISTS"
	ErrorKeywordNotFound = "KEYWORD_NOT_FOUND"
	ErrorKeywordCreation = "KEYWORD_CREATION_FAILED"
	ErrorKeywordUpdate   = "KEYWORD_UPDATE_FAILED"
	ErrorKeywordDeletion = "KEYWORD_DELETION_FAILED"
	ErrorKeywordFetch    = "KEYWORD_FETCH_FAILED"

	ErrorPositionCreation = "POSITION_CREATION_FAILED"
	ErrorPositionDeletion = "POSITION_DELETION_FAILED"
	ErrorPositionFetch    = "POSITION_FETCH_FAILED"

	ErrorGroupExists   = "GROUP_EXISTS"
	ErrorGroupNotFound = "GROUP_NOT_FOUND"
	ErrorGroupCreation = "GROUP_CREATION_FAILED"
	ErrorGroupDeletion = "GROUP_DELETION_FAILED"
	ErrorGroupFetch    = "GROUP_FETCH_FAILED"

	ErrorValidation = "VALIDATION_ERROR"
	ErrorInternal   = "INTERNAL_ERROR"
)
