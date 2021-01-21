#ifndef __VALUE_HH
#define __VALUE_HH

#include <iostream>

namespace Fract::Objects {
/// @brief Value instance.
class value {
public:
  std::string content; /* Content of value. */
  int type; /* Type of value. */
};
} // namespace Fract::Objects

#endif // __VALUE_HH
