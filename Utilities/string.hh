#ifndef __STRING_HH
#define __STRING_HH

#include <iostream>
#include <list>
#include <vector>

namespace Fract::Utilities {
/// @brief Utilities of string.
class string {
  public:
  static std::string trim_start(std::string value);
  static std::string trim_end(std::string value);
  static std::string trim(std::string value);
  static std::vector<std::string> split(std::string value, char seperator);
  static std::string to_lower(std::string value);
  static std::string to_upper(std::string value);
  static bool starts_with(std::string value, std::string start);
  static bool ends_with(std::string value, std::string end);
  static bool contains(std::string value, std::string check);
};
}  // namespace Fract::Utilities

#endif  // __STRING_HH
