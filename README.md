# Vast

VAST 4.2 : https://github.com/haxqer/vast/tree/4.2

XML/Json

:star: VAST Ad generator and parser library on GoLang.

## todo
- [x] Support for all TrackEvents

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
* [VAST Samples](https://github.com/InteractiveAdvertisingBureau/VAST_Samples)

## Installation

`go get -u github.com/haxqer/vast`



## Quick Start

```go
package main

import (
	"encoding/xml"
	"fmt"
	. "github.com/haxqer/vast"
	"time"
)

func main()  {
	v := VAST{
		Version: "3.0",
		Ads: []Ad{
			{
				ID:   "123",
				Type: "front",
				InLine: &InLine{
					AdSystem: &AdSystem{Name: "DSP"},
					AdTitle:  CDATAString{CDATA: "adTitle"},
					Impressions: []Impression{
						{ID: "11111", URI: "http://impressionv1.track.com"},
						{ID: "11112", URI: "http://impressionv2.track.com"},
					},
					Creatives: []Creative{
						{
							ID:       "987",
							Sequence: 0,
							Linear: &Linear{
								Duration: Duration(15 * time.Second),
								TrackingEvents: []Tracking{
									{Event: Event_type_start, URI: "http://track.xxx.com/q/start?xx"},
									{Event: Event_type_firstQuartile, URI: "http://track.xxx.com/q/firstQuartile?xx"},
									{Event: Event_type_midpoint, URI: "http://track.xxx.com/q/midpoint?xx"},
									{Event: Event_type_thirdQuartile, URI: "http://track.xxx.com/q/thirdQuartile?xx"},
									{Event: Event_type_complete, URI: "http://track.xxx.com/q/complete?xx"},
								},
								MediaFiles: []MediaFile{
									{
										Delivery: "progressive",
										Type:     "video/mp4",
										Width:    1024,
										Height:   576,
										URI:      "http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4",
									},
								},
							},
						},
					},
				},
			},
		},
	}
	vastXMLText, _ := xml.Marshal(v)
	fmt.Printf("%s", vastXMLText)
}

```

Result Demo
```xml
<VAST version="3.0">
    <Ad id="123" type="front">
        <InLine>
            <AdSystem><![CDATA[DSP]]></AdSystem>
            <AdTitle><![CDATA[adTitle]]></AdTitle>
            <Impression id="11111"><![CDATA[http://impressionv1.track.com]]></Impression>
            <Impression id="11112"><![CDATA[http://impressionv2.track.com]]></Impression>
            <Creatives>
                <Creative id="987">
                    <Linear>
                        <Duration>00:00:15</Duration>
                        <TrackingEvents>
                            <Tracking event="start"><![CDATA[http://track.xxx.com/q/start?xx]]></Tracking>
                            <Tracking event="firstQuartile"><![CDATA[http://track.xxx.com/q/firstQuartile?xx]]></Tracking>
                            <Tracking event="midpoint"><![CDATA[http://track.xxx.com/q/midpoint?xx]]></Tracking>
                            <Tracking event="thirdQuartile"><![CDATA[http://track.xxx.com/q/thirdQuartile?xx]]></Tracking>
                            <Tracking event="complete"><![CDATA[http://track.xxx.com/q/complete?xx]]></Tracking>
                        </TrackingEvents>
                        <MediaFiles>
                            <MediaFile delivery="progressive" type="video/mp4" width="1024" height="576"><![CDATA[http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4]]></MediaFile>
                        </MediaFiles>
                    </Linear>
                </Creative>
            </Creatives>
            <Description></Description>
            <Survey></Survey>
        </InLine>
    </Ad>
</VAST>

```

## Thanks
+ https://github.com/rs/vast
+ https://github.com/xsharp