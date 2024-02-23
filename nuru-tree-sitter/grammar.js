module.exports = grammar({
  name: "nuru",

  rules: {
    source_file: $=> seq(
      optional(repeat($.pakeji_tumia_statement)),
      optional(repeat($.declaration_statment)),
      optional(repeat($.pakeji_statement)),
      optional(repeat($.function_statement)),
      optional(repeat($.single_line_comment)),
      optional(repeat($.block_comment)),
      optional(repeat($.expression)),
      optional(repeat($.function_usage_statement)),
    ),

    declaration_statment: $ => choice(
      seq("fanya",
        field("variablename",$.identifier),
          $.equal,$.expression),
      prec(2,seq(field("functionname",$.identifier),
        $.equal,$.function_statement)),
      prec(3,seq("fanya",
        field("functionname",$.identifier),
        $.equal,$.function_statement))
    ),

    pakeji_statement: $ => seq(
      "pakeji", $.identifier, $.block
    ),

    pakeji_tumia_statement: $=> seq(
      "tumia", field("pakejiname",$.identifier)
    ),

    function_statement: $ => seq(
      "unda",$.parameter_list,$.block
    ),

    parameter_list: $ => seq("(",
      repeat(seq($.expression, optional(",")))
      ,")"),

    function_usage_statement: $=> prec(3,seq(
      field("functionname",$.expression),$.parameter_list
    )),

    block: $ => seq("{",
        optional(repeat($.declaration_statment)),
        optional(repeat($.function_statement)),
        optional(repeat($.single_line_comment)),
        optional(repeat($.block_comment)),
        prec(2,optional(repeat($.function_usage_statement))),
        optional(repeat($.expression)),
      "}"
    ),

    single_line_comment: $ => seq(
      "//", $.expression, "\n"
    ),

    block_comment: $ => seq(
      "/*", $.expression, "*/"
    ),

    ending: $=> ";",
    string_expression: $=> seq("\"",/[^\n"]*/,"\""),

    expression: $ => choice(
      $.equal,
      $.identifier,
      $.number,
      $.string_expression
    ),

    identifier: $ => /[a-z]+/,
    number: $ => /\d+/,
    equal: $=> "="
  }
})
