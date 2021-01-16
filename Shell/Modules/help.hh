#ifndef __HELP_HH
#define __HELP_HH

#include <iostream>

#include "../../Utilities/cli.hh"

namespace Fract::Shell::Modules {
/// @brief Show help menu.
class help {
  public:
  /**
   * @brief Process command in module.
   * @param cmd Command.
   */
  static void process(std::string cmd);
};
}  // namespace Fract::Shell::Modules
#endif  // __HELP_HH
