package scraper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testHTML string = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>TestHTML</title>
</head>
<body>
	<div>
		<div>
			<header>
				<div>
					<div id="testElementGetNodes">
						<span>This is the single element without tags</span>
						<p id="singleElement_OneTag">This is the single element with one tag</p>
						<p id="singleElement_MultipleTags" class="multipleTags">This is the single element with multiple tags</p>
						<div class="hasDuplicate">This is the element with a duplicate</div>
						<div class="hasDuplicate">This is the second element of the duplicate</div>
						<p id="elementToBeTrimmed"> This is the elemnt which needs some trimming </p>
						<a href="https://wikipedia.com/wiki/Wikipedia" id="websiteLink">https://wikipedia.com/wiki/Wikipedia</a>
						<p id="nestedElement">This is some <span id="insideOfNestedElement">nested text</span></p>
						<p id="emptyElement"></p>
					</div>
				</div>
			</header>
		</div>
	</div>
</body>
</html>
`

func TestGetElementNodes(t *testing.T) {
	documentNode, err := GetHTMLNode(testHTML)
	require.NoError(t, err)

	testMap := make(map[string]func(t *testing.T), 0)

	testMap["testSingleElementWithoutTag"] = func(t *testing.T) {
		testElement := HtmlElement{
			Typ: "span",
		}

		nodes, err := testElement.GetElementNodes(documentNode)
		require.NoError(t, err)
		assert.Equal(t, 2, len(nodes))
		node := nodes[0]

		assert.Equal(t, testElement.Typ, node.DataAtom.String())
	}

	testMap["testSingleElementWithOneTag"] = func(t *testing.T) {
		testElement := HtmlElement{
			Typ: "p",
			Tags: []Tag{
				{
					Typ:   "id",
					Value: "singleElement_OneTag",
				},
			},
		}

		nodes, err := testElement.GetElementNodes(documentNode)
		require.NoError(t, err)
		assert.Equal(t, 1, len(nodes))
		node := nodes[0]

		assert.Equal(t, testElement.Typ, node.DataAtom.String())
		assert.Equal(t, len(testElement.Tags), len(node.Attr))

		for k, tag := range testElement.Tags {
			assert.Equal(t, tag.Typ, node.Attr[k].Key)
			assert.Equal(t, tag.Value, node.Attr[k].Val)
		}
	}

	testMap["testSingleElementWithMultipleTags"] = func(t *testing.T) {
		testElement := HtmlElement{
			Typ: "p",
			Tags: []Tag{
				{
					Typ:   "id",
					Value: "singleElement_MultipleTags",
				},
				{
					Typ:   "class",
					Value: "multipleTags",
				},
			},
		}

		nodes, err := testElement.GetElementNodes(documentNode)
		require.NoError(t, err)
		assert.Equal(t, 1, len(nodes))
		node := nodes[0]

		assert.Equal(t, testElement.Typ, node.DataAtom.String())
		assert.Equal(t, len(testElement.Tags), len(node.Attr))

		for k, tag := range testElement.Tags {
			assert.Equal(t, tag.Typ, node.Attr[k].Key)
			assert.Equal(t, tag.Value, node.Attr[k].Val)
		}
	}

	testMap["testMultipleEquivalentElements"] = func(t *testing.T) {
		testElement := HtmlElement{
			Typ: "div",
			Tags: []Tag{
				{
					Typ:   "class",
					Value: "hasDuplicate",
				},
			},
		}
		nodes, err := testElement.GetElementNodes(documentNode)
		require.NoError(t, err)
		assert.Equal(t, 2, len(nodes))
		node := nodes[0]

		assert.Equal(t, testElement.Typ, node.DataAtom.String())
		assert.Equal(t, len(testElement.Tags), len(node.Attr))

		for k, tag := range testElement.Tags {
			assert.Equal(t, tag.Typ, node.Attr[k].Key)
			assert.Equal(t, tag.Value, node.Attr[k].Val)
		}
	}

	testMap["testReturnCustomErr"] = func(t *testing.T) {
		testElement := HtmlElement{
			Typ: "div",
			Tags: []Tag{
				{
					Typ:   "id",
					Value: "thisElementDoesNotExist",
				},
			},
		}
		_, err := testElement.GetElementNodes(documentNode)
		require.Error(t, err)
		assert.Equal(t, "missing "+testElement.Typ+" in the node tree", err.Error())
	}

	for testName, testFunc := range testMap {
		t.Run(testName, testFunc)
	}
}

func TestRenderNode(t *testing.T) {
	documentNode, err := GetHTMLNode(testHTML)
	require.NoError(t, err)

	expected := `<p id="singleElement_OneTag">This is the single element with one tag</p>`
	testElement := HtmlElement{
		Typ: "p",
		Tags: []Tag{
			{
				Typ:   "id",
				Value: "singleElement_OneTag",
			},
		},
	}

	nodes, err := testElement.GetElementNodes(documentNode)
	require.NoError(t, err)
	assert.Equal(t, 1, len(nodes))
	node := nodes[0]

	assert.Equal(t, expected, RenderNode(node))
}

func TestGetTextOfNode(t *testing.T) {
	documentNode, err := GetHTMLNode(testHTML)
	require.NoError(t, err)

	testMap := make(map[string]func(t *testing.T), 0)

	testMap["testReturnTextOfNode"] = func(t *testing.T) {
		expected := `This is some nested text`
		testElement := HtmlElement{
			Typ: "p",
			Tags: []Tag{
				{
					Typ:   "id",
					Value: "nestedElement",
				},
			},
		}

		nodes, err := testElement.GetElementNodes(documentNode)
		require.NoError(t, err)
		assert.Equal(t, 1, len(nodes))
		node := nodes[0]

		assert.Equal(t, expected, GetTextOfNode(node, false))
	}

	testMap["testReturnTextOfNode_noneRecursive"] = func(t *testing.T) {
		expected := `This is some `
		testElement := HtmlElement{
			Typ: "p",
			Tags: []Tag{
				{
					Typ:   "id",
					Value: "nestedElement",
				},
			},
		}

		nodes, err := testElement.GetElementNodes(documentNode)
		require.NoError(t, err)
		assert.Equal(t, 1, len(nodes))
		node := nodes[0]

		assert.Equal(t, expected, GetTextOfNode(node, true))
	}

	for testName, testFunc := range testMap {
		t.Run(testName, testFunc)
	}
}
