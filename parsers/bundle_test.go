package parsers

import (
	"bytes"
	"fmt"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type BundleParserSuite struct{}

var _ = Suite(&BundleParserSuite{})

var basket = `
charmworld-local:
  services:
    mongodb:
      charm: "cs:precise/mongodb-21"
      num_units: 1
      annotations:
        "gui-x": "940.5"
        "gui-y": "388.7698359714502"
      constraints: "mem=2G cpu-cores=1"
    elasticsearch:
      charm: "cs:~charming-devs/precise/elasticsearch-2"
      num_units: 1
      annotations:
        "gui-x": "490.5"
        "gui-y": "369.7698359714502"
      constraints: "mem=2G cpu-cores=1"
    charmworld:
      charm: "cs:~juju-jitsu/precise/charmworld-58"
      num_units: 1
      expose: true
      annotations:
        "gui-x": "813.5"
        "gui-y": "112.23016402854975"
      options:
        charm_import_limit: -1
        source: "lp:~bac/charmworld/ingest-local-charms"
        revno: 511
  relations:
    - - "charmworld:essearch"
      - "elasticsearch:essearch"
    - - "charmworld:database"
      - "mongodb:database"
  series: precise
`

func (s *BundleParserSuite) TestParse(c *C) {
	parser := BundleParser{}
	canvases, err := parser.Parse([]byte(basket))
	c.Assert(err, IsNil)
	c.Assert(len(canvases), Equals, 1)
	charmworld := canvases["charmworld-local"]
	var buf bytes.Buffer
	charmworld.Marshal(&buf)
	c.Assert(fmt.Sprintf("%s", buf.Bytes()), Equals,
		`<?xml version="1.0"?>
<!-- Generated by SVGo -->
<svg width="372" height="546"
     xmlns="http://www.w3.org/2000/svg" 
     xmlns:xlink="http://www.w3.org/1999/xlink">
<defs>
</defs>
<g id="relations">
<line x1="371" y1="48" x2="48" y2="305" style="stroke:black"/>
<line x1="371" y1="48" x2="498" y2="324" style="stroke:black"/>
</g>
<g id="services">
<image x="450" y="276" width="96" height="96" xlink:href="https://manage.jujucharms.com/api/3/charm/precise/mongodb-21/file/icon.svg" />
<image x="0" y="257" width="96" height="96" xlink:href="https://manage.jujucharms.com/api/3/charm/~charming-devs/precise/elasticsearch-2/file/icon.svg" />
<image x="323" y="0" width="96" height="96" xlink:href="https://manage.jujucharms.com/api/3/charm/~juju-jitsu/precise/charmworld-58/file/icon.svg" />
</g>
</svg>
`)
}
