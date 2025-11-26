package vast

import (
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	extensionCustomTracking = []byte(`<Extension type="testCustomTracking"><CustomTracking><Tracking event="event.1"><![CDATA[http://event.1]]></Tracking><Tracking event="event.2"><![CDATA[http://event.2]]></Tracking></CustomTracking></Extension>`)
	extensionData           = []byte(`<Extension type="testCustomTracking"><SkippableAdType>Generic</SkippableAdType></Extension>`)
)

func TestExtensionCustomTrackingMarshal(t *testing.T) {
	e := Extension{
		Type: "testCustomTracking",
		CustomTracking: []Tracking{
			{
				Event: "event.1",
				URI:   "http://event.1",
			},
			{
				Event: "event.2",
				URI:   "http://event.2",
			},
		},
	}

	// marshal the extension
	xmlExtensionOutput, err := xml.Marshal(e)
	assert.NoError(t, err)

	// assert the resulting marshaled extension
	assert.Equal(t, string(extensionCustomTracking), string(xmlExtensionOutput))
}

func TestExtensionCustomTracking(t *testing.T) {
	// unmarshal the Extension
	var e Extension
	assert.NoError(t, xml.Unmarshal(extensionCustomTracking, &e))

	// assert the resulting extension
	assert.Equal(t, "testCustomTracking", e.Type)
	assert.Empty(t, string(e.Data))
	if assert.Len(t, e.CustomTracking, 2) {
		// first event
		assert.Equal(t, "event.1", e.CustomTracking[0].Event)
		assert.Equal(t, "http://event.1", e.CustomTracking[0].URI)
		// second event
		assert.Equal(t, "event.2", e.CustomTracking[1].Event)
		assert.Equal(t, "http://event.2", e.CustomTracking[1].URI)
	}

	// marshal the extension
	xmlExtensionOutput, err := xml.Marshal(e)
	assert.NoError(t, err)

	// assert the resulting marshaled extension
	assert.Equal(t, string(extensionCustomTracking), string(xmlExtensionOutput))
}

func TestExtensionGeneric(t *testing.T) {
	// unmarshal the Extension
	var e Extension
	assert.NoError(t, xml.Unmarshal(extensionData, &e))

	// assert the resulting extension
	assert.Equal(t, "testCustomTracking", e.Type)
	assert.Equal(t, "<SkippableAdType>Generic</SkippableAdType>", string(e.Data))
	assert.Empty(t, e.CustomTracking)

	// marshal the extension
	xmlExtensionOutput, err := xml.Marshal(e)
	assert.NoError(t, err)

	// assert the resulting marshaled extension
	assert.Equal(t, string(extensionData), string(xmlExtensionOutput))
}

func TestViaplayFreeWheel_AdInlineExtensions(t *testing.T) {
	v, _, _, err := loadFixture("testdata/viaplay-freewheel-vast-example.xml")
	assert.NoError(t, err)

	assert.Equal(t, "4.1", v.Version)
	assert.Len(t, v.Ads, 1)
	ad := v.Ads[0]
	assert.NotNil(t, ad.InLine)
	assert.NotNil(t, ad.InLine.Extensions)
	exts := *ad.InLine.Extensions
	assert.Len(t, exts, 1)
	ext := exts[0]
	assert.Equal(t, "FreeWheel", ext.Type)
	assert.NotNil(t, ext.SSAICreativeID)
	assert.Equal(t, "252235610", ext.SSAICreativeID.CreativeID)
	assert.Equal(t, "252235610", ext.SSAICreativeID.Data)
	assert.Len(t, ext.Parameters, 1)
	param := ext.Parameters[0]
	assert.Equal(t, "_fw_4AID", param.Name)
	assert.Equal(t, "2JTVZJTWRA4G", param.Value)
}
