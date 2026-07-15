package vast

import (
	"encoding/json"
	"encoding/xml"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func assertElementsInOrder(t *testing.T, document string, elements ...string) {
	t.Helper()
	position := -1
	for _, element := range elements {
		next := strings.Index(document[position+1:], element)
		require.NotEqualf(t, -1, next, "element %s is missing from %s", element, document)
		position += next + 1
	}
}

func xmlSection(t *testing.T, document, start, end string) string {
	t.Helper()
	startAt := strings.Index(document, start)
	require.NotEqual(t, -1, startAt)
	endAt := strings.Index(document[startAt:], end)
	require.NotEqual(t, -1, endAt)
	return document[startAt : startAt+endAt+len(end)]
}

func TestCreativeResourcesPreserveMultiplicity(t *testing.T) {
	t.Run("Companion", func(t *testing.T) {
		const source = `<Companion width="300" height="250">` +
			`<HTMLResource><![CDATA[<b>one</b>]]></HTMLResource>` +
			`<HTMLResource><![CDATA[<b>two</b>]]></HTMLResource>` +
			`<IFrameResource><![CDATA[https://example.com/frame/1]]></IFrameResource>` +
			`<IFrameResource><![CDATA[https://example.com/frame/2]]></IFrameResource>` +
			`<StaticResource creativeType="image/png"><![CDATA[https://example.com/static/1]]></StaticResource>` +
			`<StaticResource creativeType="image/webp"><![CDATA[https://example.com/static/2]]></StaticResource>` +
			`</Companion>`

		var companion Companion
		require.NoError(t, xml.Unmarshal([]byte(source), &companion))
		require.Len(t, companion.HTMLResources, 2)
		require.Len(t, companion.IFrameResources, 2)
		require.Len(t, companion.StaticResources, 2)
		require.NotNil(t, companion.HTMLResource)
		require.NotNil(t, companion.IFrameResource)
		require.NotNil(t, companion.StaticResource)
		assert.Equal(t, companion.HTMLResources[0], *companion.HTMLResource)
		assert.Equal(t, companion.IFrameResources[0], *companion.IFrameResource)
		assert.Equal(t, companion.StaticResources[0], *companion.StaticResource)

		output, err := xml.Marshal(companion)
		require.NoError(t, err)
		assert.Equal(t, 2, strings.Count(string(output), "<HTMLResource>"))
		assert.Equal(t, 2, strings.Count(string(output), "<IFrameResource>"))
		assert.Equal(t, 2, strings.Count(string(output), "<StaticResource "))
		assertElementsInOrder(t, string(output), "<HTMLResource>", "<IFrameResource>", "<StaticResource ")
	})

	t.Run("NonLinear", func(t *testing.T) {
		const source = `<NonLinear width="0" height="0">` +
			`<HTMLResource>one</HTMLResource><HTMLResource>two</HTMLResource>` +
			`<IFrameResource>https://example.com/1</IFrameResource><IFrameResource>https://example.com/2</IFrameResource>` +
			`<StaticResource creativeType="image/png">https://example.com/1.png</StaticResource>` +
			`<StaticResource creativeType="image/webp">https://example.com/2.webp</StaticResource>` +
			`</NonLinear>`

		var nonLinear NonLinear
		require.NoError(t, xml.Unmarshal([]byte(source), &nonLinear))
		require.Len(t, nonLinear.HTMLResources, 2)
		require.Len(t, nonLinear.IFrameResources, 2)
		require.Len(t, nonLinear.StaticResources, 2)
		output, err := xml.Marshal(nonLinear)
		require.NoError(t, err)
		assert.Equal(t, 2, strings.Count(string(output), "<HTMLResource>"))
		assert.Equal(t, 2, strings.Count(string(output), "<IFrameResource>"))
		assert.Equal(t, 2, strings.Count(string(output), "<StaticResource "))
	})

	t.Run("Icon", func(t *testing.T) {
		const source = `<Icon offset="00:00:00" duration="00:00:00">` +
			`<HTMLResource>one</HTMLResource><HTMLResource>two</HTMLResource>` +
			`<IFrameResource>https://example.com/1</IFrameResource><IFrameResource>https://example.com/2</IFrameResource>` +
			`<StaticResource creativeType="image/png">https://example.com/1.png</StaticResource>` +
			`<StaticResource creativeType="image/webp">https://example.com/2.webp</StaticResource>` +
			`<IconClicks><IconClickThrough>https://example.com/click</IconClickThrough></IconClicks>` +
			`<IconViewTracking>https://example.com/view</IconViewTracking>` +
			`</Icon>`

		var icon Icon
		require.NoError(t, xml.Unmarshal([]byte(source), &icon))
		require.Len(t, icon.HTMLResources, 2)
		require.Len(t, icon.IFrameResources, 2)
		require.Len(t, icon.StaticResources, 2)
		assert.True(t, icon.OffsetSet)
		assert.True(t, icon.DurationSet)

		output, err := xml.Marshal(icon)
		require.NoError(t, err)
		text := string(output)
		assert.Contains(t, text, `offset="00:00:00"`)
		assert.Contains(t, text, `duration="00:00:00"`)
		assert.Equal(t, 2, strings.Count(text, "<HTMLResource>"))
		assert.Equal(t, 2, strings.Count(text, "<IFrameResource>"))
		assert.Equal(t, 2, strings.Count(text, "<StaticResource "))
		assertElementsInOrder(t, text, "<HTMLResource>", "<IFrameResource>", "<StaticResource ", "<IconClicks>", "<IconViewTracking>")
	})
}

func TestExplicitFalseAttributesRoundTrip(t *testing.T) {
	const source = `<VAST version="4.3"><Ad conditionalAd="false"><InLine>` +
		`<AdSystem>test</AdSystem><Impression>https://example.com/impression</Impression>` +
		`<AdServingId>serving-id</AdServingId><AdTitle>title</AdTitle>` +
		`<AdVerifications><Verification vendor="example.com-omid">` +
		`<JavaScriptResource apiFramework="omid" browserOptional="false">https://example.com/omid.js</JavaScriptResource>` +
		`</Verification></AdVerifications><Creatives><Creative><Linear>` +
		`<AdParameters xmlEncoded="false">key=value</AdParameters><Duration>00:00:00</Duration>` +
		`<MediaFiles><MediaFile delivery="progressive" type="video/mp4" width="1" height="1" scalable="false" maintainAspectRatio="false">https://example.com/ad.mp4</MediaFile>` +
		`<InteractiveCreativeFile variableDuration="false">data:text/html,test</InteractiveCreativeFile>` +
		`</MediaFiles></Linear></Creative></Creatives></InLine></Ad></VAST>`

	var vast VAST
	require.NoError(t, xml.Unmarshal([]byte(source), &vast))
	ad := &vast.Ads[0]
	assert.True(t, ad.ConditionalAdSet)
	assert.False(t, ad.ConditionalAd)
	verification := ad.InLine.AdVerifications.Verification[0]
	assert.True(t, verification.JavaScriptResource[0].BrowserOptionalSet)
	media := ad.InLine.Creatives[0].Linear.MediaFiles
	assert.True(t, media.MediaFile[0].ScalableSet)
	assert.True(t, media.MediaFile[0].MaintainAspectRatioSet)
	assert.True(t, media.InteractiveCreativeFile[0].VariableDurationSet)
	assert.True(t, ad.InLine.Creatives[0].Linear.AdParameters.XMLEncodedSet)

	output, err := xml.Marshal(vast)
	require.NoError(t, err)
	text := string(output)
	assert.Contains(t, text, `conditionalAd="false"`)
	assert.Contains(t, text, `browserOptional="false"`)
	assert.Contains(t, text, `scalable="false"`)
	assert.Contains(t, text, `maintainAspectRatio="false"`)
	assert.Contains(t, text, `variableDuration="false"`)
	assert.Contains(t, text, `xmlEncoded="false"`)
	assert.Contains(t, text, `<Duration>00:00:00</Duration>`)

	jsonData, err := json.Marshal(vast)
	require.NoError(t, err)
	var fromJSON VAST
	require.NoError(t, json.Unmarshal(jsonData, &fromJSON))
	assert.Equal(t, vast, fromJSON)
}

func TestIconOmittedOptionalFieldsStayOmitted(t *testing.T) {
	output, err := xml.Marshal(Icon{})
	require.NoError(t, err)
	assert.Equal(t, `<Icon></Icon>`, string(output))
}

func TestLegacySingularResourcesStillMarshal(t *testing.T) {
	static := &StaticResource{CreativeType: "image/png", URI: "https://example.com/image.png"}
	iframe := &CDATAString{CDATA: "https://example.com/frame"}
	html := &HTMLResource{HTML: "<b>creative</b>"}

	tests := []struct {
		name  string
		value any
	}{
		{name: "Companion", value: Companion{Width: 1, Height: 1, StaticResource: static, IFrameResource: iframe, HTMLResource: html}},
		{name: "NonLinear", value: NonLinear{Width: 1, Height: 1, StaticResource: static, IFrameResource: iframe, HTMLResource: html}},
		{name: "Icon", value: Icon{StaticResource: static, IFrameResource: iframe, HTMLResource: html}},
		{name: "NonLinearWrapper", value: NonLinearWrapper{NonLinearClickTracking: []CDATAString{{CDATA: "https://example.com/click"}}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, err := xml.Marshal(test.value)
			require.NoError(t, err)
			text := string(output)
			if test.name == "NonLinearWrapper" {
				assert.Contains(t, text, "<NonLinearClickTracking>")
				return
			}
			assert.Equal(t, 1, strings.Count(text, "<HTMLResource>"))
			assert.Equal(t, 1, strings.Count(text, "<IFrameResource>"))
			assert.Equal(t, 1, strings.Count(text, "<StaticResource "))
		})
	}
}

func TestParsedSingularResourceAliasesRemainWritable(t *testing.T) {
	tests := []struct {
		name      string
		source    string
		newTarget func() any
		replace   func(any)
		clear     func(any)
	}{
		{
			name: "Companion",
			source: `<Companion width="300" height="250"><HTMLResource>html-one</HTMLResource><HTMLResource>html-two</HTMLResource>` +
				`<IFrameResource>https://example.com/frame-one</IFrameResource><IFrameResource>https://example.com/frame-two</IFrameResource>` +
				`<StaticResource creativeType="image/png">https://example.com/static-one</StaticResource><StaticResource creativeType="image/png">https://example.com/static-two</StaticResource></Companion>`,
			newTarget: func() any { return &Companion{} },
			replace: func(target any) {
				value := target.(*Companion)
				value.HTMLResource = &HTMLResource{HTML: "html-replacement"}
				value.IFrameResource = &CDATAString{CDATA: "https://example.com/frame-replacement"}
				value.StaticResource = &StaticResource{CreativeType: "image/webp", URI: "https://example.com/static-replacement"}
			},
			clear: func(target any) {
				value := target.(*Companion)
				value.HTMLResource = nil
				value.IFrameResource = nil
				value.StaticResource = nil
			},
		},
		{
			name: "NonLinear",
			source: `<NonLinear width="300" height="250"><HTMLResource>html-one</HTMLResource><HTMLResource>html-two</HTMLResource>` +
				`<IFrameResource>https://example.com/frame-one</IFrameResource><IFrameResource>https://example.com/frame-two</IFrameResource>` +
				`<StaticResource creativeType="image/png">https://example.com/static-one</StaticResource><StaticResource creativeType="image/png">https://example.com/static-two</StaticResource></NonLinear>`,
			newTarget: func() any { return &NonLinear{} },
			replace: func(target any) {
				value := target.(*NonLinear)
				value.HTMLResource = &HTMLResource{HTML: "html-replacement"}
				value.IFrameResource = &CDATAString{CDATA: "https://example.com/frame-replacement"}
				value.StaticResource = &StaticResource{CreativeType: "image/webp", URI: "https://example.com/static-replacement"}
			},
			clear: func(target any) {
				value := target.(*NonLinear)
				value.HTMLResource = nil
				value.IFrameResource = nil
				value.StaticResource = nil
			},
		},
		{
			name: "Icon",
			source: `<Icon><HTMLResource>html-one</HTMLResource><HTMLResource>html-two</HTMLResource>` +
				`<IFrameResource>https://example.com/frame-one</IFrameResource><IFrameResource>https://example.com/frame-two</IFrameResource>` +
				`<StaticResource creativeType="image/png">https://example.com/static-one</StaticResource><StaticResource creativeType="image/png">https://example.com/static-two</StaticResource></Icon>`,
			newTarget: func() any { return &Icon{} },
			replace: func(target any) {
				value := target.(*Icon)
				value.HTMLResource = &HTMLResource{HTML: "html-replacement"}
				value.IFrameResource = &CDATAString{CDATA: "https://example.com/frame-replacement"}
				value.StaticResource = &StaticResource{CreativeType: "image/webp", URI: "https://example.com/static-replacement"}
			},
			clear: func(target any) {
				value := target.(*Icon)
				value.HTMLResource = nil
				value.IFrameResource = nil
				value.StaticResource = nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name+"/replace", func(t *testing.T) {
			target := test.newTarget()
			require.NoError(t, xml.Unmarshal([]byte(test.source), target))
			test.replace(target)
			output, err := xml.Marshal(target)
			require.NoError(t, err)
			text := string(output)
			assert.NotContains(t, text, "html-one")
			assert.NotContains(t, text, "frame-one")
			assert.NotContains(t, text, "static-one")
			assert.Contains(t, text, "html-replacement")
			assert.Contains(t, text, "frame-replacement")
			assert.Contains(t, text, "static-replacement")
			assert.Contains(t, text, "html-two")
			assert.Contains(t, text, "frame-two")
			assert.Contains(t, text, "static-two")
		})

		t.Run(test.name+"/clear", func(t *testing.T) {
			target := test.newTarget()
			require.NoError(t, xml.Unmarshal([]byte(test.source), target))
			test.clear(target)
			output, err := xml.Marshal(target)
			require.NoError(t, err)
			text := string(output)
			assert.NotContains(t, text, "html-one")
			assert.NotContains(t, text, "frame-one")
			assert.NotContains(t, text, "static-one")
			assert.Contains(t, text, "html-two")
			assert.Contains(t, text, "frame-two")
			assert.Contains(t, text, "static-two")
		})
	}
}

func TestParsedPluralResourceEditsRemainAuthoritative(t *testing.T) {
	const source = `<Companion width="300" height="250"><HTMLResource>original</HTMLResource></Companion>`
	t.Run("replace", func(t *testing.T) {
		var companion Companion
		require.NoError(t, xml.Unmarshal([]byte(source), &companion))
		companion.HTMLResources = []HTMLResource{{HTML: "plural-replacement"}}

		output, err := xml.Marshal(companion)
		require.NoError(t, err)
		assert.Contains(t, string(output), "plural-replacement")
		assert.NotContains(t, string(output), "original")
	})

	t.Run("clear", func(t *testing.T) {
		var companion Companion
		require.NoError(t, xml.Unmarshal([]byte(source), &companion))
		companion.HTMLResources = nil

		output, err := xml.Marshal(companion)
		require.NoError(t, err)
		assert.NotContains(t, string(output), "<HTMLResource")
	})
}

func TestOptionalFalseAttributesOnCreativeTypes(t *testing.T) {
	t.Run("HTMLResource", func(t *testing.T) {
		const source = `<HTMLResource xmlEncoded="false">content</HTMLResource>`
		var resource HTMLResource
		require.NoError(t, xml.Unmarshal([]byte(source), &resource))
		assert.True(t, resource.XMLEncodedSet)
		output, err := xml.Marshal(resource)
		require.NoError(t, err)
		assert.Contains(t, string(output), `xmlEncoded="false"`)
	})

	t.Run("NonLinear", func(t *testing.T) {
		const source = `<NonLinear width="1" height="1" scalable="false" maintainAspectRatio="false"></NonLinear>`
		var nonLinear NonLinear
		require.NoError(t, xml.Unmarshal([]byte(source), &nonLinear))
		assert.True(t, nonLinear.ScalableSet)
		assert.True(t, nonLinear.MaintainAspectRatioSet)
		output, err := xml.Marshal(nonLinear)
		require.NoError(t, err)
		assert.Contains(t, string(output), `scalable="false"`)
		assert.Contains(t, string(output), `maintainAspectRatio="false"`)
	})

	t.Run("NonLinearWrapper", func(t *testing.T) {
		const source = `<NonLinear scalable="false" maintainAspectRatio="false"></NonLinear>`
		var nonLinear NonLinearWrapper
		require.NoError(t, xml.Unmarshal([]byte(source), &nonLinear))
		assert.True(t, nonLinear.ScalableSet)
		assert.True(t, nonLinear.MaintainAspectRatioSet)
		output, err := xml.Marshal(nonLinear)
		require.NoError(t, err)
		assert.Contains(t, string(output), `scalable="false"`)
		assert.Contains(t, string(output), `maintainAspectRatio="false"`)
	})
}

func TestViewableImpressionWithoutIDOmitAttribute(t *testing.T) {
	output, err := xml.Marshal(ViewableImpression{Viewable: []CDATAString{{CDATA: "https://example.com/view"}}})
	require.NoError(t, err)
	assert.NotContains(t, string(output), `id=""`)
}

func TestXMLCodecRejectsMalformedElements(t *testing.T) {
	tests := []struct {
		name   string
		target any
		xml    string
	}{
		{name: "VAST", target: &VAST{}, xml: `<VAST><broken></VAST>`},
		{name: "Ad", target: &Ad{}, xml: `<Ad><broken></Ad>`},
		{name: "JavaScriptResource", target: &JavaScriptResource{}, xml: `<JavaScriptResource><broken></JavaScriptResource>`},
		{name: "InteractiveCreativeFile", target: &InteractiveCreativeFile{}, xml: `<InteractiveCreativeFile><broken></InteractiveCreativeFile>`},
		{name: "HTMLResource", target: &HTMLResource{}, xml: `<HTMLResource><broken></HTMLResource>`},
		{name: "AdParameters", target: &AdParameters{}, xml: `<AdParameters><broken></AdParameters>`},
		{name: "MediaFile", target: &MediaFile{}, xml: `<MediaFile><broken></MediaFile>`},
		{name: "Companion", target: &Companion{}, xml: `<Companion><broken></Companion>`},
		{name: "NonLinear", target: &NonLinear{}, xml: `<NonLinear><broken></NonLinear>`},
		{name: "NonLinearWrapper", target: &NonLinearWrapper{}, xml: `<NonLinear><broken></NonLinear>`},
		{name: "Icon", target: &Icon{}, xml: `<Icon><broken></Icon>`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.Error(t, xml.Unmarshal([]byte(test.xml), test.target))
		})
	}
}

func TestIconRejectsInvalidOffset(t *testing.T) {
	negative := Duration(-1)
	_, err := xml.Marshal(Icon{Offset: Offset{Duration: &negative}, OffsetSet: true})
	assert.EqualError(t, err, "invalid duration: -1ns")

	var icon Icon
	assert.EqualError(t, xml.Unmarshal([]byte(`<Icon offset="101%"></Icon>`), &icon), "invalid offset: 101%")
}

func TestWrapperNonLinearClickTrackingPreservesID(t *testing.T) {
	const source = `<NonLinear id="n"><NonLinearClickTracking id="track-id">https://example.com/click</NonLinearClickTracking></NonLinear>`
	var nonLinear NonLinearWrapper
	require.NoError(t, xml.Unmarshal([]byte(source), &nonLinear))
	require.Len(t, nonLinear.NonLinearClickTrackings, 1)
	assert.Equal(t, "track-id", nonLinear.NonLinearClickTrackings[0].ID)
	require.Len(t, nonLinear.NonLinearClickTracking, 1)
	assert.Equal(t, "https://example.com/click", nonLinear.NonLinearClickTracking[0].CDATA)

	output, err := xml.Marshal(nonLinear)
	require.NoError(t, err)
	assert.Contains(t, string(output), `id="track-id"`)
}

func TestParsedWrapperNonLinearClickTrackingRemainsWritable(t *testing.T) {
	const source = `<NonLinear><NonLinearClickTracking id="one">https://example.com/one</NonLinearClickTracking><NonLinearClickTracking id="two">https://example.com/two</NonLinearClickTracking></NonLinear>`

	t.Run("replace", func(t *testing.T) {
		var nonLinear NonLinearWrapper
		require.NoError(t, xml.Unmarshal([]byte(source), &nonLinear))
		nonLinear.NonLinearClickTracking = []CDATAString{{CDATA: "https://example.com/replacement"}}

		output, err := xml.Marshal(nonLinear)
		require.NoError(t, err)
		text := string(output)
		assert.Equal(t, 1, strings.Count(text, "<NonLinearClickTracking"))
		assert.Contains(t, text, "https://example.com/replacement")
		assert.NotContains(t, text, "https://example.com/one")
		assert.NotContains(t, text, `id="one"`)
	})

	t.Run("clear", func(t *testing.T) {
		var nonLinear NonLinearWrapper
		require.NoError(t, xml.Unmarshal([]byte(source), &nonLinear))
		nonLinear.NonLinearClickTracking = nil

		output, err := xml.Marshal(nonLinear)
		require.NoError(t, err)
		assert.NotContains(t, string(output), "<NonLinearClickTracking")
	})
}

func TestParsedWrapperCanonicalClickTrackingsRemainAuthoritative(t *testing.T) {
	const source = `<NonLinear><NonLinearClickTracking id="one">https://example.com/one</NonLinearClickTracking></NonLinear>`
	var nonLinear NonLinearWrapper
	require.NoError(t, xml.Unmarshal([]byte(source), &nonLinear))
	nonLinear.NonLinearClickTrackings = nil

	output, err := xml.Marshal(nonLinear)
	require.NoError(t, err)
	assert.NotContains(t, string(output), "<NonLinearClickTracking")
}

func TestVAST43SchemaElementOrder(t *testing.T) {
	vast, _, _, err := loadFixture("testdata/vast_43_inline.xml")
	require.NoError(t, err)
	output, err := xml.Marshal(vast)
	require.NoError(t, err)
	document := string(output)

	inline := xmlSection(t, document, "<InLine>", "</InLine>")
	assertElementsInOrder(t, inline,
		"<AdSystem", "<Error>", "<Extensions>", "<Impression", "<Pricing", "<ViewableImpression",
		"<AdServingId>", "<AdTitle>", "<AdVerifications>", "<Advertiser", "<Category", "<Creatives>",
		"<Description>", "<Expires>", "<Survey")

	linear := xmlSection(t, inline, "<Linear", "</Linear>")
	assertElementsInOrder(t, linear, "<Icons>", "<TrackingEvents>", "<AdParameters", "<Duration>", "<MediaFiles>", "<VideoClicks>")
	mediaFiles := xmlSection(t, linear, "<MediaFiles>", "</MediaFiles>")
	assertElementsInOrder(t, mediaFiles, "<ClosedCaptionFiles>", "<MediaFile ", "<Mezzanine ", "<InteractiveCreativeFile ")
	verification := xmlSection(t, inline, "<Verification ", "</Verification>")
	assertElementsInOrder(t, verification, "<ExecutableResource ", "<JavaScriptResource ", "<TrackingEvents>", "<VerificationParameters>")
	creative := xmlSection(t, inline, "<Creative ", "</Creative>")
	assertElementsInOrder(t, creative, "<Linear", "<UniversalAdId ")
	icon := xmlSection(t, linear, "<Icon ", "</Icon>")
	assertElementsInOrder(t, icon, "<StaticResource ", "<IconClickFallbackImages>", "<IconClickThrough>", "<IconClickTracking ", "<IconViewTracking>")
	videoClicks := xmlSection(t, linear, "<VideoClicks>", "</VideoClicks>")
	assertElementsInOrder(t, videoClicks, "<ClickTracking ", "<ClickThrough ", "<CustomClick ")

	wrapperVAST, _, _, err := loadFixture("testdata/vast_43_wrapper.xml")
	require.NoError(t, err)
	wrapperOutput, err := xml.Marshal(wrapperVAST)
	require.NoError(t, err)
	wrapper := xmlSection(t, string(wrapperOutput), "<Wrapper ", "</Wrapper>")
	assertElementsInOrder(t, wrapper,
		"<AdSystem", "<Error>", "<Extensions>", "<Impression", "<Pricing", "<ViewableImpression",
		"<AdVerifications>", "<BlockedAdCategories", "<Creatives>", "<VASTAdTagURI>")
}

func TestExtensionPreservesCustomAttributes(t *testing.T) {
	const source = `<CreativeExtension type="application/xml" vendor="example" mode="strict"><Node>value</Node></CreativeExtension>`
	var extension Extension
	require.NoError(t, xml.Unmarshal([]byte(source), &extension))
	require.Len(t, extension.Attributes, 2)
	assert.Equal(t, "vendor", extension.Attributes[0].Name.Local)
	assert.Equal(t, "example", extension.Attributes[0].Value)

	output, err := xml.Marshal(extension)
	require.NoError(t, err)
	assert.Contains(t, string(output), `vendor="example"`)
	assert.Contains(t, string(output), `mode="strict"`)
	assert.Contains(t, string(output), `<Node>value</Node>`)
}

func TestVASTPreservesRootAttributes(t *testing.T) {
	const source = `<VAST xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:noNamespaceSchemaLocation="vast.xsd" version="4.3"></VAST>`
	var vast VAST
	require.NoError(t, xml.Unmarshal([]byte(source), &vast))
	require.Len(t, vast.Attributes, 2)
	assert.Equal(t, "xmlns", vast.Attributes[0].Name.Space)
	assert.Equal(t, "xsi", vast.Attributes[0].Name.Local)
	assert.Equal(t, "http://www.w3.org/2001/XMLSchema-instance", vast.Attributes[1].Name.Space)
	assert.Equal(t, "noNamespaceSchemaLocation", vast.Attributes[1].Name.Local)

	output, err := xml.Marshal(vast)
	require.NoError(t, err)
	text := string(output)
	assert.Equal(t, source, text)
	assert.Contains(t, text, "http://www.w3.org/2001/XMLSchema-instance")
	assert.Contains(t, text, "noNamespaceSchemaLocation")
	assert.Contains(t, text, "vast.xsd")
}

func TestExtensionPreservesNamespacePrefixes(t *testing.T) {
	const source = `<Extension xmlns:vendor="urn:vendor" vendor:mode="strict"><vendor:Node>value</vendor:Node></Extension>`
	var extension Extension
	require.NoError(t, xml.Unmarshal([]byte(source), &extension))
	require.Len(t, extension.Attributes, 2)

	output, err := xml.Marshal(extension)
	require.NoError(t, err)
	assert.Equal(t, source, string(output))
}
