package vast

import (
	"encoding/json"
	"encoding/xml"
	"slices"
)

type resourceAliasState struct {
	htmlInitialized   bool
	iframeInitialized bool
	staticInitialized bool
	html              *HTMLResource
	iframe            *CDATAString
	static            *StaticResource
}

type clickTrackingAliasState struct {
	initialized bool
	values      []CDATAString
}

func clonePointer[T any](value *T) *T {
	if value == nil {
		return nil
	}
	cloned := *value
	return &cloned
}

func newResourceAliasState(
	html *HTMLResource,
	iframe *CDATAString,
	static *StaticResource,
	htmlInitialized bool,
	iframeInitialized bool,
	staticInitialized bool,
) resourceAliasState {
	return resourceAliasState{
		htmlInitialized:   htmlInitialized,
		iframeInitialized: iframeInitialized,
		staticInitialized: staticInitialized,
		html:              clonePointer(html),
		iframe:            clonePointer(iframe),
		static:            clonePointer(static),
	}
}

func equalPointers[T comparable](left, right *T) bool {
	if left == nil || right == nil {
		return left == nil && right == nil
	}
	return *left == *right
}

func reconcileResourceAlias[T comparable](plural []T, singular, original *T, initialized bool) []T {
	if !initialized {
		if len(plural) == 0 && singular != nil {
			return []T{*singular}
		}
		return plural
	}
	if equalPointers(singular, original) {
		return plural
	}
	if singular == nil {
		if len(plural) == 0 {
			return plural
		}
		return append([]T(nil), plural[1:]...)
	}
	if len(plural) == 0 {
		return []T{*singular}
	}
	reconciled := append([]T(nil), plural...)
	reconciled[0] = *singular
	return reconciled
}

func xmlAttribute(start xml.StartElement, name string) (string, bool) {
	for _, attr := range start.Attr {
		if attr.Name.Local == name {
			return attr.Value, true
		}
	}
	return "", false
}

func appendExplicitFalse(start *xml.StartElement, name string, value, set bool) {
	if set && !value {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: name}, Value: "false"})
	}
}

func appendPreservedAttributes(start *xml.StartElement, attributes []xml.Attr) {
	prefixes := make(map[string]string)
	for _, attr := range attributes {
		if attr.Name.Space == "xmlns" {
			prefixes[attr.Value] = attr.Name.Local
		}
	}
	for _, attr := range attributes {
		switch {
		case attr.Name.Space == "xmlns":
			attr.Name = xml.Name{Local: "xmlns:" + attr.Name.Local}
		case attr.Name.Space != "" && prefixes[attr.Name.Space] != "":
			attr.Name = xml.Name{Local: prefixes[attr.Name.Space] + ":" + attr.Name.Local}
		}
		start.Attr = append(start.Attr, attr)
	}
}

func (vast VAST) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type plain VAST
	appendPreservedAttributes(&start, vast.Attributes)
	return enc.EncodeElement(plain(vast), start)
}

func (vast *VAST) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type plain VAST
	var value plain
	if err := dec.DecodeElement(&value, &start); err != nil {
		return err
	}
	*vast = VAST(value)
	for _, attr := range start.Attr {
		if attr.Name.Local == "version" || (attr.Name.Space == "" && attr.Name.Local == "xmlns") {
			continue
		}
		vast.Attributes = append(vast.Attributes, attr)
	}
	return nil
}

// MarshalXML preserves an explicitly supplied conditionalAd="false" while
// retaining the historical bool field API.
func (ad Ad) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type plain Ad
	appendExplicitFalse(&start, "conditionalAd", ad.ConditionalAd, ad.ConditionalAdSet)
	return enc.EncodeElement(plain(ad), start)
}

func (ad *Ad) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type plain Ad
	var value plain
	if err := dec.DecodeElement(&value, &start); err != nil {
		return err
	}
	*ad = Ad(value)
	_, ad.ConditionalAdSet = xmlAttribute(start, "conditionalAd")
	return nil
}

func (resource JavaScriptResource) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type plain JavaScriptResource
	appendExplicitFalse(&start, "browserOptional", resource.BrowserOptional, resource.BrowserOptionalSet)
	return enc.EncodeElement(plain(resource), start)
}

func (resource *JavaScriptResource) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type plain JavaScriptResource
	var value plain
	if err := dec.DecodeElement(&value, &start); err != nil {
		return err
	}
	*resource = JavaScriptResource(value)
	_, resource.BrowserOptionalSet = xmlAttribute(start, "browserOptional")
	return nil
}

