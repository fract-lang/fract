#ifndef __EXIT_HH
#define __EXIT_HH

#include <iostream>

namespace Fract::Shell::Modules {
/// @brief Exit.
class exit {
  public:
  /**
   * @brief Process command in module.
   * @param cmd Command.
   */
  static void process(std::string cmd);
};
}  // namespace Fract::Shell::Modules
#endif  // __EXIT_HH
