// Package visitor contains walker.visitor implementations
package pkg

import (
	"fmt"
	"github.com/z7zmey/php-parser/freefloating"
	"github.com/z7zmey/php-parser/node"
	"github.com/z7zmey/php-parser/visitor"
	"github.com/z7zmey/php-parser/walker"
	"io"
	"reflect"
	"strings"
)

type AnnotationRouteDumper struct {
	Writer         io.Writer
	NsResolver     *visitor.NamespaceResolver
	depth          int
	isChildNode    bool
	isNotFirstNode bool
}

func NewAnnotationRouteDumper(Writer io.Writer, NsResolver *visitor.NamespaceResolver) *AnnotationRouteDumper {
	return &AnnotationRouteDumper{
		Writer:         Writer,
		NsResolver:     NsResolver,
		depth:          0,
		isChildNode:    false,
		isNotFirstNode: false,
	}
}

func (d *AnnotationRouteDumper) printIndent(w io.Writer) {
	for i := 0; i < d.depth; i++ {
		//fmt.Fprint(d.Writer, "  ")
	}
}

// EnterNode is invoked at every node in hierarchy
func (d *AnnotationRouteDumper) EnterNode(w walker.Walkable) bool {
	n := w.(node.Node)

	nodeType := reflect.TypeOf(n).String()

	// only care about Class and ClassMethod nodes
	if nodeType == "*stmt.Class" || nodeType == "*stmt.ClassMethod" {

		// get the node attributes
		attributes := n.Attributes()
		// make there is a PhpDocComment attribute
		if phpDocBlock, ok := attributes["PhpDocComment"]; ok {
			// cast interface{} to string
			phpDocBlockString := fmt.Sprintf("%v", phpDocBlock)

			if strings.Contains(phpDocBlockString, "@Route") {
				_, route := ExtractRoute(phpDocBlockString)
				fmt.Fprint(d.Writer, route)
				fmt.Fprintln(d.Writer, "")
			}
		}

	}

	if d.isChildNode {
		d.isChildNode = false
	} else if d.isNotFirstNode {
		//fmt.Fprint(d.Writer, ",\n")
		d.printIndent(d.Writer)
	} else {
		d.printIndent(d.Writer)
		d.isNotFirstNode = true
	}

	//fmt.Fprint(d.Writer, "{\n")
	d.depth++
	d.printIndent(d.Writer)
	//fmt.Fprintf(d.Writer, "%q: %q", "type", nodeType)

	if p := n.GetPosition(); p != nil {
		//fmt.Fprint(d.Writer, ",\n")
		d.printIndent(d.Writer)
		//fmt.Fprintf(d.Writer, "%q: {\n", "position")
		d.depth++
		d.printIndent(d.Writer)
		//fmt.Fprintf(d.Writer, "%q: %d,\n", "startPos", p.StartPos)
		d.printIndent(d.Writer)
		//fmt.Fprintf(d.Writer, "%q: %d,\n", "endPos", p.EndPos)
		d.printIndent(d.Writer)
		//fmt.Fprintf(d.Writer, "%q: %d,\n", "startLine", p.StartLine)
		d.printIndent(d.Writer)
		//fmt.Fprintf(d.Writer, "%q: %d\n", "endLine", p.EndLine)
		d.depth--
		d.printIndent(d.Writer)
		//fmt.Fprint(d.Writer, "}")
	}

	//if d.NsResolver != nil {
	//	if namespacedName, ok := d.NsResolver.ResolvedNames[n]; ok {
	//		//fmt.Fprint(d.Writer, ",\n")
	//		d.printIndent(d.Writer)
	//		//fmt.Fprintf(d.Writer, "\"namespacedName\": %q", namespacedName)
	//	}
	//}

	if !n.GetFreeFloating().IsEmpty() {
		//fmt.Fprint(d.Writer, ",\n")
		d.printIndent(d.Writer)
		//fmt.Fprint(d.Writer, "\"freefloating\": {\n")
		d.depth++
		i := 0
		for _, freeFloatingStrings := range *n.GetFreeFloating() {
			if i != 0 {
				//fmt.Fprint(d.Writer, ",\n")
			}
			i++

			d.printIndent(d.Writer)
			//fmt.Fprintf(d.Writer, "%q: [\n", key)
			d.depth++

			j := 0
			for _, freeFloatingString := range freeFloatingStrings {
				if j != 0 {
					//fmt.Fprint(d.Writer, ",\n")
				}
				j++

				d.printIndent(d.Writer)
				//fmt.Fprint(d.Writer, "{\n")
				d.depth++
				d.printIndent(d.Writer)
				switch freeFloatingString.StringType {
				case freefloating.CommentType:
					//fmt.Fprintf(d.Writer, "%q: %q,\n", "type", "freefloating.CommentType")
				case freefloating.WhiteSpaceType:
					//fmt.Fprintf(d.Writer, "%q: %q,\n", "type", "freefloating.WhiteSpaceType")
				case freefloating.TokenType:
					//fmt.Fprintf(d.Writer, "%q: %q,\n", "type", "freefloating.TokenType")
				}
				d.printIndent(d.Writer)
				//fmt.Fprintf(d.Writer, "%q: %q\n", "value", freeFloatingString.Value)
				d.depth--
				d.printIndent(d.Writer)
				//fmt.Fprint(d.Writer, "}")
			}

			d.depth--
			//fmt.Fprint(d.Writer, "\n")
			d.printIndent(d.Writer)
			//fmt.Fprint(d.Writer, "]")
		}
		d.depth--
		//fmt.Fprint(d.Writer, "\n")
		d.printIndent(d.Writer)
		//fmt.Fprint(d.Writer, "}")
	}

	if a := n.Attributes(); len(a) > 0 {
		for _, attr := range a {
			//fmt.Fprint(d.Writer, ",\n")
			d.printIndent(d.Writer)
			switch attr.(type) {
			case string:
				//fmt.Fprintf(d.Writer, "\"%s\": %q", key, attr)
			default:
				//fmt.Fprintf(d.Writer, "\"%s\": %v", key, attr)
			}
		}
	}

	return true
}

// LeaveNode is invoked after node process
func (d *AnnotationRouteDumper) LeaveNode(n walker.Walkable) {
	d.depth--
	//fmt.Fprint(d.Writer, "\n")
	d.printIndent(d.Writer)
	//fmt.Fprint(d.Writer, "}")
}

func (d *AnnotationRouteDumper) EnterChildNode(key string, w walker.Walkable) {
	//fmt.Fprint(d.Writer, ",\n")
	d.printIndent(d.Writer)
	//fmt.Fprintf(d.Writer, "%q: ", key)
	d.isChildNode = true
}

func (d *AnnotationRouteDumper) LeaveChildNode(key string, w walker.Walkable) {
	// do nothing
}

func (d *AnnotationRouteDumper) EnterChildList(key string, w walker.Walkable) {
	//fmt.Fprint(d.Writer, ",\n")
	d.printIndent(d.Writer)
	//fmt.Fprintf(d.Writer, "%q: [\n", key)
	d.depth++

	d.isNotFirstNode = false
}

func (d *AnnotationRouteDumper) LeaveChildList(key string, w walker.Walkable) {
	d.depth--
	//fmt.Fprint(d.Writer, "\n")
	d.printIndent(d.Writer)
	//fmt.Fprint(d.Writer, "]")
}