func (file InteractiveCreativeFile) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type plain InteractiveCreativeFile
	appendExplicitFalse(&start, "variableDuration", file.VariableDuration, file.VariableDurationSet)
	return enc.EncodeElement(plain(file), start)
}

func (file *InteractiveCreativeFile) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type plain InteractiveCreativeFile
	var value plain
	if err := dec.DecodeElement(&value, &start); err != nil {
		return err
	}
	*file = InteractiveCreativeFile(value)
	_, file.VariableDurationSet = xmlAttribute(start, "variableDuration")
	return nil
}

func (resource HTMLResource) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type plain HTMLResource
	appendExplicitFalse(&start, "xmlEncoded", resource.XMLEncoded, resource.XMLEncodedSet)
	return enc.EncodeElement(plain(resource), start)
}

func (resource *HTMLResource) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type plain HTMLResource
	var value plain
	if err := dec.DecodeElement(&value, &start); err != nil {
		return err
	}
	*resource = HTMLResource(value)
	_, resource.XMLEncodedSet = xmlAttribute(start, "xmlEncoded")
	return nil
}

func (parameters AdParameters) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type plain AdParameters
	appendExplicitFalse(&start, "xmlEncoded", parameters.XMLEncoded, parameters.XMLEncodedSet)
	return enc.EncodeElement(plain(parameters), start)
}

func (parameters *AdParameters) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type plain AdParameters
	var value plain
	if err := dec.DecodeElement(&value, &start); err != nil {
		return err
	}
	*parameters = AdParameters(value)
	_, parameters.XMLEncodedSet = xmlAttribute(start, "xmlEncoded")
	return nil
}

func (file MediaFile) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type plain MediaFile
	appendExplicitFalse(&start, "scalable", file.Scalable, file.ScalableSet)
	appendExplicitFalse(&start, "maintainAspectRatio", file.MaintainAspectRatio, file.MaintainAspectRatioSet)
	return enc.EncodeElement(plain(file), start)
}

func (file *MediaFile) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type plain MediaFile
	var value plain
	if err := dec.DecodeElement(&value, &start); err != nil {
		return err
	}
	*file = MediaFile(value)
	_, file.ScalableSet = xmlAttribute(start, "scalable")
	_, file.MaintainAspectRatioSet = xmlAttribute(start, "maintainAspectRatio")
	return nil
}

func (companion Companion) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type plain Companion
	value := companion
	value.populatePluralResources()
	return enc.EncodeElement(plain(value), start)
}

func (companion *Companion) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type plain Companion
	var value plain
	if err := dec.DecodeElement(&value, &start); err != nil {
		return err
	}
	*companion = Companion(value)
	companion.linkSingularResources()
	companion.captureResourceAliases()
	return nil
}

func (companion *Companion) UnmarshalJSON(data []byte) error {
	type plain Companion
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*companion = Companion(value)
	companion.captureResourceAliases()
	return nil
}

func (companion *Companion) populatePluralResources() {
	state := companion.resourceAliases
	companion.HTMLResources = reconcileResourceAlias(companion.HTMLResources, companion.HTMLResource, state.html, state.htmlInitialized)
	companion.IFrameResources = reconcileResourceAlias(companion.IFrameResources, companion.IFrameResource, state.iframe, state.iframeInitialized)
	companion.StaticResources = reconcileResourceAlias(companion.StaticResources, companion.StaticResource, state.static, state.staticInitialized)
}

func (companion *Companion) linkSingularResources() {
	if len(companion.HTMLResources) > 0 {
		companion.HTMLResource = &companion.HTMLResources[0]
	}
	if len(companion.IFrameResources) > 0 {
		companion.IFrameResource = &companion.IFrameResources[0]
	}
	if len(companion.StaticResources) > 0 {
		companion.StaticResource = &companion.StaticResources[0]
	}
}

func (companion *Companion) captureResourceAliases() {
	companion.resourceAliases = newResourceAliasState(
		companion.HTMLResource,
		companion.IFrameResource,
		companion.StaticResource,
		len(companion.HTMLResources) > 0,
		len(companion.IFrameResources) > 0,
		len(companion.StaticResources) > 0,
	)
}

