#ifndef __MAKE_HH
#define __MAKE_HH

#include <iostream>

#include "../../Include/fract.hh"
#include "../../Utilities/string.hh"
#include "../../Utilities/file_system.hh"
#include "../../Parser/parser.hh"

namespace Fract::Shell::Modules {
/// @brief Interprete Fract code.
class make {
  public:
  /**
   * @brief Process command in module.
   * @param cmd Command.
   */
  static void process(std::string cmd);

  /**
   * @brief Check is a Fract file.
   * @param value Value to check.
  */
  static bool check(std::string value);
};
}  // namespace Fract::Shell::Modules
#endif  // __MAKE_HH
