package vast

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/pquerna/ffjson/ffjson"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
	"time"
)

func TestQuickStartComplex(t *testing.T) {
	skip := Duration(5 * time.Second)
	v := VAST{
		Version: "4.2",
		Ads: []Ad{
			{
				ID:            "123",
				Type:          "front",
				AdType:        "video",
				ConditionalAd: false,
				InLine: &InLine{
					AdSystem: &AdSystem{Name: "DSP"},
					AdTitle:  PlainString{CDATA: "adTitle"},
					Impressions: []Impression{
						{ID: "11111", URI: "http://impressionv1.track.com"},
						{ID: "11112", URI: "http://impressionv2.track.com"},
					},
					Category: &[]Category{
						{Authority: "https://www.iabtechlab.com/categoryauthority", Category: "American Cuisine"},
						{Authority: "https://www.iabtechlab.com/categoryauthority", Category: "Guitar"},
					},
					Description: &CDATAString{"123"},
					ViewableImpression: &ViewableImpression{
						ID:       "1234",
						Viewable: &[]CDATAString{
							{CDATA: "http://viewable1.track.com"},
							{CDATA: "http://viewable2.track.com"},
						},
						NotViewable: &[]CDATAString{
							{CDATA: "http://notviewable1.track.com"},
							{CDATA: "http://notviewable2.track.com"},
						},
						ViewUndetermined: &[]CDATAString{
							{CDATA: "http://viewundetermined1.track.com"},
							{CDATA: "http://viewundetermined2.track.com"},
						},
					},
					Creatives: []Creative{
						{
							ID:       "987",
							Sequence: 0,
							AdID:     "12",
							UniversalAdID: &[]UniversalAdID{
								{
									IDRegistry: "Ad-ID",
									ID:         "8465",
								},
								{
									IDRegistry: "FOO-ID",
									ID:         "6666465",
								},
							},
							Linear: &Linear{
								SkipOffset: &Offset{
									Duration: &skip,
								},
								Duration: Duration(15 * time.Second),
								TrackingEvents: []Tracking{
									{Event: EventTypeStart, URI: "http://track.xxx.com/q/start?xx"},
									{Event: EventTypeFirstQuartile, URI: "http://track.xxx.com/q/firstQuartile?xx"},
									{Event: EventTypeMidpoint, URI: "http://track.xxx.com/q/midpoint?xx"},
									{Event: EventTypeThirdQuartile, URI: "http://track.xxx.com/q/thirdQuartile?xx"},
									{Event: EventTypeComplete, URI: "http://track.xxx.com/q/complete?xx"},
								},
								MediaFiles: []MediaFile{
									{
										Delivery: "progressive",
										Type:     "video/mp4",
										Width:    1024,
										Height:   576,
										URI:      "http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4",
										Label:    "123",
									},
								},
							},
						},
					},
					Extensions: &[]Extension{
						{
							Type: "ClassName",
							Data: "AdsVideoView",
						},
						{
							Type: "ExtURL",
							Data: "http://xxxxxxxx",
						},
					},
				},
			},
		},
	} //vastXMLText, _ := xml.Marshal(v)
	//fmt.Printf("%s", vastXMLText)

	out, _ := xml.MarshalIndent(v, " ", "  ")
	fmt.Println(string(out))
}

