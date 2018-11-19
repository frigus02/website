package build

import (
	"io"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/russross/blackfriday/v2"
)

type chromaRenderer struct {
	Base      blackfriday.Renderer
	Formatter *html.Formatter
	Style     *chroma.Style
}

func newChromaRenderer() *chromaRenderer {
	return &chromaRenderer{
		Base: blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
			Flags: blackfriday.CommonHTMLFlags,
		}),
		Formatter: html.New(),
		Style:     styles.MonokaiLight,
	}
}

func (r *chromaRenderer) RenderWithChroma(w io.Writer, text []byte, data blackfriday.CodeBlockData) error {
	var lexer chroma.Lexer

	if len(data.Info) > 0 {
		lexer = lexers.Get(string(data.Info))
	}

	if lexer == nil {
		lexer = lexers.Fallback
	}

	iterator, err := lexer.Tokenise(nil, string(text))
	if err != nil {
		return err
	}

	return r.Formatter.Format(w, r.Style, iterator)
}

func (r *chromaRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	switch node.Type {
	case blackfriday.CodeBlock:
		if err := r.RenderWithChroma(w, node.Literal, node.CodeBlockData); err != nil {
			return r.Base.RenderNode(w, node, entering)
		}

		return blackfriday.SkipChildren
	default:
		return r.Base.RenderNode(w, node, entering)
	}
}

func (r *chromaRenderer) RenderHeader(w io.Writer, ast *blackfriday.Node) {
	r.Base.RenderHeader(w, ast)
}

func (r *chromaRenderer) RenderFooter(w io.Writer, ast *blackfriday.Node) {
	r.Base.RenderFooter(w, ast)
}
