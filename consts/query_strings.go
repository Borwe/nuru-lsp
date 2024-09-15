package consts

import "fmt"

var PAKEJI_TAG = "pakeji_name"
var FUNCTION_TAG = "function_name"
var VARIABLE_TAG = "variable_name"

//get variables, and functions created
var TMUIA_PAEKJI_QUERY = fmt.Sprintf("(pakeji_tumia_statement pakejiname: (identifier)* @%s)",PAKEJI_TAG)
var FUNCTION_DECLARATION_QUERY = fmt.Sprintf("(declaration_statment (identifier)* @%s (equal) (function_statement))",FUNCTION_TAG)
var VARIABLE_DECLARATION_QUERY = fmt.Sprintf("(declaration_statment (identifier) @%s (equal) (expression))", VARIABLE_TAG)

// used for seeing if curent is a packaji
var HII_NI_PAKEJI = "(pakeji_statement) @pakeji"
