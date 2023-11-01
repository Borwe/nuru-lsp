package tree_paser

type Content struct {
	Load     []rune
	Position int
}

type Text struct {
	IsComent bool
	Value []rune
}

func (t *Text) IsEmpty() bool {
	return len(t.Value) == 0
}

func NewContent(text string) Content {
	return Content{
		Load: []rune(text),
		Position: 0,
	}
}

/*
Returns a string represnting the workd or comment,
with false pair if it is a word, and true if it is a
comment string
*/
func (c *Content) GetWordOrComment() Text {
	value := []rune{}
	isComment := false;

	if c.Position >= len(c.Load) {
		return Text{IsComent: false, Value: value,}
	}

	isSingleLineComment := false;
	if len(c.Load)-c.Position > 1 && c.Load[c.Position] == '/' && c.Load[c.Position+1] == '/' {
		isSingleLineComment = true
		isComment = true
	}

	isMultiLineComment := false;
	if len(c.Load)-c.Position > 3 && c.Load[c.Position] == '/' && c.Load[c.Position+1] == '*' {
		isSingleLineComment = true
		isComment = true
	}


	for c.Position < len(c.Load) {
		curr := c.Load[c.Position]
		//handle single comment
		if isSingleLineComment {
			if curr!='\n'{
				value = append(value, curr)
				continue;
			}
			break
		}

		if isMultiLineComment {
			if c.Load[c.Position] == '*' && c.Load[c.Position+1] =='/' {
				value = append(value, '*', '/' )
				c.Position += 2
				break;
			}
			value = append(value, curr)
			continue
		}

		if curr != '\t' && curr != '\n' && curr != ' ' && (curr !='/' && c.Load[c.Position+1] != '/') {
			value = append(value, curr)
		}
		c.Position +=1
	}
	return Text{IsComent: isComment, Value: value,}
}

func IsTumia(text Text) bool {
	return string(text.Value) == "tumia"
}
