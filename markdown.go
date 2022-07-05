// Package goldmark implements functions to convert markdown text to a desired format.

// modified by prr

package goldmark

import (
//	"github.com/yuin/goldmark/renderer"
//	"google/gdoc/goldmark/classes"
	"fmt"
	"google/gdoc/goldmark/parser"
	"google/gdoc/goldmark/renderer"
	"google/gdoc/goldmark/renderer/html"
	"google/gdoc/goldmark/renderer/htmlAlt"
	"google/gdoc/goldmark/text"
	"google/gdoc/goldmark/util"
	"io"
)

// DefaultParser returns a new Parser that is configured by default values.
func DefaultParser() parser.Parser {
	return parser.NewParser(parser.WithBlockParsers(parser.DefaultBlockParsers()...),
		parser.WithInlineParsers(parser.DefaultInlineParsers()...),
		parser.WithParagraphTransformers(parser.DefaultParagraphTransformers()...),
	)
}

// DefaultRenderer returns a new Renderer that is configured by default values.
func DefaultRenderer(alt bool) renderer.Renderer {
fmt.Printf("hello default renderer!\n")
fmt.Printf("renderer: new renderer!\n")

var x renderer.NodeRenderer

	if alt {
		x = htmlAlt.NewRenderer()
	} else {
		x = html.NewRenderer()
	}
fmt.Printf("html renderer alt: %t %v\n", alt, x)
	y:= renderer.WithNodeRenderers(util.Prioritized(x, 1000))
fmt.Printf("node renderer: %v\n", y)
	return renderer.NewRenderer(y)
//	return renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(x, 1000)))
}

var defaultMarkdown = NewMd()
// very confusing new() vs New() new() is reserved function


// Convert interprets a UTF-8 bytes source in Markdown and
// write rendered contents to a writer w.

// why does Convert have to be abstracted??
func Convert(source []byte, w io.Writer, opts ...parser.ParseOption) error {
	return defaultMarkdown.Convert(source, w, opts...)
}

func AltConvert(source []byte, w io.Writer, opts ...parser.ParseOption) error {
// hello func AltConvert calling defaultMarkdown
fmt.Printf("hello AltConvert!\n")
	return defaultMarkdown.AltConvert(source, w, opts...)
}



// A Markdown interface offers functions to convert Markdown text to
// a desired format.
type Markdown interface {
	// Convert interprets a UTF-8 bytes source in Markdown and write rendered
	// contents to a writer w.
	Convert(source []byte, writer io.Writer, opts ...parser.ParseOption) error

	AltConvert(source []byte, writer io.Writer, opts ...parser.ParseOption) error

	// Parser returns a Parser that will be used for conversion.
	Parser() parser.Parser

	// SetParser sets a Parser to this object.
	SetParser(parser.Parser)

	// Parser returns a Renderer that will be used for conversion.
	Renderer() renderer.Renderer

	// SetRenderer sets a Renderer to this object.
	SetRenderer(renderer.Renderer)

}

// Option is a functional option type for Markdown objects.
type Option func(*markdown)

// WithExtensions adds extensions.
func WithExtensions(ext ...Extender) Option {
	return func(m *markdown) {
		m.extensions = append(m.extensions, ext...)
	}
}

// WithParser allows you to override the default parser.
func WithParser(p parser.Parser) Option {
	return func(m *markdown) {
		m.parser = p
	}
}

// WithParserOptions applies options for the parser.
func WithParserOptions(opts ...parser.Option) Option {
	return func(m *markdown) {
		m.parser.AddOptions(opts...)
	}
}

// WithRenderer allows you to override the default renderer.
func WithRenderer(r renderer.Renderer) Option {
	return func(m *markdown) {
		m.renderer = r
	}
}

// WithRendererOptions applies options for the renderer.
func WithRendererOptions(opts ...renderer.Option) Option {
	return func(m *markdown) {
		m.renderer.AddOptions(opts...)
	}
}

type markdown struct {
	parser     parser.Parser
	renderer   renderer.Renderer
	extensions []Extender
}

// New returns a new Markdown with given options.
func NewMd(options ...Option) Markdown {
fmt.Printf("hello New Md!\n")
	alt := true

	md := &markdown{
		parser:     DefaultParser(),
		renderer:   DefaultRenderer(alt),
		extensions: []Extender{},
	}
	for _, opt := range options {
		opt(md)
	}
	for _, e := range md.extensions {
fmt.Printf("extend: %v\n", e)
		e.Extend(md)
	}
//fmt.Printf("parser: %v\n", md.parser)
fmt.Printf("renderer: %v\n", md.renderer)
	return md
}

func (m *markdown) Convert(source []byte, writer io.Writer, opts ...parser.ParseOption) error {
// method that  converts a md document into a html document

	reader := text.NewReader(source)
	// reader reads the text in the source which is the md file (document)

	doc := m.parser.Parse(reader, opts...)
	// doc is the ast tree holding the parsed md file

	return m.renderer.Render(writer, source, doc)
	// m.renderer.Render is html rendering of the parsed md file
	// not sure why we still need to read the source
}

func (m *markdown) AltConvert(source []byte, writer io.Writer, opts ...parser.ParseOption) error {
	reader := text.NewReader(source)
	doc := m.parser.Parse(reader, opts...)

// we will add the default class attributes to each ast node here! not sure anymore
//	classes.SetDefAttributes(source, doc)
fmt.Printf("***  using method AltCovert ***\n")
	x:= m.renderer
fmt.Printf("m.render: %v\n", x)

	return x.Render(writer, source, doc)
}

/*
func (m *markdown) DomConvert(source []byte, writer io.Writer, opts ...parser.ParseOption) error {
	reader := text.NewReader(source)
	doc := m.parser.Parse(reader, opts...)

	// we will add the default class attributes to each ast node here
	modDoc := m.addDefAttributes(source, doc)

	return m.renderer.ScriptRender(scriptWriter, source, doc)
}
*/

func (m *markdown) Parser() parser.Parser {
	return m.parser
}

func (m *markdown) SetParser(v parser.Parser) {
	m.parser = v
}

func (m *markdown) Renderer() renderer.Renderer {
	return m.renderer
}

func (m *markdown) SetRenderer(v renderer.Renderer) {
	m.renderer = v
}

// An Extender interface is used for extending Markdown.
type Extender interface {
	// Extend extends the Markdown.
	Extend(Markdown)
}
