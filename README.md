# Vast

:star: VAST Ad generator and parser library on GoLang. XML / JSON.

The struct model targets **VAST 4.3** (the latest version) and is backwards
compatible with VAST 2.0, 3.0, 4.0, 4.1 and 4.2.

## VAST versions

| Version | Released | Status | Highlights |
| --- | --- | --- | --- |
| 1.0 | 2008 | Deprecated | Original linear video ad template. |
| 2.0 | 2009 | Supported | Multiple media files, companion & non-linear ads, tracking events. |
| 3.0 | 2012 | Supported | Ad pods (sequenced ads), skippable ads, industry icons (AdChoices), wrapper improvements, error macros. |
| 4.0 | 2016 | Supported | Separation of media files from interactive/VPAID code, `Mezzanine` for server-side ad stitching (SSAI), viewability, `UniversalAdId`, `Category`. |
| 4.1 | 2018 | Supported | VPAID/Flash deprecation, `InteractiveCreativeFile` + SIMID/OMID, audio ads (`adType`, DAAST merged in), `AdServingId`, `Expires`, `BlockedAdCategories`. |
| 4.2 | 2019 | Supported | Multiple `UniversalAdId`, `IconClickFallbackImages`, closed-caption files, `ClosedCaptionFiles`, verification refinements. |
| **4.3** | **2022** | **Supported (default)** | Macros managed on GitHub, `[PLAYBACKMETHODS]=7` (continuous/binge play), inline **data URI** support for `InteractiveCreativeFile`, CTV-focused clarifications. |
| CTV Addendum | 2024 | N/A (backwards-compatible extension) | ACIF ad registration, DSA compliance icons, high-resolution creatives — layered on top of any VAST 4.x version. |

## Specs
* VAST 2.0 Spec: http://www.iab.net/media/file/VAST-2_0-FINAL.pdf
* VAST 3.0 Spec: http://www.iab.com/wp-content/uploads/2015/06/VASTv3_0.pdf
* VAST 4.0 Spec:
  * http://www.iab.com/wp-content/uploads/2016/01/VAST_4-0_2016-01-21.pdf
  * https://www.iab.com/wp-content/uploads/2016/04/VAST4.0_Updated_April_2016.pdf
* VAST 4.1 Spec:
  * https://iabtechlab.com/wp-content/uploads/2018/11/VAST4.1-final-Nov-8-2018.pdf
* VAST 4.2 Spec:
  * https://iabtechlab.com/wp-content/uploads/2019/06/VAST_4.2_final_june26.pdf
* VAST 4.3 Spec:
  * https://github.com/InteractiveAdvertisingBureau/VAST4.x/blob/main/4.3.md
* VAST 4.x Macros (living list): https://interactiveadvertisingbureau.github.io/vast/vast4macros/vast4-macros-latest.html
* [VAST Samples](https://github.com/InteractiveAdvertisingBureau/VAST_Samples)

## Installation

`go get -u github.com/haxqer/vast`



## Quick Start

```go
package main

import (
	"encoding/xml"
	"fmt"
	"time"

	. "github.com/haxqer/vast"
)

func main() {
	v := VAST{
		Version: "4.3",
		XMLNS:   "http://www.iab.com/VAST",
		Ads: []Ad{{
			ID:     "123",
			AdType: "video",
			Type:   "front",
			InLine: &InLine{
				AdSystem:    &AdSystem{Name: "DSP", Version: "4.3"},
				AdServingId: "DSP-123-request-456",
				AdTitle:     PlainString{CDATA: "adTitle"},
				Impressions: []Impression{
					{ID: "11111", URI: "http://impressionv1.track.com"},
					{ID: "11112", URI: "http://impressionv2.track.com"},
				},
				Creatives: []Creative{{
					ID: "987",
					UniversalAdID: &[]UniversalAdID{
						{IDRegistry: "Ad-ID", ID: "8465"},
					},
					Linear: &Linear{
						Duration: Duration(15 * time.Second),
						TrackingEvents: &TrackingEvents{Tracking: []Tracking{
							{Event: EventTypeStart, URI: "http://track.xxx.com/q/start?xx"},
							{Event: EventTypeFirstQuartile, URI: "http://track.xxx.com/q/firstQuartile?xx"},
							{Event: EventTypeMidpoint, URI: "http://track.xxx.com/q/midpoint?xx"},
							{Event: EventTypeThirdQuartile, URI: "http://track.xxx.com/q/thirdQuartile?xx"},
							{Event: EventTypeComplete, URI: "http://track.xxx.com/q/complete?xx"},
						}},
						MediaFiles: &MediaFiles{MediaFile: []MediaFile{{
							Delivery: "progressive",
							Type:     "video/mp4",
							Width:    1024,
							Height:   576,
							URI:      "http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4",
						}}},
					},
				}},
			},
		}},
	}
	vastXMLText, _ := xml.Marshal(v)
	fmt.Printf("%s", vastXMLText)
}

```

Result Demo
```xml
<VAST version="4.3" xmlns="http://www.iab.com/VAST">
    <Ad id="123" adType="video" type="front">
        <InLine>
            <AdSystem version="4.3">DSP</AdSystem>
            <Impression id="11111"><![CDATA[http://impressionv1.track.com]]></Impression>
            <Impression id="11112"><![CDATA[http://impressionv2.track.com]]></Impression>
            <AdServingId>DSP-123-request-456</AdServingId>
            <AdTitle>adTitle</AdTitle>
            <Creatives>
                <Creative id="987">
                    <Linear>
                        <TrackingEvents>
                            <Tracking event="start"><![CDATA[http://track.xxx.com/q/start?xx]]></Tracking>
                            <Tracking event="firstQuartile"><![CDATA[http://track.xxx.com/q/firstQuartile?xx]]></Tracking>
                            <Tracking event="midpoint"><![CDATA[http://track.xxx.com/q/midpoint?xx]]></Tracking>
                            <Tracking event="thirdQuartile"><![CDATA[http://track.xxx.com/q/thirdQuartile?xx]]></Tracking>
                            <Tracking event="complete"><![CDATA[http://track.xxx.com/q/complete?xx]]></Tracking>
                        </TrackingEvents>
                        <Duration>00:00:15</Duration>
                        <MediaFiles>
                            <MediaFile delivery="progressive" type="video/mp4" width="1024" height="576"><![CDATA[http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4]]></MediaFile>
                        </MediaFiles>
                    </Linear>
                    <UniversalAdId idRegistry="Ad-ID">8465</UniversalAdId>
                </Creative>
            </Creatives>
        </InLine>
    </Ad>
</VAST>

```

## Model fidelity

`Companion`, `NonLinear`, and `Icon` expose the unbounded VAST resource elements
through `HTMLResources`, `IFrameResources`, and `StaticResources`. The older
singular fields remain available as compatibility aliases for the first resource.

Optional XML booleans retain their existing `bool` fields. When an explicit
`false` must be distinguished from an omitted attribute, set the matching
`...Set` field (for example, `VariableDurationSet` or `ScalableSet`). XML parsing
sets these presence fields automatically, so XML and JSON round trips remain
lossless.

## Thanks
+ https://github.com/rs/vast
+ https://github.com/xsharp
