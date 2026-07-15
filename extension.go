package vast

import (
	"encoding/xml"
	"io"
	"strings"
)

// Extension represent arbitrary XML provided by the platform to extend the
// VAST response or by custom trackers.
type Extension struct {
	Type string `xml:"type,attr,omitempty"`
	// Attributes preserves custom attributes used by CreativeExtension and
	// vendor-defined Extension elements. The type attribute remains available
	// through Type and is not duplicated here.
	Attributes     []xml.Attr `xml:"-" json:",omitempty"`
	CustomTracking []Tracking `xml:"CustomTracking>Tracking,omitempty"  json:",omitempty"`
	Data           string     `xml:",innerxml" json:",omitempty"`
}

// the extension type as a middleware in the encoding process.
type extension Extension

type customTrackingXML struct {
	Tracking []Tracking `xml:"Tracking,omitempty"`
}

type extensionXML struct {
	Type           string             `xml:"type,attr,omitempty"`
	CustomTracking *customTrackingXML `xml:"CustomTracking,omitempty"`
	Data           string             `xml:",innerxml"`
}

// MarshalXML implements xml.Marshaler interface.
func (e Extension) MarshalXML(enc *xml.Encoder, start xml.StartElement) error {
	appendPreservedAttributes(&start, e.Attributes)
	value := extensionXML{Type: e.Type, Data: e.Data}
	if len(e.CustomTracking) > 0 {
		value.CustomTracking = &customTrackingXML{Tracking: e.CustomTracking}
	}
	return enc.EncodeElement(value, start)
}

// UnmarshalXML implements xml.Unmarshaler interface.
func (e *Extension) UnmarshalXML(dec *xml.Decoder, start xml.StartElement) error {
	// decode the extension into a temporary element from a wrapper Extension,
	// copy what we need over.
	var e2 extension
	if err := dec.DecodeElement(&e2, &start); err != nil {
		return err
	}
	// Copy the parsed fields while keeping arbitrary vendor XML independently
	// editable from the structured custom trackers.
	e.Type = e2.Type
	e.Attributes = e.Attributes[:0]
	for _, attr := range start.Attr {
		if attr.Name.Space == "" && attr.Name.Local == "type" {
			continue
		}
		e.Attributes = append(e.Attributes, attr)
	}
	e.CustomTracking = e2.CustomTracking
	data, err := extensionDataWithoutCustomTracking(e2.Data)
	if err != nil {
		return err
	}
	e.Data = data
	return nil
}

// extensionDataWithoutCustomTracking removes only top-level CustomTracking
// elements. InputOffset lets the arbitrary XML outside those elements remain
// byte-for-byte intact.
func extensionDataWithoutCustomTracking(data string) (string, error) {
	if data == "" {
		return "", nil
	}

	const prefix = "<extension-data>"
	wrapped := prefix + data + "</extension-data>"
	dec := xml.NewDecoder(strings.NewReader(wrapped))
	depth := 0
	spanStart := -1
	lastEnd := 0
	var remaining strings.Builder

	for {
		tokenStart := int(dec.InputOffset()) - len(prefix)
		token, err := dec.RawToken()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		switch token := token.(type) {
		case xml.StartElement:
			if depth == 1 && token.Name.Local == "CustomTracking" {
				spanStart = tokenStart
			}
			depth++
		case xml.EndElement:
			if spanStart >= 0 && depth == 2 {
				spanEnd := int(dec.InputOffset()) - len(prefix)
				remaining.WriteString(data[lastEnd:spanStart])
				lastEnd = spanEnd
				spanStart = -1
			}
			depth--
		}
	}

	remaining.WriteString(data[lastEnd:])
	if strings.TrimSpace(remaining.String()) == "" {
		return "", nil
	}
	return remaining.String(), nil
}
