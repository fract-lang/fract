#ifndef __CODE_LINE_HH
#define __CODE_LINE_HH

#include <iostream>

namespace Fract::Objects {
/// @brief Code line instance.
struct code_line {
  /// @brief This line is x. line.
  int line;
  /// @brief Text of line.
  std::string text;
};
} // namespace Fract::Objects

#endif // __CODE_LINE_HH
