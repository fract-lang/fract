#ifndef __TOKEN_HH
#define __TOKEN_HH

#include <iostream>

namespace Fract::Objects {
/// @brief Token instance.
struct token {
/// @brief Value of token.
std::string value;

/// @brief Type of token.
int type;

/// Line of token.
int line;

/// @brief Column of token.
int column;
};
} // namespace Fract::Objects

#endif // __TOKEN_HH