func TestQuickStart(t *testing.T) {
	d := Duration(5 * time.Second)
	v := VAST{
		Mute: true,
		Version: "3.0",
		Ads: []Ad{
			{
				ID:   "123",
				Type: "front",
				InLine: &InLine{
					AdSystem: &AdSystem{Name: "DSP"},
					AdTitle:  PlainString{CDATA: "adTitle"},
					Impressions: []Impression{
						{ID: "11111", URI: "http://impressionv1.track.com"},
						{ID: "11112", URI: "http://impressionv2.track.com"},
					},
					Creatives: []Creative{
						{
							ID:       "987",
							Sequence: 0,
							Linear: &Linear{
								SkipOffset: &Offset{
									Duration: &d,
								},
								Duration: Duration(15 * time.Second),
								TrackingEvents: []Tracking{
									{Event: EventTypeStart, URI: "http://track.xxx.com/q/start?xx"},
									{Event: EventTypeFirstQuartile, URI: "http://track.xxx.com/q/firstQuartile?xx"},
									{Event: EventTypeMidpoint, URI: "http://track.xxx.com/q/midpoint?xx"},
									{Event: EventTypeThirdQuartile, URI: "http://track.xxx.com/q/thirdQuartile?xx"},
									{Event: EventTypeComplete, URI: "http://track.xxx.com/q/complete?xx"},
								},
								MediaFiles: []MediaFile{
									{
										Delivery: "progressive",
										Type:     "video/mp4",
										Width:    1024,
										Height:   576,
										URI:      "http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4",
										Label: "123",
									},
								},
							},
						},
					},
					Extensions:  &[]Extension{
						{
							Type:           "ClassName",
							Data:           "AdsVideoView",
						},
						{
							Type:           "ExtURL",
							Data:           "http://xxxxxxxx",
						},
					},
				},
			},
		},
	}

	want := []byte(`{"Version":"3.0","Ad":[{"ID":"123","Type":"front","InLine":{"AdSystem":{"Data":"DSP"},"AdTitle":{"Data":"adTitle"},"Impressions":[{"ID":"11111","URI":"http://impressionv1.track.com"},{"ID":"11112","URI":"http://impressionv2.track.com"}],"Creatives":[{"ID":"987","Linear":{"SkipOffset":"00:00:05","Duration":"00:00:15","TrackingEvents":[{"Event":"start","URI":"http://track.xxx.com/q/start?xx"},{"Event":"firstQuartile","URI":"http://track.xxx.com/q/firstQuartile?xx"},{"Event":"midpoint","URI":"http://track.xxx.com/q/midpoint?xx"},{"Event":"thirdQuartile","URI":"http://track.xxx.com/q/thirdQuartile?xx"},{"Event":"complete","URI":"http://track.xxx.com/q/complete?xx"}],"MediaFiles":[{"Delivery":"progressive","Type":"video/mp4","Width":1024,"Height":576,"URI":"http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4","Label":"123"}]}}],"Extensions":[{"Type":"ClassName","Data":"AdsVideoView"},{"Type":"ExtURL","Data":"http://xxxxxxxx"}]}}],"Mute":true}`)
	got, err := json.Marshal(v)
	t.Logf("%s", got)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Marshal() got = %v, want %v", got, want)
	}

}

func TestEmptyVast(t *testing.T)  {
	v := VAST{
		Version: "3.0",
		Errors: []CDATAString{
			{CDATA: "http://xx.xx.com/e/error?e=__ERRORCODE__&co=__CONTENTPLAYHEAD__&ca=__CACHEBUSTING__&a=__ASSETURI__&t=__TIMESTAMP__&o=__OTHER__"},
		},
	}
	want := []byte(`{"Version":"3.0","Errors":[{"Data":"http://xx.xx.com/e/error?e=__ERRORCODE__\u0026co=__CONTENTPLAYHEAD__\u0026ca=__CACHEBUSTING__\u0026a=__ASSETURI__\u0026t=__TIMESTAMP__\u0026o=__OTHER__"}]}`)
	got, err := json.Marshal(v)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Logf("%s", got)
		t.Errorf("Marshal() got = %v, want %v", got, want)
	}

	want = []byte(`<VAST version="3.0"><Error><![CDATA[http://xx.xx.com/e/error?e=__ERRORCODE__&co=__CONTENTPLAYHEAD__&ca=__CACHEBUSTING__&a=__ASSETURI__&t=__TIMESTAMP__&o=__OTHER__]]></Error></VAST>`)
	got, err = xml.Marshal(v)
	if err != nil {
		t.Errorf("Marshal() error = %v", err)
		return
	}
	if !reflect.DeepEqual(got, want) {
		t.Logf("%s", got)
		t.Errorf("Marshal() got = %v, want %v", got, want)
	}

}

