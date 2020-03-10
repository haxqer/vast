package vast

const (
	/**
	 * not to be confused with an impression, this event indicates that an individual creative
	 * portion of the ad was viewed. An impression indicates the first frame of the ad was displayed; however
	 * an ad may be composed of multiple creative, or creative that only play on some platforms and not
	 * others. This event enables ad servers to track which ad creative are viewed, and therefore, which
	 * platforms are more common.
	 */

	Event_type_creativeView = "creativeView"
	Event_type_view         = "view"

	/**
	 * this event is used to indicate that an individual creative within the ad was loaded and playback
	 * began. As with creativeView, this event is another way of tracking creative playback.
	 */
	Event_type_start = "start"

	// the creative played for at least 25% of the total duration.
	Event_type_firstQuartile = "firstQuartile"

	// the creative played for at least 50% of the total duration.
	Event_type_midpoint = "midpoint"

	// the creative played for at least 75% of the duration.
	Event_type_thirdQuartile = "thirdQuartile"

	// The creative was played to the end at normal speed.
	Event_type_complete = "complete"

	// the user activated the mute control and muted the creative.
	Event_type_mute = "mute"

	// the user activated the mute control and unmuted the creative.
	Event_type_unmute = "unmute"

	// the user clicked the pause control and stopped the creative.
	Event_type_pause = "pause"

	// the user activated the rewind control to access a previous point in the creative timeline.
	Event_type_rewind = "rewind"

	// the user activated the resume control after the creative had been stopped or paused.
	Event_type_resume = "resume"

	// the user activated a control to extend the video player to the edges of the viewer’s screen.
	Event_type_fullscreen = "fullscreen"

	// the user activated the control to reduce video player size to original dimensions.
	Event_type_exitFullscreen = "exitFullscreen"

	// the user activated a control to expand the creative.
	Event_type_expand = "expand"

	// the user activated a control to reduce the creative to its original dimensions.
	Event_type_collapse = "collapse"

	/**
	 * the user activated a control that launched an additional portion of the
	 * creative. The name of this event distinguishes it from the existing “acceptInvitation” event described in
	 * the 2008 IAB Digital Video In-Stream Ad Metrics Definitions, which defines the “acceptInivitation”
	 * metric as applying to non-linear ads only. The “acceptInvitationLinear” event extends the metric for use
	 * in Linear creative.
	 */
	Event_type_acceptInvitationLinear = "acceptInvitationLinear"

	/**
	 * the user clicked the close button on the creative. The name of this event distinguishes it
	 * from the existing “close” event described in the 2008 IAB Digital Video In-Stream Ad Metrics
	 * Definitions, which defines the “close” metric as applying to non-linear ads only. The “closeLinear” event
	 * extends the “close” event for use in Linear creative.
	 */
	Event_type_closeLinear = "closeLinear"

	Event_type_close = "close"

	// the user activated a skip control to skip the creative, which is a
	// different control than the one used to close the creative.
	Event_type_skip = "skip"

	/**
	 * the creative played for a duration at normal speed that is equal to or greater than the
	 * value provided in an additional attribute for offset . Offset values can be time in the format
	 * HH:MM:SS or HH:MM:SS.mmm or a percentage value in the format n% . Multiple progress ev
	 */
	Event_type_progress = "progress"

	Event_type_monitor = "monitor"
)
