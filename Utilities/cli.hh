#ifndef __CLI_HH
#define __CLI_HH

#include <map>

#include "string.hh"
#include "../Objects/color.hh"

namespace Fract::Utilities {
/// @brief CLI utilities.
class cli {
  public:
  /**
   * @brief Print map as table.
   * @param map Map to print.
   */
  static void print_map_as_table(std::map<std::string, std::string> map);
};
}  // namespace Fract::Utilities

#endif  // __CLI_HH
