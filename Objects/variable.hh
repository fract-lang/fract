#ifndef __VARIABLE_HH
#define __VARIABLE_HH

#include <iostream>

namespace Fract::Objects {
/// @brief A variable instance.
class variable {
public:
  std::string name; /* Name of variable. */
  std::string value; /* Value of variable.  */
  int type; /* Type of variable */
};
} // namespace Fract::Objects

#endif // __VARIABLE_HH
