#ifndef __PARSER_HH
#define __PARSER_HH

#include <fstream>
#include <iostream>
#include <iostream>
#include <vector>

#include "tokenizer.hh"
#include "../Objects/code_file.hh"
#include "../Objects/value.hh"
#include "../Utilities/file_system.hh"

#define ptype_none -1
#define ptype_addition 0
#define ptype_multiplication 1
#define ptype_division 2
#define ptype_subtraction 3

using namespace Fract::Objects;

namespace Fract::Parser {
/// @brief Parser of Fract.
class parser {
private:
  /// @brief Parser of this file.
  code_file file;

  /// @brief Tokenizer of parser.
  tokenizer _tokenizer;

  /**
   * @brief Print value to screen.
   * @param _value Value to print.
  */
  void print_value(value _value);

  /**
   * @brief Process value from tokens.
   * @param tokens Tokens.
   * @param it Last iterator state.
   * @returns Value instance.
  */
  value process_value(std::vector<token> *tokens,
                     std::vector<token>::iterator *it);

public:
  /// @brief Type of file.
  int type;

  /**
   * @brief Create instance of code file.
   * @param path Path of file.
   * @returns Ready file.
  */
  static code_file ready_file(std::string path);

  /**
   * @brief Ready lines to process.
   * @param lines Lines to ready.
   * @returns Ready lines.
  */
  static std::vector<code_line> ready_lines(std::vector<std::string> lines);

  /**
  * @brief Exit as parser styled error.
  * @param _token Token of error.
  * @param message Message of error.
  */
  static void exit_parser_error(token _token, std::string message);

  /**
   * @brief Create new instance.
   * @param path Path of destination file.
   * @param type Type of file.
  */
  parser(std::string path, int type);

  /// @brief Parse code.
  void parse();
};
} // namespace Fract::Parser

#endif // __PARSER_HH
