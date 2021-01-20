#ifndef __PARSER_HH
#define __PARSER_HH

#include <fstream>
#include <iostream>
#include <iostream>
#include <vector>

#include "tokenizer.hh"
#include "../Objects/code_file.hh"
#include "../Utilities/file_system.hh"

using namespace Fract::Objects;

namespace Fract::Parser {
/// @brief Parser of Fract.
class parser {
private:
  /// @brief Parser of this file.
  code_file file;

  /// @brief Tokenizer of parser.
  tokenizer _tokenizer;

public:
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

  /// @brief Type of file.
  int type;

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
