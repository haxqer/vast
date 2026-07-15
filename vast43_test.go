package vast

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"aqwari.net/xml/xmltree"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This file contains a compliance test-suite for VAST 4.3.
//
// VAST 4.3 is a superset of VAST 4.0/4.1/4.2. Its only schema-level change
// relative to 4.2 is that <InteractiveCreativeFile> may carry an inline data:
// URI (spec section 3.9.3). These tests therefore verify two things:
//
//  1. The 4.3-specific behaviour (inline data URIs, string `type` attributes on
//     interactive/executable resources, verification parameters).
//  2. That the struct model round-trips the *entire* VAST 4.3 element set from
//     the human-readable schema in spec section 5, for both the InLine and the
//     Wrapper branches.
//
// The authoritative element/attribute list is taken from the VAST 4.3 spec:
// https://github.com/InteractiveAdvertisingBureau/VAST4.x/blob/main/4.3.md

// assertXMLEquivalent parses both documents and asserts they are structurally
// equal (ignoring element order, whitespace and namespace prefixes).
func assertXMLEquivalent(t *testing.T, expected, actual []byte) {
	t.Helper()
	exp, err := xmltree.Parse(expected)
	require.NoError(t, err)
	act, err := xmltree.Parse(actual)
	require.NoError(t, err)
	assert.True(t, xmltree.Equal(act, exp), "XML documents are not structurally equivalent")
}

// assertRoundTrips marshals the parsed VAST back to XML and asserts the result
// is equivalent to the source document, i.e. no element was dropped by the
// struct model.
func assertRoundTrips(t *testing.T, v *VAST, source []byte) {
	t.Helper()
	out, err := xml.Marshal(v)
	require.NoError(t, err)
	assertXMLEquivalent(t, source, out)
}

// TestVAST43Version verifies a version="4.3" document is recognised.
func TestVAST43Version(t *testing.T) {
	v, _, _, err := loadFixture("testdata/vast_43_inline.xml")
	require.NoError(t, err)
	assert.Equal(t, "4.3", v.Version)
	assert.Equal(t, "http://www.iab.com/VAST", v.XMLNS)
}

