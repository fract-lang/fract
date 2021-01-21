#ifndef __TOKENIZER_HH
#define __TOKENIZER_HH

#include <iostream>
#include <vector>

#include "../Fract/arithmetic.hh"
#include "../Objects/token.hh"
#include "../Grammar/tokens.hh"
#include "../Utilities/string.hh"
#include "../Objects/code_file.hh"

using namespace Fract::Objects;
using namespace Fract::Utilities;

namespace Fract::Parser {

#define type_entry_file 99
#define type_imported_file 100
#define type_comment 999
#define type_function 1000
#define type_equals 1001
#define type_let 1002
#define type_name 1003
#define type_dotted_name 1004
#define type_value_setter 1005
#define type_value 1006
#define type_data_type 1007
#define type_end_type 1008
#define type_return 1009
#define type_import 1010
#define type_std_import 1011
#define type_if 1012
#define type_else_if 1013
#define type_else 1014
#define type_for 1015
#define type_while 1016
#define type_delete 1017
#define type_int16 1018
#define type_int32 1019
#define type_int64 1020
#define type_int64_64 1021
#define type_unsigned_int16 1022
#define type_unsigned_int32 1023
#define type_unsigned_int64 1024
#define type_unsigned_int64_64 1025
#define type_float 1026
#define type_double 1027
#define type_boolean 1028
#define type_byte 1029
#define type_signed_byte 1030


/// @brief Tokenizer of Fract.
class tokenizer {
private:
  code_file *file;
  std::vector<code_line>::iterator line_iterator;
  int column;

public:
  /// @brief Finished tokenizing lines.
  bool finish;

  /// @brief Create new instance.
  tokenizer();

  /**
  * @brief Create new instance.
  * @param file Destination file.
  */
  tokenizer(code_file *file);

  /**
  * @brief Exit as tokenizer styled error.
  * @param message Message of error.
  */
  void exit_tokenizer_error(std::string message);

  /**
  * @brief Tokenize next token from statement.
  * @returns Returns next token.
  */
  token next_token();

  /**
  * @brief Tokenize all statement.
  * @param statement Statement to tokenize.
  * @returns Returns tokens;
  */
  std::vector<token> tokenize_next();
};
} // namespace Fract::Parser

#endif // __TOKENIZER_HH
