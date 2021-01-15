#ifndef __STRING_HH
#define __STRING_HH

#include <iostream>
#include <list>
#include <vector>

namespace Fract::Utilities {
/// @brief Utilities of string.
class string {
  public:
  static std::string trimStart(std::string value);
  static std::string trimEnd(std::string value);
  static std::string trim(std::string value);
  static std::vector<std::string> split(std::string value, char seperator);
  static bool empty(std::string value);
};
}  // namespace Fract::Utilities

#endif  // __STRING_HH