func (nonLinear NonLinear) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type plain NonLinear
	value := nonLinear
	value.populatePluralResources()
	appendExplicitFalse(&start, "scalable", value.Scalable, value.ScalableSet)
	appendExplicitFalse(&start, "maintainAspectRatio", value.MaintainAspectRatio, value.MaintainAspectRatioSet)
	return enc.EncodeElement(plain(value), start)
}

func (nonLinear *NonLinear) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type plain NonLinear
	var value plain
	if err := dec.DecodeElement(&value, &start); err != nil {
		return err
	}
	*nonLinear = NonLinear(value)
	_, nonLinear.ScalableSet = xmlAttribute(start, "scalable")
	_, nonLinear.MaintainAspectRatioSet = xmlAttribute(start, "maintainAspectRatio")
	nonLinear.linkSingularResources()
	nonLinear.captureResourceAliases()
	return nil
}

func (nonLinear *NonLinear) UnmarshalJSON(data []byte) error {
	type plain NonLinear
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*nonLinear = NonLinear(value)
	nonLinear.captureResourceAliases()
	return nil
}

func (nonLinear *NonLinear) populatePluralResources() {
	state := nonLinear.resourceAliases
	nonLinear.HTMLResources = reconcileResourceAlias(nonLinear.HTMLResources, nonLinear.HTMLResource, state.html, state.htmlInitialized)
	nonLinear.IFrameResources = reconcileResourceAlias(nonLinear.IFrameResources, nonLinear.IFrameResource, state.iframe, state.iframeInitialized)
	nonLinear.StaticResources = reconcileResourceAlias(nonLinear.StaticResources, nonLinear.StaticResource, state.static, state.staticInitialized)
}

func (nonLinear *NonLinear) linkSingularResources() {
	if len(nonLinear.HTMLResources) > 0 {
		nonLinear.HTMLResource = &nonLinear.HTMLResources[0]
	}
	if len(nonLinear.IFrameResources) > 0 {
		nonLinear.IFrameResource = &nonLinear.IFrameResources[0]
	}
	if len(nonLinear.StaticResources) > 0 {
		nonLinear.StaticResource = &nonLinear.StaticResources[0]
	}
}

func (nonLinear *NonLinear) captureResourceAliases() {
	nonLinear.resourceAliases = newResourceAliasState(
		nonLinear.HTMLResource,
		nonLinear.IFrameResource,
		nonLinear.StaticResource,
		len(nonLinear.HTMLResources) > 0,
		len(nonLinear.IFrameResources) > 0,
		len(nonLinear.StaticResources) > 0,
	)
}

func (nonLinear NonLinearWrapper) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type plain NonLinearWrapper
	value := nonLinear
	legacyChanged := value.clickTrackingAlias.initialized && !slices.Equal(value.NonLinearClickTracking, value.clickTrackingAlias.values)
	if legacyChanged || (!value.clickTrackingAlias.initialized && len(value.NonLinearClickTrackings) == 0 && len(value.NonLinearClickTracking) > 0) {
		value.NonLinearClickTrackings = make([]NonLinearClickTracking, len(value.NonLinearClickTracking))
		for i, tracking := range value.NonLinearClickTracking {
			value.NonLinearClickTrackings[i].URI = tracking.CDATA
		}
	}
	appendExplicitFalse(&start, "scalable", value.Scalable, value.ScalableSet)
	appendExplicitFalse(&start, "maintainAspectRatio", value.MaintainAspectRatio, value.MaintainAspectRatioSet)
	return enc.EncodeElement(plain(value), start)
}

func (nonLinear *NonLinearWrapper) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type plain NonLinearWrapper
	var value plain
	if err := dec.DecodeElement(&value, &start); err != nil {
		return err
	}
	*nonLinear = NonLinearWrapper(value)
	_, nonLinear.ScalableSet = xmlAttribute(start, "scalable")
	_, nonLinear.MaintainAspectRatioSet = xmlAttribute(start, "maintainAspectRatio")
	if len(nonLinear.NonLinearClickTrackings) > 0 {
		nonLinear.NonLinearClickTracking = make([]CDATAString, len(nonLinear.NonLinearClickTrackings))
		for i, tracking := range nonLinear.NonLinearClickTrackings {
			nonLinear.NonLinearClickTracking[i].CDATA = tracking.URI
		}
	}
	nonLinear.captureClickTrackingAlias()
	return nil
}

