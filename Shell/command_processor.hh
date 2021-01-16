#ifndef CMD_PROC_HH
#define CMD_PROC_HH

#include <iostream>
#include <regex>
#include <vector>

#include "../Utilities/string.hh"

namespace Fract::Shell {
/// @brief Command processor of shell.
class command_processor {
  public:
  /**
   * @brief Get namespace from command.
   * @param cmd Command.
   * @return Namespace.
   */
  static std::string getNamespace(std::string cmd);

  /**
   * @brief Remove namespace from command.
   * @param cmd Command.
   * @return Command without namespace.
   */
  static std::string removeNamespace(std::string cmd);

  /**
   * @brief Get arguments.
   * @param cmd Command.
   * @param dest Destination vector.
   * @return true if success, false if not.
   */
  static bool getArguments(std::string cmd, std::vector<std::string>* dest);

  /**
   * @brief Remove arguments.
   * @param cmd Command.
   * @return Command without arguments.
   */
  static std::string removeArguments(std::string cmd);
};
}  // namespace Fract::Shell

#endif  // __COMMAND_PROCESSOR_HH
