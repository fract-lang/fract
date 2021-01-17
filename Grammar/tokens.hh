#ifndef __TOKENS_HH
#define __TOKENS_HH

namespace Fract::Grammar {

// GENERIC TOKENS
#define token_sharp "#"
#define token_plus "+"
#define token_minus "-"
#define token_star "*"
#define token_percent "%"
#define token_slash "/"
#define token_reverse_slash "\\"
#define token_equals "="
#define token_question "?"
#define token_veritcal_bar "|"
#define token_great ">"
#define token_less "<"
#define token_semicolon ";"
#define token_colon ":"
#define token_comma ","
#define token_exclamation "!"
#define token_amper "&"
#define token_at "@"
#define token_dot "."

// KEYWORDS
#define kw_import "use"
#define kw_function "fn"
#define kw_delete "del"
#define kw_variable "let"
#define kw_block_final "end"
#define kw_return "ret"
#define kw_for_loop "for"
#define kw_while_loop "while"
#define kw_if "if"
#define kw_else_if "elif"
#define kw_else "else"

// DATA TYPES
#define dt_byte "byte"
#define dt_signed_byte "sbyte"
#define dt_16bit_integer "int16"
#define dt_32bit_integer "int32"
#define dt_64bit_integer "int64"
#define dt_6464bit_integer "int64_64"
#define dt_unsigned_16bit_integer "int16"
#define dt_unsigned_32bit_integer "int32"
#define dt_unsigned_64bit_integer "int64"
#define dt_unsigned_6464bit_integer "int64_64"
#define dt_boolean "bool"
#define dt_float "float"
#define dt_double "double"
}
#endif  // __TOKENS_HH