// TestVAST43InlineComplete parses the comprehensive InLine fixture and asserts
// every element/attribute of the InLine branch of the section-5 schema is
// mapped correctly, then verifies the document round-trips without loss.
func TestVAST43InlineComplete(t *testing.T) {
	v, source, _, err := loadFixture("testdata/vast_43_inline.xml")
	require.NoError(t, err)

	require.Equal(t, "4.3", v.Version)
	require.Len(t, v.Ads, 1)

	ad := v.Ads[0]
	assert.Equal(t, "20001", ad.ID)
	assert.Equal(t, 1, ad.Sequence)
	assert.Equal(t, "video", ad.AdType)

	inline := ad.InLine
	require.NotNil(t, inline)

	// --- InLine metadata ------------------------------------------------
	require.NotNil(t, inline.AdSystem)
	assert.Equal(t, "4.3", inline.AdSystem.Version)
	assert.Equal(t, "iabtechlab", inline.AdSystem.Name)
	assert.Equal(t, "VAST 4.3 Complete InLine Example", inline.AdTitle.CDATA)
	assert.Equal(t, "a532d16d-4d7f-4440-bd29-2ec0e693fc80", inline.AdServingId)

	require.Len(t, inline.Impressions, 1)
	assert.Equal(t, "Impression-ID", inline.Impressions[0].ID)
	assert.Equal(t, "https://example.com/track/impression", inline.Impressions[0].URI)

	require.NotNil(t, inline.Category)
	require.Len(t, *inline.Category, 1)
	assert.Equal(t, "https://www.iabtechlab.com/categoryauthority", (*inline.Category)[0].Authority)
	assert.Equal(t, "AD CONTENT description category", (*inline.Category)[0].Category)

	require.NotNil(t, inline.Description)
	assert.Contains(t, inline.Description.CDATA, "full InLine schema")

	require.NotNil(t, inline.Advertiser)
	assert.Equal(t, "advertiser-id-123", inline.Advertiser.ID)
	assert.Equal(t, "example.com", inline.Advertiser.Advertiser)

	require.NotNil(t, inline.Pricing)
	assert.Equal(t, "cpm", inline.Pricing.Model)
	assert.Equal(t, "USD", inline.Pricing.Currency)
	assert.Equal(t, "25.00", inline.Pricing.Value)

	require.NotNil(t, inline.Survey)
	assert.Equal(t, "text/javascript", inline.Survey.Type)
	assert.Equal(t, "https://example.com/survey", inline.Survey.URI)

	require.Len(t, inline.Errors, 1)
	assert.Equal(t, "https://example.com/track/error/[ERRORCODE]", inline.Errors[0].CDATA)

	require.NotNil(t, inline.Expires)
	assert.Equal(t, 3600, *inline.Expires)

	// --- ViewableImpression ---------------------------------------------
	require.NotNil(t, inline.ViewableImpression)
	assert.Equal(t, "viewable-impression-id", inline.ViewableImpression.ID)
	require.Len(t, inline.ViewableImpression.Viewable, 1)
	assert.Equal(t, "https://example.com/track/viewable", inline.ViewableImpression.Viewable[0].CDATA)
	require.Len(t, inline.ViewableImpression.NotViewable, 1)
	assert.Equal(t, "https://example.com/track/notViewable", inline.ViewableImpression.NotViewable[0].CDATA)
	require.Len(t, inline.ViewableImpression.ViewUndetermined, 1)
	assert.Equal(t, "https://example.com/track/viewUndetermined", inline.ViewableImpression.ViewUndetermined[0].CDATA)

	// --- AdVerifications ------------------------------------------------
	require.NotNil(t, inline.AdVerifications)
	require.Len(t, inline.AdVerifications.Verification, 1)
	ver := inline.AdVerifications.Verification[0]
	assert.Equal(t, "company.com-omid", ver.Vendor)
	require.Len(t, ver.JavaScriptResource, 1)
	assert.Equal(t, "omid", ver.JavaScriptResource[0].ApiFramework)
	assert.True(t, ver.JavaScriptResource[0].BrowserOptional)
	assert.Equal(t, "https://verificationcompany1.com/verification_script1.js", ver.JavaScriptResource[0].URI)
	require.Len(t, ver.ExecutableResource, 1)
	assert.Equal(t, "omid", ver.ExecutableResource[0].ApiFramework)
	assert.Equal(t, "no-op", ver.ExecutableResource[0].Type)
	require.NotNil(t, ver.TrackingEvents)
	require.Len(t, ver.TrackingEvents.Tracking, 1)
	assert.Equal(t, EventTypeVerificationNotExecuted, ver.TrackingEvents.Tracking[0].Event)
	require.NotNil(t, ver.VerificationParameters)
	assert.Equal(t, `{"key":"value"}`, ver.VerificationParameters.URI)

	// --- Extensions -----------------------------------------------------
	require.NotNil(t, inline.Extensions)
	require.Len(t, *inline.Extensions, 1)
	ext := (*inline.Extensions)[0]
	assert.Equal(t, "iab-Count", ext.Type)
	require.Len(t, ext.CustomTracking, 1)
	assert.Equal(t, EventTypeOtherAdInteraction, ext.CustomTracking[0].Event)

	// --- Creatives ------------------------------------------------------
	require.Len(t, inline.Creatives, 3)

	// Creative #1: Linear.
	linCreative := inline.Creatives[0]
	assert.Equal(t, "5480", linCreative.ID)
	assert.Equal(t, 1, linCreative.Sequence)
	assert.Equal(t, "2447226", linCreative.AdID)
	assert.Equal(t, "omid", linCreative.APIFramework)
	require.NotNil(t, linCreative.UniversalAdID)
	require.Len(t, *linCreative.UniversalAdID, 2)
	assert.Equal(t, "Ad-ID", (*linCreative.UniversalAdID)[0].IDRegistry)
	assert.Equal(t, "8465", (*linCreative.UniversalAdID)[0].ID)
	assert.Equal(t, "another-registry", (*linCreative.UniversalAdID)[1].IDRegistry)

	lin := linCreative.Linear
	require.NotNil(t, lin)
	require.NotNil(t, lin.SkipOffset)
	require.NotNil(t, lin.SkipOffset.Duration)
	assert.Equal(t, Duration(5e9), *lin.SkipOffset.Duration) // 00:00:05
	assert.Equal(t, Duration(16e9), lin.Duration)            // 00:00:16
	require.NotNil(t, lin.AdParameters)
	assert.True(t, lin.AdParameters.XMLEncoded)
	assert.Equal(t, "params=1", lin.AdParameters.Parameters)

	// MediaFiles: Mezzanine + MediaFile x2 + InteractiveCreativeFile + ClosedCaptionFiles.
	mf := lin.MediaFiles
	require.NotNil(t, mf)
	require.Len(t, mf.Mezzanine, 1)
	mez := mf.Mezzanine[0]
	assert.Equal(t, "progressive", mez.Delivery)
	assert.Equal(t, "video/mp4", mez.Type)
	assert.Equal(t, 1280, mez.Width)
	assert.Equal(t, 720, mez.Height)
	assert.Equal(t, "h264", mez.Codec)
	assert.Equal(t, "mezzanine-1", mez.ID)
	assert.Equal(t, 1048576, mez.FileSize)
	assert.Equal(t, "2D", mez.MediaType)
	assert.Equal(t, "https://example.com/media/mezzanine.mp4", mez.URI)

	require.Len(t, mf.MediaFile, 2)
	assert.Equal(t, 2000, mf.MediaFile[0].Bitrate)
	assert.True(t, mf.MediaFile[0].Scalable)
	assert.True(t, mf.MediaFile[0].MaintainAspectRatio)
	assert.Equal(t, 700, mf.MediaFile[1].MinBitrate)
	assert.Equal(t, 1500, mf.MediaFile[1].MaxBitrate)
	assert.Equal(t, "2D", mf.MediaFile[1].MediaType)
	assert.Equal(t, 524288, mf.MediaFile[1].FileSize)

	require.Len(t, mf.InteractiveCreativeFile, 1)
	icf := mf.InteractiveCreativeFile[0]
	assert.Equal(t, "SIMID", icf.ApiFramework)
	assert.Equal(t, "text/html", icf.Type)
	assert.True(t, icf.VariableDuration)
	assert.True(t, strings.HasPrefix(icf.URI, "data:text/html,"), "expected inline data URI, got %q", icf.URI)

	require.NotNil(t, mf.ClosedCaptionFiles)
	require.Len(t, *mf.ClosedCaptionFiles, 2)
	assert.Equal(t, "text/vtt", (*mf.ClosedCaptionFiles)[0].Type)
	assert.Equal(t, "en", (*mf.ClosedCaptionFiles)[0].Language)
	assert.Equal(t, "zh-TW", (*mf.ClosedCaptionFiles)[1].Language)

	// VideoClicks.
	require.NotNil(t, lin.VideoClicks)
	require.Len(t, lin.VideoClicks.ClickThroughs, 1)
	assert.Equal(t, "click-through", lin.VideoClicks.ClickThroughs[0].ID)
	require.Len(t, lin.VideoClicks.ClickTrackings, 1)
	require.Len(t, lin.VideoClicks.CustomClicks, 1)

	// Linear TrackingEvents including a progress event with an offset.
	require.NotNil(t, lin.TrackingEvents)
	require.Len(t, lin.TrackingEvents.Tracking, 7)
	var progress *Tracking
	for i := range lin.TrackingEvents.Tracking {
		if lin.TrackingEvents.Tracking[i].Event == EventTypeProgress {
			progress = &lin.TrackingEvents.Tracking[i]
		}
	}
	require.NotNil(t, progress)
	require.NotNil(t, progress.Offset)
	require.NotNil(t, progress.Offset.Duration)
	assert.Equal(t, Duration(10e9), *progress.Offset.Duration) // 00:00:10

	// Icons.
	require.NotNil(t, lin.Icons)
	require.NotNil(t, lin.Icons.Icon)
	require.Len(t, *lin.Icons.Icon, 1)
	icon := (*lin.Icons.Icon)[0]
	assert.Equal(t, "AdChoices", icon.Program)
	assert.Equal(t, 60, icon.Width)
	assert.Equal(t, 20, icon.Height)
	assert.Equal(t, "right", icon.XPosition)
	assert.Equal(t, "top", icon.YPosition)
	assert.Equal(t, "omid", icon.APIFramework)
	assert.Equal(t, "1", icon.Pxratio)
	require.NotNil(t, icon.StaticResource)
	assert.Equal(t, "image/png", icon.StaticResource.CreativeType)
	require.Len(t, icon.IconViewTracking, 1)
	require.NotNil(t, icon.IconClickThrough)
	require.Len(t, icon.IconClickTracking, 1)
	assert.Equal(t, "icon-click-tracking", icon.IconClickTracking[0].ID)
	assert.Equal(t, "https://example.com/track/iconClick", icon.IconClickTracking[0].URI)
	require.NotNil(t, icon.IconClickFallbackImages)
	require.Len(t, icon.IconClickFallbackImages.IconClickFallbackImage, 1)
	fallback := icon.IconClickFallbackImages.IconClickFallbackImage[0]
	assert.Equal(t, 400, fallback.Width)
	assert.Equal(t, 150, fallback.Height)
	assert.Equal(t, "Icon fallback alt text", fallback.AltText)
	require.NotNil(t, fallback.StaticResource)
	assert.Empty(t, fallback.StaticResource.CreativeType)
	assert.Equal(t, "https://example.com/icon/fallback.png", fallback.StaticResource.URI)

	// Creative #2: CompanionAds.
	compCreative := inline.Creatives[1]
	require.NotNil(t, compCreative.UniversalAdID)
	require.Len(t, *compCreative.UniversalAdID, 1)
	require.NotNil(t, compCreative.CompanionAds)
	assert.Equal(t, "all", compCreative.CompanionAds.Required)
	require.Len(t, compCreative.CompanionAds.Companions, 1)
	comp := compCreative.CompanionAds.Companions[0]
	assert.Equal(t, "companion-1", comp.ID)
	assert.Equal(t, 300, comp.Width)
	assert.Equal(t, 250, comp.Height)
	assert.Equal(t, 300, comp.AssetWidth)
	assert.Equal(t, 250, comp.AssetHeight)
	assert.Equal(t, 600, comp.ExpandedWidth)
	assert.Equal(t, 500, comp.ExpandedHeight)
	assert.Equal(t, "slot-1", comp.AdSlotID)
	assert.Equal(t, "1", comp.Pxratio)
	assert.Equal(t, "default", comp.RenderingMode)
	require.NotNil(t, comp.StaticResource)
	require.NotNil(t, comp.CompanionClickThrough)
	require.Len(t, comp.CompanionClickTrackings, 1)
	assert.Equal(t, "companion-click-tracking", comp.CompanionClickTrackings[0].ID)
	assert.Equal(t, "Companion alt text", comp.AltText)
	require.NotNil(t, comp.AdParameters)
	require.NotNil(t, comp.TrackingEvents)
	require.Len(t, comp.TrackingEvents.Tracking, 1)
	assert.Equal(t, EventTypeCreativeView, comp.TrackingEvents.Tracking[0].Event)

	// Creative #3: NonLinearAds.
	nlCreative := inline.Creatives[2]
	require.NotNil(t, nlCreative.UniversalAdID)
	require.Len(t, *nlCreative.UniversalAdID, 1)
	require.NotNil(t, nlCreative.NonLinearAds)
	require.Len(t, nlCreative.NonLinearAds.NonLinears, 1)
	nl := nlCreative.NonLinearAds.NonLinears[0]
	assert.Equal(t, "nonlinear-1", nl.ID)
	assert.Equal(t, 480, nl.Width)
	assert.Equal(t, 70, nl.Height)
	assert.True(t, nl.Scalable)
	assert.True(t, nl.MaintainAspectRatio)
	require.NotNil(t, nl.MinSuggestedDuration)
	assert.Equal(t, Duration(5e9), *nl.MinSuggestedDuration)
	require.NotNil(t, nl.StaticResource)
	require.NotNil(t, nl.NonLinearClickThrough)
	require.Len(t, nl.NonLinearClickTrackings, 1)
	assert.Equal(t, "nonlinear-click-tracking", nl.NonLinearClickTrackings[0].ID)
	require.NotNil(t, nlCreative.NonLinearAds.TrackingEvents)

	// --- Round-trip -----------------------------------------------------
	assertRoundTrips(t, v, source)
}

