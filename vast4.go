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