func createVastDemo() (*VAST, error) {
	adId := "123"
	adTitle := "ad title"
	assetId := "123456"
	impressionId := "456"
	impressionURI := "http://impression.track.cn"
	seconds := Duration(15 * time.Second)
	mediaType := "video/mp4"
	mediaURI := "http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4"

	v := &VAST{
		Version: "3.0",
		Ads: []Ad{
			{
				ID:   adId,
				Type: "front",
				InLine: &InLine{
					AdSystem: &AdSystem{Name: "DSP"},
					AdTitle:  PlainString{CDATA: adTitle},
					Impressions: []Impression{
						{ID: impressionId, URI: impressionURI},
					},
					Creatives: []Creative{
						{
							ID:       assetId,
							Sequence: 0,
							Linear: &Linear{
								Duration: seconds,
								TrackingEvents: []Tracking{
									{
										Event:  EventTypeStart,
										Offset: nil,
										URI:    "http://track.xxx.com/q/start?xx",
										UA:     "",
									},
								},
								MediaFiles: []MediaFile{
									{
										Delivery: "progressive",
										Type:     mediaType,
										Width:    1024,
										Height:   576,
										URI:      mediaURI,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	return v, nil
}

func BenchmarkVastMarshalXML(b *testing.B) {

	want := []byte(`<VAST version="3.0"><Ad id="123" type="front"><InLine><AdSystem><![CDATA[DSP]]></AdSystem><AdTitle><![CDATA[ad title]]></AdTitle><Impression id="456"><![CDATA[http://impression.track.cn]]></Impression><Creatives><Creative id="123456"><Linear><Duration>00:00:15</Duration><TrackingEvents><Tracking event="start"><![CDATA[http://track.xxx.com/q/start?xx]]></Tracking></TrackingEvents><MediaFiles><MediaFile delivery="progressive" type="video/mp4" width="1024" height="576"><![CDATA[http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4]]></MediaFile></MediaFiles></Linear></Creative></Creatives></InLine></Ad></VAST>`)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v, _ := createVastDemo()
		got, err := xml.Marshal(v)
		if err != nil {
			b.Errorf("Marshal() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, want) {
			b.Errorf("Marshal() got = %v, want %v", got, want)
		}
	}
}

func BenchmarkVastMarshalJson(b *testing.B) {

	want := []byte(`{"Version":"3.0","Ad":[{"ID":"123","Type":"front","InLine":{"AdSystem":{"Data":"DSP"},"AdTitle":{"Data":"ad title"},"Impressions":[{"ID":"456","URI":"http://impression.track.cn"}],"Creatives":[{"ID":"123456","Linear":{"Duration":"00:00:15","TrackingEvents":[{"Event":"start","URI":"http://track.xxx.com/q/start?xx"}],"MediaFiles":[{"Delivery":"progressive","Type":"video/mp4","Width":1024,"Height":576,"URI":"http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4"}]}}]}}]}`)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		v, _ := createVastDemo()
		got, err := ffjson.Marshal(v)
		if err != nil {
			b.Errorf("Marshal() error = %v", err)
			return
		}
		if !reflect.DeepEqual(got, want) {
			b.Errorf("Marshal() got = %v, want %v", got, want)
		}
	}
}

func TestCreateVastJson(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		{name: "testCase1", want: []byte(`{"Version":"3.0","Ad":[{"ID":"123","Type":"front","InLine":{"AdSystem":{"Data":"DSP"},"AdTitle":{"Data":"ad title"},"Impressions":[{"ID":"456","URI":"http://impression.track.cn"}],"Creatives":[{"ID":"123456","Linear":{"Duration":"00:00:15","TrackingEvents":[{"Event":"start","URI":"http://track.xxx.com/q/start?xx"}],"MediaFiles":[{"Delivery":"progressive","Type":"video/mp4","Width":1024,"Height":576,"URI":"http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4"}]}}]}}]}`),
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := createVastDemo()
			got, err := ffjson.Marshal(v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateVastXML(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		{name: "testCase1", want: []byte(`<VAST version="3.0"><Ad id="123" type="front"><InLine><AdSystem>DSP</AdSystem><AdTitle>ad title</AdTitle><Impression id="456"><![CDATA[http://impression.track.cn]]></Impression><Creatives><Creative id="123456"><Linear><Duration>00:00:15</Duration><TrackingEvents><Tracking event="start"><![CDATA[http://track.xxx.com/q/start?xx]]></Tracking></TrackingEvents><MediaFiles><MediaFile delivery="progressive" type="video/mp4" width="1024" height="576"><![CDATA[http://mp4.res.xxx.com/new_video/2020/01/14/1485/335928CBA9D02E95E63ED9F4D45DF6DF_20200114_1_1_1051.mp4]]></MediaFile></MediaFiles></Linear></Creative></Creatives></InLine></Ad></VAST>`),
			wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, _ := createVastDemo()
			got, err := xml.Marshal(v)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Marshal() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func loadFixture(path string) (*VAST, string, string, error) {
	xmlFile, err := os.Open(path)
	if err != nil {
		return nil, "", "", err
	}
	defer xmlFile.Close()
	b, _ := ioutil.ReadAll(xmlFile)

	var v VAST
	err = xml.Unmarshal(b, &v)

	res, err := xml.MarshalIndent(v, "", "  ")
	if err != nil {
		return nil, "", "", err

	}

	return &v, string(b), string(res), err
}


