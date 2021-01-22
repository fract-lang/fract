#ifndef __ARITHMETIC_HH
#define __ARITHMETIC_HH

#include <iostream>
#include "../Grammar/tokens.hh"

namespace Fract {
/// @brief Arithmetic processor.
class arithmetic {
public:
  static bool is_integer_number(std::string value);
  static bool is_floating_number(std::string value);
  static bool is_number(char ch);
  static unsigned short to_numeric(char ch);
  static bool bigger(char one, char two);
  static bool lower(char one, char two);
  static bool bigger(std::string one, std::string two);
  static bool lower(std::string one, std::string two);
  static bool equals(std::string one, std::string two);
};
} // namespace Fract

#endif // __ARITHMETIC_HH
