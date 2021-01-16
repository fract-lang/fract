#ifndef __PARSER_HH
#define __PARSER_HH

#include <iostream>
#include <vector>

#include "../Objects/code_line.hh"

using namespace Fract::Objects;

namespace Fract::Parser {
/// @brief Parser of Fract.
class parser {
public:
  /**
   * @brief Ready lines to process.
   * @param lines Lines to ready.
   * @returns Ready lines.
  */
  static std::vector<code_line> readyLines(std::vector<std::string> lines);
};
} // namespace Fract::Parser

#endif // __PARSER_HH