// TestVAST43WrapperComplete parses the comprehensive Wrapper fixture, asserts
// the Wrapper-branch elements/attributes, and verifies round-trip fidelity.
func TestVAST43WrapperComplete(t *testing.T) {
	v, source, _, err := loadFixture("testdata/vast_43_wrapper.xml")
	require.NoError(t, err)

	require.Equal(t, "4.3", v.Version)
	require.Len(t, v.Ads, 1)
	ad := v.Ads[0]
	assert.Equal(t, "20002", ad.ID)
	assert.Equal(t, "video", ad.AdType)

	w := ad.Wrapper
	require.NotNil(t, w)

	// Wrapper attributes.
	require.NotNil(t, w.FollowAdditionalWrappers)
	assert.False(t, *w.FollowAdditionalWrappers)
	require.NotNil(t, w.AllowMultipleAds)
	assert.True(t, *w.AllowMultipleAds)
	require.NotNil(t, w.FallbackOnNoAd)
	assert.False(t, *w.FallbackOnNoAd)

	require.NotNil(t, w.AdSystem)
	assert.Equal(t, "iabtechlab", w.AdSystem.Name)
	require.Len(t, w.Impressions, 1)
	assert.Equal(t, "https://example.com/vast/secondary.xml", w.VASTAdTagURI.CDATA)
	require.Len(t, w.Errors, 1)
	require.NotNil(t, w.Pricing)
	assert.Equal(t, "cpm", w.Pricing.Model)

	// ViewableImpression.
	require.NotNil(t, w.ViewableImpression)
	require.Len(t, w.ViewableImpression.Viewable, 1)
	require.Len(t, w.ViewableImpression.NotViewable, 1)
	require.Len(t, w.ViewableImpression.ViewUndetermined, 1)

	require.NotNil(t, w.AdVerifications)
	require.Len(t, w.AdVerifications.Verification, 1)

	// BlockedAdCategories is a direct child of <Wrapper> (spec section 3.19.2).
	require.Len(t, w.BlockedAdCategories, 1)
	assert.Equal(t, "https://www.iabtechlab.com/categoryauthority", w.BlockedAdCategories[0].Authority)
	assert.Equal(t, "gambling", w.BlockedAdCategories[0].Category)

	require.NotNil(t, w.Extensions)
	require.Len(t, *w.Extensions, 1)

	// Wrapped creatives: Linear, NonLinear and Companion.
	require.Len(t, w.Creatives, 3)

	linear := w.Creatives[0].Linear
	require.NotNil(t, linear)
	require.NotNil(t, linear.TrackingEvents)
	require.Len(t, linear.TrackingEvents.Tracking, 2)
	require.NotNil(t, linear.VideoClicks)
	require.Len(t, linear.VideoClicks.ClickThroughs, 1) // ClickThrough allowed in wrapper since 4.2
	require.Len(t, linear.VideoClicks.ClickTrackings, 1)
	require.Len(t, linear.VideoClicks.CustomClicks, 1)
	require.NotNil(t, linear.Icons)
	require.NotNil(t, linear.Icons.Icon)
	require.Len(t, *linear.Icons.Icon, 1)
	require.Len(t, (*linear.Icons.Icon)[0].IconClickTracking, 1)
	assert.Equal(t, "wrapper-icon-click-tracking", (*linear.Icons.Icon)[0].IconClickTracking[0].ID)
	// InteractiveCreativeFile is limited to InLine/Linear/MediaFiles by section 3.9.3.
	assert.Empty(t, linear.InteractiveCreativeFile)

	nlAds := w.Creatives[1].NonLinearAds
	require.NotNil(t, nlAds)
	require.Len(t, nlAds.NonLinears, 1)
	require.Len(t, nlAds.NonLinears[0].NonLinearClickTracking, 1)
	require.NotNil(t, nlAds.TrackingEvents)

	compAds := w.Creatives[2].CompanionAds
	require.NotNil(t, compAds)
	assert.Equal(t, "any", compAds.Required)
	require.Len(t, compAds.Companions, 1)

	// Round-trip.
	assertRoundTrips(t, v, source)
}