func (nonLinear *NonLinearWrapper) UnmarshalJSON(data []byte) error {
	type plain NonLinearWrapper
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*nonLinear = NonLinearWrapper(value)
	nonLinear.captureClickTrackingAlias()
	return nil
}

func (nonLinear *NonLinearWrapper) captureClickTrackingAlias() {
	nonLinear.clickTrackingAlias = clickTrackingAliasState{
		initialized: len(nonLinear.NonLinearClickTrackings) > 0,
		values:      append([]CDATAString(nil), nonLinear.NonLinearClickTracking...),
	}
}

type iconClicksXML struct {
	FallbackImages *IconClickFallbackImages `xml:"IconClickFallbackImages,omitempty"`
	ClickThrough   *CDATAString             `xml:"IconClickThrough,omitempty"`
	ClickTracking  []IconClickTracking      `xml:"IconClickTracking,omitempty"`
}

func (icon Icon) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	type plain Icon
	type element struct {
		plain
		Clicks       *iconClicksXML `xml:"IconClicks,omitempty"`
		ViewTracking []CDATAString  `xml:"IconViewTracking,omitempty"`
	}

	value := icon
	value.populatePluralResources()
	if value.OffsetSet || value.Offset.Duration != nil || value.Offset.Percent != 0 {
		text, err := value.Offset.MarshalText()
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "offset"}, Value: string(text)})
	}
	if value.DurationSet && value.Duration == 0 {
		text, err := value.Duration.MarshalText()
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "duration"}, Value: string(text)})
	}

	var clicks *iconClicksXML
	if value.IconClickFallbackImages != nil || value.IconClickThrough != nil || len(value.IconClickTracking) > 0 {
		clicks = &iconClicksXML{
			FallbackImages: value.IconClickFallbackImages,
			ClickThrough:   value.IconClickThrough,
			ClickTracking:  value.IconClickTracking,
		}
	}
	return enc.EncodeElement(element{plain: plain(value), Clicks: clicks, ViewTracking: value.IconViewTracking}, start)
}

func (icon *Icon) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	type plain Icon
	type element struct {
		plain
		Clicks       *iconClicksXML `xml:"IconClicks,omitempty"`
		ViewTracking []CDATAString  `xml:"IconViewTracking,omitempty"`
	}
	var value element
	if err := dec.DecodeElement(&value, &start); err != nil {
		return err
	}
	*icon = Icon(value.plain)
	if offset, ok := xmlAttribute(start, "offset"); ok {
		if err := icon.Offset.UnmarshalText([]byte(offset)); err != nil {
			return err
		}
		icon.OffsetSet = true
	}
	_, icon.DurationSet = xmlAttribute(start, "duration")
	if value.Clicks != nil {
		icon.IconClickFallbackImages = value.Clicks.FallbackImages
		icon.IconClickThrough = value.Clicks.ClickThrough
		icon.IconClickTracking = value.Clicks.ClickTracking
	}
	icon.IconViewTracking = value.ViewTracking
	icon.linkSingularResources()
	icon.captureResourceAliases()
	return nil
}

func (icon *Icon) UnmarshalJSON(data []byte) error {
	type plain Icon
	var value plain
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	*icon = Icon(value)
	icon.captureResourceAliases()
	return nil
}

func (icon *Icon) populatePluralResources() {
	state := icon.resourceAliases
	icon.HTMLResources = reconcileResourceAlias(icon.HTMLResources, icon.HTMLResource, state.html, state.htmlInitialized)
	icon.IFrameResources = reconcileResourceAlias(icon.IFrameResources, icon.IFrameResource, state.iframe, state.iframeInitialized)
	icon.StaticResources = reconcileResourceAlias(icon.StaticResources, icon.StaticResource, state.static, state.staticInitialized)
}

func (icon *Icon) linkSingularResources() {
	if len(icon.HTMLResources) > 0 {
		icon.HTMLResource = &icon.HTMLResources[0]
	}
	if len(icon.IFrameResources) > 0 {
		icon.IFrameResource = &icon.IFrameResources[0]
	}
	if len(icon.StaticResources) > 0 {
		icon.StaticResource = &icon.StaticResources[0]
	}
}

func (icon *Icon) captureResourceAliases() {
	icon.resourceAliases = newResourceAliasState(
		icon.HTMLResource,
		icon.IFrameResource,
		icon.StaticResource,
		len(icon.HTMLResources) > 0,
		len(icon.IFrameResources) > 0,
		len(icon.StaticResources) > 0,
	)
}
