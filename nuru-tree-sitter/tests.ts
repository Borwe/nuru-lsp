import assert from "assert"
import { execSync } from "child_process"

//For testing a.nr for errors, should have no errors
function testParseA_NR(){
    let data = execSync("npx tree-sitter parse ../a.nr").toString()
    assert(data.search("ERROR")==-1, "Found an ERROR when parsing ../a.nr")
}


//For testing test.nr for errors, should have no errors
function testParseTEST_NR(){
    let data = execSync("npx tree-sitter parse ../test.nr").toString()
    assert(data.search("ERROR")==-1, "Found an ERROR when parsing ../test.nr")
}

testParseA_NR()
testParseTEST_NR()