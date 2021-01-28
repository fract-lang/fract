#ifndef __ARITHMETIC_HH
#define __ARITHMETIC_HH

#include <iostream>
#include <limits>
#include <regex>

#include "../Grammar/tokens.hh"
#include "../Grammar/values.hh"
#include "../Parser/tokenizer.hh"

namespace Fract {
/// @brief Arithmetic processor.
class arithmetic {
public:
  static bool is_types_compatible(int type0, int type1);
  static bool is_integer_type(int type);
  static bool is_float_type(int type);
  static bool is_negative(std::string value);
  static unsigned short to_uint16(std::string value);
  static unsigned int to_uint32(std::string value);
  static unsigned long to_uint64(std::string value);
  static short to_int16(std::string value);
  static int to_int32(std::string value);
  static long to_int64(std::string value);
  static float to_float(std::string value);
  static double to_double(std::string value);
  static bool is_integer_number(std::string value);
  static bool is_floating_number(std::string value);
  static bool is_numeric(char ch);
  static unsigned short to_numeric(char ch);
  static bool bigger(char one, char two);
  static bool lower(char one, char two);
  static bool bigger_str(std::string one, std::string two);
  static bool lower_str(std::string one, std::string two);
  static bool equals_str(std::string one, std::string two);
};
} // namespace Fract

#endif // __ARITHMETIC_HH
