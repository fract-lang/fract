#ifndef __FORMATTER_HH
#define __FORMATTER_HH

#include <iostream>
#include <vector>

#include "rmlx_range_result.hh"
#include "../Include/parser.hh"
#include "../Grammar/tokens.hh"

namespace Fract::Parser {
/// @brief Formatter.
class formatter {
public:
  /**
  * @brief Reeturns range tokens and remove from original vector.
  * @param tokens Tokens to proecess.
  * @returns Tokens of range and replace iterator of original vector.
  */
  static rmlx_range_result rmlx_range(std::vector<token> *tokens);
};
} // namespace Fract::Parser

#endif // __FORMATTER_HH
