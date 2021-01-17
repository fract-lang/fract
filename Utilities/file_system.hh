#ifndef __FILE_SYSTEM_HH
#define __FILE_SYSTEM_HH

#include <stdio.h>

#include <fstream>
#include <iostream>
#include <iterator>
#include <vector>

#ifdef _WIN32
#include <direct.h>
#else
#include <unistd.h>
#endif  // _WIN32

namespace Fract::Utilities {
/// @brief Utilities for file system.
class file_system {
  public:
  /// @brief Working directory
  static char* _WORKING_DIR_;

  static bool existFile(std::string path);
  static std::vector<std::string> getLinesOfFile(std::string path);
};
}  // namespace Fract::Utilities

#endif  // __FILE_SYSTEM_HH