// TestVAST43InteractiveCreativeDataURI focuses on the single schema change that
// VAST 4.3 introduced over 4.2: an inline data: URI on <InteractiveCreativeFile>
// (spec section 3.9.3). It verifies both parsing and round-trip of the official
// IAB-style sample and a programmatically constructed document.
func TestVAST43InteractiveCreativeDataURI(t *testing.T) {
	const dataURI = "data:text/html,%3Chtml%3E%3Cbody%3EInteractive%3C%2Fbody%3E%3C%2Fhtml%3E"

	t.Run("sample", func(t *testing.T) {
		v, source, _, err := loadFixture("testdata/iab/vast_4.3_samples/Interactive_Creative_File_Data_URI-test.xml")
		require.NoError(t, err)
		require.Equal(t, "4.3", v.Version)

		require.Len(t, v.Ads, 1)
		require.NotNil(t, v.Ads[0].InLine)
		require.Len(t, v.Ads[0].InLine.Creatives, 1)
		require.NotNil(t, v.Ads[0].InLine.Creatives[0].Linear)
		mf := v.Ads[0].InLine.Creatives[0].Linear.MediaFiles
		require.NotNil(t, mf)
		require.Len(t, mf.InteractiveCreativeFile, 1)
		assert.True(t, strings.HasPrefix(mf.InteractiveCreativeFile[0].URI, "data:"))
		assertRoundTrips(t, v, source)
	})

	t.Run("constructed", func(t *testing.T) {
		v := VAST{
			Version: "4.3",
			Ads: []Ad{{
				ID:     "1",
				AdType: "video",
				InLine: &InLine{
					AdSystem:    &AdSystem{Name: "test"},
					AdTitle:     PlainString{CDATA: "data uri"},
					Impressions: []Impression{{URI: "https://example.com/i"}},
					AdServingId: "serving-id",
					Creatives: []Creative{{
						UniversalAdID: &[]UniversalAdID{{IDRegistry: "Ad-ID", ID: "constructed-8465"}},
						Linear: &Linear{
							Duration: Duration(16e9),
							MediaFiles: &MediaFiles{
								MediaFile: []MediaFile{{
									Delivery: "progressive",
									Type:     "video/mp4",
									Width:    1,
									Height:   1,
									URI:      "https://example.com/ad.mp4",
								}},
								InteractiveCreativeFile: []InteractiveCreativeFile{{
									ApiFramework:     "SIMID",
									Type:             "text/html",
									VariableDuration: true,
									URI:              dataURI,
								}},
							},
						},
					}},
				},
			}},
		}

		out, err := xml.Marshal(&v)
		require.NoError(t, err)
		assert.Contains(t, string(out), "data:text/html,")

		var got VAST
		require.NoError(t, xml.Unmarshal(out, &got))
		icf := got.Ads[0].InLine.Creatives[0].Linear.MediaFiles.InteractiveCreativeFile[0]
		assert.Equal(t, dataURI, icf.URI)
		assert.Equal(t, "text/html", icf.Type)
		assert.Equal(t, "SIMID", icf.ApiFramework)
		assert.True(t, icf.VariableDuration)
	})
}

