package vast

const (
	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Player Operation Metrics (for use in Linear and NonLinear Ads)

	// the user activated the mute control and muted the creative.
	EventTypeMute = "mute"
	// the user activated the mute control and unmuted the creative.
	EventTypeUnmute = "unmute"
	// the user clicked the pause control and stopped the creative.
	EventTypePause = "pause"
	// the user activated the resume control after the creative had been stopped or paused.
	EventTypeResume = "resume"
	// the user activated the rewind control to access a previous point in the creative timeline.
	EventTypeRewind = "rewind"
	// the user activated a skip control to skip the creative, which is a
	// different control than the one used to close the creative.
	EventTypeSkip = "skip"
	// the user activated a control to extend the player to a larger size. This
	// event replaces the fullscreen event per the 2014 Digital Video In-Stream Ad Metric
	// Definitions.
	EventTypePlayerExpand = "playerExpand"
	// the user activated a control to reduce player to a smaller size. This
	// event replaces the exitFullscreen event per the 2014 Digital Video In-Stream Ad
	// Metric Definitions.
	EventTypePlayerCollapse = "playerCollapse"
	// This ad was not and will not be played (e.g. it was prefetched for a particular
	// ad break but was not chosen for playback). This allows ad servers to reuse an ad earlier
	// than otherwise would be possible due to budget/frequency capping. This is a terminal
	// event; no other tracking events should be sent when this is used. Player support is
	// optional and if implemented is provided on a best effort basis as it is not technically
	// possible to fire this event for every unused ad (e.g. when the player itself is terminated
	// before playback)
	EventTypeNotUsed = "notUsed"

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Linear Ad Metrics

	// This event should be used to indicate when the player considers that
	// it has loaded and buffered the creative’s media and assets either fully or
	// to the extent that it is ready to play the media
	EventTypeLoaded = "loaded"
	// this event is used to indicate that an individual creative within the ad was loaded and playback
	//began. As with creativeView, this event is another way of tracking creative playback.
	EventTypeStart = "start"
	// the creative played for at least 25% of the total duration.
	EventTypeFirstQuartile = "firstQuartile"
	// the creative played for at least 50% of the total duration.
	EventTypeMidpoint = "midpoint"
	// the creative played for at least 75% of the duration.
	EventTypeThirdQuartile = "thirdQuartile"
	// The creative was played to the end at normal speed.
	EventTypeComplete = "complete"
	// An optional metric that can capture all other user interactions under one metric such a s hover-overs, or custom clicks.
	// It should NOT replace clickthrough events or other existing events like mute, unmute, pause, etc.
	EventTypeOtherAdInteraction = "otherAdInteraction"
	// the creative played for a duration at normal speed that is equal to or greater than the
	// value provided in an additional attribute for offset . Offset values can be time in the format
	// HH:MM:SS or HH:MM:SS.mmm or a percentage value in the format n% . Multiple progress ev
	EventTypeProgress = "progress"
	// the user clicked the close button on the creative. The name of this event distinguishes it
	// from the existing “close” event described in the 2008 IAB Digital Video In-Stream Ad Metrics
	// Definitions, which defines the “close” metric as applying to non-linear ads only. The “closeLinear” event
	// extends the “close” event for use in Linear creative.
	EventTypeCloseLinear = "closeLinear"

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// NonLinear Ad Metrics

	// Not to be confused with an impression, this event indicates that an individual creative portion of the ad was viewed.
	// An impression indicates that at least a portion of the ad was displayed;
	// however an ad may be composed of multiple creative, or creative that only play on some platforms and not others.
	// This event enables ad servers to track which ad creative are viewed, and therefore, which platforms are more common.
	EventTypeCreativeView = "creativeView"
	// The user clicked or otherwise activated a control used to pause streaming content,
	// which either expands the ad within the player’s viewable area or “takes-over” the streaming content area
	// by launching an additional portion of the ad.
	// An ad in video format ad is usually played upon acceptance, but other forms of media such as games, animation,
	// tutorials, social media, or other engaging media are also used.
	EventTypeAcceptInvitation = "acceptInvitation"
	// The user activated a control to expand the creative.
	EventTypeAdExpand = "adExpand"
	// The user activated a control to reduce the creative to its original dimensions.
	EventTypeAdCollapse = "adCollapse"
	// The user clicked or otherwise activated a control used to minimize the ad
	// to a size smaller than a collapsed ad but without fully dispatching the ad from the player environment.
	// Unlike a collapsed ad that is big enough to display it’s message,
	// the minimized ad is only big enough to offer a control that enables the user to redisplay the ad if desired.
	EventTypeMinimize = "minimize"
	// The user clicked or otherwise activated a control for removing the ad,
	// which fully dispatches the ad from the player environment in a manner that does not allow the user to re-display the ad.
	EventTypeClose = "close"
	// The time that the initial ad is displayed.
	// This time is based on the time between the impression and either the completed length of display based on the agreement between
	// transactional parties or a close, minimize, or accept invitation event.
	EventTypeOverlayViewDuration = "overlayViewDuration"

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Interactive Ad Metric

	// With VAST 4, video playback and interactive creative playback now happens in parallel.
	// Video playback and interactive creative start may not happen at the same time.
	// A separate way of tracking the interactive creative start is needed.
	// The interactive creative specification (SIMID, etc.) will define when this event should be fired.
	EventTypeInteractiveStart = "interactiveStart"

	// ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~
	// Other

	EventTypeView    = "view"
	EventTypeMonitor = "monitor"
)
