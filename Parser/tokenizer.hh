#ifndef __TOKENIZER_HH
#define __TOKENIZER_HH

#include <iostream>
#include <vector>

#include "../Include/parser.hh"
#include "../Fract/arithmetic.hh"
#include "../Objects/token.hh"
#include "../Grammar/tokens.hh"
#include "../Utilities/string.hh"
#include "../Objects/code_file.hh"

using namespace Fract::Objects;
using namespace Fract::Utilities;

namespace Fract::Parser {
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
