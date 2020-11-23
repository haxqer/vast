package vast

// The ad server may provide URIs for tracking publisher-determined view-ability
type ViewableImpression struct {
	// An ad server id for the impression.
	// Viewable impression resources of the same id should be requested at the same time,
	// or as close in time as possible, to help prevent discrepancies.
	ID string `xml:"id,attr"`
	// The <Viewable> element is used to place a URI that the player triggers if and when
	// the ad meets criteria for a viewable video ad impression.
	Viewable *[]CDATAString `xml:"Viewable"`
	// The <NotViewable> element is a container for placing a URI that the player triggers
	// if the ad is executed but never meets criteria for a viewable video ad impression.
	NotViewable *[]CDATAString `xml:"NotViewable"`
	// The <ViewUndetermined> element is a container for placing a URI that the player triggers
	// if it cannot determine whether the ad has met criteria for a viewable video ad impression.
	ViewUndetermined *[]CDATAString `xml:"ViewUndetermined"`
}

// Providing an advertiser name can help publishers prevent display of the ad with its competitors.
type Advertiser struct {
	// An (optional) identifier for the advertiser, provided by the ad server. Can be used for internal analytics.
	ID string `xml:"id,attr"`
	// A string that provides the name of the advertiser as defined by the ad serving party.
	// Recommend using the domain of the advertiser.
	Advertiser string `xml:",chardata"`
}

type AdVerifications struct {
	Verification *Verification `xml:",omitempty"`
}

// The <Verification> element contains the executable and bootstrapping data required to run the measurement code for a single verification vendor.
// Multiple <Verification> elements may be listed, in order to support multiple vendors, or if multiple API frameworks are supported.
// At least one <JavaScriptResource> or <ExecutableResource> should be provided.
// At most one of these resources should selected for execution, as best matches the technology available in the current environment.
type Verification struct {
	// An identifier for the verification vendor. The recommended format is [domain]- [useCase],
	// to avoid name collisions. For example, "company.com-omid".
	Vendor string `xml:"vendor,attr"`
	// A container for the URI to the JavaScript file used to collect verification data.
	// Some verification vendors may provide JavaScript executables which work in non-browser environments,
	// for example, in an iOS app enabled by JavaScriptCore. These resources only require methods of the API framework,
	// without relying on any browser built-in functionality.
	JavaScriptResource []JavaScriptResource
	// A reference to a non-JavaScript or custom-integration resource intended for collecting verification data via the listed apiFramework.
	ExecutableResource []ExecutableResource
	// The verification vendor may provide URIs for tracking events relating to the execution of their code during the ad session.
	TrackingEvents []Tracking `xml:"TrackingEvents>Tracking,omitempty"`
	// <VerificationParameters> contains a CDATA-wrapped string intended for bootstrapping the verification code and providing metadata about the current impression.
	// The format of the string is up to the individual vendor and should be passed along verbatim.
	VerificationParameters string `xml:",omitempty"`
	// ad categories are used in creative separation and for compliance in certain programs
	BlockedAdCategories []Category `xml:",omitempty"`
}

// A container for the URI to the JavaScript file used to collect verification data.
type JavaScriptResource struct {
	// Identifies the API needed to execute the resource file if applicable.
	ApiFramework string `xml:"apiFramework,attr"`
	// If "true", this resource is optimized and able to execute in
	// an environment without DOM and other browser built-ins (e.g. iOS' JavaScriptCore).
	BrowserOptional bool `xml:"browserOptional,attr"`
	// A CDATA-wrapped URI to a file providing Closed Caption info for the media file.
	URI string `xml:",cdata"`
}

type ExecutableResource struct {
	// Identifies the API needed to execute the resource file if applicable.
	ApiFramework string `xml:"apiFramework,attr,omitempty"`
	// Identifies the MIME type of the file provided.
	Type bool `xml:"type,attr,omitempty"`
	// A CDATA-wrapped URI to a file providing Closed Caption info for the media file.
	URI string `xml:",cdata"`
}

type Mezzanine struct {
	// Method of delivery of ad (either "streaming" or "progressive")
	Delivery string `xml:"delivery,attr"`
	// MIME type. Popular MIME types include, but are not limited to
	// “video/x-ms-wmv” for Windows Media, and “video/x-flv” for Flash
	// Video. Image ads or interactive ads can be included in the
	// MediaFiles section with appropriate Mime types
	Type string `xml:"type,attr"`
	// Pixel dimensions of video.
	Width int `xml:"width,attr"`
	// Pixel dimensions of video.
	Height int `xml:"height,attr"`
	// The codec used to produce the media file.
	Codec string `xml:"codec,attr,omitempty" json:",omitempty"`
	// Optional identifier
	ID string `xml:"id,attr,omitempty" json:",omitempty"`
	// Optional field that helps eliminate the need to calculate the size based on bitrate and duration.
	// Units - Bytes
	FileSize int `xml:"fileSize,attr,omitempty" json:",omitempty"`
	// Type of media file (2D / 3D / 360 / etc). Optional.
	// Default value = 2D
	MediaType string `xml:"mediaType,attr,omitempty" json:",omitempty"`
}

type InteractiveCreativeFile struct {
	// Identifies the API needed to execute the resource file if applicable.
	ApiFramework string `xml:"apiFramework,attr,omitempty"`
	// Identifies the MIME type of the file provided.
	Type bool `xml:"type,attr,omitempty"`
	// Useful for interactive use cases.
	// Identifies whether the ad always drops when the duration is reached,
	// or if it can potentially extend the duration by pausing the underlying video or delaying the adStopped call after adVideoComplete.
	// If it set to true the extension of the duration should be user-initiated (typically by engaging with an interactive element to view additional content).
	VariableDuration bool `xml:"variableDuration,attr,omitempty"`
	// A CDATA-wrapped URI to a file providing Closed Caption info for the media file.
	URI string `xml:",cdata"`
}

type ClosedCaptionFile struct {
	// Identifies the MIME type of the file provided.
	Type bool `xml:"type,attr,omitempty"`
	// Language of the Closed Caption File using ISO 631-1 codes.
	// An optional locale suffix can also be provided.
	// Examples - “en”, “en-US”, “zh-TW”,
	Language bool `xml:"language,attr,omitempty"`
	// A CDATA-wrapped URI to a file providing Closed Caption info for the media file.
	URI string `xml:",cdata"`
}