// TestVAST43BlockedAdCategories verifies <BlockedAdCategories> is supported in
// both placements the spec text describes: as a direct child of <Wrapper>
// (section 3.19.2) and, per the section-5 quick-reference table, nested under a
// <Verification>.
func TestVAST43BlockedAdCategories(t *testing.T) {
	t.Run("under Wrapper", func(t *testing.T) {
		const doc = `<VAST version="4.3"><Ad><Wrapper>` +
			`<AdSystem>test</AdSystem>` +
			`<Impression><![CDATA[https://example.com/i]]></Impression>` +
			`<VASTAdTagURI><![CDATA[https://example.com/vast.xml]]></VASTAdTagURI>` +
			`<BlockedAdCategories authority="iabtechlab.com">gambling</BlockedAdCategories>` +
			`</Wrapper></Ad></VAST>`
		var v VAST
		require.NoError(t, xml.Unmarshal([]byte(doc), &v))
		require.Len(t, v.Ads, 1)
		require.NotNil(t, v.Ads[0].Wrapper)
		require.Len(t, v.Ads[0].Wrapper.BlockedAdCategories, 1)
		assert.Equal(t, "iabtechlab.com", v.Ads[0].Wrapper.BlockedAdCategories[0].Authority)
		assert.Equal(t, "gambling", v.Ads[0].Wrapper.BlockedAdCategories[0].Category)
	})

	t.Run("under Verification", func(t *testing.T) {
		const doc = `<VAST version="4.3"><Ad><Wrapper>` +
			`<AdSystem>test</AdSystem>` +
			`<Impression><![CDATA[https://example.com/i]]></Impression>` +
			`<VASTAdTagURI><![CDATA[https://example.com/vast.xml]]></VASTAdTagURI>` +
			`<AdVerifications><Verification vendor="v.com-omid">` +
			`<BlockedAdCategories authority="iabtechlab.com">gambling</BlockedAdCategories>` +
			`</Verification></AdVerifications>` +
			`</Wrapper></Ad></VAST>`
		var v VAST
		require.NoError(t, xml.Unmarshal([]byte(doc), &v))
		require.NotNil(t, v.Ads[0].Wrapper)
		require.NotNil(t, v.Ads[0].Wrapper.AdVerifications)
		require.Len(t, v.Ads[0].Wrapper.AdVerifications.Verification, 1)
		require.Len(t, v.Ads[0].Wrapper.AdVerifications.Verification[0].BlockedAdCategories, 1)
		assert.Equal(t, "gambling", v.Ads[0].Wrapper.AdVerifications.Verification[0].BlockedAdCategories[0].Category)
	})
}

