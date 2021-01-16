#ifndef __VERSION_HH
#define __VERSION_HH

#include <iostream>

#include "../../Include/fract.hh"

namespace Fract::Shell::Modules {
/// @brief Show version.
class version {
  public:
  /**
   * @brief Process command in module.
   * @param cmd Command.
   */
  static void process(std::string cmd);
};
}  // namespace Fract::Shell::Modules
#endif  // __VERSION_HH
