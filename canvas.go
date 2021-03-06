package jujusvg

import (
	"image"
	"io"

	"github.com/ajstarks/svgo"
)

const (
	// iconSize is always 96px per Juju charm conventions.
	iconSize = 96
	maxInt   = int(^uint(0) >> 1)
	minInt   = -(maxInt - 1)
)

// Canvas holds the parsed form of a bundle or environment.
type Canvas struct {
	services  []*service
	relations []*serviceRelation
}

// service represents a service deployed to an environment and contains the
// point of the top-left corner of the icon, icon URL, and additional metadata.
type service struct {
	Name    string
	IconUrl string
	Point   image.Point
}

// serviceRelation represents a relation created between two services.
type serviceRelation struct {
	ServiceA *service
	ServiceB *service
}

// definition creates any necessary defs that can be used later in the SVG.
func (s *service) definition(canvas *svg.SVG) {
}

// usage creates any necessary tags for actually using the service in the SVG.
func (s *service) usage(canvas *svg.SVG) {
	canvas.Image(s.Point.X, s.Point.Y, iconSize, iconSize, s.IconUrl)
}

// definition creates any necessary defs that can be used later in the SVG.
func (r *serviceRelation) definition(canvas *svg.SVG) {
}

// usage creates any necessary tags for actually using the relation in the SVG.
func (r *serviceRelation) usage(canvas *svg.SVG) {
	canvas.Line(
		r.ServiceA.Point.X+(iconSize/2),
		r.ServiceA.Point.Y+(iconSize/2),
		r.ServiceB.Point.X+(iconSize/2),
		r.ServiceB.Point.Y+(iconSize/2),
		"stroke:black")
}

// AddService adds a new service to the canvas.
func (c *Canvas) AddService(s *service) {
	c.services = append(c.services, s)
}

// AddRelation adds a new relation to the canvas.
func (c *Canvas) AddRelation(r *serviceRelation) {
	c.relations = append(c.relations, r)
}

// layout adjusts all items so that they are positioned appropriately,
// and returns the overall size of the canvas.
func (c *Canvas) layout() (int, int) {
	minWidth := maxInt
	minHeight := maxInt
	maxWidth := minInt
	maxHeight := minInt

	for _, service := range c.services {
		if service.Point.X < minWidth {
			minWidth = service.Point.X
		}
		if service.Point.Y < minHeight {
			minHeight = service.Point.Y
		}
		if service.Point.X > maxWidth {
			maxWidth = service.Point.X
		}
		if service.Point.Y > maxHeight {
			maxHeight = service.Point.Y
		}
	}
	for _, service := range c.services {
		service.Point = service.Point.Sub(image.Point{X: minWidth, Y: minHeight})
	}
	return abs(maxWidth-minWidth) + iconSize,
		abs(maxHeight-minHeight) + iconSize
}

func (c *Canvas) definition(canvas *svg.SVG) {
	canvas.Def()
	defer canvas.DefEnd()
	for _, relation := range c.relations {
		relation.definition(canvas)
	}
	for _, service := range c.services {
		service.definition(canvas)
	}
}

func (c *Canvas) relationsGroup(canvas *svg.SVG) {
	canvas.Gid("relations")
	defer canvas.Gend()
	for _, relation := range c.relations {
		relation.usage(canvas)
	}
}

func (c *Canvas) servicesGroup(canvas *svg.SVG) {
	canvas.Gid("services")
	defer canvas.Gend()
	for _, service := range c.services {
		service.usage(canvas)
	}
}

// Marshal renders the SVG to the given io.Writer
func (c *Canvas) Marshal(w io.Writer) {
	width, height := c.layout()
	canvas := svg.New(w)
	canvas.Start(width, height)
	defer canvas.End()
	c.definition(canvas)
	c.relationsGroup(canvas)
	c.servicesGroup(canvas)
}

func abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}
