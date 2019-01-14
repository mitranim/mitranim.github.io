# Dependencies:
#
#     * lxml: `pip3 install lxml`
#
#     * atom.rng; source: http://cweiske.de/tagebuch/atom-validation.htm
#
# Note: compliance with the XML schema doesn't guarantee that the feed is
# well-formed. For example, it doesn't ensure that links have an appropriate
# base or host.

import sys
from lxml import etree

rng = etree.RelaxNG(file='atom_schema.rng')
doc = etree.parse('../public/feed.xml')
rng.assertValid(doc)
