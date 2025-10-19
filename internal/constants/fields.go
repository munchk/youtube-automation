package constants

// Phase titles used throughout the application
const (
	PhaseTitleInitialDetails    = "Initial Details"
	PhaseTitleWorkProgress      = "Work In Progress"
	PhaseTitleDefinition        = "Definition"
	PhaseTitlePostProduction    = "Post-Production"
	PhaseTitlePublishingDetails = "Publishing Details"
	PhaseTitlePostPublish       = "Post-Publish Details"
)

// Field titles used throughout the application - these must match the CLI form field titles
const (
	// Initial Details fields
	FieldTitleProjectName        = "Project Name"
	FieldTitleProjectURL         = "Project URL"
	FieldTitleSponsorshipAmount  = "Sponsorship Amount"
	FieldTitleSponsorshipEmails  = "Sponsorship Emails (comma separated)"
	FieldTitleSponsorshipBlocked = "Sponsorship Blocked Reason"
	FieldTitlePublishDate        = "Publish Date (YYYY-MM-DDTHH:MM)"
	FieldTitleDelayed            = "Delayed"
	FieldTitleGistPath           = "Gist Path (.md file)"

	// Work Progress fields
	FieldTitleCodeDone            = "Code Done"
	FieldTitleTalkingHeadDone     = "Talking Head Done"
	FieldTitleScreenRecordingDone = "Screen Recording Done"
	FieldTitleRelatedVideos       = "Related Videos (comma separated)"
	FieldTitleThumbnailsDone      = "Thumbnails Done"
	FieldTitleDiagramsDone        = "Diagrams Done"
	FieldTitleScreenshotsDone     = "Screenshots Done"
	FieldTitleFilesLocation       = "Files Location (e.g., Google Drive link)"
	FieldTitleTagline             = "Tagline"
	FieldTitleTaglineIdeas        = "Tagline Ideas"
	FieldTitleOtherLogos          = "Other Logos/Assets"

	// Definition fields
	FieldTitleTitle            = "Title"
	FieldTitleDescription      = "Description"
	FieldTitleTags             = "Tags"
	FieldTitleDescriptionTags  = "Description Tags"
	FieldTitleTweet            = "Tweet"
	FieldTitleAnimationsScript = "Animations Script"

	// Post Production fields
	FieldTitleThumbnailPath = "Thumbnail Path"
	FieldTitleMembers       = "Members (comma separated)"
	FieldTitleRequestEdit   = "Edit Request"
	FieldTitleTimecodes     = "Timecodes"
	FieldTitleMovieDone     = "Movie Done"
	FieldTitleSlidesDone    = "Slides Done"

	// Publishing fields
	FieldTitleVideoFilePath   = "Video File Path"
	FieldTitleUploadToYouTube = "Upload Video to YouTube?"
	FieldTitleCurrentVideoID  = "Current YouTube Video ID"
	FieldTitleCreateHugo      = "Create/Update Hugo Post"

	// Post Publish fields
	FieldTitleDOTPosted           = "DevOpsToolkit Post Sent (manual)"
	FieldTitleBlueSkyPosted       = "BlueSky Post Sent"
	FieldTitleLinkedInPosted      = "LinkedIn Post Sent (manual)"
	FieldTitleSlackPosted         = "Slack Post Sent"
	FieldTitleYouTubeHighlight    = "YouTube Highlight Created (manual)"
	FieldTitleYouTubeComment      = "YouTube Pinned Comment Added (manual)"
	FieldTitleYouTubeCommentReply = "Replied to YouTube Comments (manual)"
	FieldTitleGDEPosted           = "GDE Advocu Post Sent (manual)"
	FieldTitleCodeRepository      = "Code Repository URL"
	FieldTitleNotifySponsors      = "Notify Sponsors"
)

// Language constants following ISO 639-1 standard
const (
	// DefaultLanguage is the default language code for YouTube videos
	DefaultLanguage = "en"
	// LanguageEnglish is the ISO 639-1 code for English
	LanguageEnglish = "en"
	// Add other languages as needed for future expansion
)

// LanguageMap maps language codes to their full names for better readability
var LanguageMap = map[string]string{
	LanguageEnglish: "English",
	// Add other languages as needed for future expansion
}

// IsValidLanguage checks if a language code is valid according to our supported languages
func IsValidLanguage(code string) bool {
	_, exists := LanguageMap[code]
	return exists
}