// TestVAST43IABSamples parses every sample under the VAST 4.3 samples directory
// and asserts it declares version 4.3 and round-trips without loss.
func TestVAST43IABSamples(t *testing.T) {
	samples, err := filepath.Glob("testdata/iab/vast_4.3_samples/*.xml")
	require.NoError(t, err)
	require.NotEmpty(t, samples, "expected at least one VAST 4.3 sample")

	for _, sample := range samples {
		t.Run(filepath.Base(sample), func(t *testing.T) {
			v, source, _, err := loadFixture(sample)
			require.NoError(t, err)
			assert.Equal(t, "4.3", v.Version)
			assertRoundTrips(t, v, source)
		})
	}
}

// TestVAST43XMLRoundTripIdempotent proves that parsing and re-marshalling the
// comprehensive fixtures is idempotent for both the InLine and Wrapper branches:
// marshal(parse(x)) and marshal(parse(marshal(parse(x)))) are equivalent.
func TestVAST43XMLRoundTripIdempotent(t *testing.T) {
	for _, path := range []string{"testdata/vast_43_inline.xml", "testdata/vast_43_wrapper.xml"} {
		t.Run(filepath.Base(path), func(t *testing.T) {
			source, err := os.ReadFile(path)
			require.NoError(t, err)

			var first VAST
			require.NoError(t, xml.Unmarshal(source, &first))
			firstOut, err := xml.Marshal(&first)
			require.NoError(t, err)

			var second VAST
			require.NoError(t, xml.Unmarshal(firstOut, &second))
			secondOut, err := xml.Marshal(&second)
			require.NoError(t, err)

			assertXMLEquivalent(t, firstOut, secondOut)
		})
	}
}

