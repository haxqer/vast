package vast

// The ad server may provide URIs for tracking publisher-determined view-ability
type ViewableImpression struct {
	// An ad server id for the impression.
	// Viewable impression resources of the same id should be requested at the same time, or as close in time as possible, to help prevent discrepancies.
	ID string `xml:"id,attr"`
	// The <Viewable> element is used to place a URI that the player triggers if and when the ad meets criteria for a viewable video ad impression.
	Viewable *[]CDATAString `xml:"Viewable"`
	// The <NotViewable> element is a container for placing a URI that the player triggers if the ad is executed but never meets criteria for a viewable video ad impression.
	NotViewable *[]CDATAString `xml:"NotViewable"`
	// The <ViewUndetermined> element is a container for placing a URI that the player triggers if it cannot determine whether the ad has met criteria for a viewable video ad impression.
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

// The <Verification> element contains the executable and bootstrapping data required to run the measurement code for a single verification vendor.
// Multiple <Verification> elements may be listed, in order to support multiple vendors, or if multiple API frameworks are supported.
// At least one <JavaScriptResource> or <ExecutableResource> should be provided.
// At most one of these resources should selected for execution, as best matches the technology available in the current environment.
type Verification struct {
	// An identifier for the verification vendor. The recommended format is [domain]- [useCase], to avoid name collisions. For example, "company.com-omid".
	Vendor string `xml:"vendor,attr"`

}

// A container for the URI to the JavaScript file used to collect verification data.
type JavaScriptResource struct {

}