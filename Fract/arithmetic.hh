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
};
} // namespace Fract

#endif // __ARITHMETIC_HH