// TestVAST43JSONRoundTrip verifies the model also survives a JSON round-trip,
// which the library supports alongside XML.
func TestVAST43JSONRoundTrip(t *testing.T) {
	for _, path := range []string{"testdata/vast_43_inline.xml", "testdata/vast_43_wrapper.xml"} {
		t.Run(filepath.Base(path), func(t *testing.T) {
			v, _, _, err := loadFixture(path)
			require.NoError(t, err)

			jsonBytes, err := json.Marshal(v)
			require.NoError(t, err)

			var fromJSON VAST
			require.NoError(t, json.Unmarshal(jsonBytes, &fromJSON))

			// The struct decoded from JSON must equal the one decoded from XML.
			assert.Equal(t, *v, fromJSON)
		})
	}
}

// TestVAST43TrackingEvents verifies that the VAST 4.x tracking-event vocabulary
// (as exposed by the EventType* constants) parses from a <TrackingEvents> block.
// These events are all valid in a VAST 4.3 document.
func TestVAST43TrackingEvents(t *testing.T) {
	events := []string{
		EventTypeLoaded,
		EventTypeStart,
		EventTypeFirstQuartile,
		EventTypeMidpoint,
		EventTypeThirdQuartile,
		EventTypeComplete,
		EventTypeOtherAdInteraction,
		EventTypeProgress,
		EventTypeCloseLinear,
		EventTypeSkip,
		EventTypePlayerExpand,
		EventTypePlayerCollapse,
		EventTypeNotUsed,
		EventTypeInteractiveStart,
		EventTypeCreativeView,
		EventTypeAcceptInvitation,
		EventTypeAdExpand,
		EventTypeAdCollapse,
		EventTypeMinimize,
		EventTypeOverlayViewDuration,
		EventTypeVerificationNotExecuted,
	}

	var b strings.Builder
	b.WriteString(`<TrackingEvents>`)
	for _, e := range events {
		b.WriteString(`<Tracking event="` + e + `"><![CDATA[https://example.com/t/` + e + `]]></Tracking>`)
	}
	b.WriteString(`</TrackingEvents>`)

	var te TrackingEvents
	require.NoError(t, xml.Unmarshal([]byte(b.String()), &te))
	require.Len(t, te.Tracking, len(events))
	for i, e := range events {
		assert.Equal(t, e, te.Tracking[i].Event)
		assert.Equal(t, "https://example.com/t/"+e, te.Tracking[i].URI)
	}
}
